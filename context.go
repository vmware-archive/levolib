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
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	ProjectName      string
	PackageName      string
	TemplaterVersion string
	Schema           Schema
	Templates        []TemplateInfo
	Mappings         []TemplatesForModels
	Language         string
	TemplateFeatures map[string]bool
	GoAdapter        GoTemplateAdapter
}

type Schema struct {
	Project string
	Models  []Model
}

type Model struct {
	Name       string
	Parent     string
	ParentRef  *Model
	Properties []ModelProperty
}

type ModelProperty struct {
	RemoteIdentifier string
	LocalIdentifier  string
	PropertyType     string
	IsSetType        bool
}

type TemplateInfo struct {
	Language  string
	Version   string
	Directory string
	FileName  string
	Body      []byte
	Adapter   OutputAdapter
}

type TemplatesForModels struct {
	Models    []*Model
	Templates []*TemplateInfo
}

type GeneratedFile struct {
	FileName  string
	Directory string
	Body      []byte
}

func (context *Context) AddModel(model Model) (*Model, error) {
	fmt.Printf("")
	_, err := context.ModelForName(model.Name)
	if err == nil {
		return &Model{}, errors.New("Attempted to add duplicate model with name " + model.Name)
	}
	for index, _ := range model.Properties {
		if model.Properties[index].LocalIdentifier == "" {
			model.Properties[index].LocalIdentifier = model.Properties[index].RemoteIdentifier
		}
	}
	context.Schema.Models = append(context.Schema.Models, model)
	return &(context.Schema.Models[len(context.Schema.Models)-1]), nil
}

func (context *Context) AddModelWithName(name string) (*Model, error) {
	if name == "" {
		return &Model{}, errors.New("Model name must not be empty string")
	}
	model := Model{Name: name}
	return context.AddModel(model)
}

func (context *Context) AddTemplateDirectory(templateDirPath string) ([]TemplateInfo, error) {
	err := filepath.Walk(templateDirPath, context.AddTemplateFile)
	if err != nil {
		return []TemplateInfo{}, err
	}
	templateDirPath = strings.TrimSuffix(templateDirPath, "/")
	for index, template := range context.Templates {
		//Remove the first part of the directory path so that generated
		//files are relative to the working directory, not the template
		//directory
		if templateDirPath != context.Templates[index].Directory {
			context.Templates[index].Directory = template.Directory[len(templateDirPath)+1:]
		} else {
			context.Templates[index].Directory = ""
		}
	}
	return context.Templates, nil
}

func (context *Context) AddTemplateFilePath(filePath string) (TemplateInfo, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return TemplateInfo{}, err
	}
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return TemplateInfo{}, err
	}
	fileName := fileInfo.Name()
	templateInfo, err := context.AddTemplate(fileName, fileContents, "1.0", "", &context.GoAdapter)
	if err != nil {
		return TemplateInfo{}, err
	}

	return *templateInfo, nil
}

