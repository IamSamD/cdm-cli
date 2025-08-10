package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

var (
	provider          string
	resource          string
	checkName         string
	checkTemplateRepo string = "https://github.com/iamsamd/cdm_check_template"
)

// addCheckCmd represents the addCheck command
var addCheckCmd = &cobra.Command{
	Use:   "addCheck",
	Short: "Generate a template for a new check",
	Long:  ``,
	Run:   addCheck,
}

func init() {
	rootCmd.AddCommand(addCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCheckCmd.Flags().StringVarP(&provider, "provider", "p", "", "top level directory for the check, usually the cloud provider the check relates to (azure | aws | gcp | general) etc")
	addCheckCmd.Flags().StringVarP(&resource, "resource", "r", "", "second level directory for the check, usually the cloud resource the check relates to (aks | eks | cloud_compser | source_control) etc")
	addCheckCmd.Flags().StringVarP(&checkName, "check-name", "n", "", "third level directory for the check and the name of the check (upgrade | certificate_expiry | backup_check) etc")

}

func addCheck(cmd *cobra.Command, args []string) {
	// Check we are in the root dir of the cdm-checks repo
	if err := IsInChecksRepo("cdm-checks.git"); err != nil {
		log.Error(err.Error())
		return
	}

	log.Debug("checks repo validation successful")

	// Build filepath for new check
	checkPath := filepath.Join(provider, resource, checkName)

	// Make the directory structure
	if err := os.MkdirAll(checkPath, 0755); err != nil {
		log.Error(fmt.Sprintf("failed to create check directories: %v", err))
		return
	}

	log.Info("check directories created")

	// clone check template
	_, err := git.PlainClone(checkPath, &git.CloneOptions{
		URL:      checkTemplateRepo,
		Progress: nil,
		Depth:    1,
	})
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Run files through template engine and replace values for check name
	mainFile := filepath.Join(checkPath, "main.go.tmpl")
	mainFileOut := filepath.Join(checkPath, "main.go")

	modFile := filepath.Join(checkPath, "go.mod.tmpl")
	modFileOut := filepath.Join(checkPath, "go.mod")

	files := []struct {
		TemplatePath string
		OutputPath   string
	}{
		{mainFile, mainFileOut},
		{modFile, modFileOut},
	}

	data := struct {
		Module string
	}{
		Module: checkName,
	}

	for _, file := range files {
		tmpl, err := template.ParseFiles(file.TemplatePath)
		if err != nil {
			log.Error(err.Error())
			return
		}

		outFile, err := os.Create(file.OutputPath)
		if err != nil {
			log.Error(err.Error())
		}
		defer outFile.Close()

		if err := tmpl.Execute(outFile, data); err != nil {
			log.Error(err.Error())
			return
		}

		if err := os.Remove(file.TemplatePath); err != nil {
			log.Error(err.Error())
			return
		}
	}
}
