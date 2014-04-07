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
	"errors"
	"fmt"
	"os"
	"strings"
)

const templaterVersion string = "1.0"

type TemplateData struct {
	PackageName string
	ProjectName string
	PackagePath string
	Models      []Model
	Features    map[string]bool
}

type OutputAdapter interface {
	ParseTemplate(template TemplateInfo) error
	GenerateFiles(templateInfo TemplateInfo, templateData TemplateData) ([]GeneratedFile, error)
	GetFilesFromOutput(buffer *bytes.Buffer, directory string) ([]GeneratedFile, error)
}

type ConfigurationAdapter interface {
	ProcessConfigurationFile(configFile *os.File) (Context, error)
	ProcessConfigurationString(configString string) (Context, error)
}

type SchemaAdapter interface {
	ProcessSchemaFile(schemaFile *os.File) (Schema, error)
	ProcessSchemaString(schemaString string) (Schema, error)
}

func ProcessMappings(context Context) ([]GeneratedFile, error) {
	if len(context.Mappings) == 0 {
		return []GeneratedFile{}, errors.New("No mappings to process")
	}

	//Pre-parse all the templates.
	//For Go templates this will compile them all into one associated
	//set. This way templates can reference eachother.
	for _, templateInfo := range context.Templates {
		err := templateInfo.Adapter.ParseTemplate(templateInfo)
		if err != nil {
			return []GeneratedFile{}, err
		}
	}

	generatedFiles := make([]GeneratedFile, 0, 0)
	for _, mapping := range context.Mappings {
		templateModels := make([]Model, 0, 0)
		for _, model := range mapping.Models {
			templateModels = append(templateModels, *model)
		}
		for _, templateInfo := range mapping.Templates {
			newestFiles, err := processTemplate(templateInfo, context, templateModels)
			if err != nil {
				return []GeneratedFile{}, err
			}
			for _, newestFile := range newestFiles {
				generatedFiles = append(generatedFiles, newestFile)
			}
		}
	}
	return generatedFiles, nil
}

func processTemplate(templateInfo *TemplateInfo, context Context, templateModels []Model) ([]GeneratedFile, error) {
	if strings.HasSuffix(strings.ToLower(templateInfo.FileName), ".lt") == false {
		//The template is not a Levo template (.lt)
		//Treat it like a static binary file
		return getBinaryFiles(templateInfo)
	} else {
		templateData := TemplateData{PackageName: context.PackageName, ProjectName: context.ProjectName}
		templateData.PackagePath = strings.Replace(context.PackageName, ".", "/", -1)
		templateData.Models = templateModels
		templateData.Features = context.TemplateFeatures
		return templateInfo.Adapter.GenerateFiles(*templateInfo, templateData)
	}
}

func getBinaryFiles(templateInfo *TemplateInfo) ([]GeneratedFile, error) {
	bodyBuffer := bytes.NewBuffer(templateInfo.Body)
	filesForTemplate, err := templateInfo.Adapter.GetFilesFromOutput(bodyBuffer, templateInfo.Directory)
	if err != nil {
		return []GeneratedFile{}, err
	}
	if len(filesForTemplate) > 0 {
		return filesForTemplate, nil
	}
	newFile := GeneratedFile{FileName: templateInfo.FileName, Directory: templateInfo.Directory, Body: templateInfo.Body}
	return []GeneratedFile{newFile}, nil
}

func BeginContext() Context {
	fmt.Printf("")
	return Context{PackageName: "com.example", ProjectName: "ExampleProject", TemplaterVersion: templaterVersion, GoAdapter: GoTemplateAdapter{}, TemplateFeatures: make(map[string]bool, 0)}
}

func GetJSONSchemaAdapter() JSONSchemaAdapter {
	return JSONSchemaAdapter{}
}
