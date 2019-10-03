package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type session struct {
	workingDir  	string
	imports       []string
	typesAndMethods []string
	tmpCodes      []int
	code          []string
	sessionDir    string
	Writer io.Writer
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}

const moduleTemplate = `module shell

go 1.13

%s
`

func createReplaceRequireClause(moduleName, localPath string) string {
	return fmt.Sprintf(`replace %s => %s

require %s latest`, moduleName, localPath, moduleName)
}

func isShellCommand(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[0] == ':'
}

func (s *session) addImport(im string) {
	s.imports = append(s.imports, im)
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
		fmt.Println(err)
		return false
	}
	fmt.Println(matched)
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
	matched, err := regexp.Match("fmt.Print.*", []byte(code))
	if err != nil {
		panic(err)
	}
	return matched
}

func (s *session) add(code string) {
	if isShellCommand(code) {
		s.addShellCommand(len(s.code) - 1)
	} else if isImport(code) {
		s.addImport(code)
	} else if isFunc(code) || isTypeDecl(code) {
		s.typesAndMethods = append(s.typesAndMethods, code)
	} else if isPrint(code) {
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code) - 1)
	}  else {
		s.addCode(code)
	}
}
func (s *session) addCode(code string) {
	s.code = append(s.code, code)
}
func (s *session) cleanCurrentShellCommands() {
	_ = s.tmpCodes[:0]
}
func (s *session) addShellCommand(index int) {
	s.tmpCodes = append(s.tmpCodes, index)
}
func createTmpDir(workingDirectory string) (string, error) {
	sessionDir := workingDirectory + "/.goshell/sessions/" + fmt.Sprint(time.Now().Nanosecond())
	err := os.MkdirAll(sessionDir, 500)
	if err != nil {
		return sessionDir, err
	}
	return sessionDir, nil
}

func (s *session) removeTmpCodes() {
	for _,t := range s.tmpCodes {
		s.code[t] = ""
	}
	s.tmpCodes = s.tmpCodes[:0]
}
func getModuleNameOfCurrentProject(workingDirectory string) string {
	bs, err := ioutil.ReadFile(workingDirectory + "/go.mod")
	if err != nil{
		panic(err)
	}
	gomod := string(bs)
	moduleName := strings.Split(strings.Split(gomod, "\n")[0], " ")[1]
	return moduleName
}
func goGet() error {
	return exec.Command("go", "get", "-u", "./...").Run()
}
func newSession(workingDirectory string) (*session, error) {
	sessionDir, err := createTmpDir(workingDirectory)
	if err != nil {
		return nil, err
	}
	err = os.Chdir(sessionDir)
	if err != nil {
		panic(err)
	}
	session := &session{
		workingDir: "",
		imports:    []string{},
		tmpCodes:   []int{},
		code:       []string{},
		sessionDir: sessionDir,
	}
	if err = session.createModule(workingDirectory, getModuleNameOfCurrentProject(workingDirectory));err!=nil {
		return nil, err
	}
	err = goGet()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *session) removeTmpDir() error {
	return os.RemoveAll(s.sessionDir)
}

func (s *session) createModule(wd string, moduleName string) error {
	return ioutil.WriteFile("go.mod", []byte(fmt.Sprintf(moduleTemplate, createReplaceRequireClause(moduleName, wd))), 500)
}
func (s *session) writeToFile() error {
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.validGoFileFromSession()), 500)
}
func (s *session) checkIfErrIsNotDecl(err string) bool {
	return strings.Contains(err, "not used")
}
func (s *session) removeIfIsNotNotDecl(err string) {
	if !s.checkIfErrIsNotDecl(err) {
		s.code = s.code[:len(s.code)-1]
	}
}

func (s *session) run() string {
	err := os.Chdir(s.sessionDir)
	if err != nil {
		panic(err)
	}
	cmdImport := exec.Command("goimports", "-w", "main.go")
	err = cmdImport.Run()
	if err != nil {
		return err.Error()
	}
	cmdRun := exec.Command("go", "run", "main.go")
	out, err := cmdRun.CombinedOutput()
	if err != nil {
		if s.checkIfErrIsNotDecl(string(out)) {
			s.removeIfIsNotNotDecl(string(out))
			return "Note you are not using something that you define or import"
		} else {
			return fmt.Sprintf("Error:: %s", string(out))
		}
	}
	return fmt.Sprintf("%s", out)
}

func (s *session) validGoFileFromSession() string {
	code := "package main\n%s\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, strings.Join(s.imports, "\n"), strings.Join(s.typesAndMethods, "\n"), strings.Join(s.code, "\n"))
}
