# Proto-Share

Proto-Share is a simple tool for sharing prototypes between projects.
With Proto-Share you can create language agnostic schemas
(in protobuf format) and compile language specific artifacts from them.

## Supported Languages

- Java
- PHP
- TypeScript

## Configuration

In every project there has to be a top level config, e.g. `proto-share.config.yml`

### Main configuration

```yaml
projectName: # Name of the project [required]
inDir: # Directory where the proto files are located [required]
outDir: # Directory where the generated files should be placed [required]
languages: # Include which you want to compile to [at least one required]
  Java:
  PHP:
  TypeScript:
  MY_LANGUAGE:
    # protoc output 
    # used as: {config.outDir}/{language.subDir}/{language.moduleCompilePath}
    moduleCompilePath:

    # templates are copied
    # from `assets/templates/{language}/module`
    # to   `{config.outDir}/{language.subDir}/{language.moduleTemplatePath}`
    # (global templates are copied to `{config.outDir}/{language.subDir}`)
    moduleTemplatePath:

    # output subdirectory name for language
    subDir:

    # protoc command to generate code for language e.g. `java_out`, `php_out`...
    protocCommand:

    # generate publish script for language
    enablePublish:

    # additional parameters to be passed for templating
    # documented for every language
    additionalParameters:
```

### Language specific configuration and defaults

#### Java

```yaml
ModuleCompilePath: "{{ .Module.Name }}/src/main/java"
ModuleTemplatePath: "{{ .Module.Name }}"
SubDir: "java"
ProtocCommand: "java_out"
AdditionalParameters:
  # Java version
  version: 21

  # Java jar output path
  jarPath: "${rootDir}/build/libs"

  # Group id [required if enablePublish: true]
  groupId:

  # Artifact id [required if enablePublish: true]
  artifactId: "{{ .Module.Name | kebabCase }}"

  # Artifact repository url [required if enablePublish: true] 
  repositoryUrl:

  # Artifact repository username [required if enablePublish: true]
  repositoryUsername:

  # Artifact repository password [required if enablePublish: true]
  repositoryPassword: 
```

#### PHP

```yaml
ModuleCompilePath: ""
ModuleTemplatePath: ""
SubDir: "php"
ProtocCommand: "php_out"
```

#### TypeScript

```yaml
ModuleCompilePath: ""
ModuleTemplatePath: ""
SubDir: "ts"
ProtocCommand: "ts_out"
```

### Templates and pipes

As you can see in the configuration, templates are used to generate the code. Templates will be only evaluated if used
within a template file and for ModuleTemplatePath and ModuleCompilePath. For these items to be evaluated a
Context object is passed to the template. The Context object is a struct with the following fields:

```go
type Context struct {
*Config   // The top level of the configuration
*Language // The language specific configuration
*Module   // The currently evaluated module
}
```

You can access the fields of the Context object in the template by using the dot notation, e.g. `{{ .Module.Name }}`
will result in the `Name` field of the module.

You can also create recursive templates by referring to another field. The maximum depth of recursion is 10.

#### Pipes

There are also some pipe functions available for the templates. These are:

- `required` - Checks if a value is set and panics if not
- `kebabCase` - Converts a string to kebab case `(kebab-case)`
- `snakeCase` - Converts a string to snake case `(snake_case)`
- `camelCase` - Converts a string to camel case `(camelCase)`

### Modules

Modules are the building blocks of Proto-Share. One artifact will be generated per module.

Modules are usually defined in the directory of the proto files.

```yaml
name: # Name of the module
path: # Path to the proto files based on the config's inDir
hash: # A generated hash of the module - DO NOT CHANGE IT MANUALLY
version: # Version of the module - DO NOT CHANGE IT MANUALLY
```

## Usage

To generate the code for the supported languages, run the following command:

```shell
$ proto-share -config=${PATH_TO_ROOT_CONFIG}
```

When the program runs, it will generate the schema definitions for the supported languages.
If the module hash has changed since the last run, the version will be incremented and a new hash will be saved
upon successful compilation.

### Sample project

If you want to try out the tool you can check the `samples` directory. It contains a sample project with proto
files and a valid configuration. You can try to change the schema and the configs, to see what happens.

There is also a `Makefile` in the root directory, you can run `make build run` to run the program in a containerized
environment.