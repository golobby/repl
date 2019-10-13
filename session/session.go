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
	imports        ImportDatas
	types          map[string]string
	funcs          map[string]string
	vars           Vars
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
:dump => dumps current session
:file => prints go file generated from session
:vars => shows only vars of current state
:types => shows only types of current state
:funcs => shows only funcs of current state
:imports => shows only imports of current state
`
const moduleTemplate = `module shell

go 1.13

%s
`

func wrapInPrint(code string) string {
	return fmt.Sprintf(`fmt.Printf("<%%T> %%+v\n", %s, %s)`, code, code)
}
func (s *Session) importsForSource() string {
	return s.imports.String()
}

func (s *Session) addImport(im []ImportData) {
	s.imports = append(s.imports, im...)
}

func (s *Session) appendToLastCode(code string) {
	if len(s.code) == 0 {
		s.code = append(s.code, code)
		return
	}
	s.code[len(s.code)-1] += "\n" + code
	return
}

func (s *Session) addCode(t Type, code string) error {
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
	case Shell:
		return s.handleShellCommands(code)
	case Import:
		s.addImport(ExtractImportData(code))
		return nil
	case Print:
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code)-1)
		return nil
	case TypeDecl:
		s.addType(ExtractTypeName(code), code)
		return nil
	case FuncDecl:
		s.addFunc(ExtractFuncName(code), code)
		return nil
	case VarDecl:
		s.addVar(NewVar(code))
		return nil
	case Empty:
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
	err = s.addCode(typ, code)
	if err != nil {
		return err
	}
	if s.continueMode {
		return nil
	}
	if err := checkIfHasParsingError(s.String()); err != nil {
		s.removeLastCode()
		return err
	}
	return nil
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
		imports:        ImportDatas{},
		types:          map[string]string{},
		funcs:          map[string]string{},
		vars:           map[string]Var{},
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
	return session, nil
}

func (s *Session) removeTmpDir() error {
	return os.RemoveAll(s.sessionDir)
}

func (s *Session) writeToFile() error {
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.String()), 500)
}

func (s *Session) removeLastCode() {
	if len(s.code) == 0 {
		s.code = []string{}
		return
	}
	idx := len(s.code) - 1
	for tmpIdx, t := range s.tmpCodes {
		if t == idx {
			s.tmpCodes = append(s.tmpCodes[:tmpIdx], s.tmpCodes[tmpIdx+1:]...)
		}
	}
	s.code = s.code[:len(s.code)-1]
}

func checkIfErrIsNotDecl(err string) bool {
	return strings.Contains(err, "not used") && !strings.Contains(err, "evaluated")
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
		return strings.Repeat("...", s.indents)
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
	return fmt.Sprintf(code, s.importsForSource(), s.typesForSource(), s.funcsForSource(), "var(\n"+s.vars.String()+"\n)", strings.Join(s.code, "\n"))
}
