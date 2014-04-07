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

const TestBasePackage string = "com.test"
const TestLanguage string = "java"
const TestVersion string = "1.0"
const TestProjectName string = "test"
const TestModelName01 string = "People"
const TestModelName02 string = "Cats"
const TestTemplateName01 string = "_Name_.lt"
const TestPropName01 string = "hairy"
const TestPropName02 string = "bald"

var TestInvalidJSONString []byte = []byte("{this is completely {inv:alid json")
var TestValidCorrectConfigString []byte = []byte("{\"TemplaterVersion\":\"" + TestVersion + "\",\"BasePackage\":\"" + TestBasePackage + "\",\"Language\":\"" + TestLanguage + "\",\"Zip\":false,\"ModelSchemaFileName\":\"test-resources/model-schema.json\",\"TemplatesDirectory\":\"test-resources/templates\",\"Mappings\":[{\"ModelNames\":[\"" + TestModelName01 + "\",\"" + TestModelName02 + "\"],\"TemplateNames\":[\"" + TestTemplateName01 + "\"]},{\"ModelNames\":[\"" + TestModelName01 + "\"],\"TemplateNames\":[\"" + TestTemplateName01 + "\"]}]}")
var TestValidIncorrectConfigString []byte = []byte("{\"ModelSchemaFileName\":\"schema.json\",\"TemplatesDirectory\":\"templates/xl-rest_lib-android-v3.0.0\",\"Mappings\":[{\"model_name\":[\"Movie\"],\"templates\":[\"_NAME_Fragment.java\",\"list_item__name_.xml\",\"_NAME_.java\",\"_NAME_Activity.java\",\"_NAME_ContentProvider.java\",\"fragment__name_.xml\",\"activity__name_.xml\",\"_NAME_Validator.java\",\"_NAME_Application.java\",\"_NAME_Table.java\",\"_NAME_ListActivity.java\",\"_NAME_ListValidator.java\",\"Abs_NAME_.java\",\"fragment__name__list.xml\",\"_NAME_ListFragment.java\",\"activity__name__list.xml\",\"Abs_NAME_Table.java\"]},{\"model_name\":[\"Posters\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"Cast\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"Ratings\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[\"MoviesResponse\",\"ReleaseDates\"],\"templates\":[\"_NAME_.java\",\"Abs_NAME_.java\"]},{\"model_name\":[],\"templates\":[\"_NAME_Application.java\"]}]}")

var TestValidSchemaString []byte = []byte("{\"Project\":\"" + TestProjectName + "\",\"Models\":[{\"Name\":\"" + TestModelName01 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName01 + "\",\"PropertyType\":\"string\"}]},{\"Name\":\"" + TestModelName02 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName02 + "\",\"PropertyType\":\"string\"}]}]}")
var TestValidIncorrectSchemaString []byte = []byte("{\"Project\":\"" + TestProjectName + "\",\"Models\":[{\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName01 + "\",\"PropertyType\":\"string\"}]},{\"Name\":\"" + TestModelName02 + "\",\"Parent\":\"\",\"Properties\":[{\"RemoteIdentifier\":\"" + TestPropName02 + "\",\"PropertyType\":\"string\"}]}]}")

func TestParseModelSchemaString(testing *testing.T) {
	fmt.Printf("")
	//Test valid and correct Model Schema JSON
	var schemaAdapter JSONSchemaAdapter = JSONSchemaAdapter{}
	modelSchema, err := schemaAdapter.ParseModelSchemaString(TestValidSchemaString)
	models := modelSchema.Models
	if err != nil {
		testing.Errorf("Error while parsing valid JSON schema: ", err.Error())
	}
	if len(models) != 2 {
		testing.Errorf("Expecting %v models. Got %v", 2, len(models))
	} else {
		model := models[0]
		if model.Name != TestModelName01 {
			testing.Errorf("Expecting model named %v. Got %v", TestModelName01, models[0].Name)
		}
		if len(model.Properties) != 1 {
			testing.Errorf("Expecting %v properties. Got %v", 2, len(model.Properties))
		} else {
			if model.Properties[0].RemoteIdentifier != TestPropName01 {
				testing.Errorf("Expecting property named %v. Got %v", TestPropName01, model.Properties[0].RemoteIdentifier)
			}
		}
	}

	//Test invalid JSON
	schemaAdapter = JSONSchemaAdapter{}
	modelSchema, err = schemaAdapter.ParseModelSchemaString(TestInvalidJSONString)
	models = modelSchema.Models
	if err == nil {
		testing.Errorf("ParseModelSchemaString did not fail when passed invalid JSON")
	}

	//Test valid but incorrect Model Schema JSON
	schemaAdapter = JSONSchemaAdapter{}
	modelSchema, err = schemaAdapter.ParseModelSchemaString(TestValidIncorrectSchemaString)
	models = modelSchema.Models
	if err == nil {
		testing.Errorf("ParseModelSchemaString did not fail when passed incorrect JSON")
	}
}
