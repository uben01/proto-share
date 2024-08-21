package module

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

func UpdateModulesVersion(modules []*Module, inDir string) error {
	for _, module := range modules {
		moduleRoot := filepath.Join(inDir, module.Path)
		hash, err := computeModuleMD5Hash(moduleRoot)
		if err != nil {
			return err
		}

		if hash != module.Hash {
			module.Hash = hash
			module.Version += 1

			fmt.Printf("Module %s has changed. New version: %d. New hash: %s\n", module.Path, module.Version, hash)
		}
	}

	return nil
}

func WriteNewVersionToFile(modules []*Module, inDir string) error {
	for _, module := range modules {
		moduleRoot := filepath.Join(inDir, module.Path)

		moduleConfigPath := filepath.Join(moduleRoot, "module.yml")
		marshaledModule, err := yaml.Marshal(module)
		if err != nil {
			return err
		}

		if err := os.WriteFile(moduleConfigPath, marshaledModule, 0666); err != nil {
			return err
		}
	}

	return nil
}

func computeModuleMD5Hash(moduleRoot string) (string, error) {
	var concatenatedHashes string

	err := filepath.Walk(moduleRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !(filepath.Ext(path) == ".proto") {
			return nil
		}

		hash, err := computeMD5Hash(path)
		if err != nil {
			return err
		}
		concatenatedHashes += hash

		return nil
	})
	if err != nil {
		return "", err
	}

	finalHash := md5.New()
	finalHash.Write([]byte(concatenatedHashes))

	return hex.EncodeToString(finalHash.Sum(nil)), nil
}

func computeMD5Hash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	if err = file.Close(); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
