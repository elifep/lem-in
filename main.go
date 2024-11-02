package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Node grafikteki bir düğümü temsil eder
type Node struct {
	Name       string  // Düğümün adı
	X, Y       int     // Düğümün koordinatları
	Neighbours []*Node // Düğümün komşuları
}

// Graph grafiği temsil eder
type Graph struct {
	Nodes map[string]*Node // Düğümlerin haritası
}

// NewGraph yeni bir grafik oluşturur
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
	}
}

// AddNode grafiğe bir düğüm ekler
func (g *Graph) AddNode(name string, x, y int) *Node {
	if node, exists := g.Nodes[name]; exists {
		return node // Düğüm zaten varsa, mevcut düğümü döndür
	}
	node := &Node{Name: name, X: x, Y: y}
	g.Nodes[name] = node // Yeni düğüm oluştur ve grafiğe ekle
	return node
}

// AddEdge grafiğe bir kenar ekler
func (g *Graph) AddEdge(src, dst string) {
	srcNode := g.Nodes[src]
	dstNode := g.Nodes[dst]
	srcNode.Neighbours = append(srcNode.Neighbours, dstNode) // Kenarı kaynak düğüme ekle
	if len(srcNode.Neighbours) < 3 {
		dstNode.Neighbours = append(dstNode.Neighbours, srcNode) // Grafiğin yönsüz olduğunu varsayıyoruz
	}
}

// DFSAllPaths başlangıçtan sona kadar tüm yolları DFS kullanarak bulur
func (g *Graph) DFSAllPaths(start, end string) []struct {
	Path []string
} {
	var result []struct {
		Path []string
	}
	var path []string
	visited := make(map[string]bool)

	var dfs func(*Node)
	dfs = func(n *Node) {
		path = append(path, n.Name) // Düğümü mevcut yola ekle
		visited[n.Name] = true      // Düğümü ziyaret edildi olarak işaretle

		if n.Name == end {
			result = append(result, struct {
				Path []string
			}{Path: append([]string{}, path...)}) // Yolu sonuca ekle
		} else {
			for _, neighbour := range n.Neighbours {
				if !visited[neighbour.Name] {
					dfs(neighbour) // Ziyaret edilmemiş komşuları ziyaret et
				}
			}
		}

		path = path[:len(path)-1] // Düğümü yoldan çıkar
		visited[n.Name] = false   // Düğümü tekrar ziyaret edilmemiş olarak işaretle
	}

	startNode := g.Nodes[start]
	endNode := g.Nodes[end]
	if startNode == nil || endNode == nil {
		return result // Başlangıç veya bitiş düğümü tanımlanmadıysa boş döner
	}
	dfs(startNode)
	return result
}

// isNonOverlapping iki yolun örtüşüp örtüşmediğini kontrol eder
func isNonOverlapping(path1, path2 []string) bool {
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})
	for _, node := range path1[1 : len(path1)-1] { // İlk ve son düğümleri dahil etme
		set1[node] = struct{}{}
	}
	for _, node := range path2[1 : len(path2)-1] { // İlk ve son düğümleri dahil etme
		set2[node] = struct{}{}
	}
	for node := range set1 {
		if _, ok := set2[node]; ok {
			return false // İki yol arasında ortak düğüm varsa örtüşüyorlar demektir
		}
	}
	return true // Örtüşmüyorlar
}

// findMaxNonOverlappingPaths maksimum sayıda örtüşmeyen yolları bulur
func findMaxNonOverlappingPaths(paths [][]string) [][]string {
	var maxPaths [][]string
	used := make([]bool, len(paths))

	for i, path := range paths {
		if used[i] {
			continue // Bu yol zaten kullanılmışsa atla
		}
		maxPaths = append(maxPaths, path)
		for j := i + 1; j < len(paths); j++ {
			if !isNonOverlapping(path, paths[j]) {
				used[j] = true // Örtüşen yolları işaretle
			}
		}
	}
	return maxPaths
}

