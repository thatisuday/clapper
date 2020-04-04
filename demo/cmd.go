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
		registry.
			Register("").                           // root command
			AddArg("output", "").                   //
			AddFlag("force", "f", true, "").        //
			AddFlag("verbose", "v", true, "").      //
			AddFlag("version", "V", false, "").     //
			AddFlag("dir", "", false, "/var/users") // default value
	}

	// register the `info` sub-command
	registry.
		Register("info").                        // sub-command
		AddArg("username", "").                  //
		AddArg("category", "manager").           // default value
		AddFlag("verbose", "v", true, "").       //
		AddFlag("version", "V", false, "1.0.1"). // default value
		AddFlag("output", "o", false, "./")      // default value

	// register the `ghost` sub-command
	registry.
		Register("ghost")

	// parse command-line arguments
	carg, err := registry.Parse(os.Args[1:])

	// check for error
	if err != nil {
		fmt.Printf("error => %#v\n", err)
		return
	}

	// get executed sub-command name
	fmt.Printf("sub-command => %#v\n", carg.Cmd)

	// get argument values
	for _, v := range carg.Args {
		fmt.Printf("argument-value => %#v\n", v)
	}

	// get flag values
	for _, v := range carg.Flags {
		fmt.Printf("flag-value => %#v\n", v)
	}
}
