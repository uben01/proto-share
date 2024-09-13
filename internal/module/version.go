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
)

var (
	marshal   = yaml.Marshal
	writeFile = os.WriteFile
)

func UpdateModuleVersions(modules []*Module, inDir string) error {
	for _, module := range modules {
		moduleRoot := filepath.Join(inDir, module.Path)
		hash, err := computeModuleMD5Hash(os.DirFS("."), moduleRoot)
		if err != nil {
			return err
		}

		if hash != module.Hash {
			module.Hash = hash
			module.Version += 1
			module.changed = true

			fmt.Printf("Module %s has changed. New version: %d. New hash: %s\n", module.Path, module.Version, hash)
		}
	}

	return nil
}

func WriteNewVersionToFile(modules []*Module, inDir string) error {
	for _, module := range modules {
		if !module.changed {
			continue
		}

		marshaledModule, err := marshal(module)
		if err != nil {
			return err
		}

		moduleConfigPath := filepath.Join(inDir, module.Path, moduleFileName)
		if err = writeFile(moduleConfigPath, marshaledModule, 0666); err != nil {
			return err
		}
	}

	return nil
}

var computeModuleMD5Hash = func(fileSystem fs.FS, moduleRoot string) (string, error) {
	var concatenatedHashes string

	err := fs.WalkDir(fileSystem, moduleRoot, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".proto" {
			return nil
		}

		var hash string
		if hash, err = computeFileMD5Hash(fileSystem, path); err != nil {
			return err
		}
		concatenatedHashes += hash

		return nil
	})
	if err != nil {
		return "", err
	}

	if concatenatedHashes == "" {
		return "", fmt.Errorf("no proto files found in %s", moduleRoot)
	}

	finalHash := md5.New()
	finalHash.Write([]byte(concatenatedHashes))

	return hex.EncodeToString(finalHash.Sum(nil)), nil
}

var computeFileMD5Hash = func(fileSystem fs.FS, filePath string) (string, error) {
	file, err := fileSystem.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
