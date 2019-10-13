package interpreter

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func goBuild() error {
	return exec.Command("go", "build", "./...").Run()
}

func goGet() error {
	return exec.Command("go", "get", "./...").Run()
}

func goDoc(code string) ([]byte, error) {
	return exec.Command("go", "doc", code).CombinedOutput()
}

func getModuleNameOfCurrentProject(workingDirectory string) string {
	bs, err := ioutil.ReadFile(workingDirectory + "/go.mod")
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		panic(err)
	}
	gomod := string(bs)
	moduleName := strings.Split(strings.Split(gomod, "\n")[0], " ")[1]
	return moduleName
}

func createReplaceRequireClause(moduleName, localPath string) string {
	if moduleName == "" {
		return ""
	}
	return fmt.Sprintf(`replace %s => %s`, moduleName, localPath)
}

func (s *Interpreter) createModule(wd string, moduleName string) error {
	return ioutil.WriteFile("go.mod", []byte(fmt.Sprintf(moduleTemplate, createReplaceRequireClause(moduleName, wd))), 500)
}
