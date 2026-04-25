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
		// Front face (z = +3)
		{Pos: Vec3{-2, -2, 1}, Color: &Green}, // 0
		{Pos: Vec3{2, -2, 1}, Color: &Green},  // 1
		{Pos: Vec3{2, 2, 1}, Color: &Green},   // 2
		{Pos: Vec3{-2, 2, 1}, Color: &Green},  // 3

		// Back face (z = -3)
		{Pos: Vec3{-2, -2, -1}, Color: &Blue}, // 4
		{Pos: Vec3{2, -2, -1}, Color: &Red},   // 5
		{Pos: Vec3{2, 2, -1}, Color: &Red},    // 6
		{Pos: Vec3{-2, 2, -1}, Color: &Blue},  // 7
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

	screenVerts := make([]ScreenVertex, len(verts))

	fovY := float32(math.Pi / 4.0)
	aspectRatio := float32(width) / float32(height)

	cameraPos := Vec3{0, 0, 30}
	view := ViewMatrix(cameraPos)
	proj := Perspective(1, 100.0, fovY, aspectRatio)

	var model Matrix4
	MulMatrix4(&model, Translate(0, 0, 0), RotationAlongY(40))

	var temp Matrix4
	MulMatrix4(&temp, &view, &model)

	var mvp Matrix4
	MulMatrix4(&mvp, &proj, &temp)

	for i, v := range verts {
		v4 := Vec4{v.Pos.X, v.Pos.Y, v.Pos.Z, 1}

		clip := mvp.MultVec4(v4)

		if clip.W <= 0 {
			continue
		}

		divided := PerspectiveDivide(clip)

		screenVerts[i] = ScreenVertex{
			Pos: Vec3{
				X: (divided.X + 1) * 0.5 * float32(width),
				Y: (1 - divided.Y) * 0.5 * float32(height),
				Z: divided.Z,
			},
			W: clip.W,
		}

	}

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
			v1 := screenVerts[faces[i]]
			v2 := screenVerts[faces[i+1]]
			v3 := screenVerts[faces[i+2]]

			minX, maxX, minY, maxY := GetRectBounds(v1.Pos, v2.Pos, v3.Pos)

			minX = max(0, minX)
			minY = max(0, minY)
			maxX = min(width-1, maxX)
			maxY = min(height-1, maxY)

			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {

					// Create vec from center of pixel
					inTri, w0, w1, w2 := IsPixelInTriangle(Vec3{float32(x) + 0.5, float32(y) + 0.5, 0}, v1.Pos, v2.Pos, v3.Pos)
					if inTri {
						r, g, b, a := ColorFromWeights(w0, w1, w2, *verts[faces[i]].Color, *verts[faces[i+1]].Color, *verts[faces[i+2]].Color)
						coord := ToArrayCoordsYUp(x, y, width, height, 4)
						zCoord := ToArrayCoordsYUp(x, y, width, height, 1)

						interpolatedZ := w0*v1.Pos.Z + w1*v2.Pos.Z + w2*v3.Pos.Z
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
