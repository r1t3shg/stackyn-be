// Package engine provides the core deployment orchestration logic.
// It coordinates the entire deployment pipeline: Git clone -> Docker build -> Container run.
package engine

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"mvp-be/internal/apps"
	"mvp-be/internal/deployments"
	"mvp-be/internal/dockerbuild"
	"mvp-be/internal/dockerrun"
	"mvp-be/internal/gitrepo"
	"mvp-be/internal/logs"
)

// Engine orchestrates the complete deployment pipeline.
// It coordinates all the steps needed to deploy an application:
// 1. Clone the Git repository
// 2. Build a Docker image
// 3. Run a container with Traefik labels
// 4. Track status and logs in the database
type Engine struct {
	// deploymentStore provides database operations for deployments
	deploymentStore *deployments.Store

	// appStore provides database operations for apps
	appStore *apps.Store

	// cloner handles cloning Git repositories
	cloner *gitrepo.Cloner

	// builder handles building Docker images
	builder *dockerbuild.Builder

	// runner handles running Docker containers
	runner *dockerrun.Runner

	// baseDomain is the base domain for subdomain routing (e.g., "localhost" or "example.com")
	baseDomain string
}

// NewEngine creates a new Engine instance with all required dependencies.
//
// Parameters:
//   - deploymentStore: Store for deployment database operations
//   - appStore: Store for app database operations
//   - cloner: Git repository cloner
//   - builder: Docker image builder
//   - runner: Docker container runner
//   - baseDomain: Base domain for subdomain routing
//
// Returns:
//   - *Engine: A new Engine instance ready to process deployments
func NewEngine(
	deploymentStore *deployments.Store,
	appStore *apps.Store,
	cloner *gitrepo.Cloner,
	builder *dockerbuild.Builder,
	runner *dockerrun.Runner,
	baseDomain string,
) *Engine {
	return &Engine{
		deploymentStore: deploymentStore,
		appStore:        appStore,
		cloner:          cloner,
		builder:         builder,
		runner:          runner,
		baseDomain:      baseDomain,
	}
}

