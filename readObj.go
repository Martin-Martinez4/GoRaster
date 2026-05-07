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

type ObjData struct {
	Verts []Vertex
	Faces []uint32
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

type objAttribsTuple struct {
	V, Vt, Vn int
}

func CollectObjData(bufScanner *bufio.Scanner, objSD *OBJSourceData) *ObjData {
	// Keep track of position, uvs, and normals index

	verts := []Vertex{}
	faces := []uint32{}

	m := map[objAttribsTuple]int{}

	for bufScanner.Scan() {
		line := bufScanner.Text()

		if len(line) < 3 {
			continue
		}

		prefix := line[0:2]

		if prefix != "f " {
			continue
		}

		// split at space
		// for each of those split at /
		spaceSplit := strings.Split(line[2:], " ")

		if len(spaceSplit) != 3 {
			panic("only tris supported at this time")
		}

		for _, ss := range spaceSplit {

			data := strings.Split(ss, "/")

			fmt.Printf("data: %v\n", data)

			switch len(data) {
			case 1:
				d0, err := (strconv.ParseInt(data[0], 10, 0))
				if err != nil {
					panic("malformed vert index on face")
				}

				key := objAttribsTuple{V: int(d0 - 1), Vt: -1, Vn: -1}
				indx, ok := m[key]

				if !ok {
					m[key] = len(verts)
					faces = append(faces, uint32(len(verts)))
					verts = append(verts, Vertex{Pos: objSD.Positions[int(d0-1)]})
				} else {
					fmt.Println("Reused")
					faces = append(faces, uint32(indx))
				}
				// only vert
			case 2:
				// vert and texture
			case 3:
				// vert, texture, and normal
				// middle could be empty
			}
		}

		// if new combo create new record

	}

	return &ObjData{
		Faces: faces,
		Verts: nil,
	}

}

func readObjFile(path string) *ObjData {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	objSourceData := CollectAttributes(scanner)

	scanner = bufio.NewScanner(f)

	objData := CollectObjData(scanner, &objSourceData)

	return objData

}
