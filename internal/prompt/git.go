package prompt

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/go-git/go-git/v5/storage/memory"
)

func loadTemplatesFromGit(repoURL string) {
	// Clone the repository into memory
	repo, err := cloneRepo(repoURL)
	if err != nil {
		fmt.Println("Error cloning the repository:", err)
		os.Exit(1)
	}

	// Create a new template by parsing templates from the Git repository
	tmpl, err := parseTemplates(repo)
	if err != nil {
		fmt.Println("Error creating template:", err)
		os.Exit(1)
	}

	// Get user input
	fmt.Print("Enter template data (name=value): ")
	var userInput string
	fmt.Scanln(&userInput)

	// Process user input
	data := make(map[string]string)
	for _, item := range strings.Fields(userInput) {
		pair := strings.SplitN(item, "=", 2)
		if len(pair) == 2 {
			data[pair[0]] = pair[1]
		}
	}

	// Execute template with user input data
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		os.Exit(1)
	}
}

func cloneRepo(url string) (*git.Repository, error) {
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func parseTemplates(repo *git.Repository) (*template.Template, error) {
	// Get a reference to the HEAD
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get the HEAD commit
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// Get the tree of the HEAD commit
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	// Initialize an empty template
	tmpl := template.New("")

	// Add templates to the template set by iterating over the files in the tree
	err = tree.Files().ForEach(func(file *object.File) error {
		// Add the template only if it has the ".tmpl" extension
		if strings.HasSuffix(file.Name, ".tmpl") {
			// Read the file content
			reader, err := file.Reader()
			if err != nil {
				return err
			}
			content, err := io.ReadAll(reader)
			if err != nil {
				return err
			}

			// Add the template to the template set
			tmpl, err = tmpl.New(file.Name).Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
