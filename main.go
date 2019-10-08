package main

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
	"github.com/golobby/repl/engine"
)

var (
	currentSession *engine.session
)

const (
	ascii = `
   ______      __          __    __             ____  __________  __ 
  / ____/___  / /   ____  / /_  / /_  __  __   / __ \/ ____/ __ \/ / 
 / / __/ __ \/ /   / __ \/ __ \/ __ \/ / / /  / /_/ / __/ / /_/ / /  
/ /_/ / /_/ / /___/ /_/ / /_/ / /_/ / /_/ /  / _, _/ /___/ ____/ /___
\____/\____/_____/\____/_.___/_.___/\__, /  /_/ |_/_____/_/   /_____/
                                   /____/                            
`
	version = "0.0.1"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func handler(c string) {
	currentSession.removeTmpCodes()
	currentSession.add(c)
	err := currentSession.writeToFile()
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}
	if currentSession.continueMode {
		fmt.Print(engine.multiplyString("...", currentSession.indents))
		return
	}
	fmt.Println(currentSession.run())
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	currentSession, err = engine.newSession(wd)
	if err != nil {
		panic(err)
	}
	p := prompt.New(handler, completer, prompt.OptionPrefix("repl> "))
	fmt.Println(ascii)
	fmt.Printf("GoLobby Repl v%s\n", version)
	p.Run()
}
