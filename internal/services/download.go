package services

import (
	"RedhaBena/indexer/internal/logger"
	"RedhaBena/indexer/internal/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// DownloadArchive downloads an archive from the given URL and saves it to the local file.
func (client *ServiceClient) DownloadArchive(url string) error {
	logger.GlobalLogger.Debug("Downloading archive", zap.String("url", url))

	// Download the 7z archive
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download archive: %d", response.StatusCode)
	}

	// Extract the filename from the URL
	filename := filepath.Base(url)

	// Create or open the local file for writing
	archiveFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// Get the total size of the file
	totalSize := response.ContentLength
	if totalSize > 0 {
		// If the total size is known, log progress updates every 10 seconds
		progressInterval := 10 * time.Second
		var bytesRead int64
		buffer := make([]byte, 1024) // Adjust buffer size as needed

		// Create a ticker that ticks every 10 seconds
		ticker := time.NewTicker(progressInterval)
		defer ticker.Stop()

		// Create a goroutine to log progress
		go func() {
			for range ticker.C {
				logger.GlobalLogger.Info("Download progress",
					zap.String("url", url),
					zap.String("progress", utils.FormatBytes(bytesRead)),
					zap.String("size", utils.FormatBytes(totalSize)),
				)
			}
		}()

		// Copy the response body to the local file with progress updates
		for {
			n, err := response.Body.Read(buffer)
			if n > 0 {
				bytesRead += int64(n)
				_, err := archiveFile.Write(buffer[:n])
				if err != nil {
					return err
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}

		// Stop the progress logging goroutine
		ticker.Stop()
	}

	logger.GlobalLogger.Info("Download completed",
		zap.String("url", url),
		zap.String("size", utils.FormatBytes(totalSize)),
	)

	return nil
}

func (client *ServiceClient) Extract7ZArchive(path string) error {
	logger.GlobalLogger.Info("Extracting archive", zap.String("archive", path))

	cmd := exec.Command("7z", "x", "-aoa", path)

	// Capture the standard error before starting the command
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.GlobalLogger.Error("Failed to get stderr pipe", zap.Error(err))
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.GlobalLogger.Error("Failed to start command", zap.Error(err))
		return err
	}

	// Read and log the standard error
	errOutput, _ := io.ReadAll(stderr)
	if len(errOutput) > 0 {
		logger.GlobalLogger.Error("7z command stderr", zap.String("stderr", string(errOutput)))
	}

	// Wait for the command to finish
	err = cmd.Wait()
	// Check if the exit status is zero (success)
	if err != nil {
		logger.GlobalLogger.Error("7z command failed", zap.Error(err))
		return err
	}

	logger.GlobalLogger.Info("Archive uncompressed successfully")
	return nil
}
