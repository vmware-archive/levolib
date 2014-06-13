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
	"bitbucket.org/pkg/inflect"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

var CustomType map[string]map[string]string

func Lower(input string) string {
	fmt.Print("")
	return strings.ToLower(input)
}

func Upper(input string) string {
	return strings.ToUpper(input)
}

func TestEquality(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

func TestInequality(args ...interface{}) bool {
	notEqual := TestEquality(args...)
	return !notEqual
}

func Snakecase(input string) string {
	words := splitByWord(input)
	joined := strings.Join(words, "_")
	return strings.ToLower(joined)
}

func Titlecase(input string) string {
	words := splitByWord(input)
	for index, word := range words {
		if word == "" {
			continue
		}
		words[index] = strings.ToUpper(word[0:1]) + word[1:]
	}
	return strings.Join(words, "")
}

func Camelcase(input string) string {
	if input == "" {
		return ""
	}
	words := splitByWord(input)
	for index, word := range words {
		if word == "" {
			continue
		}
		words[index] = strings.ToUpper(word[0:1]) + word[1:]
	}
	words[0] = strings.ToLower(words[0])
	return strings.Join(words, "")
}

func splitByWord(input string) []string {
	words := splitBySep([]string{input})
	words = splitByLowerToUpper(words)
	words = splitByDoubleUpperToLower(words)
	words = splitByNumber(words)
	output := make([]string, 0)
	for _, word := range words {
		if word != "" {
			output = append(output, word)
		}
	}
	return output
}

func splitBySep(inputWords []string) (outputWords []string) {
	re := regexp.MustCompile("[_ !]")
	outputWords = make([]string, 0)
	for _, inputWord := range inputWords {
		words := re.Split(inputWord, -1)
		for _, word := range words {
			outputWords = append(outputWords, word)
		}
	}
	return
}

func splitByLowerToUpper(inputWords []string) (outputWords []string) {
	re := regexp.MustCompile("[a-z][A-Z]")
	outputWords = make([]string, 0)
	for _, inputWord := range inputWords {
		indexes := re.FindAllIndex([]byte(inputWord), -1)
		var splitAtIndex int = -1
		var prevIndex int = 0
		for _, match := range indexes {
			splitAtIndex = match[1] - 1
			outputWords = append(outputWords, inputWord[prevIndex:splitAtIndex])
			prevIndex = splitAtIndex
		}
		if splitAtIndex > -1 {
			outputWords = append(outputWords, inputWord[splitAtIndex:])
		} else {
			outputWords = append(outputWords, inputWord)
		}
	}
	return
}

func splitByDoubleUpperToLower(inputWords []string) (outputWords []string) {
	re := regexp.MustCompile("[A-Z][A-Z][a-z]")
	outputWords = make([]string, 0)
	for _, inputWord := range inputWords {
		indexes := re.FindAllIndex([]byte(inputWord), -1)
		var splitAtIndex int = -1
		var prevIndex int = 0
		for _, match := range indexes {
			splitAtIndex = match[1] - 2
			outputWords = append(outputWords, inputWord[prevIndex:splitAtIndex])
			prevIndex = splitAtIndex
		}
		if splitAtIndex > -1 {
			outputWords = append(outputWords, inputWord[splitAtIndex:])
		} else {
			outputWords = append(outputWords, inputWord)
		}
	}
	return
}

func splitByNumber(inputWords []string) (outputWords []string) {
	re1 := regexp.MustCompile("([^0-9][0-9])")
	re2 := regexp.MustCompile("([0-9][^0-9])")
	intermediate := make([]string, 0)
	for _, inputWord := range inputWords {
		indexes := re1.FindAllIndex([]byte(inputWord), -1)
		var splitAtIndex int = -1
		var prevIndex int = 0
		for _, match := range indexes {
			splitAtIndex = match[1] - 1
			intermediate = append(intermediate, inputWord[prevIndex:splitAtIndex])
			prevIndex = splitAtIndex
		}
		if splitAtIndex > -1 {
			intermediate = append(intermediate, inputWord[splitAtIndex:])
		} else {
			intermediate = append(intermediate, inputWord)
		}
	}
	outputWords = make([]string, 0)
	for _, inputWord := range intermediate {
		indexes := re2.FindAllIndex([]byte(inputWord), -1)
		var splitAtIndex int = -1
		var prevIndex int = 0
		for _, match := range indexes {
			splitAtIndex = match[1] - 1
			outputWords = append(outputWords, inputWord[prevIndex:splitAtIndex])
			prevIndex = splitAtIndex
		}
		if splitAtIndex > -1 {
			outputWords = append(outputWords, inputWord[splitAtIndex:])
		} else {
			outputWords = append(outputWords, inputWord)
		}
	}
	return
}

func Pluralize(input string) string {
	pluralizedString := Snakecase(input)
	pluralizedString = inflect.Pluralize(pluralizedString)
	iIndex := 0
	pIndex := 0
	re := regexp.MustCompile("[A-z0-9]")
	output := ""
	for iIndex < len(input) && pIndex < len(pluralizedString) {
		inputChar := string(input[iIndex])
		pluralChar := string(pluralizedString[pIndex])
		if strings.ToLower(inputChar) == strings.ToLower(pluralChar) {
			output = output + inputChar
			iIndex++
			pIndex++
		} else if pluralChar == "_" {
			pIndex++
		} else if re.MatchString(inputChar) == false {
			output = output + inputChar
			iIndex++
		} else if re.MatchString(pluralChar) {
			output = output + pluralChar
			iIndex++
			pIndex++
		}
	}
	if pIndex <= len(pluralizedString)-1 {
		output = output + pluralizedString[pIndex:]
	}
	return output
}

func PackageToPath(input string) string {
	return strings.Replace(input, ".", "/", -1)
}

func HasListType(model Model) bool {
	for _, property := range model.Properties {
		if property.IsSetType {
			return true
		}
	}
	return false
}

func IsSqliteType(prop ModelProperty) bool {
	_, ok := SqliteTypes()[strings.ToLower(prop.PropertyType)]
	return ok && !prop.IsSetType
}

func ToSqliteType(prop ModelProperty) string {
	if value, ok := SqliteTypes()[strings.ToLower(prop.PropertyType)]; ok {
		return value
	} else {
		return prop.PropertyType
	}
}

func IsJavaType(prop ModelProperty) bool {
	_, ok := JavaTypes()[strings.ToLower(prop.PropertyType)]
	return ok
}

func ToJavaType(prop ModelProperty) string {
	theType := ""
	if javaType, ok := JavaTypes()[strings.ToLower(prop.PropertyType)]; ok {
		theType = javaType
	} else {
		theType = prop.PropertyType
	}

	if prop.IsSetType {
		return "List<" + theType + ">"
	}
	return theType
}

func ToCoreDataType(input string) string {
	return CoreDataTypes()[strings.ToLower(input)]
}

func ToObjectiveCType(input string) string {
	return ObjectiveCTypes()[strings.ToLower(input)]
}

func ToRailsType(prop ModelProperty) string {
	theType := ""
	if railsType, ok := RailsTypes()[strings.ToLower(prop.PropertyType)]; ok {
		theType = railsType
	} else {
		theType = prop.PropertyType
	}

	if prop.IsSetType {
		return "Array"
	}
	return theType
}

func RegisterCustomType(customType string) string {
	if CustomType == nil {
		CustomType = make(map[string]map[string]string, 0)
	}

	if _, ok := CustomType[customType]; !ok {
		CustomType[customType] = make(map[string]string, 0)
	}
	return ""
}

func SetCustomType(customType string, key string, value string) string {
	if CustomType[customType] == nil {
		//throw error?
	} else {
		CustomType[customType][key] = value
	}
	return ""
}

func IsCustomType(customType string, prop ModelProperty) bool {
	if CustomType[customType] == nil {
		return false
	} else if customMap, ok := CustomType[customType]; !ok {
		return false
	} else if _, ok := customMap[prop.PropertyType]; !ok {
		return false
	}
	return true
}

func ToCustomType(customType string, prop ModelProperty) string {
	if theType, ok := CustomType[customType][prop.PropertyType]; !ok {
		return prop.PropertyType
	} else {
		return theType
	}
}

func SHA256(data string) string {
	hashWriter := sha1.New()
	io.WriteString(hashWriter, data)
	return hex.EncodeToString(hashWriter.Sum(nil))
}

func Prefix(prefix string, original string) string {
	return Concat(original, prefix)
}

func Suffix(suffix string, original string) string {
	return Concat(original, suffix)
}

func Concat(first string, second string) string {
	return first + second
}

func Truncate(targetLength int, tooLong string) string {
	return string([]byte(tooLong)[0:targetLength])
}

func IdProp(properties []ModelProperty) string {
	for i := 0; i < len(properties); i++ {
		if strings.Contains(properties[i].RemoteIdentifier, "id") {
			underscored := Snakecase(properties[i].RemoteIdentifier)
			return Upper(underscored)
		}
	}
	if len(properties) > 0 {
		underscored := Snakecase(properties[0].RemoteIdentifier)
		return Upper(underscored)
	} else {
		return "UNKNOWN_PROPERTY"
	}
}

func SqliteTypes() map[string]string {
	var dict = make(map[string]string)
	dict["int"] = "INTEGER"
	dict["integer"] = "INTEGER"
	dict["short"] = "INTEGER"
	dict["long"] = "INTEGER"
	dict["float"] = "REAL"
	dict["boolean"] = "INTEGER"
	dict["char"] = "TEXT"
	dict["character"] = "TEXT"
	dict["string"] = "TEXT"
	dict["byte"] = "TEXT"
	dict["void"] = "NONE"
	return dict
}

func JavaTypes() map[string]string {
	var dict = make(map[string]string)
	dict["int"] = "int"
	dict["integer"] = "int"
	dict["short"] = "short"
	dict["long"] = "long"
	dict["float"] = "float"
	dict["boolean"] = "boolean"
	dict["char"] = "char"
	dict["character"] = "char"
	dict["string"] = "String"
	dict["byte"] = "byte"
	dict["void"] = "NONE"
	return dict
}

func CoreDataTypes() map[string]string {
	var dict = make(map[string]string)
	dict["int"] = "Integer 32"
	dict["integer"] = "Integer 32"
	dict["short"] = "Integer 16"
	dict["long"] = "Integer 64"
	dict["float"] = "Float"
	dict["boolean"] = "Boolean"
	dict["char"] = "String"
	dict["character"] = "String"
	dict["string"] = "String"
	dict["date"] = "Date"
	dict["byte"] = "NONE"
	dict["void"] = "NONE"
	return dict
}

func ObjectiveCTypes() map[string]string {
	var dict = make(map[string]string)
	dict["int"] = "NSNumber"
	dict["integer"] = "NSNumber"
	dict["short"] = "NSNumber"
	dict["long"] = "NSNumber"
	dict["float"] = "NSNumber"
	dict["boolean"] = "NSNumber"
	dict["char"] = "NSString"
	dict["character"] = "NSString"
	dict["string"] = "NSString"
	dict["date"] = "NSDate"
	dict["byte"] = "NONE"
	dict["void"] = "NONE"
	return dict
}

func RailsTypes() map[string]string {
	var dict = make(map[string]string)
	dict["int"] = "integer"
	dict["Integer"] = "integer"
	dict["short"] = "integer"
	dict["short"] = "integer"
	dict["long"] = "integer"
	dict["Long"] = "integer"
	dict["float"] = "float"
	dict["Float"] = "float"
	dict["boolean"] = "boolean"
	dict["Boolean"] = "boolean"
	dict["char"] = "string"
	dict["Character"] = "string"
	dict["String"] = "string"
	dict["string"] = "string"
	dict["byte"] = "NONE"
	dict["Byte"] = "NONE"
	dict["void"] = "NONE"
	dict["Void"] = "NONE"
	return dict
}

func addCommonUtilitiesToTemplate(templateObject *template.Template) *template.Template {
	templateObject = templateObject.Funcs(template.FuncMap{
		"eq":                 TestEquality,
		"neq":                TestInequality,
		"lower":              Lower,
		"upper":              Upper,
		"pluralize":          Pluralize,
		"camelcase":          Camelcase,
		"titlecase":          Titlecase,
		"snakecase":          Snakecase,
		"SHA256":             SHA256,
		"concat":             Concat,
		"truncate":           Truncate,
		"registerCustomType": RegisterCustomType,
		"setCustomType":      SetCustomType,
		"isCustomType":       IsCustomType,
		"toCustomType":       ToCustomType,
	})
	return templateObject
}

func addJavaUtilitiesToTemplate(templateObject *template.Template) *template.Template {
	templateObject = templateObject.Funcs(template.FuncMap{
		"hasListType":   HasListType,
		"isSqliteType":  IsSqliteType,
		"toSqliteType":  ToSqliteType,
		"sqliteType":    SqliteTypes,
		"isJavaType":    IsJavaType,
		"toJavaType":    ToJavaType,
		"javaType":      JavaTypes,
		"idProp":        IdProp,
		"packageToPath": PackageToPath,
	})
	return templateObject
}

func addObjectiveCUtilitiesToTempalte(templateObject *template.Template) *template.Template {
	templateObject = templateObject.Funcs(template.FuncMap{
		"toCoreDataType": ToCoreDataType,
		"toObjectiveCType" : ToObjectiveCType,
	})
	return templateObject
}

func addRailsUitilitiesToTemplate(templateObject *template.Template) *template.Template {
	templateObject = templateObject.Funcs(template.FuncMap{
		"toRailsType": ToRailsType,
	})
	return templateObject
}
