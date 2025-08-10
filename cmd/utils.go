package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsInChecksRepo(repoName string) error {
	// Get working directory
	log.Debug("getting current working directory")
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	log.Debug("checking execution is in the root of a git repo")
	gitDir := filepath.Join(cwd, ".git")
	info, err := os.Stat(gitDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("command should be run in the root of the cdm-checks git repo: %v", err)
	}

	log.Debug("checking current git repos remote")
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = cwd
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check git remote: %v", err)
	}

	log.Debug("trimming remote")
	remoteURLSegments := strings.Split(strings.TrimSpace(string(out)), "/")

	log.Debug("checking executon is in cdm-checks git repo")
	if remoteURLSegments[len(remoteURLSegments)-1] != repoName {
		return fmt.Errorf("command was not run in the cdm-checks git repo")
	}

	return nil
}
