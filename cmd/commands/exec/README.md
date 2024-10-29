# Kuma runs

## Table of Contents

- [Introduction](#introduction)
- [Structure of a Run](#structure-of-a-run)
  - [Action Types](#action-types)
    - [Input](#input)
    - [Log](#log)
    - [Create](#create)
    - [Cmd](#cmd)
    - [Load](#load)
    - [Nested Run](#nested-run)
- [How to Execute a Run](#how-to-execute-a-run)
  - [Using the CLI Command](#using-the-cli-command)
  - [Interactive Run Selection](#interactive-run-selection)
- [Advanced Examples](#advanced-examples)
  - [Run that extracts variables from a swagger file](#run-that-extracts-variables-from-a-swagger-file)
- [License](#license)

## Introduction

**Runs** in Kuma are sequences of actions designed to automate repetitive tasks during project development. They enable the execution of custom pipelines, which can include everything from user inputs to terminal commands, ensuring consistency and efficiency in the workflow.

All runs must be located within the `.kuma/runs` directory

## Structure of a Run

A Run is composed of a sequence of steps that define the actions to be executed. Below are the main components and types of actions that can be included in a Run.

### Action Types

#### Input

Prompts the user for input during the execution of the Run.

```yaml
- input:
    label: "What is the name of your project's package?"
    out: packageName # Example: github.com/kuma-framework/kuma/v2-hello-world
```

**Fields:**

- `label`: The message displayed to the user.
- `out`: The variable where the entered value will be stored.

**Additional Options:**

- `options`: A list of options for selection.
- `multi`: Flag to allow selecting more than one option. Returns an array in `out`.
- `other`: If no option is selected, displays a shortcut with the **o** key to open a text input.
  **Example with Options and Multiple Selection:**

```yaml
select-runtime:
  description: "Select the runtime to use"
  steps: - input:
    label: "Select a runtime"
    multi: false
    other: false
    options:
      - label: Node
        value: node
      - label: Deno 2.0
        value: deno
      - label: Bun
        value: bun
    out: runtime
```

#### Log

Logs a message to the console.

```yaml
- log: "Creating structure for {{.data.packageName}}" # Example: Creating structure for github.com/kuma-framework/kuma/v2-hello-world
```

**Fields:**

- `log`: The message to be logged. Can include dynamic variables.

#### Create

Creates the project structure based on a defined builder.

```yaml
- create:
    from: base.yaml
```

**Fields:**

- `from`: The YAML file defining the structure and templates to be used.

#### Cmd

Executes a command in the terminal.

```yaml
- cmd: npm install
```

**Fields:**

- `cmd`: The command to be executed. Can include dynamic variables.

#### Load

Loads variables from a local file or URL.

```yaml
- load:
    from: variables.yaml
    out: vars
```

**Fields:**

- `from`: Path or URL to the JSON or YAML file containing the structure that will be stored in the `out` variable. Can include dynamic variables.
- `out`: The variable where the loaded data will be stored.

#### Nested Run

Executes one run within another. Variables from a run are automatically passed to nested runs.

```yaml
main:
  description: "main run"
  steps:
    - log: "Executing main run"
    - run: nested

nested:
  description: "nested run"
  steps:
    - log: "Executing nested run"
```

**Fields:**

- `run`: Name of the Run to be executed.

## How to Execute a Run

### Using the CLI Command

To execute a specific Run, use the `exec` command followed by the name of the Run.

```bash
kuma exec --run=initial
```

**Flags:**

- `--run`, `-r`: Name of the Run to be executed.

### Interactive Run Selection

If the name of the Run is not specified, Kuma CLI will present an interactive interface to select which Run you want to execute.

```bash
kuma exec
```

**Steps:**

1. **Run Selection:** A list of available Runs will be displayed for selection.
2. **Execution:** The selected Run will be executed based on the defined steps.

## Advanced Examples

### Run that extracts variables from a swagger file

Check out the builders and templates used by cloning this [repository](https://github.com/kuma-framework/kuma/v2-typescript-rest-services) or via the command `kuma get -t kuma-typescript-rest-services`.

```yaml
initial:
  description: "Initial run"
  steps:
    - input:
        label: "What is your project name?"
        out: projectName
    - input:
        label: "What is your project repository?"
        out: projectRepository
    - create:
        from: base.yaml
    - log: Base structure created successfully!
    - input:
        label: "Enter the local file or the file URL in the Open API 2.0 format with the data you want to generate the library:"
        out: swagger
    - load:
        from: "{{.data.swagger}}"
        out: apiData
    - create:
        from: from-swagger.yaml
    - cmd: npm i
    - cmd: npm run format
```

## License

This project is licensed under the MIT license. See the [LICENSE](LICENSE) file for more details.
