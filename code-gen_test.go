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
	"fmt"
	"testing"
)

const ExpectedBody string = "package TestPackage01.models;"

func TestProcessMappings(testing *testing.T) {
	fmt.Printf("")
	SetupContext()
	SetupTemplate()
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
		if string(generatedFiles[0].Body) != ExpectedBody {
			testing.Errorf("Was expecting file with contents:\n%v.\n\nGot:\n%v.", ExpectedBody, string(generatedFiles[0].Body))
		}
	}
}

func TestProcessTemplate(testing *testing.T) {
	SetupContext()
	SetupTemplate()
	SetupModel()

	for _, templateInfo := range context.Templates {
		templateInfo.Adapter.ParseTemplate(templateInfo)
	}

	generatedFiles, err := processTemplate(&context.Templates[0], context, []Model{context.Schema.Models[0]})
	if err != nil {
		testing.Errorf("Unexpected error: %v", err.Error())
	} else if len(generatedFiles) != 1 {
		testing.Errorf("Expecting %v templates. Got %v", 1, len(generatedFiles))
	}

	generatedFiles, err = processTemplate(&context.Templates[1], context, []Model{context.Schema.Models[0]})
	if err != nil {
		testing.Errorf("Unexpected error: %v", err.Error())
	} else if len(generatedFiles) != 1 {
		testing.Errorf("Expecting %v templates. Got %v", 1, len(generatedFiles))
	}
}
