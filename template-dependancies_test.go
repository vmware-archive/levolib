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
	"testing"
	"text/template"
)

func TestSplitByWord(testing *testing.T) {
	if output := splitByWord("One"); len(output) != 1 {
		testing.Errorf("Expecting %v. Got %v", 1, len(output))
	}
	if output := splitByWord("One Two"); len(output) != 2 {
		testing.Errorf("Expecting %v. Got %v", 2, len(output))
	}
	if output := splitByWord("One Two_Three"); len(output) != 3 {
		testing.Errorf("Expecting %v. Got %v", 3, len(output))
	}
	if output := splitByWord("RedBlue"); len(output) != 2 {
		testing.Errorf("Expecting %v. Got %v", 2, len(output))
	}
	if output := splitByWord("RedBlueGreen"); len(output) != 3 {
		testing.Errorf("Expecting %v. Got %v", 3, len(output))
	}
	if output := splitByWord("REDBLUEGREEN"); len(output) != 1 {
		testing.Errorf("Expecting %v. Got %v", 1, len(output))
	}
	if output := splitByWord("HTTPController"); len(output) != 2 {
		testing.Errorf("Expecting %v. Got %v", 2, len(output))
	}
	if output := splitByWord("Forty40forty"); len(output) != 3 {
		testing.Errorf("Expecting %v. Got %v", 3, len(output))
	}
	if output := splitByWord("IPV6Address"); len(output) != 3 {
		testing.Errorf("Expecting %v. Got %v", 3, len(output))
	}
	if output := splitByWord("Cats_AndDogs shouldNotBEUsed4eating"); len(output) != 9 {
		testing.Errorf("Expecting %v. Got %v", 9, len(output))
	}
}

func TestSnakeCase(testing *testing.T) {
	if output := Snakecase("One"); output != "one" {
		testing.Errorf("Expecting %v. Got %v", "one", output)
	}
	if output := Snakecase("One Two"); output != "one_two" {
		testing.Errorf("Expecting %v. Got %v", "one_two", output)
	}
	if output := Snakecase("One Two_Three"); output != "one_two_three" {
		testing.Errorf("Expecting %v. Got %v", "one_two_three", output)
	}
	if output := Snakecase("RedBlue"); output != "red_blue" {
		testing.Errorf("Expecting %v. Got %v", "red_blue", output)
	}
	if output := Snakecase("RedBlueGreen"); output != "red_blue_green" {
		testing.Errorf("Expecting %v. Got %v", "red_blue_green", output)
	}
	if output := Snakecase("REDBLUEGREEN"); output != "redbluegreen" {
		testing.Errorf("Expecting %v. Got %v", "redbluegreen", output)
	}
	if output := Snakecase("HTTPController"); output != "http_controller" {
		testing.Errorf("Expecting %v. Got %v", "http_controller", output)
	}
	if output := Snakecase("Forty40forty"); output != "forty_40_forty" {
		testing.Errorf("Expecting %v. Got %v", "forty_40_forty", output)
	}
	if output := Snakecase("IPV6Address"); output != "ipv_6_address" {
		testing.Errorf("Expecting %v. Got %v", "ipv_6_address", output)
	}
	if output := Snakecase("Cats_AndDogs shouldNotBEUsed4eating"); output != "cats_and_dogs_should_not_be_used_4_eating" {
		testing.Errorf("Expecting %v. Got %v", "cats_and_dogs_should_not_be_used_4_eating", output)
	}
}

