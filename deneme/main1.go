package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
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
	return &Graph{Nodes: make(map[string]*Node)}
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
	dstNode.Neighbours = append(dstNode.Neighbours, srcNode)
}

// DFSAllPaths finds all paths from start to end using DFS
func (g *Graph) DFSAllPaths(start, end string) [][]string {
	var paths [][]string
	var path []string
	visited := make(map[string]bool)

	var dfs func(*Node)
	dfs = func(n *Node) {
		path = append(path, n.Name)
		visited[n.Name] = true

		if n.Name == end {
			paths = append(paths, append([]string{}, path...))
		} else {
			for _, neighbour := range n.Neighbours {
				if !visited[neighbour.Name] {
					dfs(neighbour)
				}
			}
		}

		path = path[:len(path)-1]
		visited[n.Name] = false
	}

	if startNode, endNode := g.Nodes[start], g.Nodes[end]; startNode != nil && endNode != nil {
		dfs(startNode)
	}
	return paths
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
	// Serve HTML file with embedded JavaScript
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Serve JavaScript file
	http.Handle("/main.js", http.FileServer(http.Dir(".")))

	// Listen and serve on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	if len(os.Args) < 2 {
		log.Fatal("Dosya adı belirtilmedi")
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Dosya açma hatası:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := NewGraph()
	var startNodeName, endNodeName string
	isAntLine := true
	numAnts := 0

	for scanner.Scan() {
		line := scanner.Text()
		if isAntLine {
			numAnts, err = strconv.Atoi(line)
			if err != nil || numAnts <= 0 {
				log.Fatal("ERROR: invalid number of ants")
			}
			isAntLine = false
			continue
		}

		switch {
		case line == "##start":
			readNode(scanner, graph, &startNodeName)
		case line == "##end":
			readNode(scanner, graph, &endNodeName)
		case strings.Contains(line, "-"):
			parts := strings.Split(line, "-")
			if len(parts) == 2 && parts[0] != parts[1] {
				graph.AddEdge(parts[0], parts[1])
			} else {
				log.Fatal("ERROR: invalid edge format")
			}
		default:
			addNode(line, graph)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Dosya okuma hatası:", err)
	}

	if startNodeName == "" || endNodeName == "" {
		log.Fatal("Start veya end düğümü tanımlanmadı")
	}

	allPaths := graph.DFSAllPaths(startNodeName, endNodeName)
	fmt.Println("All paths from", startNodeName, "to", endNodeName, ":")
	for _, path := range allPaths {
		fmt.Println("Path:", strings.Join(path, " -> "))
	}

	maxNonOverlappingPaths := findMaxNonOverlappingPaths(allPaths)
	fmt.Println("\nMaximum non-overlapping paths:")
	for _, path := range maxNonOverlappingPaths {
		fmt.Println(strings.Join(path, " -> "))
	}

	simulateAnts(numAnts, maxNonOverlappingPaths)
}

func readNode(scanner *bufio.Scanner, graph *Graph, nodeName *string) {
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 3 {
			*nodeName = parts[0]
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			graph.AddNode(parts[0], x, y)
		}
	}
}

func addNode(line string, graph *Graph) {
	parts := strings.Fields(line)
	if len(parts) == 3 {
		x, _ := strconv.Atoi(parts[1])
		y, _ := strconv.Atoi(parts[2])
		graph.AddNode(parts[0], x, y)
	}
}

func simulateAnts(ants int, paths [][]string) {
	if len(paths) == 0 {
		return
	}

	antPathAssignment := make([]int, ants)
	pathUsage := make([]int, len(paths))

	for i := 0; i < ants; i++ {
		pathIndex := i % len(paths)
		antPathAssignment[i] = pathIndex
		pathUsage[pathIndex]++
	}

	for i, usage := range pathUsage {
		if usage == 0 {
			for j := 0; j < ants; j++ {
				if pathUsage[antPathAssignment[j]] > 1 {
					pathUsage[antPathAssignment[j]]--
					antPathAssignment[j] = i
					pathUsage[i]++
					break
				}
			}
		}
	}

	steps := make([][]string, 0)
	positions := make([]int, ants)
	pathLengths := make([]int, ants)
	for i, idx := range antPathAssignment {
		pathLengths[i] = len(paths[idx])
	}

	occupied := make(map[string]int)

	for {
		moves := make([]string, 0)
		anyMoved := false
		for i := 0; i < ants; i++ {
			if positions[i] < pathLengths[i]-1 {
				nextPosition := positions[i] + 1
				nextRoom := paths[antPathAssignment[i]][nextPosition]

				if occupied[nextRoom] == 0 || nextRoom == paths[antPathAssignment[i]][pathLengths[i]-1] {
					if positions[i] > 0 {
						currentRoom := paths[antPathAssignment[i]][positions[i]]
						occupied[currentRoom] = 0
					}
					positions[i] = nextPosition
					occupied[nextRoom] = i + 1
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, nextRoom))
					anyMoved = true
				}
			}
		}
		if !anyMoved {
			break
		}
		steps = append(steps, moves)
	}

	for _, step := range steps {
		fmt.Println(strings.Join(step, " "))
	}
}
