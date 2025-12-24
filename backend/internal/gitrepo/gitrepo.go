// Package gitrepo provides Git repository cloning functionality.
// It handles cloning repositories from Git URLs with support for:
//   - Specific branch selection
//   - Shallow cloning (depth=1) for faster operations
//   - Dockerfile validation
//
// The cloner creates isolated directories for each deployment
// to avoid conflicts between concurrent deployments.
package gitrepo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Cloner struct {
	WorkDir string
}

func NewCloner(workDir string) *Cloner {
	return &Cloner{WorkDir: workDir}
}

func (c *Cloner) Clone(repoURL string, deploymentID int, branch string) (string, error) {
	repoDir := filepath.Join(c.WorkDir, fmt.Sprintf("deployment-%d", deploymentID))
	log.Printf("[GIT] Cloning repository - URL: %s, Branch: %s, Target: %s", repoURL, branch, repoDir)

	// Remove directory if it exists
	if err := os.RemoveAll(repoDir); err != nil {
		log.Printf("[GIT] ERROR - Failed to clean directory %s: %v", repoDir, err)
		return "", fmt.Errorf("failed to clean directory: %w", err)
	}

	// Clone repository with specific branch
	// First clone the repository (shallow clone for the specific branch)
	log.Printf("[GIT] Executing: git clone --branch %s --single-branch --depth 1 %s %s", branch, repoURL, repoDir)
	cmd := exec.Command("git", "clone", "--branch", branch, "--single-branch", "--depth", "1", repoURL, repoDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[GIT] ERROR - Clone failed: %v, Output: %s", err, string(output))
		return "", fmt.Errorf("git clone failed: %w, output: %s", err, string(output))
	}

	log.Printf("[GIT] Repository cloned successfully to: %s", repoDir)
	return repoDir, nil
}

// CheckDockerfile checks if a Dockerfile exists in the repository directory
func CheckDockerfile(repoPath string) error {
	dockerfilePath := filepath.Join(repoPath, "Dockerfile")
	log.Printf("[GIT] Checking for Dockerfile at: %s", dockerfilePath)

	// Check if Dockerfile exists
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		log.Printf("[GIT] ERROR - Dockerfile not found at: %s", dockerfilePath)
		return fmt.Errorf("dockerfile not found in repository root directory")
	}

	log.Printf("[GIT] Dockerfile found successfully")
	return nil
}

// EnsurePackageLock handles the case where package.json exists but package-lock.json doesn't.
// This fixes the common issue where Dockerfiles use `npm ci` but the lock file is missing.
// It tries two approaches:
//   1. First, try to generate package-lock.json using npm (if Node.js is available)
//   2. If that fails, modify the Dockerfile to use `npm install` instead of `npm ci`
func EnsurePackageLock(repoPath string) error {
	packageJSONPath := filepath.Join(repoPath, "package.json")
	packageLockPath := filepath.Join(repoPath, "package-lock.json")
	dockerfilePath := filepath.Join(repoPath, "Dockerfile")

	// Check if package.json exists
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		log.Printf("[GIT] No package.json found, skipping package-lock.json check")
		return nil
	}

	// Check if package-lock.json already exists
	if _, err := os.Stat(packageLockPath); err == nil {
		log.Printf("[GIT] package-lock.json already exists, no action needed")
		return nil
	}

	log.Printf("[GIT] package.json found but package-lock.json missing")

	// Try to generate package-lock.json using npm (if Node.js is available)
	if err := generatePackageLock(repoPath); err == nil {
		log.Printf("[GIT] package-lock.json generated successfully using npm")
		return nil
	}

	log.Printf("[GIT] Could not generate package-lock.json, modifying Dockerfile to use 'npm install' instead of 'npm ci'")
	
	// Fallback: modify Dockerfile to use npm install instead of npm ci
	return fixDockerfileNpmCi(repoPath, dockerfilePath)
}

// generatePackageLock attempts to generate package-lock.json using npm
func generatePackageLock(repoPath string) error {
	// Check if npm is available
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("npm not found in PATH")
	}

	log.Printf("[GIT] Attempting to generate package-lock.json using npm...")
	cmd := exec.Command("npm", "install", "--package-lock-only")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed: %w, output: %s", err, string(output))
	}
	return nil
}

