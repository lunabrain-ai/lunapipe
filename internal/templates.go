package internal

import (
	"text/template"
	"text/template/parse"
)

func FindIndexCalls(tpl *template.Template) []string {
	var indexCalls []string
	ast := tpl.Tree

	var crawlNodes func(node parse.Node)
	crawlNodes = func(node parse.Node) {
		switch n := node.(type) {
		case *parse.ListNode:
			for _, listItem := range n.Nodes {
				crawlNodes(listItem)
			}
		case *parse.ActionNode:
			if n.Pipe != nil {
				for _, cmd := range n.Pipe.Cmds {
					if cmd.Args[0].String() == "index" {
						indexCalls = append(indexCalls, cmd.Args[2].String())
					}
				}
			}
		}
	}
	crawlNodes(ast.Root)
	return indexCalls
}