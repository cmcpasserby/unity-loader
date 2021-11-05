package main

import (
	"fmt"
	"github.com/cmcpasserby/unity-loader/unity"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func createRunCmd() *cobra.Command {
	var ( // local flags
		lFlagForceVersion string
		lFlagBuildTarget  string
	)

	cmd := &cobra.Command{
		Use:   "run [projectDirectory]",
		Short: "Launches unity and opens the selected project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := getConfig()
			if err != nil {
				return err
			}

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

			if lFlagForceVersion != "" {
				version = lFlagForceVersion
			} else {
				version, err = unity.GetVersionFromProject(path)
				if err != nil {
					return err
				}
			}

			appInstall, err := unity.GetInstallFromVersion(version, config.SearchPaths...)
			if err != nil {
				return err
			}

			return runInstalledVersion(appInstall, expandedPath, lFlagBuildTarget)
		},
	}

	cmd.Flags().StringVar(&lFlagForceVersion, "force", "", "force project to be opened with a specific Unity version")
	cmd.Flags().StringVar(&lFlagBuildTarget, "buildTarget", "", "opens project with a specific build target set")

	return cmd
}

func runInstalledVersion(installInfo unity.InstallInfo, projectPath, target string) error {
	fmt.Printf("Opening project \"%s\" in version %s\n", projectPath, installInfo.Version.String())
	return installInfo.RunWithTarget(projectPath, target)
}
