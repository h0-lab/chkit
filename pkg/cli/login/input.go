package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/containerum/chkit/pkg/context"
	"golang.org/x/crypto/ssh/terminal"
)

func InteractiveLogin(ctx *context.Context) error {
	var err error
	var username, pass string

	if strings.TrimSpace(ctx.Client.Username) == "" {
		username, err = readLogin()
		if err != nil {
			return err
		}
		if strings.TrimSpace(username) == "" {
			return ErrInvalidUsername
		}
		ctx.Client.Username = username
	}

	if strings.TrimSpace(ctx.Client.Password) == "" {
		pass, err = readPassword()
		if err != nil {
			return err
		}
		if strings.TrimSpace(pass) == "" {
			return ErrInvalidPassword
		}
		ctx.Client.Password = pass
	}
	return nil
}

func readLogin() (string, error) {
	fmt.Print("Enter your email: ")
	email, err := bufio.NewReader(os.Stdin).ReadString('\n')
	email = strings.TrimRight(email, "\r\n")
	if err != nil {
		return "", ErrUnableToReadUsername.Wrap(err)
	}
	return email, nil
}

func readPassword() (string, error) {
	fmt.Print("Enter your password: ")
	passwordB, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", ErrUnableToReadPassword.Wrap(err)
	}
	fmt.Println("")
	return string(passwordB), nil
}
