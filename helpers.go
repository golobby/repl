package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func createReplaceRequireClause(moduleName, localPath string) string {
	if moduleName == "" {
		return ""
	}
	return fmt.Sprintf(`replace %s => %s`, moduleName, localPath)
}

func isShellCommand(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[0] == ':'
}

func isTypeDecl(code string) bool {
	matched, err := regexp.Match("type .+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func isFunc(code string) bool {
	matched, err := regexp.Match("^func.+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func isImport(im string) bool {
	matched, err := regexp.Match("import .+", []byte(im))
	if err != nil {
		panic(err)
	}
	return matched
}
func isPrint(code string) bool {
	matched, err := regexp.Match("^fmt.Print.*\\(.*\\)", []byte(code))
	if err != nil {
		panic(err)
	}
	return matched
}
func goGet() error {
	return exec.Command("go", "get", "-u", "./...").Run()
}
func getModuleNameOfCurrentProject(workingDirectory string) string {
	bs, err := ioutil.ReadFile(workingDirectory + "/go.mod")
	if err != nil{
		if os.IsNotExist(err) {
			return ""
		}
		panic(err)
	}
	gomod := string(bs)
	moduleName := strings.Split(strings.Split(gomod, "\n")[0], " ")[1]
	return moduleName
}