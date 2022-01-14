// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package main

import (
	"bytes"
	"flag"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

var Template = template.Must(template.New("licenseheader").Parse(`
{{ $t := "` + "`" + `" }}
{{ .License }}

// Code generated by beats/dev-tools/cmd/license/license_generate.go - DO NOT EDIT.

package licenses

import "fmt"

{{ range $key, $value := .Licenses }}
var {{ $key }} =  {{$t}}
{{ $value }}{{$t}}
{{ end -}}

func Find(name string) (string, error) {
	switch name {
{{ range $key, $value := .Licenses }}
	case "{{ $key }}":
		return {{ $key }}, nil
{{- end -}}
	}
	return "", fmt.Errorf("unknown license: %s", name)
}
`))

var output string

type data struct {
	License  string
	Licenses map[string]string
}

func init() {
	flag.StringVar(&output, "out", "license_header.go", "output file")
}

func main() {
	Headers := make(map[string]string)
	content, err := ioutil.ReadFile("APACHE-LICENSE-2.0-header.txt")
	if err != nil {
		panic("could not read ASL2 license.")
	}
	Headers["ASL2"] = string(content)

	content, err = ioutil.ReadFile("ELASTIC-LICENSE-header.txt")
	if err != nil {
		panic("could not read Elastic license.")
	}
	Headers["Elastic"] = string(content)

	content, err = ioutil.ReadFile("ELASTIC-LICENSE-2.0-header.txt")
	if err != nil {
		panic("could not read Elastic License 2.0 license.")
	}
	Headers["Elasticv2"] = string(content)

	var buf bytes.Buffer
	Template.Execute(&buf, data{
		License:  Headers["ASL2"],
		Licenses: Headers,
	})

	bs, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if output == "-" {
		os.Stdout.Write(bs)
	} else {
		ioutil.WriteFile(output, bs, 0640)
	}
}
