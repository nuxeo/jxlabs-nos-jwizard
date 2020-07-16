package create

import (
	"fmt"
	"os"

	"github.com/jenkins-x-labs/jwizard/pkg/cmd/common"
	"github.com/jenkins-x-labs/jwizard/pkg/cmd/importcmd"
	"github.com/jenkins-x/jx/v2/pkg/cmd/create/options"

	"github.com/jenkins-x/jx/v2/pkg/cmd/helper"

	"github.com/jenkins-x/jx/v2/pkg/cmd/opts"
	"github.com/jenkins-x/jx/v2/pkg/cmd/templates"
	"github.com/jenkins-x/jx/v2/pkg/util"
	"github.com/spf13/cobra"
)

const (
	createQuickstartName = "Create new application from a Quickstart"
	createSpringName     = "Create new Spring Boot microservice"
	importDirName        = "Import existing code from a directory"
	importGitName        = "Import code from a git repository"
	importGitHubName     = "Import code from a github repository"
)

var (
	createProjectNames = []string{
		createQuickstartName,
		createSpringName,
		importDirName,
		importGitName,
		importGitHubName,
	}

	createProjectLong = templates.LongDesc(`
		Create a new project by importing code, creating a quickstart or custom wizard for spring.

`)

	createProjectExample = templates.Examples(`
		# Create a project using the wizard
		%s
	`)
)

// CreateProjectOptions contains the command line options
type CreateProjectOptions struct {
	importcmd.ImportOptions

	OutDir             string
	DisableImport      bool
	GithubAppInstalled bool
}

// CreateProjectWizardOptions the options for the command
type CreateProjectWizardOptions struct {
	options.CreateOptions
}

// NewCmdCreateProject creates a command object for the "create" command
func NewCmdCreateProject(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &CreateProjectWizardOptions{
		CreateOptions: options.CreateOptions{
			CommonOptions: commonOpts,
		},
	}

	cmd := &cobra.Command{
		Use:     "project",
		Short:   "Create a new project by importing code, creating a quickstart or custom wizard for spring",
		Long:    createProjectLong,
		Example: fmt.Sprintf(createProjectExample, common.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			setLoggingLevel(cmd)
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}

	cmd.AddCommand(NewCmdCreateQuickstart(commonOpts))
	cmd.AddCommand(NewCmdCreateSpring(commonOpts))
	cmd.AddCommand(importcmd.NewCmdImport(commonOpts))

	return cmd
}

// Run implements the command
func (o *CreateProjectWizardOptions) Run() error {
	name, err := util.PickName(createProjectNames, "Which kind of project you want to create: ",
		"there are a number of different wizards for creating or importing new projects.",
		o.GetIOFileHandles())
	if err != nil {
		return err
	}
	switch name {
	case createQuickstartName:
		return o.createQuickstart()
	case createSpringName:
		return o.createSpring()
	case importDirName:
		return o.importDir()
	case importGitName:
		return o.importGit()
	case importGitHubName:
		return o.importGithubProject()
	default:
		return fmt.Errorf("Unknown selection: %s\n", name)
	}
}

func (o *CreateProjectWizardOptions) createQuickstart() error {
	w := &CreateQuickstartOptions{}
	w.CommonOptions = o.CommonOptions
	return w.Run()
}

func (o *CreateProjectWizardOptions) createSpring() error {
	w := &CreateSpringOptions{}
	w.CommonOptions = o.CommonOptions
	return w.Run()
}

func (o *CreateProjectWizardOptions) importDir() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dir, err := util.PickValue("Which directory contains the source code: ", wd, true,
		"Please specify the directory which contains the source code you want to use for your new project", o.GetIOFileHandles())
	if err != nil {
		return err
	}
	w := &importcmd.ImportOptions{
		Dir: dir,
	}
	w.CommonOptions = o.CommonOptions
	return w.Run()
}

func (o *CreateProjectWizardOptions) importGit() error {
	repoUrl, err := util.PickValue("Which git repository URL to import: ", "", true,
		"Please specify the git URL which contains the source code you want to use for your new project", o.GetIOFileHandles())
	if err != nil {
		return err
	}

	w := &importcmd.ImportOptions{
		RepoURL: repoUrl,
	}
	w.CommonOptions = o.CommonOptions
	return w.Run()
}

func (o *CreateProjectWizardOptions) importGithubProject() error {
	w := &importcmd.ImportOptions{
		GitHub: true,
	}
	w.CommonOptions = o.CommonOptions
	return w.Run()
}

// DoImport imports the project created at the given directory
func (o *CreateProjectOptions) ImportCreatedProject(outDir string) error {
	if o.DisableImport {
		return nil
	}
	importOptions := &o.ImportOptions
	importOptions.Dir = outDir
	importOptions.DisableDotGitSearch = true
	importOptions.GithubAppInstalled = o.GithubAppInstalled
	return importOptions.Run()
}

func (o *CreateProjectOptions) addCreateAppFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.DisableImport, "no-import", "", false, "Disable import after the creation")
	cmd.Flags().StringVarP(&o.OutDir, opts.OptionOutputDir, "o", "", "Directory to output the project to. Defaults to the current directory")

	o.AddImportFlags(cmd, true)
}
