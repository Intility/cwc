# Chat With Code

<div align="center">
  <a href="https://github.com/emilkje/go-openai-toolkit">
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

### Using Go

If you have Go installed (version 1.22 or higher), you can install Chat With Code using the following command:

```sh
go install github.com/emilkje/cwc@latest
```

### Pre-built Binaries

We also provide pre-built binaries for Windows, macOS, and Linux. You can download them from the [releases page](https://github.com/emilkje/cwc/releases) on GitHub.

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
           --api-version "2023-12-01-preview" \
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

```sh
# chat across all .go files
cwc -i ".*.go"
```

```sh
# chat with everything inside src/ except .tsx files
cwc -x ".*.tsx" -p src
```

```sh
# chat with a git diff
git diff refA...refB > foo.diff
cwc -i "foo.diff"
```

## Roadmap 

> Note: these items may or may not be implemented in the future.

- [ ] tests
- [ ] support both azure and openai credentials
- [x] `cwc login` command to set up credentials
- [ ] Pull request gating
- [ ] Automatic version bumping
- [ ] Automatic release notes generation
- [ ] chat using web ui with `cwc web`
- [ ] indexing/search implementation for large codebases
- [ ] tools for dynamic context awareness
- [ ] tools for gathering external documentation

## Contributing

Please file an issue if you encounter any problems or have suggestions for improvement. We welcome contributions in the form of pull requests, bug reports, and feature requests.

## License

Chat With Code is provided under the MIT License. For more details, see the [LICENSE](LICENSE) file.

If you encounter any issues or have suggestions for improvement, please open an issue in the project's [issue tracker](https://github.com/emilkje/chat-with-code/issues).

[banner-photo-url]: ./docs/assets/yelling_at_code.webp
[screenshot-url]: ./docs/assets/screenshot.png