func TestTitleCase(testing *testing.T) {
	if output := Titlecase("One"); output != "One" {
		testing.Errorf("Expecting %v. Got %v", "One", output)
	}
	if output := Titlecase("One Two"); output != "OneTwo" {
		testing.Errorf("Expecting %v. Got %v", "OneTwo", output)
	}
	if output := Titlecase("One Two_Three"); output != "OneTwoThree" {
		testing.Errorf("Expecting %v. Got %v", "OneTwoThree", output)
	}
	if output := Titlecase("RedBlue"); output != "RedBlue" {
		testing.Errorf("Expecting %v. Got %v", "RedBlue", output)
	}
	if output := Titlecase("RedBlueGreen"); output != "RedBlueGreen" {
		testing.Errorf("Expecting %v. Got %v", "RedBlueGreen", output)
	}
	if output := Titlecase("REDBLUEGREEN"); output != "REDBLUEGREEN" {
		testing.Errorf("Expecting %v. Got %v", "REDBLUEGREEN", output)
	}
	if output := Titlecase("HTTPController"); output != "HTTPController" {
		testing.Errorf("Expecting %v. Got %v", "HTTPController", output)
	}
	if output := Titlecase("Forty40forty"); output != "Forty40Forty" {
		testing.Errorf("Expecting %v. Got %v", "Forty40Forty", output)
	}
	if output := Titlecase("IPV6Address"); output != "IPV6Address" {
		testing.Errorf("Expecting %v. Got %v", "IPV6Address", output)
	}
	if output := Titlecase("Cats_AndDogs shouldNotBEUsed4eating"); output != "CatsAndDogsShouldNotBEUsed4Eating" {
		testing.Errorf("Expecting %v. Got %v", "CatsAndDogsShouldNotBEUsed4Eating", output)
	}
}

func TestCamelCase(testing *testing.T) {
	if output := Camelcase("One"); output != "one" {
		testing.Errorf("Expecting %v. Got %v", "one", output)
	}
	if output := Camelcase("One Two"); output != "oneTwo" {
		testing.Errorf("Expecting %v. Got %v", "oneTwo", output)
	}
	if output := Camelcase("One Two_Three"); output != "oneTwoThree" {
		testing.Errorf("Expecting %v. Got %v", "oneTwoThree", output)
	}
	if output := Camelcase("RedBlue"); output != "redBlue" {
		testing.Errorf("Expecting %v. Got %v", "redBlue", output)
	}
	if output := Camelcase("RedBlueGreen"); output != "redBlueGreen" {
		testing.Errorf("Expecting %v. Got %v", "redBlueGreen", output)
	}
	if output := Camelcase("REDBLUEGREEN"); output != "redbluegreen" {
		testing.Errorf("Expecting %v. Got %v", "redbluegreen", output)
	}
	if output := Camelcase("HTTPController"); output != "httpController" {
		testing.Errorf("Expecting %v. Got %v", "httpController", output)
	}
	if output := Camelcase("Forty40forty"); output != "forty40Forty" {
		testing.Errorf("Expecting %v. Got %v", "forty40Forty", output)
	}
	if output := Camelcase("IPV6Address"); output != "ipv6Address" {
		testing.Errorf("Expecting %v. Got %v", "ipv6Address", output)
	}
	if output := Camelcase("Cats_AndDogs shouldNotBEUsed4eating"); output != "catsAndDogsShouldNotBEUsed4Eating" {
		testing.Errorf("Expecting %v. Got %v", "catsAndDogsShouldNotBEUsed4Eating", output)
	}
}

func TestPluralize(testing *testing.T) {
	if output := Pluralize("OnceUponATime"); output != "OnceUponATimes" {
		testing.Errorf("Expecting %v. Got %v", "OnceUponATimes", output)
	}
	if output := Pluralize("There was A dog"); output != "There was A dogs" {
		testing.Errorf("Expecting %v. Got %v", "There was A dogs", output)
	}
	if output := Pluralize("with!!!Long!!!tooth"); output != "with!!!Long!!!tooths" {
		testing.Errorf("Expecting %v. Got %v", "with!!!Long!!!tooths", output)
	}
	if output := Pluralize("HeLivedInABox"); output != "HeLivedInABoxes" {
		testing.Errorf("Expecting %v. Got %v", "HeLivedInABoxes", output)
	}
	if output := Pluralize("WithAHugeMouse"); output != "WithAHugeMice" {
		testing.Errorf("Expecting %v. Got %v", "WithAHugeMice", output)
	}
}

func TestHasListType(testing *testing.T) {
	modelWithList := Model{}
	modelWithList.Properties = []ModelProperty{ModelProperty{PropertyType: "int", IsSetType: true}}
	if HasListType(modelWithList) != true {
		testing.Error("Model with list returned false")
	}

	modelWithoutList := Model{}
	modelWithoutList.Properties = []ModelProperty{ModelProperty{PropertyType: "int"}}
	if HasListType(modelWithoutList) != false {
		testing.Error("Model without list returned true")
	}
}

func TestAddJavaUtilitiesToTemplate(testing *testing.T) {
	templateObject := template.New("")
	addJavaUtilitiesToTemplate(templateObject)
}

