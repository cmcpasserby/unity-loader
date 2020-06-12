package sudoer

import (
	"fmt"
	"io"
	"os/exec"
)

func (s *Sudoer) CheckRoot() bool {
	resetCmd := exec.Command("sudo", "-k")
	resetCmd.Run()

	sudoCmd := exec.Command("sudo", "-S", "whoami")
	sudoIn, _ := sudoCmd.StdinPipe()

	// todo find a better method then input a pw for all attempts
	io.WriteString(sudoIn, fmt.Sprintf("%s\n", s.password))
	io.WriteString(sudoIn, fmt.Sprintf("%s\n", s.password))
	io.WriteString(sudoIn, fmt.Sprintf("%s\n", s.password))

	err := sudoCmd.Run()
	return err == nil
}

func (s *Sudoer) RunAsRoot(command string, args ...string) error {
	sudoArgs := make([]string, 0, len(args)+2)
	sudoArgs = append(sudoArgs, "-S", command)
	for _, arg := range args {
		sudoArgs = append(sudoArgs, arg)
	}

	process := exec.Command("sudo", sudoArgs...)
	processIn, _ := process.StdinPipe()
	if _, err := io.WriteString(processIn, fmt.Sprintf("%s\n", s.password)); err != nil {
		return err
	}

	err := process.Run()
	if err != nil {
		return err
	}
	return nil
}
