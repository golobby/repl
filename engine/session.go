package engine

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golobby/repl/engine/parser"
)

type Session struct {
	imports         []string
	typesAndMethods []string
	tmpCodes        []int
	code            []string
	sessionDir      string
	Writer          io.Writer
	continueMode    bool
	indents         int
}

/*
func (s *session) add(code string) {
	if s.continueMode {
		s.code[len(s.code)-1] += code
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
	else {
		if isExpr(code) {
			s.add(wrapInPrint(code))
			return
		}
		s.addCode(code)
	}
}
*/
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

func (s *Session) addCode(typ parser.StmtType, shouldContinue bool, code string) {
	if shouldContinue {
		s.code[len(s.code)-1] += code
		if !parser.ShouldContinue(s.code[len(s.code)-1]) {
			s.continueMode = false
			code = s.code[len(s.code)-1]
			s.code = s.code[:len(s.code)-1]
			s.Add(code)
		}
	}
	switch typ {
	case parser.StmtTypeImport:
		s.addImport(code)
	case parser.StmtTypeComment:
		return
	case parser.StmtTypePrint:
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
	case parser.StmtTypeFuncDecl, parser.StmtTypeTypeDecl:
		s.typesAndMethods = append(s.typesAndMethods, code)
	case parser.StmtTypeExpr:
		s.code = append(s.code, wrapInPrint(code))
	default:
		s.code = append(s.code, code)
	}
}
func (s *Session) mergeIfShouldContinue(shouldContinue bool, code string) error {
	if !shouldContinue {
		return nil
	}

	return nil
}

func (s *Session) Add(code string) error {
	s.removeTmpCodes()
	typ, shouldContinue, err := parser.Parse(code)
	if err != nil {
		return err
	}
	s.addCode(typ, shouldContinue, code)
	return nil
}

func createTmpDir(workingDirectory string) (string, error) {
	sessionDir := workingDirectory + "/.goshell/sessions/" + fmt.Sprint(time.Now().Nanosecond())
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
}
func goGet() error {
	return exec.Command("go", "get", "-u", "./...").Run()
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
	err = goGet()
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
	_, err = cmdImport.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%s", err.Error())
	}
	cmdRun := exec.Command("go", "run", "main.go")
	out, err := cmdRun.CombinedOutput()
	if err != nil {
		if checkIfErrIsNotDecl(string(out)) {
			s.removeLastCode()
			return "Note you are not using something that you define or import"
		} else {
			return fmt.Sprintf("Error:: %s", string(out))
		}
	}
	return fmt.Sprintf("%s", out)
}

func (s *Session) validGoFileFromSession() string {
	code := "package main\n%s\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, strings.Join(s.imports, "\n"), strings.Join(s.typesAndMethods, "\n"), strings.Join(s.code, "\n"))
}
