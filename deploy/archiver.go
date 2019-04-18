// Copyright 2019 Karol Marcjan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Archiver struct {
	out string
}

func NewArchiver(cfg Config) *Archiver {
	return &Archiver{out: cfg.Out}
}

// Pack creates an archive in the temporary directory.
//
// srcDir is the path to the directory whose contents should be packaged up.
// pattern is passed to ioutil.TempFile to build the temporary file name.
//
// archive is the name of the created archive, or an empty string, if no file was created.
func (arch *Archiver) Pack(srcDir, pattern string) (archive string, err error) {
	err = arch.ensureTmpExists()
	if err != nil {
		return "", err
	}

	log.Printf("packing contents of directory %q...", srcDir)
	zipFile, err := ioutil.TempFile(arch.out, pattern)
	if err != nil {
		return "", err
	}
	log.Printf("packing into archive file %q...", zipFile.Name())

	zipW := zip.NewWriter(zipFile)

	err = filepath.Walk(srcDir, arch.addTo(zipW, srcDir))
	if err != nil {
		zipW.Close()
		os.Remove(zipFile.Name())
		return "", err
	}

	err = zipW.Close()
	if err != nil {
		os.Remove(zipFile.Name())
		return "", err
	}

	log.Printf("directory %q packed into archive %q", srcDir, zipFile.Name())
	return zipFile.Name(), nil
}

func (arch *Archiver) ensureTmpExists() error {
	err := os.MkdirAll(arch.out, 0755)
	if err != nil {
		return fmt.Errorf("cannot ensure directory %q exists: %s", arch.out, err)
	}
	return nil
}

func (*Archiver) addTo(dst *zip.Writer, srcDir string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		innerPath := strings.TrimPrefix(path, srcDir+"/")
		log.Printf("adding %q to zip file as %q", path, innerPath)
		entry, err := dst.Create(innerPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(entry, file)
		return err
	}
}
