// Package logs provides utilities for parsing and processing deployment logs.
// Currently focused on parsing Docker build logs from streams.
package logs

import (
	"bufio"
	"io"
	"strings"
)

// ParseBuildLog reads a stream of build logs and converts it to a single string.
// This is used to capture Docker build output and store it in the database.
// The reader is automatically closed when the function returns.
//
// Parameters:
//   - reader: An io.ReadCloser containing the build log stream (typically from Docker build output)
//
// Returns:
//   - string: All log lines joined with newlines, or empty string on error
//   - error: Error if reading or scanning fails
func ParseBuildLog(reader io.ReadCloser) (string, error) {
	// Ensure the reader is closed when we're done
	defer reader.Close()

	// Store all log lines in a slice
	var logLines []string
	
	// Use a scanner to read line by line (more efficient than reading all at once)
	scanner := bufio.NewScanner(reader)

	// Read each line from the stream
	for scanner.Scan() {
		line := scanner.Text()
		logLines = append(logLines, line)
	}

	// Check for scanning errors (not EOF, which is normal)
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// Join all lines with newline characters to create the full log
	return strings.Join(logLines, "\n"), nil
}

