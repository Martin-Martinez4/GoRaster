package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"testing"
)

// const red = "\033[31m"
// const reset = "\033[0m"

func Test_CollectAttributes(t *testing.T) {

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
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
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
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{1.0, 1.0, 0.0},
				{0.0, 1.0, 0.0},
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
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				{0.0, 0.0},
				{1.0, 0.0},
				{1.0, 1.0},
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
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{1.0, 1.0, 0.0},
				{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				{0.0, 0.0},
				{1.0, 0.0},
				{1.0, 1.0},
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
				{0.0, 0.0, 0.0},
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
			},
			uvs: []Tex2{},
			normals: []Vec3{
				{0, 0, 1},
				{0, 0, 1},
				{0, 1, 0},
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
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{1.0, 1.0, 0.0},
			},
			uvs: []Tex2{
				{0, 0},
				{1, 0},
				{1, 1},
			},
			normals: []Vec3{
				{0, 0, 1},
				{0, 0, 1},
				{0, 1, 0},
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
				{-1.5, 0.25, 3.14159},
				{2.0, -3.0, 0.0},
				{0.0, 0.0, -10.0},
			},
			uvs: []Tex2{
				{0.333333, 0.666666},
				{1.0, 0},
				{1, 1.0},
			},
			normals: []Vec3{
				{0.577, 0.577, 0.577},
				{-1.0, 0.0, 0.0},
				{0.0, -1.0, 0.0},
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

func Test_collectObjData(t *testing.T) {

	testCases := []struct {
		name       string
		faceString string
		positions  []Vec3
		uvs        []Tex2
		normals    []Vec3
		faces      []uint32
	}{

		{
			name:       "parses triangle face with 3 position-only vertices",
			faceString: `f 1 2 3`,

			positions: []Vec3{
				{10.0, 20.0, 30.0},
				{-5.0, 0.5, 2.0},
				{7.25, -9.0, 1.0},
			},

			uvs:     nil,
			normals: nil,
			faces:   []uint32{0, 1, 2},
		},
		{
			name: "parses multiple triangle faces",

			faceString: `f 1 2 3
f 3 2 4
`,

			positions: []Vec3{
				{0.0, 0.0, 0.0},
				{1.0, 0.0, 0.0},
				{0.0, 1.0, 0.0},
				{1.0, 1.0, 0.0},
			},

			uvs:     nil,
			normals: nil,

			faces: []uint32{
				0, 1, 2,
				2, 1, 3,
			},
		},
		{
			name: "parses multiple disconnected triangle faces",

			faceString: `f 1 2 3
f 4 5 6
`,

			positions: []Vec3{
				{0.0, 0.0, 0.0}, // 1
				{1.0, 0.0, 0.0}, // 2
				{0.0, 1.0, 0.0}, // 3
				{5.0, 5.0, 5.0}, // 4
				{6.0, 5.0, 5.0}, // 5
				{5.0, 6.0, 5.0}, // 6
			},

			uvs:     nil,
			normals: nil,

			faces: []uint32{
				0, 1, 2,
				3, 4, 5,
			},
		},
		{
			name: "parses multiple faces with reused vertices",

			faceString: `f 1 2 3
f 1 3 4
f 1 4 2
`,

			positions: []Vec3{
				{0.0, 0.0, 0.0}, // 1
				{1.0, 0.0, 0.0}, // 2
				{1.0, 1.0, 0.0}, // 3
				{0.0, 1.0, 0.0}, // 4
			},

			uvs:     nil,
			normals: nil,

			faces: []uint32{
				0, 1, 2,
				0, 2, 3,
				0, 3, 1,
			},
		},
		// {
		// 	name: "3 position vertices",

		// 	faceString: `f 1 2 3`,
		// 	positions: []Vec3{
		// 		{0.0, 0.0, 0.0},
		// 		{1.0, 0.0, 0.0},
		// 		{1.0, 1.0, 0.0},
		// 		{0.0, 1.0, 0.0},
		// 	},
		// 	uvs:     []Tex2{},
		// 	normals: []Vec3{},
		// },
		// {
		// 	name:       "3 position vertices, 3 texture coords",
		// 	faceString: `f 1/1 2/2 3/3`,
		// 	positions: []Vec3{
		// 		{0.0, 0.0, 0.0},
		// 		{1.0, 0.0, 0.0},
		// 		{1.0, 1.0, 0.0},
		// 	},
		// 	uvs: []Tex2{
		// 		{0.0, 0.0},
		// 		{1.0, 0.0},
		// 		{1.0, 1.0},
		// 	},
		// 	normals: []Vec3{},
		// },
		// {
		// 	name:       "3 position vertices, 3 texture coords",
		// 	faceString: `f 1/1 2/2 3/3`,
		// 	positions: []Vec3{
		// 		{0.0, 0.0, 0.0},
		// 		{1.0, 0.0, 0.0},
		// 		{1.0, 1.0, 0.0},
		// 		{1.0, 1.0, 0.0},
		// 	},
		// 	uvs: []Tex2{
		// 		{0.0, 0.0},
		// 		{1.0, 0.0},
		// 		{1.0, 1.0},
		// 	},
		// 	normals: []Vec3{},
		// },
		// {
		// 	name:       "position vertices and normals",
		// 	faceString: `f 1//1 2//2 3//3`,
		// 	positions: []Vec3{
		// 		{0.0, 0.0, 0.0},
		// 		{0.0, 0.0, 0.0},
		// 		{1.0, 0.0, 0.0},
		// 		{0.0, 1.0, 0.0},
		// 	},
		// 	uvs: []Tex2{},
		// 	normals: []Vec3{
		// 		{0, 0, 1},
		// 		{0, 0, 1},
		// 		{0, 1, 0},
		// 	},
		// },
		// {
		// 	name:       "position texture and normals",
		// 	faceString: `f 1/1/1 2/2/2 3/3/3`,
		// 	positions: []Vec3{
		// 		{0.0, 0.0, 0.0},
		// 		{1.0, 0.0, 0.0},
		// 		{1.0, 1.0, 0.0},
		// 	},
		// 	uvs: []Tex2{
		// 		{0, 0},
		// 		{1, 0},
		// 		{1, 1},
		// 	},
		// 	normals: []Vec3{
		// 		{0, 0, 1},
		// 		{0, 0, 1},
		// 		{0, 1, 0},
		// 	},
		// },
		// {
		// 	name:       "position texture and normals",
		// 	faceString: `f 1/1/1 2/2/2 3/3/3`,
		// 	positions: []Vec3{
		// 		{-1.5, 0.25, 3.14159},
		// 		{2.0, -3.0, 0.0},
		// 		{0.0, 0.0, -10.0},
		// 	},
		// 	uvs: []Tex2{
		// 		{0.333333, 0.666666},
		// 		{1.0, 0},
		// 		{1, 1.0},
		// 	},
		// 	normals: []Vec3{
		// 		{0.577, 0.577, 0.577},
		// 		{-1.0, 0.0, 0.0},
		// 		{0.0, -1.0, 0.0},
		// 	},
		// },
	}

	for _, tt := range testCases {

		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.faceString))

			objSD := &OBJSourceData{
				Positions: tt.positions,
				UVs:       tt.uvs,
				Normals:   tt.normals,
			}

			fmt.Println(tt.faceString)
			got := CollectObjData(scanner, objSD)

			for i, f := range got.Faces {
				if f != tt.faces[i] {
					log.Fatalf("Test #%d %s: faces do not match; want: %v; got: %v", i, tt.name, tt.faces, got.Faces)
				}
			}

		})

	}

}
