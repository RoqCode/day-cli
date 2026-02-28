package scope

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrAborted = errors.New("aborted by user")

func getCustomScope() (string, error) {
	fmt.Print("scope: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func GetScope(s []string) (string, error) {
	chosenScope, err := fzf(s)
	if err != nil {
		return "", err
	}

	if chosenScope == "CUSTOM" {
		chosenScope, err = getCustomScope()
		if err != nil {
			return "", err
		}

		if chosenScope == "" {
			fmt.Println("scope can't be empty")
			chosenScope, err = getCustomScope()
			if err != nil {
				return "", err
			}
			if chosenScope == "" {
				return "", ErrAborted
			}
		}
	}

	return chosenScope, nil
}

func fzf(s []string) (string, error) {
	var result strings.Builder
	cmd := exec.Command("fzf")
	cmd.Stdout = &result
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(strings.Join(s, "\n"))

	err := cmd.Start()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return getCustomScope()
		}
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code := exitErr.ExitCode()
			if code == 1 || code == 130 {
				return "", ErrAborted
			}
		}
		return "", fmt.Errorf("fzf failed: %w", err)
	}
	return strings.TrimSpace(result.String()), nil
}
