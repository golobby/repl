package session

import (
	"fmt"
	p "go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Session struct {
	shellCmdOutput string
	imports        []string
	types          map[string]string
	funcs          map[string]string
	vars           map[string]string
	tmpCodes       []int
	code           []string
	sessionDir     string
	Writer         io.Writer
	continueMode   bool
	indents        int
}

const helpText = `
List of REPL commands:
:help => shows help
:doc => shows go documentation of package/function
:e => evaluates expression
:pop => pop latest code from session
:dump => shows all codes in the session
`
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

const outputTemplte = `import(
%s
)
types(
%s
)
funcs(
%s
)
var(
%s
)
main(
%s
)`

func (s *Session) dump() string {
	return fmt.Sprintf(outputTemplte, strings.Join(s.imports, "\n"), s.typesAsString(), s.funcsAsString(), s.varsString(), strings.Join(s.code, "\n"))
}

func (s *Session) addCode(t StmtType, code string) error {
	if s.continueMode {
		s.appendToLastCode(code)
		indents, shouldContinue := ShouldContinue(s.code[len(s.code)-1])
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
	indents, shouldContinue := ShouldContinue(code)
	s.indents = indents
	if shouldContinue {
		s.continueMode = true
		s.code = append(s.code, code)
		return nil
	}
	switch t {
	case StmtShell:
		return s.handleShellCommands(code)
	case StmtTypeImport:
		s.addImport(code)
		return nil
	case StmtTypePrint:
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
		return nil
	case StmtTypeTypeDecl:
		s.addType(ExtractTypeName(code), code)
		return nil
	case StmtTypeFuncDecl:
		s.addFunc(ExtractFuncName(code), code)
		return nil
	case StmtVarDecl:
		varName, value := ExtractNameAndValueFromVarInit(code)
		s.addVar(varName, value)
		return nil
	case StmtEmpty:
		return nil
	default:
		s.code = append(s.code, code)
		return nil
	}
}

func (s *Session) Add(code string) error {
	s.removeTmpCodes()
	typ, err := Parse(code)
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
		shellCmdOutput: "",
		imports:        []string{},
		types:          map[string]string{},
		funcs:          map[string]string{},
		vars:           map[string]string{},
		tmpCodes:       []int{},
		code:           []string{},
		sessionDir:     sessionDir,
		Writer:         nil,
		continueMode:   false,
		indents:        0,
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
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.String()), 500)
}

func (s *Session) removeLastCode() {
	idx := len(s.code) - 1
	for tmpIdx, t := range s.tmpCodes {
		if t == idx {
			s.tmpCodes = append(s.tmpCodes[:tmpIdx], s.tmpCodes[tmpIdx+1:]...)
		}
	}
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
	return strings.Contains(err, "not used") && !strings.Contains(err, "evaluated")
}
func multiplyString(s string, n int) string {
	return strings.Repeat(s, n)
}
func checkIfHasParsingError(code string) error {
	fs := token.NewFileSet()
	_, err := p.ParseFile(fs, "", code, p.AllErrors)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Eval() string {
	if s.shellCmdOutput != "" {
		output := s.shellCmdOutput
		s.shellCmdOutput = ""
		return output + "\n"
	}
	if len(s.code) == 0 {
		return ""
	}
	if s.continueMode {
		return multiplyString("...", s.indents)
	}
	if err := checkIfHasParsingError(s.String()); err != nil {
		s.removeLastCode()
		return err.Error() + "\n"
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
		return fmt.Sprintf("%s %s\n", string(out), err.Error())
	}
	cmdRun := exec.Command("go", "run", "main.go")
	out, err = cmdRun.CombinedOutput()
	if err != nil {
		if checkIfErrIsNotDecl(string(out)) {
			return fmt.Sprintf("%s %s\n", string(out), err.Error())
		} else {
			s.removeLastCode()
			return fmt.Sprintf("%s %s\n", string(out), err.Error())
		}
	}
	return fmt.Sprintf("%s", out)
}

func (s *Session) String() string {
	code := "package main\n%s\n%s\n%s\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, strings.Join(s.imports, "\n"), s.typesForSource(), s.funcsForSource(), s.varsForSource(), strings.Join(s.code, "\n"))
}
