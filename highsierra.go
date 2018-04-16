package mic

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/stephen-fox/mic/internal/createinstallmediautil"
	"github.com/stephen-fox/mic/internal/hdiutilw"
)

type HighSierra struct {
	name             string
	isLoggingEnabled bool
}

func (o HighSierra) Name() string {
	return o.name
}

func (o *HighSierra) SetLogging(enabled bool) {
	o.isLoggingEnabled = enabled
}

func (o HighSierra) CreateIso(isoDestinationPath string, installerApplicationPath string) error {
	cim, err := createinstallmediautil.Get(installerApplicationPath + "/Contents/Resources/createinstallmedia")
	if err != nil {
		return err
	}

	cdr, err := ioutil.TempFile("", "macos-installer-cdr-")
	if err != nil {
		return errors.New("Failed to create empty .cdr file - " + err.Error())
	}
	cdr.Close()
	os.Remove(cdr.Name())

	if o.isLoggingEnabled {
		log.Println("Creating empty .cdr...")
	}

	cdrFinalFilePath, err := hdiutilw.CreateCdr(cdr.Name(), 5500)
	if err != nil {
		return err
	}
	defer os.Remove(cdrFinalFilePath)

	cdrMountPath, err := ioutil.TempDir("", "mount-cdr-")
	if err != nil {
		return errors.New("Failed to create .cdr mount point - " + err.Error())
	}

	err = hdiutilw.Attach(cdrFinalFilePath, cdrMountPath)
	if err != nil {
		return errors.New("Failed to mount .cdr - " + err.Error())
	}
	defer hdiutilw.Detach(cdrMountPath)

	if o.isLoggingEnabled {
		log.Println("Writing installer files to image...")
	}

	err = cim.CreateDmg(cdrMountPath)
	if err != nil {
		return errors.New("Failed to create installer - " + err.Error())
	}

	err = hdiutilw.Detach("/Volumes/Install macOS High Sierra")
	if err != nil {
		return err
	}

	iso, err := ioutil.TempFile("", "macos-installer-iso-")
	if err != nil {
		return errors.New("Failed to create empty .iso file - " + err.Error())
	}
	iso.Close()
	os.Remove(iso.Name())

	if o.isLoggingEnabled {
		log.Println("Converting image to an .iso...")
	}

	err = hdiutilw.Detach(cdrMountPath)
	if err != nil {
		return err
	}

	isoFinalFilePath, err := hdiutilw.ConvertImageToIso(cdrFinalFilePath, iso.Name())
	if err != nil {
		return err
	}

	err = os.Rename(isoFinalFilePath, isoDestinationPath)
	if err != nil {
		return err
	}

	return nil
}