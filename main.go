package main

import (
	"math"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

func interpolateZValue(w0, w1, w2, z0, z1, z2 float32) float32 {
	return w0*z0 + w1*z1 + w2*z2
}

func main() {
	defer binsdl.Load().Unload()
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	width := 1080
	height := 900

	window, renderer, err := sdl.CreateWindowAndRenderer("Hello World!", width, height, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	renderer.SetDrawColor(255, 255, 255, 255)

	pixels := make([]byte, width*height*4)
	// use for double buffer setup later
	// pixels2 := make([]byte, width*height*4)

	zbuffer := make([]float32, width*height)
	for i := range zbuffer {
		zbuffer[i] = float32(math.Inf(-1))
	}

	Red := Vec4{255, 0, 0, 255}
	Green := Vec4{0, 255, 0, 255}
	Blue := Vec4{0, 0, 255, 255}

	verts := []Vertex{
		// Front face (z = 1)
		{Pos: Vec3{X: 200, Y: 100, Z: 1}, Color: &Green}, // 0
		{Pos: Vec3{X: 600, Y: 100, Z: 1}, Color: &Green}, // 1
		{Pos: Vec3{X: 600, Y: 500, Z: 1}, Color: &Green}, // 2
		{Pos: Vec3{X: 200, Y: 500, Z: 1}, Color: &Green}, // 3

		// Back face (z = 0)
		{Pos: Vec3{X: 200, Y: 100, Z: 0}, Color: &Blue}, // 4
		{Pos: Vec3{X: 600, Y: 100, Z: 0}, Color: &Red},  // 5
		{Pos: Vec3{X: 600, Y: 500, Z: 0}, Color: &Red},  // 6
		{Pos: Vec3{X: 200, Y: 500, Z: 0}, Color: &Blue}, // 7
	}

	faces := []uint32{
		// Front face
		0, 1, 2,
		0, 2, 3,

		// Back face
		5, 4, 7,
		5, 7, 6,

		// Left face
		4, 0, 3,
		4, 3, 7,

		// Right face
		1, 5, 6,
		1, 6, 2,

		// Top face
		4, 5, 1,
		4, 1, 0,

		// Bottom face
		3, 2, 6,
		3, 6, 7,
	}

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		width, height,
	)

	sdl.RunLoop(func() error {
		var event sdl.Event

		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
		}

		for i := 0; i < len(pixels); i += 4 {
			pixels[i] = 0
			pixels[i+1] = 0
			pixels[i+2] = 0
			pixels[i+3] = 255
		}

		for i := 0; i < len(faces); i += 3 {
			v1 := verts[faces[i]]
			v2 := verts[faces[i+1]]
			v3 := verts[faces[i+2]]

			minX, maxX, minY, maxY := GetRectBounds(v1.Pos, v2.Pos, v3.Pos)

			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {

					if minX < 0 || minY < 0 {
						continue
					}
					if maxX >= height || maxY >= width {
						continue
					}

					// Create vec from center of pixel
					inTri, w0, w1, w2 := IsPixelInTriangle(Vec3{float32(x) + 0.5, float32(y) + 0.5, 0}, v1.Pos, v2.Pos, v3.Pos)
					if inTri {
						r, g, b, a := ColorFromWeights(w0, w1, w2, *v1.Color, *v2.Color, *v3.Color)
						coord := ToArrayCoordsYUp(x, y, width, height, 4)
						zCoord := ToArrayCoordsYUp(x, y, width, height, 1)

						interpolatedZ := interpolateZValue(w0, w1, w2, v1.Pos.Z, v2.Pos.Z, v3.Pos.Z)
						if interpolatedZ >= zbuffer[zCoord] {

							zbuffer[zCoord] = interpolatedZ

							pixels[coord] = r
							pixels[coord+1] = g
							pixels[coord+2] = b
							pixels[coord+3] = a
						}

					}
				}
			}

		}

		texture.Update(nil, pixels, int32(width*4))

		renderer.Clear()
		renderer.RenderTexture(texture, nil, nil)
		// renderer.DebugText(50, 50, "Hello World")
		renderer.Present()

		return nil
	})
}
