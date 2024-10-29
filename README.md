<p align="center">
  <img src="https://github.com/user-attachments/assets/c023465c-132c-4fef-b4b4-4f30552148fb" />
</p>

### This README is not being updated with the latest releases. I am working on new documentation

Kuma is a powerful framework designed to generate scaffolds for any programming language, based on [Go templates](https://pkg.go.dev/text/template). It streamlines the process of setting up new projects by automating the creation of directories, files, and base code, ensuring consistency and saving valuable development time. Additionally, Kuma features a customizable TUI, providing an intuitive and efficient experience both for those creating scaffolds and those using them, making the process accessible and seamless for developers of all levels.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Create Your Own Scaffolds](#create-your-own-scaffolds)
  - [Builders](#builders)
  - [Templates](#templates)
  - [Runs](#runs)
- [Terminal Commands](#terminal-commands)
  - [Create a Scaffold](#create-a-scaffold)
  - [Execute a Run](#execute-a-run)
  - [Get Templates from GitHub](#get-templates-from-github)
    - [Official Templates](#official-templates)
- [Contribution](#contribution)
- [License](#license)

## Features

- Customize your project’s directory and file structures through [Go templates](https://pkg.go.dev/text/template).
- GitHub integration to download pre-defined templates from the community or for personal use via private repositories.
- Ability to create custom CLI/TUI command workflows through a YAML file using runs.
- Dynamic variable usage to be applied to templates. Variables can be extracted from a local YAML or JSON file or fetched from a public URL. They can also be obtained from user input during the execution of a [run](cmd/commands/exec).

## Installation

### Requirements

- [Go](https://golang.org/dl/) version 1.23 or higher.
- Git installed and configured on your system.

### Step by Step

1. **Run the installation command:**

   ```bash
   go install github.com/kuma-framework/kuma/v2@latest
   ```

2. **Add the Go bin directory to your PATH (if not already included):**

   Add the following line to your shell configuration file (`.bashrc`, `.zshrc`, etc.):

   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

   Then, reload your shell or run:

   ```bash
   source ~/.bashrc
   ```

   _Replace `.bashrc` with your shell configuration file if necessary._

3. **Verify if $GOPATH is set correctly:**

   Run the following command to display the current value of $GOPATH:

   ```bash
   echo $GOPATH
   ```

   **Expected Results:**

   - **If $GOPATH is set correctly:** The command will return the path to the GOPATH directory, usually something like `/home/user/go` on Linux or `C:Users/user/go` on Windows.
   - **If $GOPATH is empty or incorrect:** You will need to set it by adding the following line to your shell configuration file:

     ```bash
     export GOPATH=$(go env GOPATH)
     ```

     Then, reload your shell or run:

     ```bash
     source ~/.bashrc
     ```

     _Replace `.bashrc` with your shell configuration file if necessary._

4. **Verify the installation:**

   ```bash
   kuma --help
   ```

   You should see the Kuma CLI help, confirming that the installation was successful.

## Create Your Own Scaffolds

For Kuma to work, all framework-related files must be inside the `.kuma` folder.

The framework uses a parser for [Go templates](https://pkg.go.dev/text/template) divided into three file types:

### Builders

YAML-formatted templates that contain the folder and file structure to be created.

```yaml
# .kuma/base.yaml

# Global variables to be used in all templates.
global:
  packageName: "{{.data.packageName}}"

# Directory and file structure to be generated
structure:
  # main.go
  main.go:
    # The template that will be used as the basis for creating the main.go file
    template: templates/Main.go

    # Variables that will be used inside the template
    data:
      # msg: Hello, Kuma!
      msg: "{{ .data.msg }}"
```

### Templates

Individual [Go templates](https://pkg.go.dev/text/template) for the files that will be created.

```go
package main

import (
  "fmt"
)

func main() {
  // fmt.Print("Hello, Kuma!")
  fmt.Println("{{ .data.msg }}")
}
```

### Runs

A YAML file containing a sequence of actions that will be executed when calling a `run`. It includes logs, terminal commands, HTTP calls, text input, and multiple-choice prompts, along with actions to create folders and files based on builders and templates.

all runs must be located within the `.kuma/runs` directory

**Check the full documentation [here](cmd/commands/exec).**

```yaml
# Name of the run that will be executed as soon as the repository is obtained via the `kuma get` command
initial:
  # Description of the run
  description: "Initial run to set up the project"

  # Steps that will be executed when the run is called
  steps:
    # Input action
    - input:
        label: "What is the package name of your project?"
        out: packageName # Example: github.com/kuma-framework/kuma/v2-hello-world

    # Another input action
    - input:
        label: "What message would you like to print?"
        out: msg # Example: Hello, Kuma!

    # Log message
    - log: "Creating structure for {{.data.packageName}}" # Example: Creating structure for github.com/kuma-framework/kuma/v2-hello-world

    # Create the project structure using the base.yaml builder
    - create:
        from: base.yaml

    # Success log
    - log: "Base structure created successfully!"

    # Initialize the Go module
    - cmd: go mod init {{.data.packageName}} # Example: go mod init github.com/kuma-framework/kuma/v2-hello-world

    # Install dependencies
    - cmd: go mod tidy

    # Run the main.go file
    - cmd: go run main.go # Example: Hello, Kuma!
```

### Additional Notes

- This project uses [sprout](https://github.com/go-sprout/sprout) as a dependency, which contains hundreds of extremely useful functions for working with Go templates. In addition to our official functions, further enhancing your experience with our framework.

  Sproute: [Official Docs](https://docs.atom.codes/sprout)

  Kuma functions: [Read me](internal/functions)
  &nbsp;

## Terminal Commands

![render1728444387729](https://github.com/user-attachments/assets/54f74beb-cd85-47b0-87c2-4e7bd471cb54)

### Create a Scaffold

The `create` command is used to create a scaffold based on the builders and templates inside the `.kuma` folder and a JSON or YAML file containing the variables to replace in the templates.

```bash
kuma create --variables=swagger.json --project=. --from=base.yaml
```

**Flags:**

- `--variables`, `-v`: Path or URL to the variables file.
- `--project`, `-p`: Path to the project where the scaffold will be created.
- `--from`, `-f`: Path to the YAML file with the structure and templates.

### Execute a Run

The `exec` command is used to start the process of a [run](cmd/commands/exec).

```bash
kuma exec --run=initial
```

**Flags:**

- `--run`, `-r`: Name of the run to be executed.

### Get Templates from GitHub

Fetch templates and runs from a GitHub repository.

You can get templates from any repository using the command:

```bash
kuma get --repo=kuma-framework/kuma-typescript-rest-services
```

Or use one of our official templates with:

```bash
kuma get --template=kuma-typescript-rest-services
```

**Flags:**

- `--repo`, `-r`: Name of the GitHub repository.
- `--template`, `-t`: Name of the official template.

#### Official Templates

- **[Hello World](https://github.com/kuma-framework/kuma/v2-hello-world):** A simple Hello World in Go.
- **[OpenAPI 2.0 TypeScript Services](https://github.com/kuma-framework/kuma/v2-typescript-rest-services):** Create a TypeScript library with typed services for all endpoints described in an Open API 2.0 file.
- **[Changelog Generator](https://github.com/kuma-framework/kuma/v2-changelog-generator):** Helper to write a good changelog to your project.

## Contribution

Contributions are welcome! Feel free to open issues or submit pull requests to improve Kuma.

1. **Fork the repository.**
2. **Create a branch for your feature:**

   ```bash
   git checkout -b feature/my-new-feature
   ```

3. **Commit your changes:**

   ```bash
   git commit -m "Add new feature"
   ```

4. **Push to the branch:**

   ```bash
   git push origin feature/my-new-feature
   ```

5. **Open a Pull Request.**

## License

This project is licensed under the MIT license. See the [LICENSE](LICENSE) file for more details.
