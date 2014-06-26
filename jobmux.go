package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

type job struct {
	cmd    string
	stdout chan<- []byte
	stderr chan<- []byte
}

// flag variables
var numWorkers int
var shell string

// global variables
var jobs chan job
var stdouts chan (<-chan []byte)
var stderrs chan (<-chan []byte)
var finish chan bool

func worker() {
	for j := range jobs {
		var stdout, stderr bytes.Buffer
		cmd := exec.Cmd{
			Path:   shell,
			Args:   []string{shell},
			Stdin:  bytes.NewBufferString(j.cmd),
			Stdout: &stdout,
			Stderr: &stderr,
		}
		cmd.Run()
		go func(j job) {
			j.stdout <- stdout.Bytes()
			j.stderr <- stderr.Bytes()
			close(j.stdout)
			close(j.stderr)
		}(j)
	}
}

func inputReader(reader io.Reader) {
	token := make(chan bool, 1)
	token <- true
	bufreader := bufio.NewReader(reader)
	for line, err := bufreader.ReadBytes('\n'); err == nil; line, err = bufreader.ReadBytes('\n') {
		stdout := make(chan []byte)
		stderr := make(chan []byte)
		jobs <- job{cmd: string(line), stdout: stdout, stderr: stderr}
		prev := token
		next := make(chan bool, 1)
		go func() {
			<-prev
			stdouts <- stdout
			stderrs <- stderr
			next <- true
		}()
		token = next
	}
	go func() {
		<-token
		close(stdouts)
		close(stderrs)
	}()
}

func outputWriter(writer io.Writer, chans <-chan (<-chan []byte)) {
	for c := range chans {
		for data := range c {
			writer.Write(data)
		}
	}
	finish <- true
}

func init() {
	flag.IntVar(&numWorkers, "n", -1, "number of workers. -1 means the number of logical CPUs.")
	flag.StringVar(&shell, "shell", "", "Absolute path to your shell.")
}

func main() {
	flag.Parse()
	if numWorkers < 0 {
		if numWorkers != -1 {
			fmt.Fprintln(os.Stderr, "Invalid number of workers.")
			os.Exit(-1)
		}
		numWorkers = runtime.NumCPU()
	}
	if len(shell) == 0 {
		shell = os.Getenv("SHELL")
		if len(shell) == 0 {
			fmt.Fprintln(os.Stderr, "No available shells.")
			os.Exit(-1)
		}
	}
	jobs = make(chan job)
	stdouts = make(chan (<-chan []byte))
	stderrs = make(chan (<-chan []byte))
	finish = make(chan bool, 2)
	go inputReader(os.Stdin)
	go outputWriter(os.Stdout, stdouts)
	go outputWriter(os.Stderr, stderrs)
	for i := 0; i < numWorkers; i += 1 {
		go worker()
	}
	<-finish
	<-finish
}
