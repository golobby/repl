package session

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golobby/repl/parser"
)

type Session struct {
	shellCmdOutput  string
	imports         []string
	typesAndMethods []string
	tmpCodes        []int
	code            []string
	sessionDir      string
	Writer          io.Writer
	continueMode    bool
	indents         int
}

const moduleTemplate = `module shell

go 1.13

%s
`

func wrapInPrint(code string) string {
	return fmt.Sprintf(`fmt.Printf("<%%T> %%+v\n", %s, %s)`, code, code)
}

func (s *Session) addImport(im string) {
	s.imports = append(s.imports, im)
}

func (s *Session) appendToLastCode(code string) {
	if len(s.code) == 0 {
		s.code = append(s.code, code)
		return
	}
	s.code[len(s.code)-1] += "\n" + code
	return
}

func (s *Session) handleShellCommands(code string) error {
	typ, data := parser.ParseCmd(code)
	switch typ {
	case parser.REPLCmdDoc:
		output, err := goDoc(data)
		if err != nil {
			return err
		}
		s.shellCmdOutput = string(output)
	default:
		return nil
	}
	return nil
}
func (s *Session) addCode(t parser.StmtType, code string) error {
	if s.continueMode {
		s.appendToLastCode(code)
		indents, shouldContinue := parser.ShouldContinue(s.code[len(s.code)-1])
		s.indents = indents
		if !shouldContinue {
			s.continueMode = false
			code = s.code[len(s.code)-1]
			s.code = s.code[:len(s.code)-1]
			s.Add(code)
			return nil
		}
		return nil
	}
	indents, shouldContinue := parser.ShouldContinue(code)
	s.indents = indents
	if shouldContinue {
		s.continueMode = true
		s.code = append(s.code, code)
		return nil
	}
	switch t {
	case parser.StmtShell:
		return s.handleShellCommands(code)
	case parser.StmtTypeImport:
		s.addImport(code)
		return nil
	case parser.StmtTypeComment:
		return nil
	case parser.StmtTypePrint:
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
		return nil
	case parser.StmtTypeFuncDecl, parser.StmtTypeTypeDecl:
		s.typesAndMethods = append(s.typesAndMethods, code)
		return nil
	case parser.StmtTypeExpr:
		return s.Add(wrapInPrint(code))
	case parser.StmtEmpty:
		return nil
	default:
		s.code = append(s.code, code)
		return nil
	}
}

func (s *Session) Add(code string) error {
	s.removeTmpCodes()
	typ, err := parser.Parse(code)
	if err != nil {
		return err
	}
	return s.addCode(typ, code)
}

func createTmpDir(workingDirectory string) (string, error) {
	sessionDir := workingDirectory + "/.repl/sessions/" + fmt.Sprint(time.Now().Nanosecond())
	err := os.MkdirAll(sessionDir, 500)
	if err != nil {
		return sessionDir, err
	}
	return sessionDir, nil
}

func (s *Session) removeTmpCodes() {
	for _, t := range s.tmpCodes {
		s.code[t] = ""
	}
	s.tmpCodes = s.tmpCodes[:0]
	for idx, c := range s.code {
		if c == "" {
			s.code = append(s.code[:idx], s.code[idx+1:]...)
		}
	}
}
func goBuild() error {
	return exec.Command("go", "build", "./...").Run()
}

func goGet() error {
	return exec.Command("go", "get", "./...").Run()
}

func goDoc(code string) ([]byte, error) {
	return exec.Command("go", "doc", code).CombinedOutput()
}

func NewSession(workingDirectory string) (*Session, error) {
	sessionDir, err := createTmpDir(workingDirectory)
	if err != nil {
		return nil, err
	}
	err = os.Chdir(sessionDir)
	if err != nil {
		panic(err)
	}
	session := &Session{
		imports:    []string{},
		tmpCodes:   []int{},
		code:       []string{},
		sessionDir: sessionDir,
	}
	currentModule := getModuleNameOfCurrentProject(workingDirectory)
	if err = session.createModule(workingDirectory, currentModule); err != nil {
		return nil, err
	}
	err = session.writeToFile()
	if err != nil {
		return nil, err
	}
	err = goGet()
	if err != nil {
		return nil, err
	}
	err = goBuild()
	if err != nil {
		return nil, err
	}
	return session, nil
}
func createReplaceRequireClause(moduleName, localPath string) string {
	if moduleName == "" {
		return ""
	}
	return fmt.Sprintf(`replace %s => %s`, moduleName, localPath)
}
func (s *Session) removeTmpDir() error {
	return os.RemoveAll(s.sessionDir)
}

func (s *Session) createModule(wd string, moduleName string) error {
	return ioutil.WriteFile("go.mod", []byte(fmt.Sprintf(moduleTemplate, createReplaceRequireClause(moduleName, wd))), 500)
}
func (s *Session) writeToFile() error {
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.validGoFileFromSession()), 500)
}

func (s *Session) removeLastCode() {
	s.code = s.code[:len(s.code)-1]
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

func checkIfErrIsNotDecl(err string) bool {
	return strings.Contains(err, "not used")
}
func multiplyString(s string, n int) string {
	var out string
	if n == 0 {
		return ""
	}
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func (s *Session) Eval() string {
	if s.shellCmdOutput != "" {
		output := s.shellCmdOutput
		s.shellCmdOutput = ""
		return output
	}
	if len(s.code) == 0 {
		return ""
	}
	if s.continueMode {
		return multiplyString("...", s.indents)
	}
	err := s.writeToFile()
	if err != nil {
		return err.Error()
	}
	err = os.Chdir(s.sessionDir)
	if err != nil {
		panic(err)
	}
	cmdImport := exec.Command("goimports", "-w", "main.go")
	out, err := cmdImport.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s %s", string(out), err.Error())
	}
	cmdRun := exec.Command("/bin/sh", "-c", "go build -o repl_session; ./repl_session")
	out, err = cmdRun.CombinedOutput()
	if err != nil {
		if checkIfErrIsNotDecl(string(out)) {
			return "Note you are not using something that you define or import"
		} else {
			s.removeLastCode()
			return fmt.Sprintf("%s %s", string(out), err.Error())
		}
	}
	return fmt.Sprintf("%s", out)
}

func (s *Session) validGoFileFromSession() string {
	code := "package main\n%s\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, strings.Join(s.imports, "\n"), strings.Join(s.typesAndMethods, "\n"), strings.Join(s.code, "\n"))
}
