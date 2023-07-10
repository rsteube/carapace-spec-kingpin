package spec

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"gopkg.in/yaml.v3"
)

func Register(app *kingpin.Application) {
	cmd := app.GetCommand("_carapace")
	if cmd == nil {
		cmd = app.Command("_carapace", "")
		cmd.Hidden()
	}

	specCmd := cmd.Command("spec", "")
	specCmd.Action(func(pc *kingpin.ParseContext) error {
		fmt.Println(Scrape(app))
		return nil
	})
}

func Scrape(app *kingpin.Application) string {
	cmd := command(&kingpin.CmdModel{
		Name:           app.Name,
		Help:           app.Help,
		FlagGroupModel: app.Model().FlagGroupModel,
		CmdGroupModel:  app.Model().CmdGroupModel,
	}, true)
	m, err := yaml.Marshal(cmd)
	if err != nil {
		panic(err.Error())
	}
	return string(m)
}

func command(c *kingpin.CmdModel, root bool) Command {
	cmd := Command{
		Name:            c.Name,
		Aliases:         c.Aliases,
		Description:     c.Help,
		Hidden:          c.Hidden,
		Flags:           make(map[string]string),
		PersistentFlags: make(map[string]string),
		Commands:        make([]Command, 0),
	}
	cmd.Completion.Flag = make(map[string][]string)

	// TODO groups
	// if group := node.Group; group != nil {
	// 	cmd.Group = group.Key
	// }

	for _, flag := range c.Flags {
		formatted := ""

		if flag.Short != 0 {
			formatted += fmt.Sprintf("-%v, ", string(flag.Short))
		}
		formatted += fmt.Sprintf("--%v", flag.Name)

		if flag.Hidden {
			formatted += "&"
		}

		if flag.Required {
			formatted += "!"
		}

		// 	if flag.IsCounter() || flag.IsCumulative() { // TODO
		// 		formatted += "*"
		// 	}

		switch {
		case flag.IsBoolFlag():
		//case optionalArgument: // TODO
		//	formatted += "?"
		default:
			formatted += "="
		}

		flags := cmd.Flags
		if root {
			flags = cmd.PersistentFlags
		}
		flags[formatted] = flag.Help
		if flag.IsBoolFlag() {
			flags["--no-"+flag.Name] = flag.Help
		}
	}

	for _, subcmd := range c.Commands {
		if subcmd.Name != "_carapace" {
			cmd.Commands = append(cmd.Commands, command(subcmd, false))
		}
	}
	return cmd
}
