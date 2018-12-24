package commands

type Command struct {
	Name     string
	HelpText string
	Action   func(...string) error
}

var CommandOrder = [...]string{"run", "version", "list", "update", "install", "uninstall", "cleanup", "repair"}

var Commands = map[string]Command{
	"run": {
		"run",
		"run the passed in project with an auto detected version of unity",
		run,
	},
	"version": {
		"version",
		"check what version of unity a project is using",
		version,
	},
	"list": {
		"list",
		"list all installed unity versions",
		list,
	},
	"update": {
		"update",
		"update the package index",
		update,
	},
	"install": {
		"install",
		"installed the specified version of unity",
		install,
	},
	"uninstall": {
		"uninstall",
		"uninstall one or multiple versions of Unity",
		uninstall,
	},
	"cleanup": {
		"cleanup",
		"removes unused unity versons",
		cleanup,
	},
	"repair": {
		"repair",
		"fix paths to unity installs",
		repair,
	},
}
