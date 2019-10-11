package session

import "os/exec"

func goBuild() error {
	return exec.Command("go", "build", "./...").Run()
}

func goGet() error {
	return exec.Command("go", "get", "./...").Run()
}

func goDoc(code string) ([]byte, error) {
	return exec.Command("go", "doc", code).CombinedOutput()
}
