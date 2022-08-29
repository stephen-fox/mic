package mic

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/stephen-fox/cabinet"
)

const (
	osxBaseSystemMountPath = "/Volumes/OS X Base System"
)

type Installer interface {
	CreateIso(isoDestinationPath string, installerApplicationPath string) error
	Name() string
	SetLogging(enabled bool)
}

func Get(installerPath string) (Installer, error) {
	plist := installerPath + "/Contents/Info.plist"

	info, err := os.Stat(plist)
	if err != nil {
		return &PreHighSierra{}, err
	}

	if info.Size() > 10000000 {
		return &PreHighSierra{}, errors.New("Failed to parse installer's Info.plist - it is too big")
	}

	raw, err := ioutil.ReadFile(plist)
	if err != nil {
		return &PreHighSierra{}, err
	}

	for _, l := range strings.Split(string(raw), "\n") {
		if strings.Contains(strings.ToLower(l),"install macos high sierra") {
			return &HighSierra{
				name: "High Sierra",
			}, nil
		}
	}

	return &PreHighSierra{
		name: "Pre-High Sierra",
	}, nil
}

func Validate(isoPath string, installerApplicationPath string) error {
	installerInfo, err := os.Stat(installerApplicationPath)
	if err != nil {
		return errors.New("The specified installer application does not exist - " + err.Error())
	}

	if !installerInfo.IsDir() {
		return errors.New("The specified installer application is not a directory")
	}

	isoInfo, statErr := os.Stat(isoPath)
	if statErr == nil {
		if !isoInfo.IsDir() {
			return errors.New("The specified .iso already exists at '" + isoPath + "'")
		}
	}

	return nil
}

func copyInstallerResources(sourceMountPath string, destinationMountPath string) error {
	tempInstallation := destinationMountPath + "/System/Installation"

	err := os.Remove(tempInstallation + "/Packages")
	if err != nil {
		return errors.New("Failed to remove packages symlink - " + err.Error())
	}

	err = cabinet.CopyDirectory(sourceMountPath + "/Packages", tempInstallation, false)
	if err != nil {
		return errors.New("Failed to copy installer packages into .iso - " + err.Error())
	}

	err = cabinet.CopyFile(sourceMountPath + "/BaseSystem.chunklist", destinationMountPath, false)
	if err != nil {
		return errors.New("Failed to copy installer chunklist into .iso - " + err.Error())
	}

	err = cabinet.CopyFile(sourceMountPath + "/BaseSystem.dmg", destinationMountPath, false)
	if err != nil {
		return errors.New("Failed to copy installer BaseSystem.dmg into .iso - " + err.Error())
	}

	return nil
}