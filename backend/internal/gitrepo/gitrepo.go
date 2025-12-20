// Package gitrepo provides functionality for cloning Git repositories.
// It handles the Git clone operation needed before building Docker images.
package gitrepo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Cloner handles cloning Git repositories to the local filesystem.
// Each deployment gets its own directory to avoid conflicts.
type Cloner struct {
	// WorkDir is the base directory where repositories will be cloned.
	// Each deployment gets a subdirectory named "deployment-{id}" within this directory.
	WorkDir string
}

// NewCloner creates a new Cloner instance with the specified working directory.
//
// Parameters:
//   - workDir: The base directory path where repositories should be cloned
//     (e.g., "/tmp/mvp-deployments")
//
// Returns:
//   - *Cloner: A new Cloner instance ready to clone repositories
func NewCloner(workDir string) *Cloner {
	return &Cloner{WorkDir: workDir}
}

// Clone clones a Git repository to a deployment-specific directory.
// It first removes any existing directory with the same name to ensure a clean clone.
//
// The cloned repository will be located at: {WorkDir}/deployment-{deploymentID}
//
// Parameters:
//   - repoURL: The Git repository URL to clone (e.g., "https://github.com/user/repo.git")
//   - deploymentID: The unique deployment ID used to create a unique directory name
//
// Returns:
//   - string: The absolute path to the cloned repository directory, or empty string on error
//   - error: Error if directory cleanup fails, git command fails, or repository is inaccessible
func (c *Cloner) Clone(repoURL string, deploymentID int) (string, error) {
	// Create a unique directory name for this deployment
	// Format: deployment-{id} (e.g., "deployment-123")
	repoDir := filepath.Join(c.WorkDir, fmt.Sprintf("deployment-%d", deploymentID))

	// Remove directory if it exists to ensure a clean clone
	// This handles cases where a previous deployment with the same ID failed
	if err := os.RemoveAll(repoDir); err != nil {
		return "", fmt.Errorf("failed to clean directory: %w", err)
	}

	// Execute git clone command
	// This will clone the repository into the repoDir directory
	cmd := exec.Command("git", "clone", repoURL, repoDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Include the git output in the error for debugging
		return "", fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	// Return the path to the cloned repository
	return repoDir, nil
}

// CheckDockefile checks if a Dockerfile exists in the repository directory
func CheckDockerfile(repoPath string) error {
	dockerfilePath := filepath.Join(repoPath, "Dockerfile")

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("Dockerfile not found in repository root directory")
	}

	return nil
}
