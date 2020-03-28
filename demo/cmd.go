package main

import (
	"fmt"
	"os"

	"github.com/thatisuday/clapper"
)

func main() {

	// create new registry
	registry := clapper.NewRegistry()

	// register the root command
	registry.
		Register("").
		AddArg("output").
		AddFlag("force", "f", true, "").
		AddFlag("verbose", "v", true, "").
		AddFlag("version", "V", false, "").
		AddFlag("dir", "", false, "/var/users")

	// register the `info` command
	registry.
		Register("info").
		AddArg("username").
		AddArg("category").
		AddFlag("verbose", "v", true, "").
		AddFlag("version", "V", false, "1.0.1").
		AddFlag("output", "o", false, "./")

	// register the `ghost` command
	registry.
		Register("ghost")

	// parse command line arguments
	carg, err := registry.Parse(os.Args[1:])

	// check for error
	if err != nil {
		fmt.Printf("Error => %#v\n", err)
		return
	}

	// get executed command name
	fmt.Printf("Command Name => %#v\n", carg.Cmd)

	// get argument values
	for _, v := range carg.Args {
		fmt.Printf("Arg => %#v\n", v)
	}

	// get flag values
	for _, v := range carg.Flags {
		fmt.Printf("Flag => %#v\n", v)
	}
}
