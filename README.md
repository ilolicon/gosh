# gosh

## usage example

- single command

```go
package main

import (
    "fmt"

    command "github.com/ilolicon/gosh"
)

func main() {
    cmd := command.NewCommand("echo 'hello, gosh'")
    _ = cmd.Run()
    fmt.Printf("stdout:%s stderr:%s\n", cmd.Stdout(), cmd.Stderr())
}

```

- multi commands

```go
package main

import (
    "fmt"
    command "github.com/ilolicon/gosh"    
)

func main() {
    cmds := []string{
        "echo 'task test'",
        "uptime",
        "date",
    }
    task := command.NewTask(cmds, 10)
    task.Run(false)
    if err := task.GetError(); err != nil {
        fmt.Printf("task running error: %s\n", err.Error())
    }
}

```
