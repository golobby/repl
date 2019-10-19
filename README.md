[![Build Status](https://travis-ci.org/golobby/repl.png?branch=master)](https://travis-ci.org/golobby/repl)
[![Coverage State](https://coveralls.io/repos/github/golobby/repl/badge.png?branch=master)](https://coveralls.io/github/golobby/repl)
# repl
## Getting Started

### Requirements
- go 1.13
- goimports
- gofmt
```bash
go get golang.org/x/tools/cmd/goimports
```

## Installing repl
```bash
go get github.com/golobby/repl
go install github.com/golobby/repl
```
####
## Why
repl is a REPL, Read-Eval-Print-Loop. In scripting languages like PHP or Python there is an environment called repl, which is 
basically a command line interface that is connected to the language interpreter and it would show instantly the result of 
user input. In Golang we don't have this feature by default but it does'nt mean that we cannot benefit the power that repl gives
us.

## How
repl parses the input and put the code you entered in one of these categories:
- Imports
- Var declaration/assignment
- type definition
- function definition
- function calls

then repl creates a file from state of the session running.

#### repl Pipeline
`repl Prompt -> goimports -> go compiler -> shows result of compile`

## Features

### Global Access to vars
In a repl you need to have access to the vars defined no matter what scope you are in, repl provides this feature by
defining all vars in a global scope.
```go
repl> func someFunc() string {
...repl>     return fmt.Sprint(a)
...repl> }
repl> a = 2
repl> someFunc() // 2
```

### Variable Redefine/assignment
repl does not care about either type or value of a variable, so you can redefine or change type of variable with ease.
```go
repl> x := 3
repl> x := "amirreza"
repl> x := someType{}
repl> var x = 5
// all above codes are valid in a repl interpreter
```

### Instant Expression Evaluation
repl can easily evaluate any valid Go expression.
```go
repl> 2 // <int> 2
repl> 2*3*(2+1) // <int> 18
repl> "Hello World" // <string> "Hello World"
```

### Automated Imports
repl uses the power of goimports so it can almost identify all packages you use and automatically import them for you.
```go
repl> fmt.Println("Hello") // no need for manual import, goimports will take care of that
```

### Helpful Error Messages
repl does not suppress any error message so you can see exact error message from go toolchain

### Go module discovery
repl build from ground up to be compatible with go modules, so when you fire up a repl in go project which is a go module project
repl instantly identifies your module and imports it as a module with local path, so you can access all your public 
types and functions.(watch Demo)

### Shell Commands

#### Go Doc
go doc is available as a shell command so you can access any document about any package with repl.
```go
repl> :doc fmt
repl> :doc json.Marshal
```
#### Eval
```go
repl> :e 1+2
// <int> 3
repl> :e "HelloWorld"
// <string> "HelloWorld"
```
#### Dump
shows formatted view of current session state.
```go
repl> :dump
```
#### Pop
pops latest entered code from session.
```go
repl> :pop
```

#### file
file shows exactly code that will be generated from current state of session.
```go
repl> :file
```
## Demo
[![asciicast](https://asciinema.org/a/273628.svg)](https://asciinema.org/a/273628)
## TODO
##### Code completion (with gocode)
