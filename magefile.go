// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"go.uber.org/multierr"

	devtools "github.com/elastic/elastic-agent-poc/dev-tools/mage"
	"github.com/elastic/elastic-agent-poc/dev-tools/mage/gotool"
)

// Fmt formats code and adds license headers.
func Fmt() {
	mg.Deps(devtools.GoImports, devtools.PythonAutopep8)
	mg.Deps(AddLicenseHeaders)
}

// AddLicenseHeaders adds ASL2 headers to .go files outside of x-pack and
// add Elastic headers to .go files in x-pack.
func AddLicenseHeaders() error {
	fmt.Println(">> fmt - go-licenser: Adding missing headers")

	mg.Deps(devtools.InstallGoLicenser)

	licenser := gotool.Licenser

	return multierr.Combine(
		licenser(
			licenser.Check(),
			licenser.License("ASL2"),
			licenser.Exclude("elastic-agent"),
			licenser.Exclude("generator/_templates/beat/{beat}"),
			licenser.Exclude("generator/_templates/metricbeat/{beat}"),
		),
		licenser(
			licenser.License("Elastic"),
			licenser.Path("elastic-agent"),
		),
	)
}

// CheckLicenseHeaders checks ASL2 headers in .go files outside of x-pack and
// checks Elastic headers in .go files in x-pack.
func CheckLicenseHeaders() error {
	fmt.Println(">> fmt - go-licenser: Checking for missing headers")

	mg.Deps(devtools.InstallGoLicenser)

	licenser := gotool.Licenser

	return multierr.Combine(
		licenser(
			licenser.Check(),
			licenser.License("ASL2"),
			licenser.Exclude("elastic-agent"),
			licenser.Exclude("generator/_templates/beat/{beat}"),
			licenser.Exclude("generator/_templates/metricbeat/{beat}"),
		),
		licenser(
			licenser.Check(),
			licenser.License("Elastic"),
			licenser.Path("elastic-agent"),
		),
	)
}

// DumpVariables writes the template variables and values to stdout.
func DumpVariables() error {
	return devtools.DumpVariables()
}

