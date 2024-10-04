package module

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
)

var (
	marshal   = yaml.Marshal
	writeFile = os.WriteFile
)

func UpdateModuleVersions(modules []*Module, inDir string) {
	for _, module := range modules {
		moduleRoot := filepath.Join(inDir, module.Path)
		hash := computeModuleMD5Hash(os.DirFS("."), moduleRoot)

		if hash != module.Hash {
			module.Hash = hash
			module.Version += 1
			module.Changed = true

			fmt.Printf("Module %s has Changed. New version: %d. New hash: %s\n", module.Path, module.Version, hash)
		}
	}
}

func WriteNewVersionToFile(modules []*Module, inDir string) {
	for _, module := range modules {
		if !module.Changed {
			continue
		}

		marshaledModule, err := marshal(module)
		if err != nil {
			log.Panicf("Error marshaling module %s: %s", module.Path, err)
		}

		moduleConfigPath := filepath.Join(inDir, module.Path, moduleFileName)
		if err = writeFile(moduleConfigPath, marshaledModule, 0666); err != nil {
			log.Panicf("Error writing module %s: %s", module.Path, err)
		}
	}
}

var computeModuleMD5Hash = func(fileSystem fs.FS, moduleRoot string) string {
	var concatenatedHashes string

	err := fs.WalkDir(fileSystem, moduleRoot, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".proto" {
			return nil
		}

		concatenatedHashes += computeFileMD5Hash(fileSystem, path)

		return nil
	})
	if err != nil {
		log.Panicf("Error walking through %s: %s", moduleRoot, err)
	}

	if concatenatedHashes == "" {
		log.Panicf("No proto files found in %s", moduleRoot)
	}

	finalHash := md5.New()
	finalHash.Write([]byte(concatenatedHashes))

	return hex.EncodeToString(finalHash.Sum(nil))
}

var computeFileMD5Hash = func(fileSystem fs.FS, filePath string) string {
	file, err := fileSystem.Open(filePath)
	if err != nil {
		log.Panicf("Error opening %s: %s", filePath, err)
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		log.Panicf("Error reading %s: %s", filePath, err)
	}

	return hex.EncodeToString(hash.Sum(nil))
}
