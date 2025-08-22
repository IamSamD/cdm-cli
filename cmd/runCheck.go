/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// runCheckCmd represents the runCheck command
var runCheckCmd = &cobra.Command{
	Use:   "runCheck",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runCheck,
}

func init() {
	rootCmd.AddCommand(runCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCheck(cmd *cobra.Command, args []string) {
	// Load env vars from .env file if it exists
	godotenv.Load()

	// Parse orchestrator level env vars
	config, err := parseOrchestratorLevelEnvVars()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check if we need to skip the check
	skipCheck, err := isSkipCheck()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if skipCheck {
		fmt.Printf("Skipping check until: %v\n", os.Getenv("SKIP_UNTIL"))
		fmt.Println("##vso[task.logissue type=warning]Check Skipped")
		os.Exit(0)
	}

	// Download the check
	if err := downloadCheck(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Run the check
	if err := runCheckBinary(config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func runCheckBinary(config map[string]string) error {
	cmd := exec.Command(fmt.Sprintf("./%s", config["BINARY_NAME"]))

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run check: %s : %v", output, err)
	}

	// Print the output
	fmt.Println(string(output))
	return nil
}

func downloadCheck(config map[string]string) error {
	// Replace check name for binary name

	url := fmt.Sprintf("https://github.com/IamSamD/cdm-checks/releases/download/%s_%s/%s", config["BINARY_NAME"], config["CHECK_VERSION"], config["BINARY_NAME"])

	// Download the check binary
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download check binary: %v", err)
	}
	defer resp.Body.Close()

	binary, err := os.Create(config["BINARY_NAME"])
	if err != nil {
		return fmt.Errorf("failed to create binary file: %v", err)
	}

	_, err = io.Copy(binary, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body to binary file: %v", err)
	}

	binary.Close()

	err = os.Chmod(config["BINARY_NAME"], 0755)
	if err != nil {
		return fmt.Errorf("failed to set permissions on binary file: %v", err)
	}

	return nil
}

func parseOrchestratorLevelEnvVars() (map[string]string, error) {
	if os.Getenv("CHECK_NAME") == "" || os.Getenv("CHECK_VERSION") == "" {
		return nil, fmt.Errorf("CHECK_NAME and CHECK_VERSION must be set")
	}

	checkName := os.Getenv("CHECK_NAME")
	checkVersion := os.Getenv("CHECK_VERSION")
	binaryName := strings.ReplaceAll(checkName, "/", "_")

	return map[string]string{
		"CHECK_NAME":    checkName,
		"CHECK_VERSION": checkVersion,
		"BINARY_NAME":   binaryName,
	}, nil
}

func isSkipCheck() (bool, error) {
	// Check the skip
	if os.Getenv("SKIP_UNTIL") != "" {
		timeFormat := "02/01/2006 15:04" // Go's time formatting string
		skipUntilDatetime, err := time.Parse(timeFormat, os.Getenv("SKIP_UNTIL"))
		if err != nil {
			return false, fmt.Errorf(fmt.Sprintf("failed to parse skip until datetime string: %v", err))
		}

		// Check if we need to skip the check
		if skipUntilDatetime.After(time.Now()) {
			return true, nil
		}
	}
	return false, nil
}
