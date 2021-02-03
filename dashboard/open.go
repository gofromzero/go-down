package dashboard

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Cmd struct {
	Path string
	Args []string
}

var commands = map[string]Cmd{
	"windows": {Path: "cmd", Args: []string{"/c", "start"}},
	"darwin":  {Path: "open", Args: []string{}},
	"linux":   {Path: "xdg-open", Args: []string{}},
}

// Open calls the OS default program for uri
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	run.Args = append(run.Args, uri)

	cmd := exec.Command(run.Path, run.Args...)
	return cmd.Start()
}
