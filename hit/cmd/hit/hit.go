package main

import (
	"fmt"
	"io"
	"os"
)

const logo = `
 __  __     __     ______  
/\ \_\ \   /\ \   /\__  _\ 
\ \  __ \  \ \ \  \/_/\ \/ 
 \ \_\ \_\  \ \_\    \ \_\ 
  \/_/\/_/   \/_/     \/_/`

func main() {
	e := &env{
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   os.Args[1:],
	}
	if err := run(e); err != nil {
		os.Exit(1)
	}
}

type env struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	dryrun bool
}

func run(e *env) error {
	c := config{n: 100, c: 10}

	if err := parseArgs(&c, e.args[1:], e.stderr); err != nil {
		return err
	}
	fmt.Fprintf(
		e.stdout,
		"%s\n\nSendidng %d requeests to %q (concurrency: %d)\n",
		logo,
		c.n,
		c.url,
		c.c,
	)
	return nil
}
