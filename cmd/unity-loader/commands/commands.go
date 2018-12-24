package commands

type Command struct {
	Name     string
	HelpText string
	Action   func(...string) error
}

var CommandOrder = [...]string{"run", "version", "list", "update"}

var Commands = map[string]Command{
	"update": {
		"update",
		"update the package index",
		update,
	},
	"list": {
		"list",
		"list all installed unity versions",
		list,
	},
	"version": {
		"version",
		"check what version of unity a project is using",
		version,
	},
}
