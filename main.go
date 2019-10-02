package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
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
func prompt() {
	fmt.Print(fmt.Sprintf("%s > ", time.Now().Format("15:04:05")))
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	session, err := newSession(wd)
	if err != nil {
		panic(err)
	}
	fmt.Println(ascii)
	fmt.Printf("GoLobby Repl v%s\n", version)
	r := bufio.NewReader(os.Stdin)
	for {
		session.removeTmpCodes()
		prompt()
		code, err := r.ReadString(';')
		if err != nil {
			panic(err)
		}
		session.add(code)

		err = session.writeToFile()
		if err != nil {
			fmt.Printf("Err: %v\n",err)
			continue
		}
		session.run(os.Stdout, os.Stdout)
	}
}
