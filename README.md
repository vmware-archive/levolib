# Code Generation With Levolib

At the most basic level, this library accepts a *context* and produces *source code files*. The information within a context acts like instructions; it provides levolib with the models and templates it will need, instructions regarding which model(s) to use when filling each template, and a few other things as well.

Levolib provides **BeginContext** which returns an initalized context. A context object provides methods for adding new templates and models and mappings to itself. Once the context is complete, it can be passed to levolib's **ProcessMappings** method, which will start generating source code files.

Building a context piece by piece isn't necessary. Levolib also provides **GetJSONConfigurationAdapter**, which returns an adapter specialized in converting a json configuration file into a context. Simply pass the file path to the adapter's **ProcessConfigurationFile** method.

The [levo project](https://github.com/cfmobile/levo) is a useful example of how the above methods can be used to instrument this library. It also provides example files, include an example json configuration.

## License
This library is licensed under the Apache License, Version 2.0 [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

## Code-Generator Context Structure
```go
type Context struct {
	ProjectName      string
	PackageName      string
	TemplaterVersion string
	Schema           Schema
	Templates        []TemplateInfo
	Mappings         []TemplatesForModels
	Language         string
	Zip              bool
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
	RemoteIdentifier   string
	LocalIdentifier string
	PropertyType   string
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

type OutputAdapter interface {
	GenerateFiles(templateInfo TemplateInfo, templateData TemplateData) ([]GeneratedFile, error)
}

type GeneratedFile struct {
	FileName string
	Body     string
}

func (context *Context) AddModelWithName(name string) (*Model, error)
//Initalize a new Model with a specific name. Add that Model to the Context

func (context *Context) AddModel(model Model) (*Model, error)
//Add a Model object to the Context

func (context *Context) AddTemplateDirectory(templateDirPath string) error
//Parse a directory and all it's subdirectories. Update the Context with a new Template
//for every file found.

func (context *Context) AddTemplateFilePath(filePath string) (TemplateInfo, error)
//Add the file at filePath to the Context as a Template

func (context *Context) AddTemplate(fileName string, body []byte, version string, directory string, adapter OutputAdapter) (*TemplateInfo, error)
//Initalize a template with the information provided and add that TemplateInfo to the Context
//fileName          : Exactly what it seems like
//body              : The contents of the template file
//version           : The template version. Used to ensure that levolib is able to process the template.
//directory         : The path to the template file (can be a relative path). This directory will be
//                    used to determine where the output of the template goes.
//adapter           : This library comes with a GoTemplateAdapter for processing go templates.
//                    Other adapters can be created. For example, adapters for Mako templates or
//                    xslt templates. 

func (context *Context) AddTemplatesForModelsMapping(templateFileNames []string, modelNames []string) error
//Specify which models should be used to fill a template (or a set of templates).

func (context *Context) ModelForName(name string) (*Model, error)
//Simple getter method

func (context *Context) TemplateForFileName(fileName string) (*TemplateInfo, error)
//Simple getter method
```

## Code-Generatior Structure
```go
func BeginContext(packageName string, language string, version string) (Context, error)
//Get an initialized context object.

func ProcessMappings(context Context) ([]GeneratedFile, error)
//Create source code files using the contents of this context

func GetJSONConfigurationAdapter() JSONConfigAdapter
//An adapter for converting JSON configuration files into Contexts. The adapter provides
//**ProcessConfigurationFile** and **ProcessConfigurationString**, each of which returns a Context
//based on the JSON configuration information they are fed

func GetJSONSchemaAdapter() JSONSchemaAdapter
//An adapter for converting JSON schema files into Contexts. The adapter provides
//**ProcessSchemaFile** and **ProcessSchemaString**, each of which returns a Context
//based on the JSON schema information they are fed
```
