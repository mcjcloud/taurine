package util

import (
  "fmt"

  "github.com/mcjcloud/taurine/ast"
)

// ImportGraph represents a directed graph of imports
type ImportGraph struct {
  Matrix map[string][]string
  Nodes  map[string]*ImportNode
}

func NewImportGraph() *ImportGraph {
  return &ImportGraph{
    Matrix: make(map[string][]string),
    Nodes:  make(map[string]*ImportNode),
  }
}

// ImportNode represents a single node in an ImportGraph
type ImportNode struct {
  Path     string
  Children []*ImportNode
  Ast      *ast.Ast
}

// String implements string interface
func (n ImportNode) String() string {
  return n.Path
}

// SetAst sets the ast for the node
func (n *ImportNode) SetAst(tree *ast.Ast) {
  n.Ast = tree
}

// Print prints the graph starting with the given node
func (g ImportGraph) Print(start string) {
  if n, ok := g.Nodes[start]; ok {
    visited := make(map[string]bool)
    visited[n.Path] = true
    queue := []*ImportNode{n}

    // loop until the queue is empty
    for len(queue) > 0 {
      n = queue[0]
      // queue unvisited children
      for _, c := range n.Children {
        if v, ok := visited[c.Path]; !ok || !v {
          queue = append(queue, c)
        }
      }

      // print current top of queue
      fmt.Println(n)
      queue = queue[1:]
    }

  } else {
    fmt.Println("error: could not find start node")
  }
}

// Add creates a directed connection from src to dest. Returns the destination node
func (g *ImportGraph) Add(src, dest string) *ImportNode {
  if _, ok := g.Nodes[src]; !ok {
    g.Nodes[src] = &ImportNode{
      Path: src,
      Children: make([]*ImportNode, 0),
    }
  }
  if _, ok := g.Nodes[dest]; !ok {
    g.Nodes[dest] = &ImportNode{
      Path: dest,
      Children: make([]*ImportNode, 0),
    }
  }

  // create connection in matrix
  if ok := g.addMatrixConnection(src, dest); ok {
    g.Nodes[src].Children = append(g.Nodes[src].Children, g.Nodes[dest])
  }
  return g.Nodes[dest]
}

func (g *ImportGraph) addMatrixConnection(src, dest string) bool {
  if _, ok := g.Matrix[src]; !ok {
    g.Matrix[src] = make([]string, 0)
  }
  if _, ok := g.Matrix[dest]; !ok {
    g.Matrix[dest] = make([]string, 0)
  }

  for _, item := range g.Matrix[src] {
    if item == dest {
      return false
    }
  }

  g.Matrix[src] = append(g.Matrix[src], dest)
  return true
}

