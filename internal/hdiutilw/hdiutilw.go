package hdiutilw

import (
	"errors"
	"os"
	"strconv"

	"github.com/stephen-fox/mic/internal/executil"
)

const (
	defaultHdiutil = "hdiutil"
	attachArg      = "attach"
	detachArg      = "detach"
	createArg      = "create"
	compactArg     = "compact"
	convertArg     = "convert"
	noVerifyArg    = "-noverify"
	noBrowseArg    = "-nobrowse"
	mountPointArg  = "-mountpoint"
	outputArg      = "-o"
	sizeArg        = "-size"
	layoutArg      = "-layout"
	fileSystemArg  = "-fs"
	typeArg        = "-type"
	formatArg      = "-format"
)

var (
	ExePath = defaultHdiutil
)

func Attach(filePath string, mountPath string) error {
	args := []string{
		attachArg, filePath, noVerifyArg, noBrowseArg,
		mountPointArg, mountPath,
	}

	_, err := executil.Run(ExePath, args)
	if err != nil {
		return errors.New("Failed to mount file - " + err.Error())
	}

	return nil
}

func Detach(mountPointPath string) error {
	args := []string{
		detachArg, mountPointPath,
	}

	_, err := executil.Run(ExePath, args)
	return err
}

func Compact(filePath string) error {
	args := []string{
		compactArg, filePath,
	}

	_, err := executil.Run(ExePath, args)
	return err
}

func CreateSparseImage(filePath string, sizeMb int) (string, error) {
	args := []string{
		createArg,
		outputArg, filePath,
		sizeArg, strconv.Itoa(sizeMb) + "m",
		layoutArg, "SPUD",
		fileSystemArg, "HFS+J",
		typeArg, "SPARSE",
	}

	_, err := executil.Run(ExePath, args)
	if err != nil {
		return "", errors.New("Failed to populate sparseimage - " + err.Error())
	}

	sparseImagePath := filePath + ".sparseimage"
	_, statErr := os.Stat(sparseImagePath)
	if statErr == nil {
		return sparseImagePath, nil
	}

	return "", errors.New("Failed to create sparseimage")
}

func CreateCdr(filePath string, sizeMb int) (string, error) {
	args := []string{
		createArg,
		outputArg, filePath,
		sizeArg, strconv.Itoa(sizeMb) + "m",
		layoutArg, "SPUD",
		fileSystemArg, "HFS+J",
	}

	_, err := executil.Run(ExePath, args)
	if err != nil {
		return "", errors.New("Failed to populate cdr - " + err.Error())
	}

	sparseImagePath := filePath + ".dmg"
	_, statErr := os.Stat(sparseImagePath)
	if statErr == nil {
		return sparseImagePath, nil
	}

	return "", errors.New("Failed to create cdr")
}

func ConvertImageToIso(imageFilePath string, destinationPath string) (string, error) {
	args := []string{
		convertArg, imageFilePath,
		formatArg, "UDTO",
		outputArg, destinationPath,
	}

	_, err := executil.Run(ExePath, args)
	if err != nil {
		return "", errors.New("Failed to convert image file to a .iso - " + err.Error())
	}

	resultingFilePath := destinationPath + ".cdr"
	_, statErr := os.Stat(resultingFilePath)
	if statErr == nil {
		return resultingFilePath, nil
	}

	return "", errors.New("Failed to convert image file to a .iso")
}