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

				var d0Int int
				if d0 < 0 {
					d0Int = int(d0) + int(int64(len(objSD.Positions)))
				} else {
					d0Int = int(d0 - 1)
				}

				key := objAttribsTuple{V: d0Int, Vt: -1, Vn: -1}
				indx, ok := m[key]

				if !ok {
					m[key] = len(verts)
					faces = append(faces, uint32(len(verts)))
					verts = append(verts, Vertex{Pos: objSD.Positions[d0Int]})
				} else {
					faces = append(faces, uint32(indx))
				}
				// only vert
			case 2:
				// vert and texture
				d0, err := (strconv.ParseInt(data[0], 10, 0))
				if err != nil {
					panic("malformed vert index on face")
				}

				var d0Int int
				if d0 < 0 {
					d0Int = int(d0) + int(int64(len(objSD.Positions)))
				} else {
					d0Int = int(d0 - 1)
				}

				d1, err := (strconv.ParseInt(data[1], 10, 0))
				if err != nil {
					panic("malformed vert index on face")
				}
				var d1Int int
				if d1 < 0 {
					d1Int = int(d1) + int(int64(len(objSD.UVs)))
				} else {
					d1Int = int(d1 - 1)
				}

				fmt.Printf("d0: %d, d1: %d\n", d0, d1)

				key := objAttribsTuple{V: d0Int, Vt: d1Int, Vn: -1}
				indx, ok := m[key]

				if !ok {
					m[key] = len(verts)
					faces = append(faces, uint32(len(verts)))
					verts = append(verts, Vertex{Pos: objSD.Positions[d0Int], UV: &objSD.UVs[d1Int]})
				} else {
					faces = append(faces, uint32(indx))
				}
			case 3:
				// vert, texture, and normal
				// middle could be empty
				d0, err := strconv.ParseInt(data[0], 10, 0)
				if err != nil {
					panic("malformed vert index on face")
				}

				var d0Int int
				if d0 < 0 {
					d0Int = int(d0) + len(objSD.Positions)
				} else {
					d0Int = int(d0 - 1)
				}

				var d1Int = -1
				if data[1] != "" {
					d1, err := strconv.ParseInt(data[1], 10, 0)
					if err != nil {
						panic("malformed uv index on face")
					}

					if d1 < 0 {
						d1Int = int(d1) + len(objSD.UVs)
					} else {
						d1Int = int(d1 - 1)
					}
				}

				var d2Int = -1
				if len(data) > 2 && data[2] != "" {
					d2, err := strconv.ParseInt(data[2], 10, 0)
					if err != nil {
						panic("malformed normal index on face")
					}

					if d2 < 0 {
						d2Int = int(d2) + len(objSD.Normals)
					} else {
						d2Int = int(d2 - 1)
					}
				}

				var uv *Tex2
				if d1Int >= 0 {
					uv = &objSD.UVs[d1Int]
				}

				var normal *Vec3
				if d2Int >= 0 {
					normal = &objSD.Normals[d2Int]
				}

				key := objAttribsTuple{V: d0Int, Vt: d1Int, Vn: d2Int}
				indx, ok := m[key]

				if !ok {
					m[key] = len(verts)
					faces = append(faces, uint32(len(verts)))

					verts = append(verts, Vertex{
						Pos:    objSD.Positions[d0Int],
						UV:     uv,
						Normal: normal,
					})
				} else {
					faces = append(faces, uint32(indx))
				}
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
