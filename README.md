[![Build Status](https://travis-ci.org/golobby/gshell.png?branch=master)](https://travis-ci.org/golobby/gshell)
[![Coverage State](https://coveralls.io/repos/github/golobby/gshell/badge.png?branch=master)](https://coveralls.io/github/golobby/gshell)
# gshell
## Getting Started

### Requirements
- go 1.13
- goimports
- gofmt
```bash
go get golang.org/x/tools/cmd/goimports
```

## Installing gshell
```bash
go get github.com/golobby/gshell
go install github.com/golobby/gshell
```
####
## Why
gshell stands for Read-Eval-Print-Loop. In scripting languages like PHP or Python there is an environment called gshell, which is 
basically a command line interface that is connected to the language interpreter and it would show instantly the result of 
user input. In Golang we don't have this feature by default but it does'nt mean that we cannot benefit the power that gshell gives
us.

## How
gshell parses the input and put the code you entered in one of these categories:
- Imports
- Var declaration/assignment
- type definition
- function definition
- function calls

then gshell creates a file from state of the session running.

#### gshell Pipeline
`gshell Prompt -> goimports -> go compiler -> shows result of compile`

## Features

### Global Access to vars
In a gshell you need to have access to the vars defined no matter what scope you are in, gshell provides this feature by
defining all vars in a global scope.
```go
gshell> func someFunc() string {
...gshell>     return fmt.Sprint(a)
...gshell> }
gshell> a = 2
gshell> someFunc() // 2
```

### Variable Redefine/assignment
gshell does not care about either type or value of a variable, so you can redefine or change type of variable with ease.
```go
gshell> x := 3
gshell> x := "amirreza"
gshell> x := someType{}
gshell> var x = 5
// all above codes are valid in a gshell interpreter
```

### Instant Expression Evaluation
gshell can easily evaluate an expression for you with a simple built-in command.

### Automated Imports
gshell uses the power of goimports so it can almost identify all packages you use and automatically import them for you.
```go
gshell> fmt.Println("Hello") // no need for manual import, goimports will take care of that
```

### Helpful Error Messages
gshell does not suppress any error message so you can see exact error message from go toolchain

### Go module discovery
gshell build from ground up to be compatible with go modules, so when you fire up a gshell in go project which is a go module project
gshell instantly identifies your module and imports it as a module with local path, so you can access all your public 
types and functions.(watch Demo)

### Shell Commands

#### Go Doc
go doc is available as a shell command so you can access any document about any package with gshell.
```go
gshell> :doc fmt
gshell> :doc json.Marshal
```
#### Eval
```go
gshell> :e 1+2
// <int> 3
gshell> :e "HelloWorld"
// <string> "HelloWorld"
```
#### Dump
shows formatted view of current session state.
```go
gshell> :dump
```
#### Pop
pops latest entered code from session.
```go
gshell> :pop
```

#### file
file shows exactly code that will be generated from current state of session.
```go
gshell> :file
```
## Demo
[![asciicast](https://asciinema.org/a/273628.svg)](https://asciinema.org/a/273628)
## TODO
##### Code completion (with gocode)
