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
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type JSONSchemaAdapter struct{}

func (self *JSONSchemaAdapter) ProcessSchemaFile(schemaPath string) (Schema, error) {
	fmt.Printf("")
	fileContents, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return Schema{}, err
	}
	return self.ParseModelSchemaString(fileContents)
}

func (self *JSONSchemaAdapter) ParseModelSchemaString(schemaString []byte) (Schema, error) {
	var schema Schema
	err := json.Unmarshal(schemaString, &schema)
	if err != nil {
		return Schema{}, err
	}
	if err := schema.validate(); err != nil {
		return Schema{}, err
	}
	for _, model := range schema.Models {
		for _, prop := range model.Properties {
			if prop.LocalIdentifier == "" {
				prop.LocalIdentifier = prop.RemoteIdentifier
			}
		}
	}
	return schema, nil
}
