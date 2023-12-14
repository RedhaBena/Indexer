package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"RedhaBena/indexer/internal/config"
	"RedhaBena/indexer/internal/logger"
	"RedhaBena/indexer/internal/services"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	blockExit := false

	ctx := context.Background()
	// Logger
	if err := logger.InitGlobalLogger(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create global logger: %v", err)
		os.Exit(1)
	}

	defer func() {
		if rerr := recover(); rerr != nil {
			logger.GlobalLogger.Error("Fatal panic", zap.String("panic", fmt.Sprintf("%+v", rerr)))
		}
	}()

	// Context
	ctx, cancelFn := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			logger.GlobalLogger.Warn("Cancel signal triggered")
			cancelFn()
		}
	}()

	// CLI
	app := &cli.App{
		Name:  "indexer",
		Usage: "populate a Neo4j database with json file data",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       config.DefaultDatabaseHost,
				Usage:       "Database Host",
				Category:    "Neo4j Database",
				EnvVars:     []string{"DATABASE_HOST"},
				Destination: &config.GlobalConfig.DatabaseConfig.Host,
			},
			&cli.StringFlag{
				Name:        "user",
				Value:       config.DefaultDatabaseUser,
				Usage:       "Database User",
				Category:    "Neo4j Database",
				EnvVars:     []string{"DATABASE_USER"},
				Destination: &config.GlobalConfig.DatabaseConfig.User,
			},
			&cli.StringFlag{
				Name:        "pass",
				Value:       config.DefaultDatabasePass,
				Usage:       "Database Pass",
				Category:    "Neo4j Database",
				EnvVars:     []string{"DATABASE_PASS"},
				Destination: &config.GlobalConfig.DatabaseConfig.Pass,
			},
			&cli.StringFlag{
				Name:        "file",
				Value:       config.DefaultFilePath,
				Usage:       "File Path",
				Category:    "File",
				EnvVars:     []string{"FILE_PATH"},
				Destination: &config.GlobalConfig.FileConfig.LocalPath,
			},
			&cli.StringFlag{
				Name:        "download-file",
				Usage:       "7z archive url",
				Category:    "File",
				EnvVars:     []string{"DOWNLOAD_FILE"},
				Destination: &config.GlobalConfig.FileConfig.DownloadPath,
			},
			&cli.UintFlag{
				Name:        "size",
				Value:       config.DefaultBatchSize,
				Usage:       "Batch Size",
				EnvVars:     []string{"BATCH_SIZE"},
				Destination: &config.GlobalConfig.IndexerConfig.BatchSize,
			},
			&cli.BoolFlag{
				Value:       false,
				EnvVars:     []string{"BLOCK_EXIT"},
				Destination: &blockExit,
			},
		},
		Before: func(cCtx *cli.Context) error {
			download_url := cCtx.String("download-file")
			if download_url != "" {
				err := services.Client.DownloadArchive(download_url)
				if err != nil {
					return fmt.Errorf("Failed to download archive: %v", err)
				}
				if err := services.Client.Extract7ZArchive(filepath.Base(download_url)); err != nil {
					return fmt.Errorf("Failed to extract archive: %v", err)
				}

				config.GlobalConfig.FileConfig.LocalPath = "dblpv13.json" // TODO: hardcoded value, get it from extracted file.
			}

			if err := services.InitServiceClient(cCtx.Context); err != nil {
				return fmt.Errorf("Failed to create service client: %v", err)
			}

			return nil
		},
		Action: func(cCtx *cli.Context) error {
			services.Client.ReadFile(cCtx.Context)
			return nil
		},
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		logger.GlobalLogger.Fatal("Fatal error", zap.Error(err))
	}

	if blockExit {
		for {
			time.Sleep(time.Second)
		}
	}

	logger.GlobalLogger.Info("Exiting")
}
