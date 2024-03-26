# **C**hat **W**ith **C**ode

<div align="center">
  <a href="https://github.com/intility/cwc">
    <img src="docs/assets/yelling_at_code.webp" alt="Logo">
  </a>
</div>

## Overview

Chat With Code is yet another command-line application that bridges the gap between coding and conversation. This tool lets you engage with your codebases using natural language prompts, offering a fresh approach to code exploration and problem-solving.

**Why does this tool exist?**

I was frequently disappointed with Github Copilot's chat context discovery, as it often missed relevant files for accurate answers. 
CWC improves this by allowing you to specify include and exclude patterns across your codebase, giving you complete control over 
the context window during chats. Additionally, its terminal-based operation makes it independent of IDEs like VSCode, 
enhancing flexibility with other text editors not natively supporting Github Copilot.

## Features

- **Interactive Chat Sessions**: Start a dialogue with your codebase to learn about its structure, get summaries of different parts, or even debug issues.
- **Intelligent Context-Aware Responses**: Powered by OpenAI, Chat With Code understands the context of your project, providing meaningful insights and relevant code snippets.
- **Customizable File Inclusion**: Filter the files you want the tool to consider using regular expressions, ensuring focused and relevant chat interactions.
- **Gitignore Awareness**: Exclude files listed in `.gitignore` from the chat context to maintain confidentiality and relevance.
- **Simplicity**: A simple and intuitive interface that requires minimal setup to get started.

## Installation

### Using homebrew

Intility provides a shared Homebrew tap with all our formulae. Install Chat With Code using:

```sh
brew tap intility/tap
brew install cwc
```

### Using Go

If you have Go installed (version 1.22 or higher), you can install Chat With Code using the following command:

```sh
go install github.com/intility/cwc@latest
```

### Pre-built Binaries

