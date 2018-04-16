package createinstallmediautil

import (
	"errors"
	"os"
	"os/user"

	"github.com/stephen-fox/mic/internal/executil"
)

const (
	volumeArg        = "--volume"
	noInteractionArg = "--nointeraction"
)

type Cim struct {
	exePath string
}

func Get(exePath string) (Cim, error) {
	err := isRoot()
	if err != nil {
		return Cim{}, err
	}

	info, err := os.Stat(exePath)
	if err != nil {
		return Cim{}, err
	}

	if info.IsDir() {
		return Cim{}, errors.New("The specified createinstallmedia path is a directory")
	}

	return Cim{
		exePath: exePath,
	}, nil
}

func isRoot() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	if u.Name == "root" || u.Name == "System Administrator" {
		return nil
	}

	return errors.New("Please execute this application using 'sudo'")
}

func (o Cim) CreateDmg(imageMountPath string) error {
	args := []string{
		volumeArg, imageMountPath,
		noInteractionArg,
	}

	_, err := executil.Run(o.exePath, args)
	return err
}