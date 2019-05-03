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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type SpecLoader struct {
	filename string
}

type DeploymentSpec struct {
	Functions []FuncDeploymentSpec `json:"functions"`
}

type FuncDeploymentSpec struct {
	Name        string `json:"name"`
	Entrypoint  string `json:"entrypoint"`
	Src         string `json:"src"`
	Description string `json:"description"`
	Memory      int    `json:"memory"`
}

func NewSpecLoader(cfg Config) *SpecLoader {
	return &SpecLoader{
		filename: cfg.SpecFile,
	}
}

func (p *SpecLoader) Load() (DeploymentSpec, error) {
	log.Printf("reading deployment spec from %s...", p.filename)
	rawSpec, err := ioutil.ReadFile(p.filename)
	if err != nil {
		return DeploymentSpec{}, err
	}

	log.Printf("deserializing deployment spec...")
	spec := DeploymentSpec{}
	err = json.Unmarshal(rawSpec, &spec)
	if err != nil {
		return DeploymentSpec{}, err
	}

	log.Printf("verifying deployment spec...")
	err = p.verify(spec)
	if err != nil {
		return DeploymentSpec{}, err
	}

	return spec, nil
}

func (sl *SpecLoader) verify(spec DeploymentSpec) error {
	for _, fun := range spec.Functions {
		err := sl.verifyFunc(fun)
		if err != nil {
			return fmt.Errorf("function named %q: %s", fun.Name, err)
		}
	}

	return nil
}

func (*SpecLoader) verifyFunc(fun FuncDeploymentSpec) error {
	if strings.TrimSpace(fun.Name) == "" {
		return fmt.Errorf("invalid function name")
	}

	if strings.TrimSpace(fun.Entrypoint) == "" {
		return fmt.Errorf("invalid function entrypoint %q", fun.Entrypoint)
	}

	if fun.Memory != 0 && (fun.Memory < 128 || fun.Memory > 2048) {
		return fmt.Errorf("invalid function memory size %d (not in [%d;%d] range)", fun.Memory, 128, 2048)
	}

	return nil
}