We also provide pre-built binaries for Windows, macOS, and Linux. You can download them from the [releases page](https://github.com/intility/cwc/releases) on GitHub.
You can install the latest release with the following command using a bash shell (git bash or WSL on Windows):

```sh
bash <(curl -sSL https://raw.githubusercontent.com/intility/cwc/main/scripts/install.sh)

# move the binary to a directory in your PATH
mv cwc /usr/local/bin
```

## Getting Started

After installing Chat With Code, you're just a few steps away from a conversational coding experience. Get everything up and running with these instructions:

1. **Launch Chat With Code**: Open your terminal and enter `cwc` to start an interactive session with your codebase.
   If you are not already signed you will be prompted to configure your credentials.

2. **Authenticate**: To enable communication with your code, authenticate using your Azure credentials by executing:

    ```sh
    cwc login
    ```

   *For a seamless login experience, follow the non-interactive authentication method below:*

    1. Safeguard your API Key by storing it in a variable (avoid direct command-line input to protect the key from your history logs):

         ```sh
         read -s API_KEY
         ```

    2. Authenticate securely using the following command:

         ```sh
         cwc login \
           --api-key=$API_KEY \
           --endpoint "https://your-endpoint.openai.azure.com/" \
           --deployment-model "gpt-4-turbo"
         ```

   > **Security Notice**: Never input your API key directly into the command-line arguments to prevent potential exposure in shell history and process listings. The API key is securely stored in your personal keyring.

After completing these steps, you will have established a secure session, ready to explore and interact with your codebase in the most natural way.

![screenshot][screenshot-url]

Need a more tailored experience? Try customizing your session. Use the `--include`, `--exclude` flag to filter for specific file patterns or `--paths` to set directories included in the session. Discover all the available options with:

```sh
cwc --help
```

## Example usage

The simplest example would be to chat with a single file or output from a command. This use-case is easy using a pipe:

```sh
cat README.md | cwc "help me rewrite the getting started section"
```

If you have multiple files you want to include in the context you can provide a regular expression matching your criteria for inclusion using the `-i` flag:

```sh
# chat across all .go files
cwc -i ".*.go"

# chat with README and test files
cwc -i "README.md|.*_test.go"
```

The include flag can also be combined with exclusion expressions, these work exactly the same as the inclusion patterns, but takes priority:

```sh
# chat with all .ts files, excluding a large .ts file
cwc -i ".*.ts$" -x "large_file.ts"
```

In addition to include and exclude expressions you can also scope the search space to a particular directory. Multiple paths can be provided by a comma separated list or by providing multiple instances of the `-p` flag.

```sh
# chat with everything inside src/ except .tsx files
cwc -x ".*.tsx" -p src

# chat with all yaml files in prod and lab
cwc -i ".*.ya?ml" -p prod,lab
```

The result output from cwc can also be piped to other commands as well. This example automates the creation of a conventional commit based on the current git diff.

```sh
# generate a commit message for current changes
PROMPT="please write me a conventional commit for these changes"
git diff HEAD | cwc $PROMPT | git commit -e --file -
```

## Configuration

Managing your configuration is simple with the `cwc config` command. This command allows you to view and set configuration options for cwc.
To view the current configuration, use:

```sh
cwc config get
```

To set a configuration option, use:

```sh
cwc config set key1=value1 key2=value2 ...
```

For example, to disable the gitignore feature and the git directory exclusion, use:

```sh
cwc config set useGitignore=false excludeGitDir=false
```

To reset the configuration to default values use `cwc login` to re-authenticate.

## Templates

### Overview

Chat With Code (CWC) introduces the flexibility of custom templates to enhance the conversational coding experience. Templates are pre-defined system messages and prompts that tailor interactions with your codebase. A template envelops default prompts, system messages and variables, allowing for easier access to common tasks.

### Template Schema

Each template follows a specific YAML schema defined in `templates.yaml`. 
Here's an outline of the schema for a CWC template:

```yaml
templates:
  - name: template_name
    description: A brief description of the template's purpose
    defaultPrompt: An optional default prompt to use if none is provided
    systemMessage: |
      The system message that details the instructions and context for the chat session.
      This message supports placeholders for {{ .Context }} which is the gathered file context,
      as well as custom variables `{{ .Variables.variableName }}` fed into the session with cli args.
    variables:
      - name: variableName
        description: Description of the variable
        defaultValue: Default value for the variable
```

### Placement

Templates may be placed within the repository or under the user's configuration directory, adhering to the XDG Base Directory Specification:

1. **In the Repository Directory**: To include the templates specifically for a repository, place a `templates.yaml` inside the `.cwc` directory at the root of your repository:

   ```
   .
   ├── .cwc
   │   └── templates.yaml
   ...

2. **In the User XDG CWC Config Directory**: For global user templates, place the `templates.yaml` within the XDG configuration directory for CWC, which is typically `~/.config/cwc/` on Unix-like systems:

   ```
   $XDG_CONFIG_HOME/cwc/templates.yaml
   ```

   If `$XDG_CONFIG_HOME` is not set, it defaults to `~/.config`.

### Example Usage

You can specify a template using the `-t` flag and pass variables with the `-v` flag in the terminal. These flags allow you to customize the chat session based on the selected template and provided variables.

#### Selecting a Template

To begin a chat session using a specific template, use the `-t` flag followed by the template name:

```sh
cwc -t my_template
```

This command will start a conversation with the system message and default prompt defined in the template named `my_template`.

#### Passing Variables to a Template

You can pass variables to a template using the `-v` flag followed by a key-value pair:

```sh
cwc -t my_template -v personality="a helpful assistant",name="Juno"
```

Here, the `my_template` template is used. The `personality` variable is set to "a helpful coding assistant", and
the `name` variable is set to "Juno". These variables will be fed into the template's system message where placeholders are specified.

The template supporting these variables might look like this:

```yaml
name: my_template
description: A custom template with modifiable personality and name
systemMessage: |
  You are {{ .Variables.personality }} named {{ .Variables.name }}. 
  Using the following context you will be able to help the user.

  Context:
  {{ .Context }}
   
  Please keep in mind your personality when responding to the user.
  If the user asks for your name, you should respond with {{ .Variables.name }}.
variables:
  - name: personality
    description: The personality of the assistant. e.g. "a helpful assistant"
    defaultValue: a helpful assistant
  - name: name
    description: The name of the assistant. e.g. "Juno"
    defaultValue: Juno
```

> Notice that the `personality` and `name` variables have default values, which will be used if no value is provided in the `-v` flag.

## Roadmap 

These items may or may not be implemented in the future.

- [ ] tests
- [ ] support both azure and openai credentials
- [ ] customizable tools

## Contributing

Please file an issue if you encounter any problems or have suggestions for improvement. We welcome contributions in the form of pull requests, bug reports, and feature requests.

## License

Chat With Code is provided under the MIT License. For more details, see the [LICENSE](LICENSE) file.

If you encounter any issues or have suggestions for improvement, please open an issue in the project's [issue tracker](https://github.com/intility/chat-with-code/issues).

[banner-photo-url]: ./docs/assets/yelling_at_code.webp
[screenshot-url]: ./docs/assets/screenshot.png
