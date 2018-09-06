package sudoer

import (
    "fmt"
    "gopkg.in/AlecAivazis/survey.v1"
    "io"
    "os/exec"
)

type Sudoer struct {
    password string
}

func (s *Sudoer) AskPass() bool {
    password := ""
    pwPrompt := &survey.Password{
        Message: "enter admin password",
    }
    fmt.Println("admin access is required")
    survey.AskOne(pwPrompt, &password, nil)

    return s.CheckRoot()
}

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
    sudoArgs := make([]string, 0, len(args) + 2)
    sudoArgs = append(sudoArgs, "-S", command)
    for _, arg := range args {
        sudoArgs = append(sudoArgs, arg)
    }

    process := exec.Command("sudo", sudoArgs...)
    processIn, _ := process.StdinPipe()
    io.WriteString(processIn, fmt.Sprintf("%s\n", s.password))

    err := process.Run()
    if err != nil {return err}
    return nil
}
