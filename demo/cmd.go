package main

import (
	"fmt"
	"os"

	"github.com/thatisuday/clapper"
)

func main() {

	// create a new registry
	registry := clapper.NewRegistry()

	// register the root command
	if _, ok := os.LookupEnv("NO_ROOT"); !ok {
		rootCommand, _ := registry.Register("")             // root command
		rootCommand.AddArg("output", "")                    //
		rootCommand.AddFlag("force", "f", true, "")         // --force, -f | default value: "false"
		rootCommand.AddFlag("verbose", "v", true, "")       // --verbose, -v | default value: "false"
		rootCommand.AddFlag("version", "V", false, "")      // --version, -V <value>
		rootCommand.AddFlag("dir", "", false, "/var/users") // --dir <value> | default value: "/var/users"
	}

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info")         // sub-command
	infoCommand.AddArg("category", "manager")           // default value: manager
	infoCommand.AddArg("username", "")                  //
	infoCommand.AddArg("subjects...", "")               // variadic argument
	infoCommand.AddFlag("verbose", "v", true, "")       // --verbose, -v | default value: "false"
	infoCommand.AddFlag("version", "V", false, "1.0.1") // --version, -V <value> | default value: "1.0.1"
	infoCommand.AddFlag("output", "o", false, "./")     // --output, -o <value> | default value: "./"
	infoCommand.AddFlag("no-clean", "", true, "")       // --no-clean | default value: "true"

	// register the `ghost` sub-command
	registry.Register("ghost")

	/*----------------*/

	// parse command-line arguments
	command, err := registry.Parse(os.Args[1:])

	/*----------------*/

	// check for error
	if err != nil {
		fmt.Printf("error => %#v\n", err)
		return
	}

	// get executed sub-command name
	fmt.Printf("sub-command => %#v\n", command.Name)

	// get argument values
	for _, v := range command.Args {
		fmt.Printf("argument-value => %#v\n", v)
	}

	// get flag values
	for _, v := range command.Flags {
		fmt.Printf("flag-value => %#v\n", v)
	}
}
