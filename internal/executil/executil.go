package executil

import (
	"errors"
	"log"
	"os/exec"
	"strings"
)

var (
	Logger *log.Logger
)

func Run(exe string, args []string) ([]string, error) {
	command := exec.Command(exe, args...)

	if Logger != nil {
		Logger.Printf("executing: '%s'...", command.Args)
	}

	raw, err := command.CombinedOutput()
	outputStr := string(raw)
	if err != nil {
		res := "Failed to execute command: '" + exe

		if len(args) > 0 {
			res = res + " " + strings.Join(args, " ")
		}

		res = res + "'. "

		if len(strings.TrimSpace(outputStr)) > 0 {
			res = res + "Output:\n" + outputStr
		} else {
			res = res + "No additional output"
		}

		return []string{}, errors.New(res)
	}

	return strings.Split(outputStr, "\n"), nil
}
