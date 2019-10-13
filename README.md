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
REPL parses the input and put the code you entered in one of these categories:
- Imports
- Var declaration/assignment
- type definition
- function definition
- function calls

then REPL creates a file from state of the session running.

#### REPL Pipeline
`REPL Prompt -> goimports -> go compiler -> shows result of compile`

## Features

### Global Access to vars
In a REPL you need to have access to the vars defined no matter what scope you are in, REPL provides this feature by
defining all vars in a global scope.
```go
REPL> func someFunc() string {
...REPL>     return fmt.Sprint(a)
...REPL> }
REPL> a = 2
REPL> someFunc() // 2
```

### Variable Redefine/assignment
REPL does not care about either type or value of a variable, so you can redefine or change type of variable with ease.
```go
REPL> x := 3
REPL> x := "amirreza"
REPL> x := someType{}
REPL> var x = 5
// all above codes are valid in a REPL interpreter
```

### Instant Expression Evaluation
REPL can easily evaluate an expression for you with a simple built-in command.

### Automated Imports
REPL uses the power of goimports so it can almost identify all packages you use and automatically import them for you.
```go
REPL> fmt.Println("Hello") // no need for manual import, goimports will take care of that
```

### Helpful Error Messages
REPL does not suppress any error message so you can see exact error message from go toolchain

### Go module discovery
REPL build from ground up to be compatible with go modules, so when you fire up a REPL in go project which is a go module project
REPL instantly identifies your module and imports it as a module with local path, so you can access all your public 
types and functions.(watch Demo)

### Shell Commands

#### Go Doc
go doc is available as a shell command so you can access any document about any package with REPL.
```go
REPL> :doc fmt
REPL> :doc json.Marshal
```
#### Eval
```go
REPL> :e 1+2
// <int> 3
REPL> :e "HelloWorld"
// <string> "HelloWorld"
```
#### Dump
shows formatted view of current session state.
```go
REPL> :dump
```
#### Pop
pops latest entered code from session.
```go
REPL> :pop
```

#### file
file shows exactly code that will be generated from current state of session.
```go
REPL> :file
```
## Demo
[![asciicast](https://asciinema.org/a/273628.svg)](https://asciinema.org/a/273628)
## TODO
##### Code completion (with gocode)
