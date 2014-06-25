jobmux
======

A lightweight job multiplexer.

`jobmux` reads jobs from stdin, one job each line. It executes the jobs
concurrently, while keeping stdout and stderr of the jobs sorted.

## Installation

If you are new to Go, you need to setup your Go environment. Please follow the
instructions in [http://golang.org/doc/install].

After setting up Go and and its environment, installation of jobmux is pretty
simple:
```
cd $GOROOT
go get github.com/kaicheng/jobmux
```

Please make sure you have internet connection during installation process.

## How to use

If we have a file(input) with indenpendent jobs:

```
job1 arg1 arg2
job2 arga argb argc
job3 argi argii argiii
```

```
jobmux <input >stdout 2>stderr
```
will be equivalent to
```
$(SHELL) <input >stdout 2>stderr
```

jobmux only keep stdout and stdin sorted. If the jobs interference in
environment variables, files, etc, the behavior is undefined.

### Arguments
Check `jobmux -h` for arguments.