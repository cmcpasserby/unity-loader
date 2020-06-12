package sudoer

type InvalidPasswordError struct{}

func (err InvalidPasswordError) Error() string {
	return "invalid password\n"
}
