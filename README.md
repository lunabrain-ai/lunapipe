# aicli
Use ChatGPT in your terminal. Kind of feels like a bash utility?

## Installation
```bash
go get github.com/lunabrain-ai/aicli
```
or if you want the binary release, go to [releases](https://github.com/lunabrain-ai/aicli/releases/).

## Usage
````bash
export OPENAI_API_KEY="<your openai api key>"
aicli "Write me a go function that prints 'Hello World'"
Here's an example Go function that prints "Hello World" to the console:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```
````

### Pipe
You can pipe text into aicli. For example, if you have a file called `main.go` that contains the following code:
```go
ls | aicli "Based on the files, what language is this repo?"
This repo is written in Go (also known as Golang).
```

#### Templates
You can use templates to generate code. For example, if you want to generate a go function that prints "Hello World", you can use the following template:
```bash
aicli -t code -p language=go "Read values from a map"

# create an alias to code even faster
alias aigo="aicli -t function -p language=go"
aigo "Read values from a map"
```
To see all available templates, check out the available [prompt templates](prompts).

To define your own templates, you can pass a directory to where your templates are:
```bash
mkdir my_prompt_templates
echo "This is a test" > my_prompt_templates/test.tmpl
aicli --prompts my_prompt_templates -t test "Is this thing working?"

# If you want to add parameters to your template
echo 'This is a test with params: {{ index .Params "testparam" }} ' > my_prompt_templates/test_with_params.tmpl
aicli --prompts my_prompt_templates -t test -p testparam="Hello, world!" "Is this thing working?"
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
