package mic

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/stephen-fox/mic/internal/asrutil"
	"github.com/stephen-fox/mic/internal/hdiutilw"
)

type PreHighSierra struct {
	name             string
	isLoggingEnabled bool
}

func (o PreHighSierra) Name() string {
	return o.name
}

func (o *PreHighSierra) SetLogging(enabled bool) {
	o.isLoggingEnabled = enabled
}

func (o PreHighSierra) CreateIso(isoDestinationPath string, installerApplicationPath string) error {
	_, statErr := os.Stat(osxBaseSystemMountPath)
	if statErr == nil {
		return errors.New("An OS-X base system is already mounted at '" +
			osxBaseSystemMountPath + "'. You must umount it first")
	}

	esdMountPath, err := ioutil.TempDir("", "mount-esd-")
	if err != nil {
		return errors.New("Failed to create mount point - " + err.Error())
	}

	err = hdiutilw.Attach(installerApplicationPath + "/Contents/SharedSupport/InstallESD.dmg", esdMountPath)
	if err != nil {
		return errors.New("Failed to mount installer application's .dmg - " + err.Error())
	}
	defer hdiutilw.Detach(esdMountPath)

	sparseImage, err := ioutil.TempFile("", "macos-installer-sparseimage-")
	if err != nil {
		return errors.New("Failed to create empty sparseimage - " + err.Error())
	}
	sparseImage.Close()
	os.Remove(sparseImage.Name())

	sparseImagePath, err := hdiutilw.CreateSparseImage(sparseImage.Name(), 8000)
	if err != nil {
		return err
	}
	defer os.Remove(sparseImagePath)

	sparseImageMountPath, err := ioutil.TempDir("", "mount-sparseimage-")
	if err != nil {
		return errors.New("Failed to create mount point - " + err.Error())
	}

	err = hdiutilw.Attach(sparseImagePath, sparseImageMountPath)
	if err != nil {
		return err
	}
	defer hdiutilw.Detach(sparseImageMountPath)

	if o.isLoggingEnabled {
		log.Println("Adding installer files to .iso sparseimage...")
	}

	err = asrutil.Restore(esdMountPath + "/BaseSystem.dmg", sparseImageMountPath)
	if err != nil {
		return err
	}

	if o.isLoggingEnabled {
		log.Println("Copying installer files...")
	}

	err = copyInstallerResources(esdMountPath, osxBaseSystemMountPath)
	if err != nil {
		return err
	}

	err = hdiutilw.Detach(osxBaseSystemMountPath)
	if err != nil {
		return err
	}

	if o.isLoggingEnabled {
		log.Println("Shrinking .iso sparseimage...")
	}

	err = hdiutilw.Compact(sparseImagePath)
	if err != nil {
		return err
	}

	if o.isLoggingEnabled {
		log.Println("Converting .iso sparseimage to a .iso file...")
	}

	iso, err := ioutil.TempFile("", "macos-installer-iso-")
	if err != nil {
		return errors.New("Failed to create empty .iso file - " + iso.Name())
	}
	iso.Close()
	os.Remove(iso.Name())

	isoFinalFilePath, err := hdiutilw.ConvertImageToIso(sparseImagePath, iso.Name())
	if err != nil {
		return err
	}

	err = os.Rename(isoFinalFilePath, isoDestinationPath)
	if err != nil {
		return err
	}

	return nil
}