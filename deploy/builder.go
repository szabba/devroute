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

const _BackendArchivePattern = "backend-*.zip"

type Builder struct {
	archiver *Archiver

	backend struct {
		dir     string
		archive string
	}

	err error
}

type BuildResult struct {
	BackendArchive string
}

func NewBuilder(cfg Config, archiver *Archiver) *Builder {
	bld := new(Builder)
	bld.backend.dir = cfg.BackendModule
	bld.archiver = archiver
	return bld
}

func (bld *Builder) Build() (BuildResult, error) {
	if bld.err != nil {
		return bld.result()
	}
	bld.buildBackend()
	return bld.result()
}

func (bld *Builder) result() (BuildResult, error) {
	if bld.err != nil {
		return BuildResult{}, bld.err
	}
	return BuildResult{
		BackendArchive: bld.backend.archive,
	}, nil
}

func (bld *Builder) buildBackend() {
	bld.backend.archive, bld.err = bld.archiver.Pack(bld.backend.dir, _BackendArchivePattern)
}
