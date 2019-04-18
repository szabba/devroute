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
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/googleapi"
)

const (
	_ApplicationZip = "application/zip"
	_Runtime        = "go111"
	_Location       = "europe-west1"
)

type Deployer struct {
	projectID   string
	codeBucket  string
	location    string
	codeStorage *CodeStorage
	functions   *cloudfunctions.Service
}

func NewDeployer(config Config, codeStorage *CodeStorage, functions *cloudfunctions.Service) *Deployer {
	return &Deployer{
		projectID:   config.CloudProject,
		location:    _Location,
		codeStorage: codeStorage,
		functions:   functions,
	}
}

func (depl *Deployer) Deploy(build BuildResult, spec DeploymentSpec) error {
	log.Printf("uploading the backend source code archive...")
	archiveURL, err := depl.codeStorage.Upload(build.BackendArchive)
	if err != nil {
		return err
	}
	log.Printf("backend code at %q", archiveURL)

	for _, fn := range spec.Functions {
		err := depl.deployFunction(fn, archiveURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func (depl *Deployer) deployFunction(fn FuncDeploymentSpec, archiveURL string) error {
	location := fmt.Sprintf("projects/%s/locations/%s", depl.projectID, depl.location)
	apiFn := depl.toApiFn(fn, location, archiveURL)

	_, err := depl.functions.Projects.Locations.Functions.Get(apiFn.Name).Do()
	if err, ok := err.(*googleapi.Error); ok && err.Code == http.StatusNotFound {
		return depl.deployNew(location, &apiFn)
	}
	if err != nil {
		return err
	}
	return depl.deployUpdated(&apiFn)
}

func (depl *Deployer) toApiFn(fn FuncDeploymentSpec, location, archiveURL string) cloudfunctions.CloudFunction {
	name := fmt.Sprintf("%s/functions/%s", location, fn.Name)
	return cloudfunctions.CloudFunction{
		Name:             name,
		Runtime:          _Runtime,
		EntryPoint:       fn.Entrypoint,
		Description:      fn.Description,
		SourceArchiveUrl: archiveURL,
		HttpsTrigger:     &cloudfunctions.HttpsTrigger{},
	}
}

func (depl *Deployer) deployNew(location string, fn *cloudfunctions.CloudFunction) error {
	log.Printf("deploying function %q for the first time...", fn.Name)
	_, err := depl.functions.Projects.Locations.Functions.Create(location, fn).Do()
	return err
}

func (depl *Deployer) deployUpdated(fn *cloudfunctions.CloudFunction) error {
	log.Printf("deploying updated function %q...", fn.Name)
	_, err := depl.functions.Projects.Locations.Functions.Patch(fn.Name, fn).Do()
	return err
}
