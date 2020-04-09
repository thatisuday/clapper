# clapper
A simple but powerful Go package to parse command-line arguments [_getopt(3) style_](http://man7.org/linux/man-pages/man3/getopt.3.html). Designed especially for making CLI based libraries with ease. It has built-in support for sub-commands, long and short flag name combination (_for example `--version` <==> `-v`_), `--flag=<value>` syntax, inverted flag (_for example `--no-clean`_), variadic arguments, etc.

> [**Commando**](https://github.com/thatisuday/commando) CLI application builder library uses this package to parse command-line arguments.

![go-version](https://img.shields.io/github/go-mod/go-version/thatisuday/clapper?label=Go%20Version) &nbsp;
![Build](https://github.com/thatisuday/clapper/workflows/CI/badge.svg?style=flat-square)

![logo](/assets/clapper-logo.png)

## Documentation
[**pkg.go.dev**](https://pkg.go.dev/github.com/thatisuday/clapper?tab=doc)

## Installation
```
$ go get "github.com/thatisuday/clapper"
```

## Usage

```go
// cmd.go
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
	for argName, argValue := range command.Args {
		fmt.Printf("argument(%s) => %#v\n", argName, argValue)
	}

	// get flag values
	for flagName, flagValue := range command.Flags {
		fmt.Printf("flag(%s) => %#v\n", flagName, flagValue)
	}
}
```

In the above example, we have registred a **root** command and an `info` command. The `registry` can parse arguments passed to the command that executed this program.

#### Example 1
When the **root command** is executed with no command-line arguments.

```
$ go run cmd.go

sub-command => ""
argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:""}
flag(dir) => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:""}
flag(force) => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:""}
```

#### Example 2
When the **root command** is executed but not registered.

```
$ NO_ROOT=TRUE go run cmd.go

error => clapper.ErrorUnknownCommand{Name:""}
```

#### Example 3
When the **root command** is executed with short/long flag names as well as by changing the positions of the arguments.

```
$ go run cmd.go userinfo -V 1.0.1 -v --force --dir ./sub/dir
$ go run cmd.go -V 1.0.1 --verbose --force userinfo --dir ./sub/dir
$ go run cmd.go -V 1.0.1 -v --force --dir ./sub/dir userinfo
$ go run cmd.go --version 1.0.1 --verbose --force --dir ./sub/dir userinfo

sub-command => ""
argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:"userinfo"}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:"1.0.1"}
flag(dir) => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:"./sub/dir"}
flag(force) => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
```

#### Example 4
When an **unregistered flag** is provided in the command-line arguments.

```
$ go run cmd.go userinfo -V 1.0.1 -v --force -d ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"-d"}

$ go run cmd.go userinfo -V 1.0.1 -v --force --d ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"--d"}

$ go run cmd.go userinfo -V 1.0.1 -v --force --directory ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"--directory"}

$ go run cmd.go info student --dump
error => clapper.ErrorUnknownFlag{Name:"--dump"}

$ go run cmd.go info student --clean
error => clapper.ErrorUnknownFlag{Name:"--clean"}
```


#### Example 5
When `information` was intended to be a sub-command but not registered and the root command accepts arguments.

```
$ go run cmd.go information --force

sub-command => ""
argument(output) => &clapper.Arg{Name:"output", IsVariadic:false, DefaultValue:"", Value:"information"}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"", Value:""}
flag(dir) => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, IsInverted:false, DefaultValue:"/var/users", Value:""}
flag(force) => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:""}
```

#### Example 6
When a **sub-command** is executed.

```
$ go run cmd.go info student -V -v --output ./opt/dir

sub-command => "info"
argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}
argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}
argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}
flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt/dir"}
flag(clean) => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:""}
```

#### Example 7
When a command is executed with an **inverted** flag (flag that starts with `--no-` prefix).

```
$ go run cmd.go info student -V -v --output ./opt/dir --no-clean

sub-command => "info"
argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:""}
argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:""}
argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:""}
flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:"./opt
```

#### Example 8
When the position of argument values are changed and variadic arguments are provided.

```
$ go run cmd.go info -v student -V 2.0.0 thatisuday math science physics
$ go run cmd.go info student -v --version=2.0.0 thatisuday math science physics
$ go run cmd.go info student thatisuday math science -v physics -V=2.0.0

sub-command => "info"
argument(category) => &clapper.Arg{Name:"category", IsVariadic:false, DefaultValue:"manager", Value:"student"}
argument(username) => &clapper.Arg{Name:"username", IsVariadic:false, DefaultValue:"", Value:"thatisuday"}
argument(subjects) => &clapper.Arg{Name:"subjects", IsVariadic:true, DefaultValue:"", Value:"math,science,physics"}
flag(output) => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, IsInverted:false, DefaultValue:"./", Value:""}
flag(clean) => &clapper.Flag{Name:"clean", ShortName:"", IsBoolean:true, IsInverted:true, DefaultValue:"true", Value:""}
flag(verbose) => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, IsInverted:false, DefaultValue:"false", Value:"true"}
flag(version) => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, IsInverted:false, DefaultValue:"1.0.1", Value:"2.0.0"}
```

#### Example 9
When a **sub-command** is registered without any flags.

```
$ go run cmd.go ghost -v thatisuday -V 2.0.0 teachers

error => clapper.ErrorUnknownFlag{Name:"-v"}
```

#### Example 10
When a **sub-command** is registered without any arguments.

```
$ go run cmd.go ghost
$ go run cmd.go ghost thatisuday extra

sub-command => "ghost
```

#### Example 11
When the **root command** is not registered or the **root command** is registered with no arguments.

```
$ NO_ROOT=TRUE go run cmd.go information
error => clapper.ErrorUnknownCommand{Name:"information"}

$ go run cmd.go ghost
sub-command => "ghost"
```

#### Example 12
When unsupported flag format is provided.

```
$ go run cmd.go ---version 
error => clapper.ErrorUnsupportedFlag{Name:"---version"}

$ go run cmd.go ---v=1.0.0 
error => clapper.ErrorUnsupportedFlag{Name:"---v"}

$ go run cmd.go -version 
error => clapper.ErrorUnsupportedFlag{Name:"-version"}
```

## Contribution
A lot of improvements can be made to this library, one of which is the support for combined short flags, like `-abc`. If you are willing to contribute, create a pull request and mention your bug fixes or enhancements in the comment.