func (context *Context) AddTemplateFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info == nil {
		return errors.New("Nil file info when adding template file " + path)
	}
	if info.IsDir() {
		if info.Name() == ".git" || info.Name() == ".hg" {
			return filepath.SkipDir
		}
	} else if info.Name() == ".DS_Store" {
		//do nothing
	} else {
		//it's a template! We should add it!
		fileContents, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fileName := info.Name()
		directory := path[0 : len(path)-len(fileName)]
		_, err = context.AddTemplate(fileName, fileContents, "1.0", directory, &context.GoAdapter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *Context) AddTemplate(fileName string, body []byte, version string, directory string, adapter OutputAdapter) (*TemplateInfo, error) {
	if fileName == "" {
		return &TemplateInfo{}, errors.New("TemplateInfo must have a filename")
	} else if version == "" {
		return &TemplateInfo{}, errors.New("TemplateInfo must have a version")
	}

	if strings.HasSuffix(fileName, ".lt") == false {
		encodedBody := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
		base64.StdEncoding.Encode(encodedBody, body)
		prefix := []byte("<<levobase64>>")
		body = append(prefix, encodedBody...)
	}

	_, err := context.FindTemplate(fileName, directory)
	if err == nil {
		return &TemplateInfo{}, errors.New("Attempted to add duplicate template with name " + fileName)
	}
	templateInfo := TemplateInfo{FileName: fileName, Body: body, Directory: directory, Version: version, Adapter: adapter, Language: context.Language}
	context.Templates = append(context.Templates, templateInfo)
	return &templateInfo, nil
}

func (context *Context) FindTemplate(fileName string, directory string) (*TemplateInfo, error) {
	for _, template := range context.Templates {
		if template.FileName == fileName && template.Directory == directory {
			return &template, nil
		}
	}
	return &TemplateInfo{}, errors.New("Template not found: " + fileName)
}

func (context *Context) TemplateForFileName(fileName string) ([]*TemplateInfo, error) {
	tempTemplates := make([]*TemplateInfo, 0)
	for index, template := range context.Templates {
		if template.FileName == fileName {
			tempTemplates = append(tempTemplates, &context.Templates[index])
		}
	}
	if len(tempTemplates) > 0 {
		return tempTemplates, nil
	}
	return tempTemplates, errors.New("Template not found: " + fileName)
}

func (context *Context) ModelForName(name string) (*Model, error) {
	for _, model := range context.Schema.Models {
		if model.Name == name {
			return &model, nil
		}
	}
	return &Model{}, errors.New("Model not found: " + name)
}

func (context *Context) AddTemplatesForModelsMapping(templateFileNames []string, modelNames []string) error {
	templates := make([]*TemplateInfo, 0)
	models := make([]*Model, 0)
	for _, fileName := range templateFileNames {
		templateInfos, err := context.TemplateForFileName(fileName)
		if err != nil {
			return err
		}
		for _, templateInfo := range templateInfos {
			templates = appendIfUnique(templates, templateInfo)
		}
	}
	for _, modelName := range modelNames {
		model, err := context.ModelForName(modelName)
		if err == nil {
			models = append(models, model)
		} else {
			return err
		}
	}
	if len(templates) <= 0 {
		return errors.New("Mapping must have at least one template")
	}

	context.Mappings = append(context.Mappings, TemplatesForModels{Models: models, Templates: templates})
	return nil
}

func (context *Context) AddTemplateFeature(feature string) {
	context.TemplateFeatures[strings.ToLower(feature)] = true
}

func (context *Context) RemoveTemplateFeature(feature string) {
	context.TemplateFeatures[strings.ToLower(feature)] = false
}

func (model *Model) AddProperty(remoteIdentifier string, localIdentifier string, propertyType string) (*ModelProperty, error) {
	if remoteIdentifier == "" && localIdentifier == "" {
		return &ModelProperty{}, errors.New("Properties must have an identifier")
	}

	for _, property := range model.Properties {
		if property.LocalIdentifier == localIdentifier {
			return &ModelProperty{}, errors.New("Attempted to add duplicate properties with name " + localIdentifier)
		}
	}

	if propertyType == "" {
		return &ModelProperty{}, errors.New("Properties must have a type")
	}
	isSetType := false
	if strings.HasPrefix(propertyType, "[]") {
		propertyType = propertyType[2:]
		isSetType = true
	} else if strings.HasSuffix(propertyType, "[]") {
		propertyType = propertyType[:len(propertyType)-2]
		isSetType = true
	}
	property := ModelProperty{RemoteIdentifier: remoteIdentifier, LocalIdentifier: localIdentifier, PropertyType: propertyType, IsSetType: isSetType}
	model.Properties = append(model.Properties, property)
	return &(model.Properties[len(model.Properties)-1]), nil
}

func (self *Schema) validate() error {
	for _, model := range self.Models {
		if model.Name == "" {
			return errors.New("At least one model missing Remote Identifier")
		}
		for _, property := range model.Properties {
			if property.RemoteIdentifier == "" {
				return errors.New("Model " + model.Name + " has at least one property missing it's Remote Identifier. Type: " + property.PropertyType)
			}
		}
	}
	//TODO fill this out more?
	return nil
}

func appendIfUnique(slice []*TemplateInfo, item *TemplateInfo) []*TemplateInfo {
	for _, existingTemplate := range slice {
		if existingTemplate.FileName == item.FileName && existingTemplate.Directory == item.Directory {
			return slice
		}
	}
	return append(slice, item)
}
