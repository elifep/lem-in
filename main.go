package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Room struct represents a room in the layout
type Room struct {
	Name string
	X, Y int
}

// Link struct represents a connection between rooms
type Link struct {
	From, To string
}

func main() {
	// Dosyadan veriyi oku
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var ants int
	var startRoom, endRoom string
	var rooms []Room
	var links []Link

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "##start") {
			parts := strings.Fields(line)
			startRoom = parts[1]
		} else if strings.HasPrefix(line, "##end") {
			parts := strings.Fields(line)
			endRoom = parts[1]
		} else if strings.HasPrefix(line, "room") {
			parts := strings.Fields(line)
			name := parts[0]
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			rooms = append(rooms, Room{Name: name, X: x, Y: y})
		} else if strings.HasPrefix(line, "link") {
			parts := strings.Fields(line)
			from := parts[1]
			to := parts[2]
			links = append(links, Link{From: from, To: to})
		} else {
			ants, _ = strconv.Atoi(line)
		}
	}

	// Ant hareketini hesapla
	fmt.Println(ants)
	fmt.Println(startRoom, endRoom)
	for _, room := range rooms {
		fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
	}
	for _, link := range links {
		fmt.Printf("%s %s\n", link.From, link.To)
	}

	// Antları hedefe taşı
	moveAnts(startRoom, endRoom, rooms, links, ants)
}

// moveAnts fonksiyonu, antları hedef odasına taşır
func moveAnts(startRoom, endRoom string, rooms []Room, links []Link, ants int) {
	// Antlar başlangıç odasında başlar
	currentRoom := startRoom
	fmt.Println("Ants' movements:")
	for i := 0; i < ants; i++ {
		// Hedefe ulaşana kadar antı taşı
		for currentRoom != endRoom {
			// Antın bulunduğu odadan bir sonraki odaya geç
			for _, link := range links {
				if link.From == currentRoom {
					currentRoom = link.To
					break
				}
			}
			// Antın hareketini yazdır
			fmt.Printf("Ant%d: %s -> %s\n", i+1, currentRoom, endRoom)
		}
		// Hedef odaya ulaşıldığında, antı başlangıç odayasına geri taşı
		currentRoom = startRoom
	}
}
