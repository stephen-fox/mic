package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stephen-fox/mic"
)

const (
	applicationName = "mic"

	installerApplicationPathArg = "i"
	isoFileOutputArg            = "o"
	forceHighSierraArg          = "high-sierra-strategy"

	noLogArg    = "q"
	versionArg  = "v"
	helpArg     = "h"
	examplesArg = "x"

	examples = `Create a macOS installer .iso:
	` + applicationName + ` -` + installerApplicationPathArg + ` '/Applications/Install macOS High Sierra.app -` + isoFileOutputArg + ` ~/Desktop/macos-high-sierra.iso`
)

var (
	version string

	installerApplicationPath = flag.String(installerApplicationPathArg, "", "The path to the macOS installer application")
	isoFileOutputPath        = flag.String(isoFileOutputArg, "", "The path to save the .iso to")
	forceHighSierra          = flag.Bool(forceHighSierraArg, false, "Force .iso creation to use the High Sierra strategy")

	noLog         = flag.Bool(noLogArg, false, "Do not print log output")
	printVersion  = flag.Bool(versionArg, false, "Prints the version")
	printExamples = flag.Bool(examplesArg, false, "Print usage examples")
	printHelp     = flag.Bool(helpArg, false, "Prints this help page")
)

func main() {
	flag.Parse()

	if len(os.Args) <= 1 || *printHelp {
		fmt.Println(applicationName, version, "\n")

		fmt.Println("[ABOUT]\nApplication for creating macOS installer .iso files (macOS ISO Creator).")
		fmt.Println("This application supports macOS El Capitan through High Sierra.")
		fmt.Println("Note: You must first download a macOS installer from the App Store.\n")

		fmt.Println("[USAGE]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *printExamples {
		fmt.Println(examples)
		os.Exit(0)
	}

	if len(strings.TrimSpace(*installerApplicationPath)) == 0 {
		log.Fatal("Please specify the path to the macOS installer application using '-",
			installerApplicationPathArg, " </path/to/application.app>'")
	}

	if len(strings.TrimSpace(*isoFileOutputPath)) == 0 {
		log.Fatal("Please specify where to create the macOS installer .iso using '-",
			isoFileOutputArg, " </path/to/macos.iso>'")
	}

	err := mic.Validate(*isoFileOutputPath, *installerApplicationPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	var installer mic.Installer

	if *forceHighSierra {
		log.Println("Installer .iso will be created using the High Sierra strategy")

		installer = &mic.HighSierra{}
	} else {
		installer, err = mic.Get(*installerApplicationPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Determined installer as", installer.Name())
	}

	installer.SetLogging(!*noLog)

	err = installer.CreateIso(*isoFileOutputPath, *installerApplicationPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Created installer .iso at '" + *isoFileOutputPath + "'")
}