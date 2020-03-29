# clapper
A simple but powerful Go package to parse command-line arguments [_getopt(3) style_](http://man7.org/linux/man-pages/man3/getopt.3.html). Designed especially for making CLI based libraries with ease.

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

    // parse command-line arguments
    carg, err := registry.Parse(os.Args[1:])

    // check for error
    if err != nil {
        fmt.Printf("Error => %#v\n", err)
        return
    }

    // get executed sub-command name
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
```

In the above example, we have registred a **root** command and an `info` command. The `registry` can parse arguments passed to the command that executed this program.

#### Example 1
```
$ go run cmd.go

Command Name => ""
Arg => &clapper.Arg{Name:"output", Value:""}
Flag => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:""}
Flag => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:""}
Flag => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:""}
Flag => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:""}
```

#### Example 2
```
$ go run cmd.go userinfo -V 1.0.1 -v --force --dir ./sub/dir
$ go run cmd.go -V 1.0.1 --verbose --force userinfo --dir ./sub/dir
$ go run cmd.go -V 1.0.1 -v --force --dir ./sub/dir userinfo
$ go run cmd.go --version 1.0.1 --verbose --force -dir ./sub/dir userinfo

Command Name => ""
Arg => &clapper.Arg{Name:"output", Value:"userinfo"}
Flag => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
Flag => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:"1.0.1"}
Flag => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:"./sub/dir"}
Flag => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:"true"}
```

#### Example 3
```
$ go run cmd.go userinfo -V 1.0.1 -v --force -d ./sub/dir
Error => clapper.ErrorUnknownFlag{Name:"d", IsShort:true}

$ go run cmd.go userinfo -V 1.0.1 -v --force --d ./sub/dir
Error => clapper.ErrorUnknownFlag{Name:"d", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force -di ./sub/dir
Error => clapper.ErrorUnknownFlag{Name:"di", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force --directory ./sub/dir
Error => clapper.ErrorUnknownFlag{Name:"directory", IsShort:false}

$ go run cmd.go userinfo -V 1.0.1 -v --force -directory ./sub/dir
Error => clapper.ErrorUnknownFlag{Name:"directory", IsShort:false}
```


#### Example 4
```
$ go run cmd.go information --force

Command Name => ""
Arg => &clapper.Arg{Name:"output", Value:"information"}
Flag => &clapper.Flag{Name:"force", ShortName:"f", IsBoolean:true, DefaultValue:"false", Value:"true"}
Flag => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:""}
Flag => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"", Value:""}
Flag => &clapper.Flag{Name:"dir", ShortName:"", IsBoolean:false, DefaultValue:"/var/users", Value:""}
```

#### Example 5
```
$ go run cmd.go info thatisuday teachers -V -v --output ./opt/dir extra

Command Name => "info"
Arg => &clapper.Arg{Name:"username", Value:"thatisuday"}
Arg => &clapper.Arg{Name:"category", Value:"teachers"}
Flag => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
Flag => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:""}
Flag => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:"./opt/dir"}
```

#### Example 6
```
$ go run cmd.go info -v thatisuday -V 2.0.0  teachers extra
$ go run cmd.go info thatisuday -v -V 2.0.0  teachers extra
$ go run cmd.go info thatisuday teachers extra -v -V 2.0.0

Command Name => "info"
Arg => &clapper.Arg{Name:"username", Value:"thatisuday"}
Arg => &clapper.Arg{Name:"category", Value:"teachers"}
Flag => &clapper.Flag{Name:"verbose", ShortName:"v", IsBoolean:true, DefaultValue:"false", Value:"true"}
Flag => &clapper.Flag{Name:"version", ShortName:"V", IsBoolean:false, DefaultValue:"1.0.1", Value:"2.0.0"}
Flag => &clapper.Flag{Name:"output", ShortName:"o", IsBoolean:false, DefaultValue:"./", Value:""}
```

#### Example 7
```
$ go run cmd.go ghost -v thatisuday -V 2.0.0 teachers extra

Error => clapper.ErrorUnknownFlag{Name:"v", IsShort:true}
```

#### Example 8
```
$ go run cmd.go ghost
$ go run cmd.go ghost thatisuday extra

Command Name => "ghost
```

## Contribution
I am looking for some contributors for writing better test cases and documentation.