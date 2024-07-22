Sure, here's a comprehensive README.md for your Go CLI tool:

# Code Generator CLI Tool

This CLI tool is designed for code generation in Go, similar to the "rails g" command, but more general and flexible. It allows you to define templates and transformations in your project and generate code based on those definitions.

## Features

- Generates code based on templates.

- Supports custom transformations using JavaScript.

- Can be used to generate anything, however, automatically formats generated Go code using `gopls` if installed.

- Customizable configuration through YAML and JavaScript.

## Installation

On mac or linux, run:

```sh
brew install quinn/tap/g
```

if you're using ohmyzsh, add to the bottom of your .zshrc:

```sh
unalias g
```

## Usage

### Basic Command

```sh
g -path <target-directory> <generator-name> [args...]
```

- -path: Specifies the target directory containing the .g directory.
- <generator-name>: The name of the generator to run.
- [args...]: Arguments required by the generator.

### Example

```sh
g my-generator arg1 arg2
```

### Configuration

The configuration is defined in a g.yaml file located in the root directory specified by -path.

### g.yaml Structure

```yaml
version: "1.0"
generators:
  - name: "my-generator"
    args:
      - "arg1"
      - "arg2"
transforms:
  - myTransformFunction: "path/to/file"
```

- version: The version of the configuration file.
- generators: A list of generators.
- name: The name of the generator.
- args: A list of arguments required by the generator.
- transforms: A list of transformations to apply.
- myTransformFunction: The JavaScript function to apply.
- path/to/file: The path to the file to transform.

### Template Directory

Each generator should have a corresponding directory under .g/<generator-name>/tpl containing the template files.

### JavaScript Configuration

JavaScript files can be used to define transformations. The config.js file should define a config function that takes G_CONFIG_INPUT and returns additional configuration values.

```js
function config(input) {
  return {
    additionalKey: "additionalValue",
  };
}
```

### Example Project Structure

```
my-project/
├── g.yaml
├── .g/
│ └── my-generator/
│ ├── tpl/
│ │ └── template.go.tpl
│ └── config.js
└── main.go
```

### Template File

Template files use Go's text/template syntax and can access variables from the configuration.

template.go.tpl:

```go
package main

// imports can be omitted, and will be added automatically.

func main() {
    fmt.Println("Generated with arg1: {{ .arg1 }} and arg2: {{ .arg2 }}")
}
```

### JavaScript Transformations

Transformations allow you to manipulate files using JavaScript functions.

config.js:

```js
function config(input) {
  return {
    ...input,
    additionalKey: "additionalValue",
  };
}

function myTransformFunction(fileContent, config) {
  return fileContent.replace("PLACEHOLDER", config.additionalKey);
}
```

### Development

### Build

To build the CLI tool, run:

go build -o g

### Testing

To test the CLI tool, create test projects with the appropriate structure and run the tool with different configurations.

### Contributing

Feel free to open issues or submit pull requests. Any contributions are welcome!

### License

This project is licensed under the MIT License.
