package main

import (
	"embed"
	"flag"

	log "github.com/sirupsen/logrus"

	. "github.com/uben01/proto-share/internal/compiler"
	. "github.com/uben01/proto-share/internal/config"
	. "github.com/uben01/proto-share/internal/context"
	_ "github.com/uben01/proto-share/internal/logger"
	. "github.com/uben01/proto-share/internal/module"
	. "github.com/uben01/proto-share/internal/renderer"
)

//go:embed assets/templates
var embedFileSystem embed.FS

func main() {
	silent := flag.Bool("silent", false, "Set log level to Error")
	verbose := flag.Bool("verbose", false, "Set log level to Debug")

	configPath := flag.String("config", "", "Path to the configuration file. If not provided, read from stdin.")
	flag.Parse()
	if *configPath == "" {
		*configPath = "/dev/stdin"
	}
	if *silent && *verbose {
		log.Fatal("Cannot set both silent and verbose flags")
	}

	if *silent {
		log.SetLevel(log.ErrorLevel)
	} else if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	config := ParseConfig(*configPath)
	CTX.Config = config

	modules := DiscoverModules(config.InDir)
	config.Modules = modules

	UpdateModuleVersions(config.Modules, config.InDir)

	if !config.AnyModuleChanged() && !config.ForceGeneration {
		log.Warn("No changes detected. Exiting.")

		return
	}

	RenderTemplates(embedFileSystem, config)
	CompileModules(config)
	WriteNewVersionToFile(config.Modules, config.InDir)
}
