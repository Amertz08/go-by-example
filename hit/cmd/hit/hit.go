package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Amertz08/go-by-example/hit"
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
		args:   os.Args,
	}
	if err := run(e); err != nil {
		os.Exit(1)
	}
}

type env struct {
	stdout io.Writer
	stderr io.Writer
	args   []string
	dryRun bool
}

func run(e *env) error {
	c := config{n: 100, c: 10}

	if err := parseArgs(&c, e.args[1:], e.stderr); err != nil {
		return err
	}
	fmt.Fprintf(
		e.stdout,
		"%s\n\nSending %d requests to %q (concurrency: %d)\n",
		logo,
		c.n,
		c.url,
		c.c,
	)
	if e.dryRun {
		return nil
	}
	if err := runHit(&c, e.stdout); err != nil {
		fmt.Fprintf(e.stderr, "\nerror occurred: %v\n", err)
		return err
	}
	return nil
}

func runHit(c *config, stdout io.Writer) error {
	req, err := http.NewRequest(http.MethodGet, c.url, http.NoBody)

	if err != nil {
		return fmt.Errorf("creating a new request: %w", err)
	}
	results, err := hit.SendN(
		c.n, req, hit.Options{Concurrency: c.c, RPS: c.rps},
	)
	if err != nil {
		return fmt.Errorf("sending requests: %w", err)
	}
	sum := hit.Summarize(results)
	printSummary(sum, stdout)
	return nil
}

func printSummary(sum hit.Summary, stdout io.Writer) {
	fmt.Fprintf(stdout, `
Summary:
	Success: 	%.0f%%
	RPS: 		%.1f
	Requests: 	%d
	Errors: 	%d
	Bytes: 		%d
	Duration: 	%s
	Avg. Latency: %s
	Fastest: 	%s
	Slowest: 	%s
`,
		sum.Success,
		math.Round(sum.RPS),
		sum.Requests,
		sum.Errors,
		sum.Bytes,
		sum.Duration.Round(time.Millisecond),
		(sum.Duration / time.Duration(sum.Requests)).Round(time.Millisecond),
		sum.Fastest.Round(time.Millisecond),
		sum.Slowest.Round(time.Millisecond),
	)
}