// fixDockerfileNpmCi modifies the Dockerfile to replace `npm ci` with `npm install`
// when package-lock.json is missing
func fixDockerfileNpmCi(repoPath, dockerfilePath string) error {
	// Read Dockerfile
	file, err := os.Open(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	modified := false

	for scanner.Scan() {
		line := scanner.Text()
		// Check if line contains `npm ci` (case-insensitive, handles variations)
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "npm ci") || strings.Contains(lowerLine, "npmci") {
			// Replace npm ci with npm install
			// Preserve the original formatting and any flags
			originalLine := line
			line = strings.ReplaceAll(line, "npm ci", "npm install")
			line = strings.ReplaceAll(line, "npmci", "npm install")
			line = strings.ReplaceAll(line, "npm  ci", "npm install")
			// Also handle case variations
			line = strings.ReplaceAll(line, "NPM CI", "npm install")
			line = strings.ReplaceAll(line, "Npm Ci", "npm install")
			
			if line != originalLine {
				log.Printf("[GIT] Modified Dockerfile line: %s -> %s", originalLine, line)
				modified = true
			}
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read Dockerfile: %w", err)
	}

	if !modified {
		log.Printf("[GIT] Dockerfile does not contain 'npm ci', no modification needed")
		return nil
	}

	// Write modified Dockerfile
	if err := os.WriteFile(dockerfilePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write modified Dockerfile: %w", err)
	}

	log.Printf("[GIT] Dockerfile modified successfully to use 'npm install' instead of 'npm ci'")
	return nil
}

// DetectPortFromDockerfile attempts to detect the port from the Dockerfile's EXPOSE directive.
// Also checks for common port patterns in CMD/RUN directives (e.g., uvicorn --port, flask run --port).
// Returns the first port found, or attempts to detect from common patterns, or 8080 as default.
func DetectPortFromDockerfile(repoPath string) int {
	dockerfilePath := filepath.Join(repoPath, "Dockerfile")
	
	file, err := os.Open(dockerfilePath)
	if err != nil {
		log.Printf("[GIT] WARNING - Failed to open Dockerfile for port detection: %v, using default port 8080", err)
		return 8080
	}
	defer file.Close()

	// Regex patterns for port detection
	exposeRegex := regexp.MustCompile(`(?i)^\s*EXPOSE\s+(\d+)`)
	// Common Python patterns: uvicorn --port 8000, gunicorn -b :8000, flask run --port 5000
	uvicornRegex := regexp.MustCompile(`(?i)uvicorn.*--port\s+(\d+)`)
	gunicornRegex := regexp.MustCompile(`(?i)gunicorn.*-b\s+[:\d.]*:(\d+)`)
	flaskRegex := regexp.MustCompile(`(?i)flask.*--port\s+(\d+)`)
	// Generic port patterns: -p 8000, :8000, port=8000
	genericPortRegex := regexp.MustCompile(`(?i)(?:-p|--port|port=)[\s:=]+(\d+)`)
	
	scanner := bufio.NewScanner(file)
	var detectedPort int
	foundExpose := false
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// First, check for EXPOSE directive (highest priority)
		matches := exposeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			port, err := strconv.Atoi(matches[1])
			if err == nil && port > 0 && port < 65536 {
				log.Printf("[GIT] Detected port %d from Dockerfile EXPOSE directive", port)
				return port
			}
		}
		
		// If no EXPOSE found yet, check for common Python patterns
		if !foundExpose {
			// Check for uvicorn (FastAPI)
			if matches := uvicornRegex.FindStringSubmatch(line); len(matches) > 1 {
				if port, err := strconv.Atoi(matches[1]); err == nil && port > 0 && port < 65536 {
					detectedPort = port
					log.Printf("[GIT] Detected port %d from uvicorn command", port)
				}
			}
			// Check for gunicorn
			if matches := gunicornRegex.FindStringSubmatch(line); len(matches) > 1 {
				if port, err := strconv.Atoi(matches[1]); err == nil && port > 0 && port < 65536 {
					detectedPort = port
					log.Printf("[GIT] Detected port %d from gunicorn command", port)
				}
			}
			// Check for flask
			if matches := flaskRegex.FindStringSubmatch(line); len(matches) > 1 {
				if port, err := strconv.Atoi(matches[1]); err == nil && port > 0 && port < 65536 {
					detectedPort = port
					log.Printf("[GIT] Detected port %d from flask command", port)
				}
			}
			// Check for generic port patterns
			if matches := genericPortRegex.FindStringSubmatch(line); len(matches) > 1 {
				if port, err := strconv.Atoi(matches[1]); err == nil && port > 0 && port < 65536 {
					detectedPort = port
					log.Printf("[GIT] Detected port %d from generic port pattern", port)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[GIT] WARNING - Error reading Dockerfile: %v, using default port 8080", err)
		return 8080
	}

	// If we detected a port from patterns, use it
	if detectedPort > 0 {
		log.Printf("[GIT] Using detected port %d from Dockerfile patterns", detectedPort)
		return detectedPort
	}

	log.Printf("[GIT] No EXPOSE directive or port patterns found in Dockerfile, using default port 8080")
	return 8080
}
