package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type OBJSourceData struct {
	Positions []Vec3
	UVs       []Tex2
	Normals   []Vec3
}

func CollectAttributes(bufScanner *bufio.Scanner) OBJSourceData {
	var positions []Vec3
	var uvs []Tex2
	var normals []Vec3

	for bufScanner.Scan() {
		line := bufScanner.Text()
		// populate attribute slices

		if len(line) < 3 {
			continue
		}
		prefix := line[0:2]

		switch prefix {
		case "v ":
			// vertex
			pos := strings.Split(line[2:], " ")

			fmt.Printf("pos: %v\n", pos)
			f0, err := strconv.ParseFloat(pos[0], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex position")
				log.Fatal("Could not parse string to vertex position")
			}
			f1, err := strconv.ParseFloat(pos[1], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex position")
				log.Fatal("Could not parse string to vertex position")
			}
			f2, err := strconv.ParseFloat(pos[2], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex position")
				log.Fatal("Could not parse string to vertex position")
			}

			positions = append(positions, Vec3{float32(f0), float32(f1), float32(f2)})
			fmt.Printf("positions: %v\n", positions[len(positions)-1])

		case "vn":
			// normals
			pos := strings.Split(line[3:], " ")

			f0, err := strconv.ParseFloat(pos[0], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex normal")
				log.Fatal("Could not parse string to vertex normal")
			}
			f1, err := strconv.ParseFloat(pos[1], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex normal")
				log.Fatal("Could not parse string to vertex normal")
			}
			f2, err := strconv.ParseFloat(pos[2], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex normal")
				log.Fatal("Could not parse string to vertex normal")
			}

			normals = append(normals, Vec3{float32(f0), float32(f1), float32(f2)})

		case "vt":
			// texture coords
			vt := strings.Split(line[3:], " ")

			f0, err := strconv.ParseFloat(vt[0], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex texture")
				log.Fatal("Could not parse string to vertex texture")
			}
			f1, err := strconv.ParseFloat(vt[1], 32)
			if err != nil {
				fmt.Println("Could not parse string to vertex texture")
				log.Fatal("Could not parse string to vertex texture")
			}

			uvs = append(uvs, Tex2{float32(f0), float32(f1)})
		}
	}

	return OBJSourceData{
		Positions: positions,
		UVs:       uvs,
		Normals:   normals,
	}
}

func readObjFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	CollectAttributes(scanner)

}