func main() {
	startTime := time.Now() // Programın başlangıç zamanını kaydet

	filename := os.Args[1]         // Komut satırından dosya adını al
	file, err := os.Open(filename) // Dosyayı aç (Örneğin: example02.txt, example04.txt, example05.txt)
	if err != nil {
		log.Fatal("Dosya açma hatası:", err) // Dosya açılamazsa hata ver ve çık
	}
	defer file.Close() // Program bittiğinde dosyayı kapat

	scanner := bufio.NewScanner(file)
	graf := NewGraph()                        // Yeni bir grafik oluştur
	var startNodeName, endNodeName string     // Başlangıç ve bitiş düğümleri için değişkenler
	var readingStartNode, readingEndNode bool // Başlangıç ve bitiş düğümlerini okuma bayrakları
	isAnt := true                             // Karınca sayısını okuma bayrağı
	numAnt := 0                               // Karınca sayısı
	for scanner.Scan() {
		line := scanner.Text() // Satırı oku
		if isAnt {
			numAnt, _ = strconv.Atoi(scanner.Text()) // Karınca sayısını al
			if err != nil || numAnt <= 0 {
				fmt.Println("ERROR: invalid data format") // Geçersiz veri formatı hatası
				return
			}
			isAnt = false // Karınca sayısı okundu, bayrağı kapat
		}
		if line == "##start" {
			readingStartNode = true // Başlangıç düğümünü okuma moduna geç
			continue
		} else if line == "##end" {
			readingEndNode = true // Bitiş düğümünü okuma moduna geç
			continue
		}

		if readingStartNode {
			parts := strings.Fields(line) // Satırı parçalara ayır
			if len(parts) == 3 {
				startNodeName = parts[0]       // Başlangıç düğümünün adı
				x, _ := strconv.Atoi(parts[1]) // X koordinatı
				y, _ := strconv.Atoi(parts[2]) // Y koordinatı
				graf.AddNode(parts[0], x, y)   // Düğümü grafiğe ekle
			}
			readingStartNode = false // Başlangıç düğümü okundu, bayrağı kapat
		} else if readingEndNode {
			parts := strings.Fields(line) // Satırı parçalara ayır
			if len(parts) == 3 {
				endNodeName = parts[0]         // Bitiş düğümünün adı
				x, _ := strconv.Atoi(parts[1]) // X koordinatı
				y, _ := strconv.Atoi(parts[2]) // Y koordinatı
				graf.AddNode(parts[0], x, y)   // Düğümü grafiğe ekle
			}
			readingEndNode = false // Bitiş düğümü okundu, bayrağı kapat
		} else if strings.Contains(line, "-") {
			parts := strings.Split(line, "-") // Kenarı parçalara ayır
			if len(parts) != 2 || parts[0] == parts[1] {
				fmt.Println("ERROR: invalid data format") // Geçersiz veri formatı hatası
				return
			}
			graf.AddEdge(parts[0], parts[1]) // Kenarı grafiğe ekle
		} else {
			parts := strings.Fields(line) // Satırı parçalara ayır
			if len(parts) == 3 {
				x, _ := strconv.Atoi(parts[1]) // X koordinatı
				y, _ := strconv.Atoi(parts[2]) // Y koordinatı
				graf.AddNode(parts[0], x, y)   // Düğümü grafiğe ekle
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Dosya okuma hatası:", err) // Dosya okuma hatası varsa bildir ve çık
	}

	if startNodeName == "" || endNodeName == "" {
		log.Fatal("Start veya end düğümü tanımlanmadı") // Başlangıç veya bitiş düğümü yoksa hata ver
	}

	allPaths := graf.DFSAllPaths(startNodeName, endNodeName) // Başlangıçtan sona kadar tüm yolları bul
	fmt.Println("All paths from", startNodeName, "to", endNodeName, ":")
	for _, path := range allPaths {
		fmt.Printf("Path: %s\n", strings.Join(path.Path, " -> ")) // Tüm yolları yazdır
	}

	// Maksimum örtüşmeyen yolları bul
	pathStrings := make([][]string, len(allPaths))
	for i, path := range allPaths {
		pathStrings[i] = path.Path
	}
	maxNonOverlappingPaths := findMaxNonOverlappingPaths(pathStrings)

	fmt.Println("\nMaximum non-overlapping paths:")
	for _, path := range maxNonOverlappingPaths {
		fmt.Println(strings.Join(path, " -> ")) // Maksimum örtüşmeyen yolları yazdır
	}
	simulateAnts(numAnt, maxNonOverlappingPaths) // Karıncaları simüle et
	elapsed := time.Since(startTime)             // Geçen zamanı hesapla
	fmt.Println()
	fmt.Printf("Execution time: %.8f seconds\n", elapsed.Seconds()) // Geçen süreyi yazdır
}

// Karıncaları simüle eden fonksiyon
func simulateAnts(ants int, paths [][]string) {
	if len(paths) == 0 {
		return
	}

	// En kısa yolu bul
	shortestPathIndex := 0
	shortestPathLength := len(paths[0])
	for i := 1; i < len(paths); i++ {
		if len(paths[i]) < shortestPathLength {
			shortestPathIndex = i
			shortestPathLength = len(paths[i])
		}
	}

	// Karıncaları yollara dağıt
	antPathAssignment := make([]int, ants)
	pathUsage := make([]int, len(paths))

	// İlk karınca en kısa yoldan gidecek
	antPathAssignment[0] = shortestPathIndex
	pathUsage[shortestPathIndex]++

	// Diğer karıncaları yollara dengeli bir şekilde dağıt
	for i := 1; i < ants; i++ {
		pathIndex := i % len(paths)
		antPathAssignment[i] = pathIndex
		pathUsage[pathIndex]++
	}

	// Her yolun en az bir karınca tarafından kullanıldığını kontrol et
	for i, usage := range pathUsage {
		if usage == 0 { // Karınca atanmayan bir yol bul ve bir karınca ata
			for j := 1; j < ants; j++ {
				if pathUsage[antPathAssignment[j]] > 1 {
					pathUsage[antPathAssignment[j]]--
					antPathAssignment[j] = i
					pathUsage[i]++
					break
				}
			}
		}
	}

	// Simülasyonu başlat
	steps := make([][]string, 0)
	positions := make([]int, ants)
	pathLengths := make([]int, ants)
	for i, idx := range antPathAssignment {
		pathLengths[i] = len(paths[idx])
	}

	occupied := make(map[string]int) // Hangi karıncanın hangi odayı işgal ettiğini takip et

	endRoom := paths[0][len(paths[0])-1]

	for {
		moves := make([]string, 0)
		anyMoved := false
		for i := 0; i < ants; i++ {
			if positions[i] < pathLengths[i]-1 {
				nextPosition := positions[i] + 1
				nextRoom := paths[antPathAssignment[i]][nextPosition]

				// Karıncanın end odasına girmesini kontrol et
				if nextRoom == endRoom || occupied[nextRoom] == 0 { //end odası bir sonraki oda ve bos mu yani 0 mı ?
					if positions[i] > 0 {
						currentRoom := paths[antPathAssignment[i]][positions[i]]
						occupied[currentRoom] = 0 //şuan bulduğu odayı boş yap
					}
					positions[i] = nextPosition
					occupied[nextRoom] = i + 1
					//fmt.Print("i:", i)
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

	// Sonucu yazdır
	fmt.Println("Tur sayısı:", len(steps))
	for _, step := range steps {
		// if i == len(steps)-1 {
		// 	return
		// }
		fmt.Println(strings.Join(step, " "))

	}

}
