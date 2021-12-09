// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/elastic/elastic-agent-poc/licenses"
	"github.com/elastic/elastic-agent-poc/elastic-agent/pkg/packer"
)

var (
	input   string
	output  string
	license string
)

func init() {
	flag.StringVar(&input, "in", "", "config to embed")
	flag.StringVar(&output, "out", "-", "Output path. \"-\" means writing to stdout")
	flag.StringVar(&license, "license", "Elastic", "License header for generated file.")
}

var tmpl = template.Must(template.New("cfg").Parse(`
{{ .License }}
// Code generated by dev-tools/cmd/buildfleetcfg/buildfleetcfg.go - DO NOT EDIT.

package application

import "github.com/elastic/elastic-agent-poc/elastic-agent/pkg/packer"

// DefaultAgentFleetConfig is the content of the default configuration when we enroll a beat, the elastic-agent.yml
// will be replaced with this variables.
var DefaultAgentFleetConfig []byte

func init() {
	// Packed File
	{{ range $i, $f := .Files -}}
	// {{ $f }}
	{{ end -}}
	unpacked := packer.MustUnpack("{{ .Pack }}")
	raw, ok := unpacked["_meta/elastic-agent.fleet.yml"]
	if !ok {
		// ensure we have something loaded.
		panic("elastic-agent.fleet.yml is not included in the binary")
	}
	DefaultAgentFleetConfig = raw
}
`))

func main() {
	flag.Parse()

	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "Invalid input source")
		os.Exit(1)
	}

	l, err := licenses.Find(license)
	if err != nil {
		fmt.Fprintf(os.Stderr, "problem to retrieve the license, error: %+v", err)
		os.Exit(1)
	}

	data, err := gen(input, l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while generating the file, err: %+v\n", err)
		os.Exit(1)
	}

	if output == "-" {
		os.Stdout.Write(data)
		return
	}

	ioutil.WriteFile(output, data, 0640)
	return
}

func gen(path string, l string) ([]byte, error) {
	pack, files, err := packer.Pack(input)
	if err != nil {
		return nil, err
	}

	if len(files) > 1 {
		return nil, fmt.Errorf("Can only embed a single configuration file")
	}

	var buf bytes.Buffer
	tmpl.Execute(&buf, struct {
		Pack    string
		Files   []string
		License string
	}{
		Pack:    pack,
		Files:   files,
		License: l,
	})

	return format.Source(buf.Bytes())
}
