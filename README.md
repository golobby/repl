[![Build Status](https://travis-ci.org/golobby/repl.svg?branch=master)](https://travis-ci.org/golobby/repl)
[![Coverage State](https://coveralls.io/repos/github/golobby/repl/badge.svg?branch=master)](https://coveralls.io/github/golobby/repl)
# REPL
## Why
REPL stands for READ-EVAL-PRINT-LOOP. In scripting languages like PHP or Python there is an environment called REPL, which is 
basically a command line interface that is connected to the language interpreter and it would show instantly the result of 
user input. In Golang we don't have this feature by default but it does'nt mean that we cannot benefit the power that REPL gives
us.

## How
golobby/repl basically creates a new golang project and gives a direct interface to the project, every line that you type into the 
console will get directly to the go project and will be compiled instantly and you will see the result.

## Demo
[![asciicast](https://asciinema.org/a/14.png)](https://asciinema.org/a/14)
## TODO
##### Code completion (with gocode)
