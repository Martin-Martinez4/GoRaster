package main

import (
	"bufio"
	"log"
	"strings"
	"testing"
)

// const red = "\033[31m"
// const reset = "\033[0m"

func Test_readObjFile(t *testing.T) {

	testCases := []struct {
		name      string
		objString string
		positions []Vec3
		uvs       []Tex2
		normals   []Vec3
	}{

		{
			name: "3 position vertices",
			objString: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.0 1.0 0.0

f 1 2 3`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{0.0, 1.0, 0.0},
			},
			uvs:     []Tex2{},
			normals: []Vec3{},
		},
		{
			name: "4 position vertices",

			objString: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
v 0 1 0

f 1 2 3 4`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{1.0, 1.0, 0.0},
				Vec3{0.0, 1.0, 0.0},
			},
			uvs:     []Tex2{},
			normals: []Vec3{},
		},
		{
			name: "3 position vertices, 3 texture coords",
			objString: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0

vt 0.0 0.0
vt 1.0 0.0
vt 1.0 1.0

f 1/1 2/2 3/3`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				Tex2{0.0, 0.0},
				Tex2{1.0, 0.0},
				Tex2{1.0, 1.0},
			},
			normals: []Vec3{},
		},
		{
			name: "3 position vertices, 3 texture coords",
			objString: `v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
v 1 1 0


vt 0.0 0.0
vt 1.0 0.0
vt 1.0 1.0

f 1/1 2/2 3/3`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{1.0, 1.0, 0.0},
				Vec3{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				Tex2{0.0, 0.0},
				Tex2{1.0, 0.0},
				Tex2{1.0, 1.0},
			},
			normals: []Vec3{},
		},
		{
			name: "position vertices and normals",
			objString: `v 0.0 0.0 0.0
v 0.0 0.0 0.0
v 1 0 0
v 0 1 0


vn 0 0 1
vn 0 0 1
vn 0 1 0

f 1//1 2//2 3//3`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{0.0, 1.0, 0.0},
			},
			uvs: []Tex2{},
			normals: []Vec3{
				Vec3{0, 0, 1},
				Vec3{0, 0, 1},
				Vec3{0, 1, 0},
			},
		},
		{
			name: "position texture and normals",
			objString: `v 0 0 0
v 1 0 0
v 1 1 0

vt 0 0
vt 1 0
vt 1 1

vn 0 0 1
vn 0 0 1
vn 0 1 0

f 1/1/1 2/2/2 3/3/3`,
			positions: []Vec3{
				Vec3{0.0, 0.0, 0.0},
				Vec3{1.0, 0.0, 0.0},
				Vec3{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				Tex2{0, 0},
				Tex2{1, 0},
				Tex2{1, 1},
			},
			normals: []Vec3{
				Vec3{0, 0, 1},
				Vec3{0, 0, 1},
				Vec3{0, 1, 0},
			},
		},
		{
			name: "position texture and normals",
			objString: `v -1.5 0.25 3.14159
v 2.0 -3.0 0.0
v 0.0 0.0 -10.0

vt 0.333333 0.666666
vt 1.0 0.0
vt 0.0 1.0

vn 0.577 0.577 0.577
vn -1.0 0.0 0.0
vn 0.0 -1.0 0.0

f 1/1/1 2/2/2 3/3/3`,
			positions: []Vec3{
				Vec3{-1.5, 0.25, 3.14159},
				Vec3{2.0, -3.0, 0.0},
				Vec3{0.0, 0.0, -10.0},
			},
			uvs: []Tex2{
				Tex2{0.333333, 0.666666},
				Tex2{1.0, 0},
				Tex2{1, 1.0},
			},
			normals: []Vec3{
				Vec3{0.577, 0.577, 0.577},
				Vec3{-1.0, 0.0, 0.0},
				Vec3{0.0, -1.0, 0.0},
			},
		},
	}

	for i, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.objString))
			objAttrs := CollectAttributes(scanner)

			if len(objAttrs.Positions) != len(tt.positions) {
				log.Fatalf("Test #%d %s: positions len mismatch", i, tt.name)
			}
			if len(objAttrs.UVs) != len(tt.uvs) {
				log.Fatalf("Test #%d %s: uvs len mismatch", i, tt.name)
			}
			if len(objAttrs.Normals) != len(tt.normals) {
				log.Fatalf("Test #%d %s: normals len mismatch", i, tt.name)
			}

			for indx, p := range tt.positions {
				if objAttrs.Positions[indx].X != p.X {
					log.Fatalf("Test #%d %s: positions do not match X; objAttrs.Positions[i]: %v  p: %v", i, tt.name, objAttrs.Positions[indx], p)
				}
				if objAttrs.Positions[indx].Y != p.Y {
					log.Fatalf("Test #%d %s: positions do not match Y; objAttrs.Positions[i]: %v  p: %v", i, tt.name, objAttrs.Positions[indx], p)
				}
				if objAttrs.Positions[indx].Z != p.Z {
					log.Fatalf("Test #%d %s: positions do not match Z", i, tt.name)
				}
			}
		})

	}

}
