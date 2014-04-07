/* Copyright (C) 2014 Pivotal Software, Inc.

All rights reserved. This program and the accompanying materials
are made available under the terms of the under the Apache License,
Version 2.0 (the "License”); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/
package levo

import (
	"bytes"
	"fmt"
	"testing"
)

func SetupMapping() {
	fileNames := []string{templateFileName, "binary.template"}
	modelNames := []string{modelName}
	context.AddTemplatesForModelsMapping(fileNames, modelNames)
}

func TestParseTemplates(testing *testing.T) {
	SetupContext()
	SetupTemplate()

	err := context.GoAdapter.ParseTemplate(context.Templates[0])
	if err != nil {
		testing.Errorf("Unexpected Error: %v", err.Error())
	} else if len(context.GoAdapter.ParsedTemplates.Templates()) != 1 {
		testing.Errorf("Expected %v template. Got %v.", 1, len(context.GoAdapter.ParsedTemplates.Templates()))
	}
	parsedTemplate := context.GoAdapter.ParsedTemplates.Lookup("Template01.lt")
	if parsedTemplate == nil {
		testing.Errorf("Did not find template with name _Name_.lt")
	}

	err = context.GoAdapter.ParseTemplate(context.Templates[1])
	if err != nil {
		testing.Errorf("Unexpected Error: %v", err.Error())
	} else if len(context.GoAdapter.ParsedTemplates.Templates()) != 2 {
		testing.Errorf("Expected %v template. Got %v.", 2, len(context.GoAdapter.ParsedTemplates.Templates()))
	}
	parsedTemplate = context.GoAdapter.ParsedTemplates.Lookup("binary.template")
	if parsedTemplate == nil {
		testing.Errorf("Did not find template with name binary.template")
	}
}

func TestGenerateFiles(testing *testing.T) {
	fmt.Printf("")
	SetupContext()
	context.AddTemplate(templateFileName, []byte(templateBody), testTemplaterVersion, "template", &context.GoAdapter)
	context.AddTemplateFilePath("test-resources/templates/binary.template")
	context.Templates[0].Adapter = &context.GoAdapter
	SetupModel()
	SetupMapping()

	generatedFiles, err := ProcessMappings(context)
	if err != nil {
		testing.Error(err)
	}
	if len(generatedFiles) != 2 {
		testing.Errorf("Was expecting %v generated files. Got %v", 2, len(generatedFiles))
	} else {
		if generatedFiles[0].FileName != sourceFileName {
			testing.Errorf("Was expecting file named %v. Got %v", sourceFileName, generatedFiles[0].FileName)
		}
	}
}

func TestGetFilesFromOutput(testing *testing.T) {
	firstBuffer := bytes.NewBufferString("<<levo filename:first>>\nContents\n<<levo>>")
	directory := "default"

	adapter := GoTemplateAdapter{}
	files, err := adapter.GetFilesFromOutput(firstBuffer, directory)
	if err != nil {
		testing.Errorf("Unexpected Error: %v", err.Error())
	} else if len(files) != 1 {
		testing.Errorf("Expecting %v file. Got %v.", 1, len(files))
	} else if files[0].FileName != "first" {
		testing.Errorf("Expecting filename %v. Got %v.", "first", files[0].FileName)
	} else if files[0].Directory != "default" {
		testing.Errorf("Expecting directory %v. Got %v.", "default", files[0].Directory)
	}

	secondBuffer := bytes.NewBufferString("<<levo filename:first directory:subdir>>\nContents\n<<levo>>\n<<levo filename:second directory:subdir>>\nContents\n<<levo>>")
	files, err = adapter.GetFilesFromOutput(secondBuffer, directory)
	if err != nil {
		testing.Errorf("Unexpected Error: %v", err.Error())
	} else if len(files) != 2 {
		testing.Errorf("Expecting %v file. Got %v.", 2, len(files))
	} else if files[1].FileName != "second" {
		testing.Errorf("Expecting filename %v. Got %v.", "second", files[1].FileName)
	} else if files[1].Directory != "subdir" {
		testing.Errorf("Expecting directory %v. Got %v.", "subdir", files[1].Directory)
	}
}
