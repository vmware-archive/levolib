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
	"reflect"
	"testing"
)

const testPackageName string = "TestPackage01"
const testProjectName string = "TestProject01"
const testLanguage string = "test"
const testTemplaterVersion string = "1.0.0"
const modelName string = "TestModel01"
const brokenModelName string = "NonExistentModelººº"
const propertyRemoteIdent string = "TestPropRemote01"
const propertyLocalIdent string = "TestPropLocal01"
const propertyPropertyType string = "string"
const templateFileName string = "Template01.lt"
const sourceFileName string = "TestProject01.generic"
const obviouslyBrokenFileName string = "::*&thisºººcan'twork.&&&"
const manyToOneRelationship string = "ManyToOne"
const manyToManyRelationship string = "ManyToMany"
const templateBody string = "<<levo filename:{{.ProjectName}}.generic>>\npackage {{.PackageName}}.models;\n<<levo>>\n"

var context Context
var templateInfo TemplateInfo

func SetupContext() Context {
	fmt.Printf("")
	context = BeginContext()
	context.PackageName = testPackageName
	context.ProjectName = testProjectName
	context.Language = testLanguage
	context.TemplaterVersion = testTemplaterVersion
	return context
}

func SetupTemplate() {
	context.AddTemplate(templateFileName, []byte(templateBody), testTemplaterVersion, "template", &context.GoAdapter)
	context.AddTemplate("binary.template", []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."), testTemplaterVersion, "template", &context.GoAdapter)
}

func SetupModel() {
	context.AddModelWithName(modelName)
}

func TestContextAddModelWithName(testing *testing.T) {
	SetupContext()
	originalModels := context.Schema.Models
	context.AddModelWithName(modelName)

	//check if there exactly one more model than before.
	if (len(originalModels) + 1) != len(context.Schema.Models) {
		testing.Errorf("Model count did not increase by exactly one, change was %v", len(context.Schema.Models)-len(originalModels)+1)
	}

	//check that the added model has the correct data (name)
	if context.Schema.Models[len(context.Schema.Models)-1].Name != modelName {
		testing.Errorf("Expecting ModelName '%s'. Got '%s'", modelName, context.Schema.Models[0].Name)
	}

	//check that all the other models still exist
	for modelIndex, model := range originalModels {
		if model.Name != context.Schema.Models[modelIndex].Name {
			testing.Errorf("Unexpected model change when adding. the changed model name was %v", model.Name)
		}
	}

	//check that adding a model with the same name as one that already exists fails
	duplicateModelCountCheck := len(context.Schema.Models)
	_, err := context.AddModelWithName(modelName)
	if err == nil {
		testing.Errorf("No error returned when adding model with duplicate name %v", modelName)
	} else if duplicateModelCountCheck != len(context.Schema.Models) {
		testing.Errorf("Model with identical name %v was added (it shouldn't have been)", modelName)
	}

	//check that a model without a name fails
	unnamedModelCountCheck := len(context.Schema.Models)
	_, err = context.AddModelWithName("")
	if err == nil {
		testing.Errorf("No error returned when adding model with no name")
	} else if unnamedModelCountCheck != len(context.Schema.Models) {
		testing.Errorf("Model with no name was added (it shouldn't have been)")
	}
}

func TestModelAddProperty(testing *testing.T) {
	SetupContext()
	SetupModel()

	model := context.Schema.Models[0]

	originalProperties := model.Properties
	property, err := model.AddProperty(propertyRemoteIdent, propertyLocalIdent, propertyPropertyType)
	if err != nil {
		testing.Errorf("Error adding valid property: %v", err.Error())
	}

	//check if there is exactly one more property than before
	if len(originalProperties)+1 != len(model.Properties) {
		testing.Errorf("Property count did not increase by exactly one, change was %v", len(model.Properties)-len(originalProperties)+1)
	}

	//check that the added property has the correct data (api identifier)
	if property.RemoteIdentifier != propertyRemoteIdent {
		testing.Errorf("Expecting RemoteIdentifier '%s'. Got '%s'", propertyRemoteIdent, property.RemoteIdentifier)
	}

	//check that adding a property with duplicate data fails
	duplicatePropertyCountCheck := len(model.Properties)
	_, err = model.AddProperty(propertyRemoteIdent, propertyLocalIdent, propertyPropertyType)
	if err == nil {
		testing.Errorf("No error returned when adding property with duplicate data")
	} else if duplicatePropertyCountCheck != len(model.Properties) {
		testing.Errorf("Property with identical data was added (it shouldn't have been)")
	}

	//Test adding a no-name property
	if property, err = model.AddProperty("", "", propertyPropertyType); err == nil {
		testing.Errorf("No error returned when adding a property with no name")
	} else if reflect.DeepEqual(*property, ModelProperty{}) == false {
		testing.Errorf("Non-empty property returned when adding a property with no name")
	}

	//Test adding a un-typed property
	property, err = model.AddProperty(propertyRemoteIdent+"02", propertyLocalIdent+"02", "")
	if err == nil {
		testing.Errorf("No error returned when adding a property with no type")
	} else if reflect.DeepEqual(*property, ModelProperty{}) == false {
		testing.Errorf("Non-empty property returned when adding a property with no type")
	}

	//TODO add test for invalid type, like "potato"
}

func TestContextAddTemplateDirectory(testing *testing.T) {
	SetupContext()
	templates, err := context.AddTemplateDirectory("test-resources/templates")
	if err != nil {
		testing.Errorf("Unexpected error: %v", err.Error())
	}
	for _, templateInfo := range templates {
		if len(templateInfo.Directory) >= 6 && templateInfo.Directory[0:6] != "subdir" {
			testing.Errorf("Template directory not properly trimmed: %v", templateInfo.Directory)
		}
	}
}

func TestContextAddFilePath(testing *testing.T) {
	SetupContext()
	template, err := context.AddTemplateFilePath("test-resources/templates/_Name_.lt")
	if err != nil {
		testing.Errorf("Unexpected error: %v", err.Error())
	}
	if template.FileName != "_Name_.lt" {
		testing.Errorf("Expecting template named %v. Got %v.", "_Name_.lt", template.FileName)
	}
	template, err = context.AddTemplateFilePath("test-resources/templates/NotARealFile.lt")
	if err == nil {
		testing.Errorf("No error when trying to add non-existing file")
	}
}

func TestContextAddTemplate(testing *testing.T) {
	SetupContext()
	originalTemplates := context.Templates
	contextCopy := context

	if template, err := contextCopy.AddTemplate("", []byte(templateBody), testTemplaterVersion, "template", &context.GoAdapter); err == nil {
		testing.Error("No error when adding template without filename")
	} else if reflect.DeepEqual(*template, TemplateInfo{}) == false {
		testing.Error("Non empty template returned when adding template without filename")
	}

	SetupTemplate()
	//check if there is exactly one more template than before
	if len(originalTemplates)+2 != len(context.Templates) {
		testing.Errorf("Template count did not increase by exactly one, change was %v", len(context.Templates))
	}

	//check that the added template has the correct data (file name)
	if len(context.Templates) != 2 {
		testing.Errorf("Expecting exactly 1 template. Got '%v'", len(context.Templates))
	}

	//check that adding a template with duplicate data fails
	duplicateTemplateCountCheck := len(context.Templates)
	_, err := context.AddTemplate(templateFileName, []byte(templateBody), "1.0.0", "template", &context.GoAdapter)
	if err == nil {
		testing.Errorf("No error returned when adding template with duplicate data")
	} else if duplicateTemplateCountCheck != len(context.Templates) {
		testing.Errorf("Template with identical data was added (it shouldn't have been)")
	}
}

func TestContextAddTemplatesForModelsMapping(testing *testing.T) {
	SetupContext()
	SetupTemplate()
	SetupModel()
	fileNames := []string{templateFileName}
	modelNames := []string{modelName}

	//Test valid template and model mapping
	err := context.AddTemplatesForModelsMapping(fileNames, modelNames)

	if err != nil {
		testing.Error(err)
	}
	if len(context.Mappings) == 0 {
		testing.Errorf("No Mappings")
	} else {
		if len(context.Mappings[0].Models) < 1 {
			testing.Errorf("No Models")
		}
		if len(context.Mappings[0].Templates) == 0 {
			testing.Errorf("No Templates")
		}
		if context.Mappings[0].Models[0].Name != modelName {
			testing.Errorf("Expecting modelname '%s'. Got '%s'", modelName, context.Mappings[0].Models[0].Name)
		}
		if context.Mappings[0].Templates[0].FileName != templateFileName {
			testing.Errorf("Expecting modelname '%s'. Got '%s'", templateFileName, context.Mappings[0].Templates[0].FileName)
		}
	}

	//Test mapping with empty models
	if err = context.AddTemplatesForModelsMapping(fileNames, make([]string, 0, 0)); err == nil {
		//We might actually accept templates with no models
		//testing.Errorf("No Templates")
	}

	//Test mapping with empty templates
	if err = context.AddTemplatesForModelsMapping(make([]string, 0, 0), modelNames); err == nil {
		testing.Errorf("No error returned when mapping with zero templates is added")
	}

	//Test mapping with empty models and templates
	if err = context.AddTemplatesForModelsMapping(make([]string, 0, 0), make([]string, 0, 0)); err == nil {
		testing.Errorf("No error returned when mapping with zero templates and zero models is added")
	}

	//Test mapping with invalid model
	if err = context.AddTemplatesForModelsMapping(fileNames, append(modelNames, brokenModelName)); err == nil {
		testing.Errorf("No error returned when mapping with broken models is added")
	}

	//Test mapping with invalid template
	if err = context.AddTemplatesForModelsMapping(append(fileNames, obviouslyBrokenFileName), modelNames); err == nil {
		testing.Errorf("No error returned when mapping with broken models is added")
	}
}

func TestContextTemplateForFileName(testing *testing.T) {
	SetupContext()
	SetupTemplate()

	//test a working filename
	if templates, err := context.TemplateForFileName(templateFileName); err != nil {
		testing.Errorf("Error returned for valid template filename: %v", err.Error())
	} else if reflect.DeepEqual(*templates[0], TemplateInfo{}) == true {
		testing.Errorf("Empty template returned for valid filename")
	}

	//test a blank filename
	if _, err := context.TemplateForFileName(""); err == nil {
		testing.Errorf("No error returned for blank template filename")
	}
	//test a broken filename
	if _, err := context.TemplateForFileName(obviouslyBrokenFileName); err == nil {
		testing.Errorf("No error returned for broken template filename")
	}
}

func TestContextModelForName(testing *testing.T) {
	SetupContext()
	SetupModel()

	//test a working model name
	if model, err := context.ModelForName(modelName); err != nil {
		testing.Errorf("Error returned for valid model name: %v", err.Error())
	} else if reflect.DeepEqual(*model, Model{}) == true {
		testing.Errorf("Empty model returned for valid model name")
	}
	//test a blank model name
	if model, err := context.ModelForName(""); err == nil {
		testing.Errorf("No error returned for blank model name")
	} else if reflect.DeepEqual(*model, Model{}) == false {
		testing.Errorf("Non-empty model returned for blank model name")
	}
	//test an obviously broken model name
	if model, err := context.ModelForName(brokenModelName); err == nil {
		testing.Errorf("No error returned for broken model name")
	} else if reflect.DeepEqual(*model, Model{}) == false {
		testing.Errorf("Non-empty model returned for broken model name")
	}
}
