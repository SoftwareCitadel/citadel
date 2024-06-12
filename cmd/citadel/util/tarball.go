package util

import (
	"bufio"
	"io"
	"os"

	"github.com/docker/docker/pkg/archive"
)

func MakeTarball() (io.ReadCloser, error) {
	paths := listIgnorePaths()

	ar, err := archive.TarWithOptions(".", &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: paths,
	})
	if err != nil {
		return nil, err
	}

	return ar, nil
}

func listIgnorePaths() []string {
	paths := []string{
		"./.git",
		"./node_modules",
		".git/*",
		".git/**/*",
		"node_modules/*",
		"node_modules/**/*",
	}

	readFile, err := os.Open(".dockerignore")
	if err != nil {
		return paths
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		paths = append(paths, fileScanner.Text())
	}

	readFile.Close()

	return paths
}
