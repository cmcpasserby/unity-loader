package main

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"path/filepath"
)

func createRunCmd() *cobra.Command {
	var ( // local flags
		lFlagForceVersion bool
		lFlagBuildTarget  string
	)

	cmd := &cobra.Command{
		Use:   "run [projectDirectory]",
		Short: "Launches unity and opens the selected project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string

			if len(args) > 0 {
				path = args[0]
			} else {
				wd, err := os.Getwd()
				if err != nil {
					return err
				}
				path = wd
			}

			expandedPath, _ := filepath.Abs(path)

			var version string
			var err error

			if lFlagForceVersion { // TODO use this flow or direct version number
				installs, err := unity.GetInstalls()
				if err != nil {
					return err
				}

				options := make([]string, 0, len(installs))
				for _, install := range installs {
					options = append(options, install.Version.String())
				}

				prompt := &survey.Select{
					Message:  "Select Unity version to run project with",
					Options:  options,
					PageSize: 10,
				}

				if err := survey.AskOne(prompt, &version, nil); err != nil {
					return err
				}
			} else {
				version, err = unity.GetVersionFromProject(path)
				if err != nil {
					return err
				}
			}

			appInstall, err := unity.GetInstallFromVersion(version)
			if err != nil {
				// TODO if versionNotFoundError offer to install version
				return err
			}

			return runInstalledVersion(appInstall, expandedPath, lFlagBuildTarget)
		},
	}

	cmd.Flags().BoolVar(&lFlagForceVersion, "force", false, "force project to be opened with a specific Unity version") // TODO should let version be passed into flag
	cmd.Flags().StringVar(&lFlagBuildTarget, "buildTarget", "", "opens project with a specific build target set")

	return cmd
}

func runInstalledVersion(installInfo unity.InstallInfo, projectPath, target string) error {
	fmt.Printf("Opening project %q in version %s\n", projectPath, installInfo.Version.String())
	return installInfo.RunWithTarget(projectPath, target)
}
