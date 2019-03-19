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
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Timeout       time.Duration
	Out           string
	BackendModule string
	CloudProject  string
	CodeBucket    string
	SpecFile      string
}

func ParseConfig() Config {
	out := flag.CommandLine.Output()
	var cfg Config

	flag.Usage = func() {
		fmt.Fprintf(out, "%s is a tool for deploying the project\n\n", os.Args[0])
		fmt.Fprintf(out, "usage:\n\n")
		fmt.Fprintf(out, "\t%s [FLAG]... PROJECT_NAME\n\n", os.Args[0])
		fmt.Fprintf(out, "flag defaults:\n\n")
		flag.PrintDefaults()
	}

	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "the Google API request timeout")
	flag.StringVar(&cfg.Out, "out", "out", "the output diretory to use")
	flag.StringVar(&cfg.CodeBucket, "code-bucket", "devroute-dev-code", "the Google Cloud Storage bucket to use for code uploads")
	flag.StringVar(&cfg.BackendModule, "backend-module", "backend/functions/hello", "the directory containing the Go module with Cloud Functions")
	flag.StringVar(&cfg.SpecFile, "spec", "backend/deploy.json", "the file specifying the deployment details")
	flag.Parse()

	cfg.CloudProject = flag.Arg(0)
	if cfg.CloudProject == "" {
		flag.Usage()
		fmt.Fprintf(out, "\nmissing PROJECT_NAME\n")
		os.Exit(1)
	}

	return cfg
}
