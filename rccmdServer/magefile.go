//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(BuildLinux)
	mg.Deps(BuildWindows)
	fmt.Println("Compilation finished")
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

func BuildLinux() error {
	fmt.Println("Building Linux executable...")
	cmd := exec.Command("go", "build", "-o", "./bin/linux/rccmdServer", "./bin/linux")
	return cmd.Run()
}

func BuildWindows() error {
	fmt.Println("Building Windows executable...")
	cmd := exec.Command("go", "build", "-o", "./bin/windows/rccmdServer.exe", "./bin/windows")
	// Add env variable to cross-compile for windows
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOOS=windows", "GOARCH=amd64")

	return cmd.Run()
}
