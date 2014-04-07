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
	"regexp"
	"text/template"
)

type GoTemplateAdapter struct {
	ParsedTemplates *template.Template
}

func (self *GoTemplateAdapter) ParseTemplate(templateInfo TemplateInfo) error {
	if templateInfo.Version != LibraryVersion {
		return errors.New("Expecting templates with version " + LibraryVersion + ". Template " + templateInfo.FileName + " has version " + templateInfo.Version)
	}
	newTemplate := template.New(templateInfo.FileName)
	addCommonUtilitiesToTemplate(newTemplate)
	addJavaUtilitiesToTemplate(newTemplate)
	addObjectiveCUtilitiesToTempalte(newTemplate)
	addRailsUitilitiesToTemplate(newTemplate)

	if _, err := newTemplate.Parse(string(templateInfo.Body)); err != nil {
		return err
	}

	if self.ParsedTemplates == nil {
		//Parse the first template
		self.ParsedTemplates = newTemplate
	} else {
		if _, err := self.ParsedTemplates.AddParseTree(templateInfo.FileName, newTemplate.Tree); err != nil {
			return err
		}
	}
	return nil
}

func (self *GoTemplateAdapter) GenerateFiles(templateInfo TemplateInfo, templateData TemplateData) ([]GeneratedFile, error) {
	fmt.Printf("")
	err := self.cleanTemplateData(&templateData)
	if err != nil {
		return []GeneratedFile{}, err
	}

	buffer := bytes.NewBufferString("")
	err = self.ParsedTemplates.ExecuteTemplate(buffer, templateInfo.FileName, templateData)
	if err != nil {
		return []GeneratedFile{}, err
	}

	return self.GetFilesFromOutput(buffer, templateInfo.Directory)
}

func (self *GoTemplateAdapter) cleanTemplateData(data *TemplateData) error {
	data.PackageName = self.cleanPackageName(data.PackageName)
	data.ProjectName = self.cleanName(data.ProjectName)
	models := make([]Model, 0)
	for _, model := range data.Models {
		models = append(models, self.cleanModel(model))
	}
	data.Models = models
	return nil
}

func (self *GoTemplateAdapter) cleanPackageName(packageName string) string {
	if packageName == "" {
		return "com.example"
	}

	re := regexp.MustCompile("[^A-z0-9.]")
	value := re.ReplaceAll([]byte(packageName), []byte(""))
	return string(value)
}

func (self *GoTemplateAdapter) cleanName(input string) string {
	if input == "" {
		return ""
	}
	reNotAlphaNumeric := regexp.MustCompile("[^A-z0-9_]")
	return string(reNotAlphaNumeric.ReplaceAll([]byte(input), []byte("")))
}

func (self *GoTemplateAdapter) cleanModel(model Model) Model {
	model.Name = self.cleanName(model.Name)
	model.Parent = self.cleanName(model.Parent)
	properties := make([]ModelProperty, 0)
	for _, property := range model.Properties {
		properties = append(properties, self.cleanProperty(property))
	}
	model.Properties = properties
	return model
}

func (self *GoTemplateAdapter) cleanProperty(property ModelProperty) ModelProperty {
	property.LocalIdentifier = self.cleanName(property.LocalIdentifier)
	return property
}

func (self *GoTemplateAdapter) GetFilesFromOutput(buffer *bytes.Buffer, directory string) ([]GeneratedFile, error) {
	generatedFiles := make([]GeneratedFile, 0)
	templateContentsRegex := regexp.MustCompile("<<levo filename:(.*?)( directory:(.*?))?>>((.*\n?)*?)<<levo>>")
	replaceWhitespaceRegex := regexp.MustCompile("[\t ]*!>\n")
	separatedFiles := templateContentsRegex.FindAllStringSubmatch(buffer.String(), -1)
	for _, fileContents := range separatedFiles {
		generatedFile := GeneratedFile{FileName: fileContents[1]}
		if fileContents[3] == "" {
			generatedFile.Directory = directory
		} else {
			generatedFile.Directory = fileContents[3]
		}
		generatedFile.Body = replaceWhitespaceRegex.ReplaceAll([]byte(fileContents[4]), []byte(""))
		generatedFile.Body = bytes.Trim(generatedFile.Body, "\n	 ")
		generatedFiles = append(generatedFiles, generatedFile)
	}
	return generatedFiles, nil
}
