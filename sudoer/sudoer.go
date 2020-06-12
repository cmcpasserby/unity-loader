package sudoer

import (
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Sudoer struct {
	password string
}

func (s *Sudoer) AskPass() error {
	pwPrompt := &survey.Password{
		Message: "enter admin password",
	}
	fmt.Println("admin access is required")
	if err := survey.AskOne(pwPrompt, &s.password, nil); err != nil {
		return err
	}

	if !s.CheckRoot() {
		return InvalidPasswordError{}
	}

	return nil
}