func TestPackageToPath(testing *testing.T) {
	if output := PackageToPath("com.example"); output != "com/example" {
		testing.Errorf("Expecting %v. Got %v", "com/example", output)
	}
}

func TestSqliteType(testing *testing.T) {
	goodProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string"}
	badProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "potato"}

	if isSql := IsSqliteType(goodProp); isSql != true {
		testing.Errorf("Expecting %v. Got %v", true, isSql)
	}
	if isSql := IsSqliteType(badProp); isSql != false {
		testing.Errorf("Expecting %v. Got %v", false, isSql)
	}
	if sqlType := ToSqliteType(goodProp); sqlType != "TEXT" {
		testing.Errorf("Expecting %v. Got %v", "TEXT", sqlType)
	}
	if sqlType := ToSqliteType(badProp); sqlType != "potato" {
		testing.Errorf("Expecting %v. Got %v", "potato", sqlType)
	}
}

func TestJavaType(testing *testing.T) {
	goodProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string"}
	goodPropArray := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string", IsSetType: true}
	badProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "potato"}

	if isSql := IsJavaType(goodProp); isSql != true {
		testing.Errorf("Expecting %v. Got %v", true, isSql)
	}
	if isSql := IsJavaType(badProp); isSql != false {
		testing.Errorf("Expecting %v. Got %v", false, isSql)
	}
	if sqlType := ToJavaType(goodProp); sqlType != "String" {
		testing.Errorf("Expecting %v. Got %v", "String", sqlType)
	}
	if sqlType := ToJavaType(goodPropArray); sqlType != "List<String>" {
		testing.Errorf("Expecting %v. Got %v", "List<String>", sqlType)
	}
	if sqlType := ToJavaType(badProp); sqlType != "potato" {
		testing.Errorf("Expecting %v. Got %v", "potato", sqlType)
	}
}

func TestCoreDataType(testing *testing.T) {
	goodProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string"}
	goodPropArray := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string", IsSetType: true}
	badProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "potato"}

	if sqlType := ToCoreDataType(goodProp.PropertyType); sqlType != "String" {
		testing.Errorf("Expecting %v. Got %v", "String", sqlType)
	}
	if sqlType := ToCoreDataType(goodPropArray.PropertyType); sqlType != "String" {
		testing.Errorf("Expecting %v. Got %v", "String", sqlType)
	}
	if sqlType := ToCoreDataType(badProp.PropertyType); sqlType != "" {
		testing.Errorf("Expecting %v. Got %v", "", sqlType)
	}
}

func TestRailsType(testing *testing.T) {
	goodProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string"}
	goodPropArray := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string", IsSetType: true}
	badProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "potato"}

	if sqlType := ToRailsType(goodProp); sqlType != "string" {
		testing.Errorf("Expecting %v. Got %v", "string", sqlType)
	}
	if sqlType := ToRailsType(goodPropArray); sqlType != "Array" {
		testing.Errorf("Expecting %v. Got %v", "Array", sqlType)
	}
	if sqlType := ToRailsType(badProp); sqlType != "potato" {
		testing.Errorf("Expecting %v. Got %v", "potato", sqlType)
	}
}

func TestCustomType(test *testing.T) {
	goodProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "string"}
	badProp := ModelProperty{RemoteIdentifier: "Prop01", PropertyType: "potato"}

	RegisterCustomType("TestType1")
	SetCustomType("TestType1", "string", "String1")
	RegisterCustomType("TestType2")
	SetCustomType("TestType2", "string", "String2")

	if customType := ToCustomType("TestType1", goodProp); customType != "String1" {
		test.Errorf("Expecting %v. Got %v", "String1", customType)
	}
	if customType := ToCustomType("TestType2", goodProp); customType != "String2" {
		test.Errorf("Expecting %v. Got %v", "String2", customType)
	}
	if customType := ToCustomType("TestType1", badProp); customType != "potato" {
		test.Errorf("Expecting %v. Got %v", "potato", customType)
	}
}

func TestSHA256(testing *testing.T) {
	if output := SHA256("test"); output != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		testing.Errorf("Expecting %v. Got %v", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", output)
	}
}

func TestXcodeHash(testing *testing.T) {
	if output := SHA256("test"); len(output) != 40 {
		testing.Errorf("Expecting %v. Got %v", 40, len(output))
	}
}
