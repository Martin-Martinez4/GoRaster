package main

import (
	"fmt"
	"math"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

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
	// Black := Vec4{0, 0, 0, 255}

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

	// screenVerts := make([]ScreenVertex, len(verts))

	// aspectY := float32(height) / float32(width)
	aspectX := float32(width) / float32(height)

	fovY := float32(math.Pi / 4.0)
	fovx := float32(math.Atan(math.Tan(float64(fovY/2))*float64(aspectX)) * 2.0)
	zNear := float32(0.1)
	zFar := float32(100.0)

	// aspectRatio := float32(width) / float32(height)

	cameraPos := Vec3{0, 0, 8}
	view := ViewMatrix(cameraPos)

	proj := Perspective(zNear, zFar, fovY, aspectX) // correct

	initFrustumPlanes(fovY, fovx, zNear, zFar)

	var lastTime uint64
	var fps float32

	frameCountFrequency := 10
	var frameCount int
	var fpsAcc float32

	var rotationY float32 = 40.0

	viewVerts := make([]Vec4, len(verts))
	clipSpaceVerts := make([]Vec4, len(verts))

	// might delete later
	vertNormals := make([]Vec3, len(verts))
	lightValues := make([]float32, len(verts))

	for i := 0; i < len(faces); i += 3 {

		v0 := verts[faces[i]].Pos
		v1 := verts[faces[i+1]].Pos
		v2 := verts[faces[i+2]].Pos

		edge1 := Vec3{v1.X - v0.X, v1.Y - v0.Y, v1.Z - v0.Z}
		edge2 := Vec3{v2.X - v0.X, v2.Y - v0.Y, v2.Z - v0.Z}

		faceNormal := edge1.Cross(edge2)

		vertNormals[faces[i]].X += faceNormal.X
		vertNormals[faces[i]].Y += faceNormal.Y
		vertNormals[faces[i]].Z += faceNormal.Z

		vertNormals[faces[i+1]].X += faceNormal.X
		vertNormals[faces[i+1]].Y += faceNormal.Y
		vertNormals[faces[i+1]].Z += faceNormal.Z

		vertNormals[faces[i+2]].X += faceNormal.X
		vertNormals[faces[i+2]].Y += faceNormal.Y
		vertNormals[faces[i+2]].Z += faceNormal.Z
	}

	for i := range vertNormals {
		vertNormals[i] = vertNormals[i].Normalize()

	}

	// Baked in lighting; may remove later
	lightDir := Vec3{0.5, 1, 0.5}.Normalize()
	for i := range verts {

		intensity := max(float32(0), vertNormals[i].Dot(lightDir))
		// 0.2 ambient, 0.8 max
		lightValues[i] = 0.2 + intensity*0.8
	}

	sdl.RunLoop(func() error {
		currentTime := sdl.Ticks()
		deltaTime := float32(currentTime-lastTime) / 1000.0
		lastTime = currentTime

		// cap at 100ms to avoid first frame spike
		if deltaTime > 0.1 {
			deltaTime = 0.1
		}
		fpsAcc += 1.0 / deltaTime
		frameCount++

		var event sdl.Event

		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
		}

		// object space

		// Create world space
		var model Matrix4
		MulMatrix4(&model, Translate(0, 0, 0), RotationAlongY(rotationY))

		// Create view space
		var temp Matrix4
		MulMatrix4(&temp, &view, &model)

		for i, v := range verts {
			v4 := Vec4{v.Pos.X, v.Pos.Y, v.Pos.Z, 1}
			viewVerts[i] = temp.MultVec4(v4)

		}

		// Create Clip Space?
		// mvp is proj * view * model combined
		var mvp Matrix4
		MulMatrix4(&mvp, &proj, &temp)

		for i, v := range verts {
			v4 := Vec4{v.Pos.X, v.Pos.Y, v.Pos.Z, 1}
			clipSpaceVerts[i] = mvp.MultVec4(v4)
		}

		// copy instead of fill per pixel later
		for i := 0; i < len(pixels); i += 4 {
			pixels[i] = 0
			pixels[i+1] = 0
			pixels[i+2] = 0
			pixels[i+3] = 255
		}

		for i := 0; i < len(faces); i += 3 {

			// back face culling
			vs0 := viewVerts[faces[i]]
			vs1 := viewVerts[faces[i+1]]
			vs2 := viewVerts[faces[i+2]]

			l0 := lightValues[faces[i]]
			l1 := lightValues[faces[i+1]]
			l2 := lightValues[faces[i+2]]

			edge1 := Vec3{vs1.X - vs0.X, vs1.Y - vs0.Y, vs1.Z - vs0.Z}
			edge2 := Vec3{vs2.X - vs0.X, vs2.Y - vs0.Y, vs2.Z - vs0.Z}

			normal := edge1.Cross(edge2)

			// in view space camera is at origin so camera ray is just -vertex
			cameraRay := Vec3{-vs0.X, -vs0.Y, -vs0.Z}

			if normal.Dot(cameraRay) <= 0 {
				continue
			}

			v0 := clipSpaceVerts[faces[i]]
			v1 := clipSpaceVerts[faces[i+1]]
			v2 := clipSpaceVerts[faces[i+2]]

			col0 := *verts[faces[i]].Color
			col1 := *verts[faces[i+1]].Color
			col2 := *verts[faces[i+2]].Color

			poly := CreatePolygonFromTriangle(
				v0, v1, v2,
				Tex2{}, Tex2{}, Tex2{},
				col0, col1, col2,
			)
			ClipPolygon(&poly)

			trianglesAfterClipping := make([]Triangle, MAX_NUM_POLY_TRIANGLES)
			numTrianglesAfterClipping := 0
			TriangleFromPolygon(&poly, trianglesAfterClipping, &numTrianglesAfterClipping)

			// NDC Space stuff
			for j := 0; j < numTrianglesAfterClipping; j++ {

				col0 := trianglesAfterClipping[j].colors[0]
				col1 := trianglesAfterClipping[j].colors[1]
				col2 := trianglesAfterClipping[j].colors[2]

				tri := trianglesAfterClipping[j].points

				divided0 := PerspectiveDivide(tri[0])
				divided1 := PerspectiveDivide(tri[1])
				divided2 := PerspectiveDivide(tri[2])

				sv1 := Vec4{
					X: (divided0.X + 1) * 0.5 * float32(width),
					Y: (1 - divided0.Y) * 0.5 * float32(height),
					Z: 1.0 / tri[0].W,
					W: tri[0].W,
				}
				sv2 := Vec4{
					X: (divided1.X + 1) * 0.5 * float32(width),
					Y: (1 - divided1.Y) * 0.5 * float32(height),
					Z: 1.0 / tri[1].W,
					W: tri[1].W,
				}
				sv3 := Vec4{
					X: (divided2.X + 1) * 0.5 * float32(width),
					Y: (1 - divided2.Y) * 0.5 * float32(height),
					Z: 1.0 / tri[2].W,
					W: tri[2].W,
				}

				// divide each color channel by W at each vertex
				r0 := float32(col0.X) / sv1.W
				g0 := float32(col0.Y) / sv1.W
				b0 := float32(col0.Z) / sv1.W

				r1 := float32(col1.X) / sv2.W
				g1 := float32(col1.Y) / sv2.W
				b1 := float32(col1.Z) / sv2.W

				r2 := float32(col2.X) / sv3.W
				g2 := float32(col2.Y) / sv3.W
				b2 := float32(col2.Z) / sv3.W

				oneOverW0 := 1.0 / sv1.W
				oneOverW1 := 1.0 / sv2.W
				oneOverW2 := 1.0 / sv3.W

				// perspective correct prep for light, same as color
				l0OverW := l0 / sv1.W
				l1OverW := l1 / sv2.W
				l2OverW := l2 / sv3.W

				vec31 := Vec3{
					X: sv1.X,
					Y: sv1.Y,
					Z: sv1.Z,
				}

				vec32 := Vec3{
					X: sv2.X,
					Y: sv2.Y,
					Z: sv2.Z,
				}

				vec33 := Vec3{
					X: sv3.X,
					Y: sv3.Y,
					Z: sv3.Z,
				}

				minX, maxX, minY, maxY := GetRectBounds(vec31, vec32, vec33)

				minX = max(0, minX)
				minY = max(0, minY)
				maxX = min(width-1, maxX)
				maxY = min(height-1, maxY)

				// screen space stuff
				for y := minY; y <= maxY; y++ {
					for x := minX; x <= maxX; x++ {

						// Create vec from center of pixel
						inTri, w0, w1, w2 := IsPixelInTriangle(Vec3{float32(x) + 0.5, float32(y) + 0.5, 0}, vec31, vec32, vec33)
						if inTri {

							base := ToArrayCoordsYUp(x, y, width, height, 1)
							coord := base * 4
							zCoord := base

							// inside pixel loop, interpolate and recover
							interpR := w0*r0 + w1*r1 + w2*r2
							interpG := w0*g0 + w1*g1 + w2*g2
							interpB := w0*b0 + w1*b1 + w2*b2
							interpW := w0*oneOverW0 + w1*oneOverW1 + w2*oneOverW2

							// lighting
							interpL := w0*l0OverW + w1*l1OverW + w2*l2OverW
							finalLight := interpL / interpW

							finalR := byte(interpR / interpW * finalLight)
							finalG := byte(interpG / interpW * finalLight)
							finalB := byte(interpB / interpW * finalLight)

							interpolatedZ := w0*sv1.Z + w1*sv2.Z + w2*sv3.Z
							if interpolatedZ >= zbuffer[zCoord] {

								zbuffer[zCoord] = interpolatedZ

								pixels[coord] = byte(finalR)
								pixels[coord+1] = byte(finalG)
								pixels[coord+2] = byte(finalB)
								pixels[coord+3] = 255
							}

						}
					}
				}
			}

			// drawLineZ(v1.Pos, v2.Pos, width, height, Black, pixels, zbuffer)

		}

		texture.Update(nil, pixels, int32(width*4))

		renderer.Clear()
		renderer.RenderTexture(texture, nil, nil)

		if frameCount >= frameCountFrequency {
			fps = fpsAcc / float32(frameCount)
			fpsAcc = 0
			frameCount = 0
		}

		renderer.DebugText(10, 10, fmt.Sprintf("FPS: %.0f", fps))
		renderer.Present()

		rotationY += 1.0 * deltaTime
		for i := range zbuffer {
			zbuffer[i] = float32(math.Inf(-1))
		}

		return nil
	})
}
