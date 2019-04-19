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
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"go.uber.org/dig"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/option"
)

var _GoogleCloudOAuth2Scopes = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/devstorage.read_write",
}

func main() {
	dic := dig.New()
	dic.Provide(ParseConfig)
	dic.Provide(NewArchiver)
	dic.Provide(NewStorageClient)
	dic.Provide(NewCloudFunctionsService)
	dic.Provide(NewSpecLoader)
	dic.Provide(NewCodeStorage)
	dic.Provide(NewDeployer)
	err := dic.Invoke(Run)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func Run(specLoader *SpecLoader, depl *Deployer) error {
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("working at %s", workdir)

	spec, err := specLoader.Load()
	if err != nil {
		return err
	}

	return depl.Deploy(spec)
}

func NewCloudFunctionsService() (*cloudfunctions.Service, error) {
	return cloudfunctions.NewService(context.Background(), option.WithScopes(_GoogleCloudOAuth2Scopes...))
}

func NewStorageClient() (*storage.Client, error) {
	return storage.NewClient(context.Background(), option.WithScopes(_GoogleCloudOAuth2Scopes...))
}
