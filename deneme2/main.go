package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Node represents a node in the graph
type Node struct {
	Name       string
	X, Y       int
	Neighbours []*Node
}

// Graph represents a graph
type Graph struct {
	Nodes map[string]*Node
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(name string, x, y int) *Node {
	if node, exists := g.Nodes[name]; exists {
		return node
	}
	node := &Node{Name: name, X: x, Y: y}
	g.Nodes[name] = node
	return node
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(src, dst string) {
	srcNode := g.Nodes[src]
	dstNode := g.Nodes[dst]
	srcNode.Neighbours = append(srcNode.Neighbours, dstNode)
	dstNode.Neighbours = append(dstNode.Neighbours, srcNode) // Assuming the graph is undirected
}

// Mesafe calculates the Euclidean distance between two nodes
func Mesafe(a, b *Node) float64 {
	return math.Sqrt(float64((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y)))
}

// DFSAllPaths finds all paths from start to end using DFS
func (g *Graph) DFSAllPaths(start, end string) []struct {
	Path     []string
	Distance float64
} {
	var result []struct {
		Path     []string
		Distance float64
	}
	var path []string
	var distance float64
	visited := make(map[string]bool)

	var dfs func(*Node)
	dfs = func(n *Node) {
		path = append(path, n.Name)
		visited[n.Name] = true

		if n.Name == end {
			result = append(result, struct {
				Path     []string
				Distance float64
			}{Path: append([]string{}, path...), Distance: distance})
		} else {
			for _, neighbour := range n.Neighbours {
				if !visited[neighbour.Name] {
					distance += Mesafe(n, neighbour)
					dfs(neighbour)
					distance -= Mesafe(n, neighbour)
				}
			}
		}

		path = path[:len(path)-1]
		visited[n.Name] = false
	}

	startNode := g.Nodes[start]
	endNode := g.Nodes[end]
	if startNode == nil || endNode == nil {
		return result // Return empty if start or end is not defined
	}
	dfs(startNode)
	return result
}

// isNonOverlapping checks if two paths are non-overlapping
func isNonOverlapping(path1, path2 []string) bool {
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})
	for _, node := range path1[1 : len(path1)-1] {
		set1[node] = struct{}{}
	}
	for _, node := range path2[1 : len(path2)-1] {
		set2[node] = struct{}{}
	}
	for node := range set1 {
		if _, ok := set2[node]; ok {
			return false
		}
	}
	return true
}

// findMaxNonOverlappingPaths finds the maximum number of non-overlapping paths
func findMaxNonOverlappingPaths(paths [][]string) [][]string {
	var maxPaths [][]string
	used := make([]bool, len(paths))

	for i, path := range paths {
		if used[i] {
			continue
		}
		maxPaths = append(maxPaths, path)
		for j := i + 1; j < len(paths); j++ {
			if !isNonOverlapping(path, paths[j]) {
				used[j] = true
			}
		}
	}
	return maxPaths
}

func main() {
	file, err := os.Open("example00.txt")
	if err != nil {
		log.Fatal("Dosya açma hatası:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graf := NewGraph()
	var startNodeName, endNodeName string
	var readingStartNode, readingEndNode bool
	isAnt := true
	numAnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		if isAnt {
			numAnt, _ = strconv.Atoi(scanner.Text())
			isAnt = false
		}
		if line == "##start" {
			readingStartNode = true
			continue
		} else if line == "##end" {
			readingEndNode = true
			continue
		}

		if readingStartNode {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				startNodeName = parts[0]
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graf.AddNode(parts[0], x, y)
			}
			readingStartNode = false
		} else if readingEndNode {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				endNodeName = parts[0]
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graf.AddNode(parts[0], x, y)
			}
			readingEndNode = false
		} else if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			graf.AddEdge(parts[0], parts[1])
		} else {
			parts := strings.Fields(line)
			if len(parts) == 3 {
				x, _ := strconv.Atoi(parts[1])
				y, _ := strconv.Atoi(parts[2])
				graf.AddNode(parts[0], x, y)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Dosya okuma hatası:", err)
	}

	if startNodeName == "" || endNodeName == "" {
		log.Fatal("Start veya end düğümü tanımlanmadı")
	}

	allPaths := graf.DFSAllPaths(startNodeName, endNodeName)
	fmt.Println("All paths from", startNodeName, "to", endNodeName, ":")
	for _, path := range allPaths {
		fmt.Printf("Path: %s, Total Distance: %.2f\n", strings.Join(path.Path, " -> "), path.Distance)
	}

	// Find maximum non-overlapping paths
	pathStrings := make([][]string, len(allPaths))
	for i, path := range allPaths {
		pathStrings[i] = path.Path
	}
	maxNonOverlappingPaths := findMaxNonOverlappingPaths(pathStrings)

	fmt.Println("\nMaximum non-overlapping paths:")
	for _, path := range maxNonOverlappingPaths {
		fmt.Println(strings.Join(path, " -> "))
	}
	simulateAnts(numAnt, maxNonOverlappingPaths)
}

func simulateAnts(ants int, paths [][]string) {
	steps := make([][]string, 0)
	positions := make([]int, ants)
	pathLengths := make([]int, ants)

	for i := 0; i < ants; i++ {
		pathLengths[i] = len(paths[i%len(paths)])
	}

	occupied := make(map[string]bool)

	for {
		moves := make([]string, 0)
		moved := false
		for i := 0; i < ants; i++ {
			if positions[i] < pathLengths[i]-1 {
				nextPosition := positions[i] + 1
				nextRoom := paths[i%len(paths)][nextPosition]

				if !occupied[nextRoom] {
					if positions[i] > 0 {
						occupied[paths[i%len(paths)][positions[i]]] = false
					}
					positions[i] = nextPosition
					occupied[nextRoom] = true
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, nextRoom))
					moved = true
				}
			}
		}
		if !moved {
			break
		}
		steps = append(steps, moves)
	}

	for _, step := range steps {
		fmt.Println(strings.Join(step, " "))
	}
}
