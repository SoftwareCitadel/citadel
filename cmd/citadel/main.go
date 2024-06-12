package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
)

const (
	version = "0.1.30"
)

func main() {
	log.SetOutput(io.Discard)

	if err := update(); err != nil {
		fmt.Printf("error occurred while updating binary: %v\n", err)
		os.Exit(1)
	}

	execute(version)
}

func update() error {
	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("softwarecitadel/cli"))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if latest.Version() == version {
		return nil
	}

	fmt.Printf("Current version (%s) is not the latest\nUpdating to version %s...\n", version, latest.Version())

	exe, err := os.Executable()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	fmt.Printf("Successfully updated to version %s.\n", latest.Version())

	return nil
}
