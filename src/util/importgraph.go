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

func NewImportGraph(root string) *ImportGraph {
  return &ImportGraph{
    Matrix: map[string][]string{root: make([]string, 0)},
    Nodes: map[string]*ImportNode{root: {
      Path: root,
      Children: make([]*ImportNode, 0),
    }},
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

// FindCycles returns an ordered list representing the import cycle
func (graph *ImportGraph) FindCycles() []string {
  w := make([]*ImportNode, len(graph.Nodes))
  g := make([]*ImportNode, 0)
  b := make([]*ImportNode, 0)
  dfsMap := make(map[string]string)

  // populate unvisited set
  var i int
  for _, n := range graph.Nodes {
    w[i] = n
    i++
  }

  for len(w) > 0 {
    // add the first element to the stack
    n := w[0]
    w = w[1:]
    g = append(g, n)
    dfsMap[n.Path] = ""

    for _, c := range n.Children {
      if _, ok := dfsMap[c.Path]; !ok {
        dfsMap[c.Path] = n.Path
      }
      l := findCycles(c, &w, &g, &b, dfsMap)
      if len(l) > 0 {
        return l
      }
      // take c out of g and into b
      vs := remove(&g, c)
      g = *vs
      b = append(b, c)
    }
  }
  return make([]string, 0)
}

func findCycles(node *ImportNode, w, g, b *[]*ImportNode, dfsMap map[string]string) []string {
  // take node out of w
  remove(w, node)
  // add node to g
  *g = append(*g, node)
  // go over the children
  for _, c := range node.Children {
    if _, ok := dfsMap[c.Path]; !ok {
      dfsMap[c.Path] = node.Path
    }
    // check b
    if contains(b, c) {
      continue
    }
    // check g
    if contains(g, c) {
      // cycle detected
      l := []string{c.Path}
      x := node.Path
      for x != "" {
        l = append(l, x)
        x = dfsMap[x]
      }
      return reverse(l)
    }

    // recursive call
    l := findCycles(c, w, g, b, dfsMap)
    if len(l) > 0 {
      return l
    }
    // take c out of g
    g = remove(g, c)
    // add c to b
    *b = append(*b, c)
  }
  return make([]string, 0)
}

func contains(l *[]*ImportNode, v *ImportNode) bool {
  for _, n := range *l {
    if n == v {
      return true
    }
  }
  return false
}

func remove(l *[]*ImportNode, v *ImportNode) *[]*ImportNode {
  for i, n := range *l {
    if n == v {
      *l =  append((*l)[:i], (*l)[i+1:]...)
      return l
    }
  }
  return l
}

func reverse(l []string) []string {
  res := make([]string, len(l))
  for i, n := range l {
    res[len(res)-i-1] = n
  }
  return res
}