// ProcessDeployment executes the complete deployment pipeline for a single deployment.
// It performs all steps sequentially: clone -> build -> run -> update status.
//
// Deployment pipeline:
//   1. Update status to "building"
//   2. Clone the Git repository
//   3. Build Docker image from the repository
//   4. Parse and store build logs
//   5. Run container with Traefik labels
//   6. Update status to "running"
//
// If any step fails, the deployment status is updated to "failed" with an error message.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - deploymentID: The unique identifier of the deployment to process
//
// Returns:
//   - error: Error if any step in the pipeline fails
func (e *Engine) ProcessDeployment(ctx context.Context, deploymentID int) error {
	// Step 0: Load deployment and app information from database
	deployment, err := e.deploymentStore.GetByID(deploymentID)
	if err != nil {
		return fmt.Errorf("failed to get deployment: %w", err)
	}

	// Get the app associated with this deployment
	app, err := e.appStore.GetByID(deployment.AppID)
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err)
	}

	log.Printf("Processing deployment %d for app %s", deploymentID, app.Name)

	// Step 1: Update status to "building" and clone the Git repository
	// This marks the deployment as in-progress
	if err := e.deploymentStore.UpdateStatus(deploymentID, deployments.StatusBuilding); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Clone the repository to a local directory
	// The directory will be: {WorkDir}/deployment-{deploymentID}
	repoPath, err := e.cloner.Clone(app.RepoURL, deploymentID)
	if err != nil {
		// Record the error in the database before returning
		e.deploymentStore.UpdateError(deploymentID, fmt.Sprintf("Git clone failed: %v", err))
		return fmt.Errorf("git clone failed: %w", err)
	}

	// Step 2: Build Docker image from the cloned repository
	// Image name format: mvp-{app-name}:{deployment-id}
	imageName := fmt.Sprintf("mvp-%s:%d", strings.ToLower(app.Name), deploymentID)
	builtImage, buildLogReader, err := e.builder.Build(ctx, repoPath, imageName)
	if err != nil {
		// Record the error in the database before returning
		e.deploymentStore.UpdateError(deploymentID, fmt.Sprintf("Docker build failed: %v", err))
		return fmt.Errorf("docker build failed: %w", err)
	}

	// Parse the Docker build output and store it in the database
	// This allows users to see what happened during the build
	buildLog, err := logs.ParseBuildLog(buildLogReader)
	if err != nil {
		log.Printf("Warning: failed to parse build log: %v", err)
	} else {
		// Store the build log (non-critical, so we log warnings but don't fail)
		if err := e.deploymentStore.UpdateBuildLog(deploymentID, buildLog); err != nil {
			log.Printf("Warning: failed to update build log: %v", err)
		}
	}

	// Update the deployment record with the image name
	if err := e.deploymentStore.UpdateImage(deploymentID, builtImage); err != nil {
		return fmt.Errorf("failed to update image name: %w", err)
	}

	// Step 3: Run the container with Traefik routing labels
	// Subdomain format: {app-name}-{deployment-id}
	subdomain := fmt.Sprintf("%s-%d", strings.ToLower(app.Name), deploymentID)
	containerID, err := e.runner.Run(ctx, builtImage, subdomain, e.baseDomain)
	if err != nil {
		// Record the error in the database before returning
		e.deploymentStore.UpdateError(deploymentID, fmt.Sprintf("Container run failed: %v", err))
		return fmt.Errorf("container run failed: %w", err)
	}

	// Update the deployment record with container information
	if err := e.deploymentStore.UpdateContainer(deploymentID, containerID, subdomain); err != nil {
		return fmt.Errorf("failed to update container info: %w", err)
	}

	// Step 4: Mark deployment as successfully running
	if err := e.deploymentStore.UpdateStatus(deploymentID, deployments.StatusRunning); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	log.Printf("Deployment %d completed successfully. Container: %s, Subdomain: %s.%s",
		deploymentID, containerID, subdomain, e.baseDomain)

	return nil
}

// RunLoop runs the deployment engine in a continuous loop.
// It polls for pending deployments and processes them one by one.
// The loop continues until the context is cancelled (e.g., on application shutdown).
//
// Polling behavior:
//   - Fetches all pending deployments from the database
//   - Processes each deployment sequentially
//   - Waits 2 seconds before checking for new pending deployments
//   - Responds to context cancellation for graceful shutdown
//
// This is a simple polling implementation. In production, you might want to use:
//   - Database triggers/notifications
//   - Message queues (RabbitMQ, Redis, etc.)
//   - Webhooks from the API
//
// Parameters:
//   - ctx: Context for cancellation control (when cancelled, the loop stops gracefully)
func (e *Engine) RunLoop(ctx context.Context) {
	log.Println("Deployment engine started")

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled (e.g., SIGTERM signal received)
			log.Println("Deployment engine stopped")
			return
		default:
			// Get all pending deployments from the database
			// These are deployments that need to be processed
			pending, err := e.deploymentStore.GetPending()
			if err != nil {
				log.Printf("Error fetching pending deployments: %v", err)
				// Continue to next iteration instead of exiting
				// This makes the engine resilient to temporary database issues
				continue
			}

			// Process each pending deployment
			// They are processed sequentially (one at a time) to avoid resource contention
			for _, deployment := range pending {
				// Process the deployment (clone, build, run)
				if err := e.ProcessDeployment(ctx, deployment.ID); err != nil {
					// Log the error but continue processing other deployments
					// The error has already been recorded in the database by ProcessDeployment
					log.Printf("Error processing deployment %d: %v", deployment.ID, err)
				}
			}

			// Wait before checking for new pending deployments again
			// This prevents a busy loop and reduces database load
			select {
			case <-ctx.Done():
				// Check if context was cancelled during the wait
				return
			case <-time.After(2 * time.Second):
				// Wait 2 seconds before next poll
				// This is a simple polling interval - adjust based on your needs
			}
		}
	}
}

