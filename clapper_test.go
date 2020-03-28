package clapper

import (
	"testing"
)

// return a new registry
func getTestRegistry() Registry {

	// create new registry
	registry := NewRegistry()

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
