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

	// Note: Creating individual temp files and removing them
	// seems to cause some kind of file system race. It was
	// unnecessary to begin with - but "Resource temporarily
	// unavailable" errors sealed the deal.
	tempDirPath, err := ioutil.TempDir("", "macos-installer-")
	if err != nil {
		return errors.New("Failed to create empty .cdr file - " + err.Error())
	}

	installerAppSizeBytes, err := dirSizeBytes(installerApplicationPath)
	if err != nil {
		return fmt.Errorf("failed to get installer application size - %s", err.Error())
	}

	// TODO: Used to be 500 mb - now 5gb.
	//
	// Investigate auto-sizing based on hdiutil error:
	// Failed to create .dmg image - Failed to execute command:
	// '.../createinstallmedia --volume /tmp/mount-cdr-12345 --nointeraction'. Output:
	// /tmp/mount-cdr-12345 is not large enough for install media. An additional 1.34 GB is needed.
	installerAppSizeBytes += 5000000000

	if o.isLoggingEnabled {
		log.Printf("Creating empty %d gb .cdr...", installerAppSizeBytes/1000000000)
	}

	cdrFinalFilePath, err := hdiutilw.CreateCdr(path.Join(tempDirPath, "cdr"), installerAppSizeBytes)
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

	ejectInstallVolumeFn := func() {
		const volumesPath = "/Volumes"
		infos, err := ioutil.ReadDir(volumesPath)
		if err != nil {
			if o.isLoggingEnabled {
				log.Printf("unable to unmount installer volume - failed to read volumes directory - %s",
					err.Error())
			}
			return
		}

		installerAppName := strings.TrimSuffix(path.Base(installerApplicationPath), ".app")
		for _, info := range infos {
			if strings.Contains(info.Name(), installerAppName) {
				installVolume := path.Join(volumesPath, info.Name())
				err := hdiutilw.ForceEject(installVolume)
				if err != nil && o.isLoggingEnabled {
					log.Printf("failed to detach installer volume '%s' - %s",
						installVolume, err.Error())
				}
			}
		}
	}

	err = cim.CreateDmg(cdrMountPath)
	ejectInstallVolumeFn()
	if err != nil {
		return errors.New("Failed to create .dmg image - " + err.Error())
	}

	if o.isLoggingEnabled {
		log.Println("Converting image to an .iso...")
	}

	isoFinalFilePath, err := hdiutilw.ConvertImageToIso(cdrFinalFilePath, path.Join(tempDirPath, "installer-iso"))
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
