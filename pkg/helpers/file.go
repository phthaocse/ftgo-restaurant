package helpers

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

func ProjectPath() (path string) {
	var delimeter string
	var listAncientDirs []string
	if runtime.GOOS == "windows" {
		delimeter = "\\"
	} else {
		delimeter = "/"
	}
	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	path = string(bytes.TrimSpace(stdout))
	if path != "" {
		pathSeparated := strings.Split(path, delimeter)
		listAncientDirs = pathSeparated[:len(pathSeparated)-1]
		path = strings.Join(listAncientDirs, delimeter) + delimeter
		return
	}
	return
}
