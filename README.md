# lunapipe
Use ChatGPT in your terminal. It can:

* Prompt ChatGPT and get a streamed response
* Pipe stdin/stdout to and from ChatGPT. Ex:`cat README.md | lunapipe "What does LunaPipe do?" > output.txt`
* Chat back and forth with ChatGPT in your terminal instead of having to open OpenAI's web UI. Ex: `lunapipe chat`
* Use templates for more advanced flows (or write your own template)


## Usage Examples

#### Direct prompts

```bash
lunapipe "Write me a go function that prints 'Hello World'" 
```

#### Pipe in (and out) anything

```bash
cat somefile.go | lunapipe "Write documentation for this code" > README.md
grep -r "hello" . | lunapipe "which of these files looks like a hello world?"
ls | lunapipe "Based on the files, what language is this repo?"
```

#### Chat

```bash
lunapipe chat                                                                                                                                 
> Starting chat, close with ctrl+D...
> hi
> Hello! How may I assist you today?
```

#### Use a template
For example, here we are using the built-in `function` template the generates a function in the language you choose. 
```bash
lunapipe -t function -p language=go "print hello world"
> func main() {
>     fmt.Println("Hello World")
> }
```

## Video Demonstration of LunaPipe
[![screenshot-www youtube com-2023 04 14-20_52_03](http://img.youtube.com/vi/2Y4i3rtFvAI/0.jpg)](https://www.youtube.com/watch?v=2Y4i3rtFvAI)

## Setup

### Installation
Use the automatic install script for any system:
```bash
curl https://raw.githubusercontent.com/lunabrain-ai/lunapipe/main/scripts/install.sh | sh
```
or install as a go package
```shell
go install github.com/lunabrain-ai/lunapipe@latest
```
or if you are looking for other releases, go to [releases](https://github.com/lunabrain-ai/lunapipe/releases/).
### Input your OpenAI API key

Don't have an API key? Sign up [here](https://platform.openai.com/overview) and generate an API key [here](https://platform.openai.com/account/api-keys).

Once you have your key, call `lunapipe configure` to be prompted for your key, which will be stored for future use
```shell
lunapipe configure
> Enter your API key: <your openai api key>
> Wrote API key to ~/.lunapipe/config.yaml
```
Or put your key on the environment variable `OPENAI_API_KEY`
```bash
export OPENAI_API_KEY="<your openai api key>"
```

## Advanced Usage

#### Using GPT-4
`gpt-3.5-turbo` (ChatGPT) is used by default, but GPT-4 is also available by using the `-m` flag.
```bash
lunapipe -m gpt-4 "What are the improvements in GPT-4 compared to previous GPT models?"
```

#### Templates
Templates will compose a more complex message to the LLM. Different templates take different arguments. Several 
templates are included by default, or you can write your own. 

For example, the `function` template writes a function in the specified language.
```bash
lunapipe -t function -p language=go "Print 'Hello World'"
```
The included templates are `code, function, getinfo, markdown, rubberduck, shell, and summarizedir`.
To read more about templates and see the definitions for these templates, look in [prompt templates](prompts).

Templates are written in to go templating format, `.tmpl`. To expose your own templates to LunaPipe, you can pass a 
directory to where your templates are:
```bash
mkdir my_prompt_templates
echo "This is a test" > my_prompt_templates/test.tmpl
lunapipe --prompts my_prompt_templates -t test "Is this thing working?"

# If you want to add parameters to your template
echo 'This is a test with params: {{ index .Params "testparam" }} ' > my_prompt_templates/test_with_params.tmpl
lunapipe --prompts my_prompt_templates -t test -p testparam="Hello, world!" "Is this thing working?"
```



## Tips for contributors

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


