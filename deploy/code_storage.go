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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
)

type CodeStorage struct {
	projectID string
	config    struct {
		bucketName string
		location   string
	}
	storage *storage.Client
}

func NewCodeStorage(config Config, storage *storage.Client) *CodeStorage {
	cs := new(CodeStorage)
	cs.projectID = config.CloudProject
	cs.config.bucketName = config.CodeBucket
	cs.config.location = _Location
	cs.storage = storage
	return cs
}

func (cs *CodeStorage) Upload(archivePath string) (url string, err error) {
	content, checksum, err := cs.readContent(archivePath)
	if err != nil {
		return "", err
	}

	err = cs.ensureBucketExists()
	if err != nil {
		return "", err
	}

	return cs.upload(content, checksum)
}

func (cs *CodeStorage) readContent(path string) (content []byte, checksum string, err error) {
	content, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	rawChecksum := sha256.Sum256(content)
	checksum = base64.RawURLEncoding.EncodeToString(rawChecksum[:])

	return content, checksum, nil
}

func (cs *CodeStorage) ensureBucketExists() error {
	ok, err := cs.bucketExists()
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return cs.createBucket()
}

func (cs *CodeStorage) bucketExists() (bool, error) {
	_, err := cs.bucket().Attrs(context.Background())
	if err == storage.ErrBucketNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (cs *CodeStorage) createBucket() error {
	attrs := cs.bucketAttrs()
	log.Printf("creating bucket %q", attrs.Name)
	return cs.
		bucket().
		Create(context.Background(), cs.projectID, &attrs)
}

func (cs *CodeStorage) upload(content []byte, checksum string) (url string, err error) {
	name := fmt.Sprintf("%s.zip", checksum)
	w := cs.bucketWriter(name)
	defer w.Close()

	_, err = io.Copy(w, bytes.NewBuffer(content))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("gs://%s/%s", cs.config.bucketName, name), nil
}

func (cs *CodeStorage) bucketAttrs() storage.BucketAttrs {
	return storage.BucketAttrs{
		Name:         cs.config.bucketName,
		StorageClass: "REGIONAL",
		Location:     cs.config.location,
	}
}

func (cs *CodeStorage) bucketWriter(name string) *storage.Writer {
	w := cs.
		bucket().
		Object(name).
		If(storage.Conditions{DoesNotExist: true}).
		NewWriter(context.Background())
	w.ContentType = _ApplicationZip
	return w
}

func (cs *CodeStorage) bucket() *storage.BucketHandle {
	return cs.storage.Bucket(cs.config.bucketName)
}
