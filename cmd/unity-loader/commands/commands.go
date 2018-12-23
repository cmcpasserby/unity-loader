package commands

type Command struct {
	Name     string
	HelpText string
	Action   func(...string) error
}

var CommandOrder = [...]string{"update"}

var Commands = map[string]Command{

	"update": {
		"update",
		"update the package index",
		update,
	},
}
