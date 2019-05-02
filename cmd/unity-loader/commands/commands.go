package commands

import "flag"

type Command struct {
	Name     string
	HelpText string
	Flags    *flag.FlagSet
	Action   func(...string) error
}

var CommandOrder = [...]string{"run", "version", "list", "update", "install", "uninstall", "repair", "config"}

var Commands = map[string]Command{

	"run": {
		"run",
		"run the passed in project with an auto detected version of unity",
		func() *flag.FlagSet {
			fs := flag.NewFlagSet("run", flag.ExitOnError)
			fs.Bool("force", false, "force a certain version to be used for running")
			fs.String("buildTarget", "", "Allows the selection of an active build target before loading a project")
			return fs
		}(),
		run,
	},

	"version": {
		"version",
		"check what version of unity a project is using",
		nil,
		version,
	},

	"list": {
		"list",
		"list all installed unity versions",
		nil,
		list,
	},

	"update": {
		"update",
		"update the package index",
		nil,
		update,
	},

	"install": {
		"install",
		"installed the specified version of unity",
		nil,
		install,
	},

	"uninstall": {
		"uninstall",
		"uninstall one or multiple versions of Unity",
		nil,
		uninstall,
	},

	// "cleanup": {
	// 	"cleanup",
	// 	"removes unused unity versions",
	// 	cleanup,
	// },

	"repair": {
		"repair",
		"fix paths to unity installs",
		nil,
		repair,
	},

	"config": {
		"config",
		"open the config file",
		nil,
		config,
	},
}
