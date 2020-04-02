# clapper
A simple but powerful Go package to parse command-line arguments [_getopt(3) style_](http://man7.org/linux/man-pages/man3/getopt.3.html). Designed especially for making CLI based libraries with ease.

![go-version](https://img.shields.io/github/go-mod/go-version/thatisuday/clapper?label=Go%20Version&style=flat-square)

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

import "github.com/thatisuday/clapper"

func main() {
    
    // create a new registry
    registry := clapper.NewRegistry()

    // register the root command
    registry.
        Register("").
        AddArg("output").
        AddFlag("force", "f", true, "").
        AddFlag("verbose", "v", true, "").
        AddFlag("version", "V", false, "").
        AddFlag("dir", "", false, "/var/users")

    // register the `info` sub-command
    registry.
        Register("info").
        AddArg("username").
        AddArg("category").
        AddFlag("verbose", "v", true, "").
        AddFlag("version", "V", false, "1.0.1").
        AddFlag("output", "o", false, "./")

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
```

In the above example, we have registred a **root** command and an `info` command. The `registry` can parse arguments passed to the command that executed this program.

#### Example 1
When the **root command** is executed with no command-line arguments.

```
$ go run cmd.go

sub-command => ""
argument-value => &clapper.Arg{Name:"output", Value:""}
flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:""}
flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:""}
flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:""}
flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:""}
```

#### Example 2
When the **root command** is executed but not registered.

```
$ go run cmd.go

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
argument-value => &clapper.Arg{Name:"output", Value:"userinfo"}
flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:"1.0.1"}
flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:"./sub/dir"}
flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:"true"}
```

#### Example 4
When an **unregistered flag** is provided in the command-line arguments.

```
$ go run cmd.go userinfo -V 1.0.1 -v --force -d ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"d", IsShort:true}

$ go run cmd.go userinfo -V 1.0.1 -v --force --d ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"d", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force -di ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"di", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force --directory ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"directory", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force -directory ./sub/dir
error => clapper.ErrorUnknownFlag{Name:"directory", IsShort:false}
```


#### Example 5
When `information` was intended to be a sub-command but not registered.

```
$ go run cmd.go information --force

sub-command => ""
argument-value => &clapper.Arg{Name:"output", Value:"information"}
flag-value => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:"true"}
flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:""}
flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:""}
flag-value => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:""}
```

#### Example 6
When a **sub-command** is executed.

```
$ go run cmd.go info thatisuday teachers -V -v --output ./opt/dir

sub-command => "info"
argument-value => &clapper.Arg{Name:"username", Value:"thatisuday"}
argument-value => &clapper.Arg{Name:"category", Value:"teachers"}
flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:""}
flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:"./opt/dir"}
```

#### Example 7
When the position of arguments' values are changed and extra argument values are provided.

```
$ go run cmd.go info -v thatisuday -V 2.0.0  teachers extra
$ go run cmd.go info thatisuday -v -V 2.0.0  teachers extra
$ go run cmd.go info thatisuday teachers extra -v -V 2.0.0

sub-command => "info"
argument-value => &clapper.Arg{Name:"username", Value:"thatisuday"}
argument-value => &clapper.Arg{Name:"category", Value:"teachers"}
flag-value => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
flag-value => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:"2.0.0"}
flag-value => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:""}
```

#### Example 8
When a **sub-command** is registered without any flags.

```
$ go run cmd.go ghost -v thatisuday -V 2.0.0 teachers extra

error => clapper.ErrorUnknownFlag{Name:"v", IsShort:true}
```

#### Example 9
When a **sub-command** is registered without any arguments.

```
$ go run cmd.go ghost
$ go run cmd.go ghost thatisuday extra

sub-command => "ghost
```

#### Example 10
When the **root command** is not registered or the **root command** is registered with no arguments.

```
$ go run cmd.go information
error => clapper.ErrorUnknownCommand{Name:"information"}

$ go run cmd.go ghost
sub-command => "ghost"
```

## Contribution
I am looking for some contributors for writing better test cases and documentation.