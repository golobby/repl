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
	imports       []string
	tmpCodes      []int
	code          []string
	sessionDir    string
	errorDisplay  io.Writer
	outputDisplay io.Writer
}

func isShellCommand(code string) bool {
	return code[0] == ':'
}

func (s *session) addImport(im string) {
	s.imports = append(s.imports, im)
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
	} else if isPrint(code) {
		s.code = append(s.code, code)
		s.tmpCodes = append(s.tmpCodes, len(s.code) - 1)
	}else {
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
	for idx,c := range s.code{
		if c == "" {
			s.code = append(s.code[:idx] ,s.code[idx+1:]...)
		}
	}
}

func newSession(workingDirectory string) (*session, error) {
	sessionDir, err := createTmpDir(workingDirectory)
	if err != nil {
		return nil, err
	}
	return &session{
		imports:       []string{},
		tmpCodes:      []int{},
		code:          []string{},
		sessionDir:    sessionDir,
		errorDisplay:  os.Stdout,
		outputDisplay: os.Stdout,
	}, nil
}

func (s *session) removeTmpDir() error {
	return os.RemoveAll(s.sessionDir)
}

func (s *session) writeToFile() error {
	return ioutil.WriteFile(s.sessionDir+"/main.go", []byte(s.validGoFileFromSession()), 500)
}

func (s *session) run(stdOut, stdErr io.Writer) error {
	err := os.Chdir(s.sessionDir)
	if err != nil {
		return err
	}
	cmdImport := exec.Command("goimports", "-w", "main.go")
	cmdImport.Stdout = s.outputDisplay
	cmdImport.Stderr = s.errorDisplay
	err = cmdImport.Run()
	if err != nil {
		return err
	}
	cmdRun := exec.Command("go", "run", "main.go")
	cmdRun.Stdout = s.outputDisplay
	cmdRun.Stderr = s.errorDisplay
	err = cmdRun.Run()
	if err != nil {
		return err
	}
	return nil
}
func (s *session) displayError(err error) error {
	_, err = s.errorDisplay.Write([]byte(err.Error()+"\n"))
	return err
}
func (s *session) displayOutput(output string) error {
	_, err := s.errorDisplay.Write([]byte(output + "\n"))
	return err
}

func (s *session) validGoFileFromSession() string {
	code := "package main\n%s\nfunc main() {\n%s\n}"
	return fmt.Sprintf(code, strings.Join(s.imports, "\n"), strings.Join(s.code, "\n"))
}
