package clapper

import (
	"os"
	"testing"
)

// return a new registry
func getTestRegistry() Registry {

	// create a new registry
	registry := NewRegistry()

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
		Register("info").                        // `info` command
		AddArg("username", "").                  //
		AddArg("category", "manager").           // default value
		AddFlag("verbose", "v", true, "").       //
		AddFlag("version", "V", false, "1.0.1"). // default value
		AddFlag("output", "o", false, "./")      // default value

	// register the `ghost` sub-command
	registry.
		Register("ghost")

	// return registry
	return registry

}

func TestRootWhenNoCmdArgs(t *testing.T) {

	// registry and parse
	cmdArgs1 := []string{}
	carg1, err1 := getTestRegistry().Parse(cmdArgs1)

	// error
	if err1 != nil || carg1.Cmd != "" {
		t.Error("root command not found")
	}

	/*-------------*/

	// create registry and parse
	cmdArgs2 := []string{"ls", "-f"}
	carg2, err2 := getTestRegistry().Parse(cmdArgs2)

	// error
	if err2 != nil || carg2.Cmd != "" {
		t.Error("root command not found")
	}
}
