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
			AddFlag("force", "f", true, "").        // --force, -f | default value: "false"
			AddFlag("verbose", "v", true, "").      // --verbose, -v | default value: "false"
			AddFlag("version", "V", false, "").     // --version, -V <value>
			AddFlag("dir", "", false, "/var/users") // --dir <value> | default value: "/var/users"
	}

	// register the `info` sub-command
	registry.
		Register("info").                        // sub-command
		AddArg("category", "manager").           // default value: manager
		AddArg("username", "").                  //
		AddArg("subjects...", "").               // variadic argument
		AddFlag("verbose", "v", true, "").       // --verbose, -v | default value: "false"
		AddFlag("version", "V", false, "1.0.1"). // --version, -V <value> | default value: "1.0.1"
		AddFlag("output", "o", false, "./").     // --output, -o <value> | default value: "./"
		AddFlag("no-clean", "", true, "")        // --no-clean | default value: "true"

	// register the `ghost` sub-command
	registry.
		Register("ghost")

	// parse command-line arguments
	carg, err := registry.Parse(os.Args[1:])

	/*----------------*/

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
