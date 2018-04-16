package asrutil

import (
	"errors"

	"github.com/stephen-fox/mic/internal/executil"
)

const (
	defaultAsr  = "asr"
	noVerifyArg = "-noverify"
	restoreArg  = "restore"
	sourceArg   = "-source"
	targetArg   = "-target"
	noPromptArg = "-noprompt"
	eraseArg    = "-erase"
)

var (
	ExePath = defaultAsr
)

func Restore(installerBaseSystemDmgPath string, sparseImageMountPath string) error {
	args := []string{
		restoreArg,
		sourceArg, installerBaseSystemDmgPath,
		targetArg, sparseImageMountPath,
		noPromptArg,
		//noVerifyArg,
		eraseArg,
	}

	_, err := executil.Run(ExePath, args)
	if err != nil {
		return errors.New("Failed to copy installer base system into .iso - " + err.Error())
	}

	return nil
}