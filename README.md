# gosh

## usage example

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
