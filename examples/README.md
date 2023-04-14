# Examples directory for lunapipe

This directory contains examples of commands executed using lunapipe. 

## Example commands

- `lunapipe "Write me json for a recipe" > recipe.json`: Generates json for a recipe
- `cat recipe.json | aigo "Parse recipe, stored in 'recipe.json' with struct" > parse.go`: Generates a go file to parse the recipe json
- `cat parse.go | aigo "Write me a test for the function ParseRecipe" > parse_test.go`: Generates a go test file for the ParseRecipe function
- `go run main.go -t code -p language=ruby -p type=struct "Recipe"`: Generates code for a Ruby struct named Recipe
- `go run main.go -t code -p language=cpp -p type=struct "Recipe"`: Generates code for a C++ struct named Recipe
- `go run main.go -t code -p language=c# -p type=struct "Recipe"`: Generates code for a C# struct named Recipe
- `lunapipe "Write me a python discord bot. Say hello when the user mentions the bot."`: Executes a task to write a python discord bot that says hello when mentioned
- `lunapipe -t code -p language=python -p type=file "Write me a discord bot. Say hello when the user mentions the bot. Include intents." > discord_bot.py`: Generates
- `pip install discord`: Installs the discord module
- `python discord_bot.py`: Runs the python file for the discord bot again after installing the discord module
