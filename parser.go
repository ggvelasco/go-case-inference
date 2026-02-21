package main

import (
	"regexp"
	"strings"
)

// Representa um nó do grafo
type Node struct {
	Name   string
	Result string // ex: "approved=true,segment=prime"
}

// Representa uma aresta (seta)
type Edge struct {
	From string // nó de origem
	To   string // nó de destino
	Cond string // condição (ex: "age>=18")
}

// Representa o grafo completo
type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge
}

func parseDOT(dot string) (*Graph, error) {
	graph := &Graph{
		Nodes: make(map[string]*Node),
		Edges: []*Edge{},
	}

	// Remove quebras de linha e espaços extras
	dot = strings.ReplaceAll(dot, "\n", " ")
	dot = strings.ReplaceAll(dot, "\t", " ")

	// Regex pra encontrar nós: nome [result="..."]
	nodeRegex := regexp.MustCompile(`(\w+)\s*\[result="([^"]*)"\]`)
	nodeMatches := nodeRegex.FindAllStringSubmatch(dot, -1)

	for _, match := range nodeMatches {
		nodeName := match[1]
		result := match[2]

		graph.Nodes[nodeName] = &Node{
			Name:   nodeName,
			Result: result,
		}
	}

	// Regex pra encontrar arestas: from -> to [cond="..."]
	edgeRegex := regexp.MustCompile(`(\w+)\s*->\s*(\w+)\s*\[cond="([^"]*)"\]`)
	edgeMatches := edgeRegex.FindAllStringSubmatch(dot, -1)

	for _, match := range edgeMatches {
		from := match[1]
		to := match[2]
		cond := match[3]

		graph.Edges = append(graph.Edges, &Edge{
			From: from,
			To:   to,
			Cond: cond,
		})
	}

	return graph, nil
}