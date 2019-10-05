package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	mapChars = map[string]string{
		"(": ")",
		"{": "}",
	}
)

type session struct {
	workingDir      string
	imports         []string
	typesAndMethods []string
	tmpCodes        []int
	code            []string
	sessionDir      string
	Writer          io.Writer
	continueMode    bool
	indents         int
}

func (s *session) shouldContinue(code string) bool {
	var stillOpenChars string
	for _, c := range code {
		if c == '{' || c == '(' {
			stillOpenChars += string(c)
			continue
		}
		if c == '}' || c == ')' {
			idx := strings.Index(stillOpenChars, mapChars[string(c)])
			if idx >= 0 {
				if len(stillOpenChars) == 0 {
					return false
				}
				stillOpenChars = stillOpenChars[:idx] + stillOpenChars[idx+1:]
			}
		}
	}
	if len(stillOpenChars) > 0 {
		s.indents = len(stillOpenChars)
		return true
	}
	return false
}

const moduleTemplate = `module shell

go 1.13

%s
`

func (s *session) addImport(im string) {
	s.imports = append(s.imports, im)
}

func (s *session) add(code string) {
	if s.continueMode {
		s.code[len(s.code)-1] += "\n" + code
		if !s.shouldContinue(s.code[len(s.code)-1]) {
			s.continueMode = false
			code = s.code[len(s.code)-1]
			s.code = s.code[:len(s.code)-1]
			s.add(code)
			return
		}
		return
	}

	if s.shouldContinue(code) {
		s.continueMode = true
		s.code = append(s.code, code)
		return
	}

	if isShellCommand(code) {
		s.addShellCommand(len(s.code) - 1)
	} else if isImport(code) {
		s.addImport(code)
	} else if isFunc(code) || isTypeDecl(code) {
		s.typesAndMethods = append(s.typesAndMethods, code)
	} else if isPrint(code) {
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
	} else if isComment(code) {
		s.code = append(s.code, code)
	} else {
		if isExpr(code) {
			s.add(wrapInPrint(code))
			return
		}
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
func goCompilerPretifyOutput(output string) string {
	return ""
}

func (s *session) removeTmpCodes() {
	for _, t := range s.tmpCodes {
		s.code[t] = ""
	}
	s.tmpCodes = s.tmpCodes[:0]
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
	currentModule := getModuleNameOfCurrentProject(workingDirectory)
	if err = session.createModule(workingDirectory, currentModule); err != nil {
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
	out1, err := cmdImport.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", out1)
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
