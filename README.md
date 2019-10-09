[![Build Status](https://travis-ci.org/golobby/repl.png?branch=master)](https://travis-ci.org/golobby/repl)
[![Coverage State](https://coveralls.io/repos/github/golobby/repl/badge.png?branch=master)](https://coveralls.io/github/golobby/repl)
# REPL
## Getting Started

### Requirements
- go 1.13
- goimports
- gofmt
```bash
go get golang.org/x/tools/cmd/goimports
```

## Installing REPL
```bash
go get github.com/golobby/repl
go install github.com/golobby/repl
```
####
## Why
REPL stands for Read-Eval-Print-Loop. In scripting languages like PHP or Python there is an environment called REPL, which is 
basically a command line interface that is connected to the language interpreter and it would show instantly the result of 
user input. In Golang we don't have this feature by default but it does'nt mean that we cannot benefit the power that REPL gives
us.

## How
golobby/repl basically creates a new golang project and gives a direct interface to the project, every line that you type into the 
console will get directly to the go project and will be compiled instantly and you will see the result.
### REPL Pipeline
`REPL Prompt -> gofmt -> goimports -> go compiler -> shows result of compile`

## Features

### Instant Expression Evaluation
REPL automatically understands if you have typed an expression and evaluates it and shows to you type and value of result.
### Automated Imports
REPL uses the power of goimports so it can almost identify all packages you use and automatically import them for you.

### Helpful Error Messages
REPL does not suppress any error message so you can see exact error message from go toolchain

### Go module discovery
REPL build from ground up to be compatible with go modules, so when you fire up a REPL in go project which is a go module project
REPL instantly identifies your module and imports it as a module with local path, so you can access all your public types and functions.


## Demo
[![asciicast](https://asciinema.org/a/272640.svg)](https://asciinema.org/a/272640)

## TODO
##### Code completion (with gocode)
