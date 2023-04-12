# Prompt Templates

## Example

Here is an example of a template that generates a function in a given language:
```gotemplate
{{ $lang := index .Params "language" }}
Return a function{{if $lang }}, written in {{ $lang }},{{end}} based on the following requirements:
```

Running this template would look like:
```shell
go run main.go -t function -p language=go "Read values from a map"
```

The formatted prompt will look like:
```text
Return a function, written in go, based on the following requirements:
Read values from a map
```

Notice that `language` in the template is replaced with `go` in the prompt. 

We can edit this example so that it is more generic:
```gotemplate
{{ $lang := index .Params "language" }}
{{ $type := index .Params "type" }}
Return a {{ if $type}}{{$type}}{{end}}{{if $lang }}, written in {{ $lang }},{{end}} based on the following requirements:
```

You can pass multiple parameters to the template like so:
```shell
go run main.go -t code -p language=python -p type=class "Read values from a map"
class MapReader:
    def __init__(self, map):
        self.map = map
    
    def get_value(self, key):
        """
        Returns the value associated with the given key.
        If the key is not present in the map, returns None.
        """
        return self.map.get(key)
```

## Making your own
By creating a new template in the `prompts` directory, you can create your own prompt templates. The name of the file will be the name of the template. For example, if you create a file called `mytemplate.tmpl`, you can use it like so:
```shell
go run main.go -t mytemplate "Read values from a map"
```
