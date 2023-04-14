# lunapipe
Use ChatGPT in your terminal. Kind of feels like a bash utility?

## Installation
```bash
curl https://raw.githubusercontent.com/lunabrain-ai/lunapipe/main/scripts/install.sh | sh
```
or
```shell
go install github.com/lunabrain-ai/lunapipe@latest
```
or if you are looking for other releases, go to [releases](https://github.com/lunabrain-ai/lunapipe/releases/).

## Usage
````bash
export OPENAI_API_KEY="<your openai api key>"
lunapipe "Write me a go function that prints 'Hello World'"
Here's an example Go function that prints "Hello World" to the console:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```
````

To have your API key persist, you can use:
```shell
lunapipe configure
Enter your API key: <your openai api key>
Wrote API key to ~/.lunapipe/config.yaml
```

Don't have an API key? Sign up [here](https://platform.openai.com/overview) and generate an API key [here](https://platform.openai.com/account/api-keys).

### Pipe
You can pipe text into lunapipe. For example, if you have a file called `main.go` that contains the following code:
```go
ls | lunapipe "Based on the files, what language is this repo?"
This repo is written in Go (also known as Golang).
```

#### Templates
You can use templates to generate code. For example, if you want to generate a go function that prints "Hello World", you can use the following template:
```bash
lunapipe -t function -p language=go "Read values from a map"

# create an alias to code even faster
alias aigo="lunapipe -t function -p language=go"
aigo "Read values from a map"
```
To see all available templates, check out the available [prompt templates](prompts).

To define your own templates, you can pass a directory to where your templates are:
```bash
mkdir my_prompt_templates
echo "This is a test" > my_prompt_templates/test.tmpl
lunapipe --prompts my_prompt_templates -t test "Is this thing working?"

# If you want to add parameters to your template
echo 'This is a test with params: {{ index .Params "testparam" }} ' > my_prompt_templates/test_with_params.tmpl
lunapipe --prompts my_prompt_templates -t test -p testparam="Hello, world!" "Is this thing working?"
```

## Hack

### Debug
```bash
LOG_LEVEL=debug go run main.go "Test prompt"
```

### Run locally
```bash
go run main.go
```

### Build
```bash
go build main.go
```
