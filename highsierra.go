package mic

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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

	installerAppSizeBytes, err := dirSizeBytes(installerApplicationPath)
	if err != nil {
		return fmt.Errorf("failed to get installer application size - %s", err.Error())
	}

	// Add another 500 mb.
	installerAppSizeBytes += 500000000

	if o.isLoggingEnabled {
		log.Println("Creating empty .cdr...")
	}

	cdrFinalFilePath, err := hdiutilw.CreateCdr(cdr.Name(), installerAppSizeBytes)
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

	// First, unmount the installer volume.
	volumesPath := "/Volumes"
	infos, err := ioutil.ReadDir(volumesPath)
	if err != nil {
		if o.isLoggingEnabled {
			log.Printf("unable to unmount installer volume - failed to read volumes directory - %s", err.Error())
		}
	} else {
		installerAppName := strings.TrimSuffix(path.Base(installerApplicationPath), ".app")
		for _, info := range infos {
			if strings.Contains(info.Name(), installerAppName) {
				err := hdiutilw.Detach(path.Join(volumesPath, info.Name()))
				if err != nil && o.isLoggingEnabled {
					log.Printf("failed to detach installer volumes - %s", err.Error())
				}
			}
		}
	}

	// Process the error for creating the .dmg *after* unmounting
	// the installer volume.
	if err != nil {
		return errors.New("Failed to create .dmg image - " + err.Error())
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

// https://stackoverflow.com/a/32482941
func dirSizeBytes(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
