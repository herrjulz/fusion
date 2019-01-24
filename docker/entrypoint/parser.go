package entrypoint

import (
	"errors"
	"io"
	"strings"

	"github.com/JulzDiverse/fusion/docker"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func Parse(dockerfile io.Reader) (docker.Entrypoint, error) {
	result, err := parser.Parse(dockerfile)
	if err != nil {
		return docker.Entrypoint{}, err
	}

	for _, c := range result.AST.Children {
		if c.Value == "entrypoint" {
			currentNode := c.Next
			if currentNode.Next == nil {
				return parseShellForm(currentNode.Value), nil
			} else {
				return parseExecForm(currentNode.Value, currentNode.Next), nil
			}
		}
	}
	return docker.Entrypoint{}, errors.New("Could not find entrypoint")
}

func parseShellForm(params string) (entrypoint docker.Entrypoint) {
	result := strings.Split(params, " ")
	entrypoint.Executable = result[0]
	entrypoint.Args = result[1:]
	return
}

func parseExecForm(executable string, node *parser.Node) (entrypoint docker.Entrypoint) {
	entrypoint.Executable = executable
	for node != nil {
		entrypoint.Args = append(entrypoint.Args, node.Value)
		node = node.Next
	}
	return
}
