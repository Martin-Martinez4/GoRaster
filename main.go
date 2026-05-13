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

	width := 1280
	height := 720

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

	// Red := Vec4{255, 0, 0, 255}
	// White := Vec4{0, 255, 0, 255}
	// Blue := Vec4{0, 0, 255, 255}
	White := Vec4{255, 255, 255, 255}
	// Black := Vec4{0, 0, 0, 255}

	textureImg := LoadTexture("./assets/textures/mossy_brick_diff_4k.jpg")

	objData := ReadObjFile("./assets/obj/maze1.obj")
	faces := objData.Faces
	verts := objData.Verts

	for i := range verts {
		verts[i].Color = &White

		// verts[i].Pos.X *= 4
		// verts[i].Pos.Y *= 8
		// verts[i].Pos.Z *= 4
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
	zFar := float32(1000.0)

	proj := Perspective(zNear, zFar, fovY, aspectX)

	initFrustumPlanes(fovY, fovx, zNear, zFar)

	var lastTime uint64
	var fps float32

	frameCountFrequency := 10
	var frameCount int
	var fpsAcc float32

	// var rotationY float32 = 40.0

	viewVerts := make([]Vec4, len(verts))
	clipSpaceVerts := make([]Vec4, len(verts))

	// might delete later
	vertNormals := make([]Vec3, len(verts))
	lightValues := make([]Vec3, len(verts))

	for i := range vertNormals {
		vertNormals[i] = verts[i].Normal.Normalize()

	}

	// Baked in lighting; may remove later
	lightDir := Vec3{10, 10, 10}.Normalize()
	lightColor := Vec3{0.1, 0.5, 0.2}
	ambient := Vec3{0.4, 0.4, 0.8}
	for i := range verts {

		intensity := max(float32(0), vertNormals[i].Dot(lightDir))
		// 0.2 ambient, 0.8 max
		lightValues[i].X = ambient.X + intensity*0.8*lightColor.X
		lightValues[i].Y = ambient.Y + intensity*0.8*lightColor.Y
		lightValues[i].Z = ambient.Z + intensity*0.8*lightColor.Z
	}

	// rotation := IdentityMatrix()
	rotationAmount := 5 * float32(math.Pi/180)
	speed := float32(5.0)
	// yaw := float32(0.0)

	camera := Camera{Position: Vec3{-127, -16, -230}, Yaw: 0.0, Pitch: 0.0}

	// cameraPos := Vec3{0, 0, 75}

	// translate := Vec3{0, 0, 0}

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
			switch event.Type {
			case sdl.EVENT_QUIT:
				return sdl.EndLoop

			case sdl.EVENT_KEY_DOWN:
				switch event.KeyboardEvent().Key {
				case sdl.K_A:
					camera.Yaw += rotationAmount
				case sdl.K_D:
					camera.Yaw -= rotationAmount

				case sdl.K_UP:
					camera.Pitch += rotationAmount
				case sdl.K_DOWN:
					camera.Pitch -= rotationAmount

				case sdl.K_W:
					camera.Position.Z -= float32(math.Cos((float64(camera.Yaw)))) * speed
					camera.Position.X -= float32(math.Sin((float64(camera.Yaw)))) * speed
				case sdl.K_S:
					camera.Position.Z += float32(math.Cos((float64(camera.Yaw)))) * speed
					camera.Position.X += float32(math.Sin((float64(camera.Yaw)))) * speed
				}

			}

		}

		// object space

		// Create world space
		// angle := 90 * float32(math.Pi/180)
		translation := ViewMatrix(camera.Position)
		var rotation Matrix4
		MulMatrix4(&rotation, RotationAlongY(camera.Yaw), RotationAlongX(camera.Pitch))

		// -- Do collision here

		var view Matrix4
		MulMatrix4(&view, &rotation, &translation)

		// var model Matrix4
		// MulMatrix4(&model, &idMatrix, &idMatrix)
		var model = IdentityMatrix()

		// // dynamic lighting
		CalculateSimpleLighting(vertNormals, lightValues, lightColor, ambient, model)

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
			if cullBackFace(viewVerts[faces[i]], viewVerts[faces[i+1]], viewVerts[faces[i+2]]) {
				continue
			}

			poly := CreatePolygonFromTriangle(
				clipSpaceVerts[faces[i]], clipSpaceVerts[faces[i+1]], clipSpaceVerts[faces[i+2]],
				verts[faces[i]].UV, verts[faces[i+1]].UV, verts[faces[i+2]].UV,
				*verts[faces[i]].Color, *verts[faces[i+1]].Color, *verts[faces[i+2]].Color,
				lightValues[faces[i]], lightValues[faces[i+1]], lightValues[faces[i+2]],
			)
			ClipPolygon(&poly)

			trianglesAfterClipping := make([]Triangle, MAX_NUM_POLY_TRIANGLES)
			numTrianglesAfterClipping := 0
			TriangleFromPolygon(&poly, trianglesAfterClipping, &numTrianglesAfterClipping)

			// NDC Space stuff
			for j := 0; j < numTrianglesAfterClipping; j++ {

				tris := trianglesAfterClipping[j].points

				sv1, sv2, sv3 := ClipToScreenSpace(tris[0], tris[1], tris[2], width, height)

				col0 := trianglesAfterClipping[j].colors[0]
				col1 := trianglesAfterClipping[j].colors[1]
				col2 := trianglesAfterClipping[j].colors[2]

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

				uv0 := trianglesAfterClipping[j].textcoords[0]
				uv1 := trianglesAfterClipping[j].textcoords[1]
				uv2 := trianglesAfterClipping[j].textcoords[2]

				u0OverW := uv0.U / sv1.W
				v0OverW := uv0.V / sv1.W

				u1OverW := uv1.U / sv2.W
				v1OverW := uv1.V / sv2.W

				u2OverW := uv2.U / sv3.W
				v2OverW := uv2.V / sv3.W

				// perspective correct prep for light, same as color
				l0OverW, l1OverW, l2OverW := GetLightOverW(trianglesAfterClipping[j].lights[0], trianglesAfterClipping[j].lights[1], trianglesAfterClipping[j].lights[2], sv1.W, sv2.W, sv3.W)

				oneOverW0 := 1.0 / sv1.W
				oneOverW1 := 1.0 / sv2.W
				oneOverW2 := 1.0 / sv3.W

				// screen space stuff

				area := edge(vec31, vec32, vec33)
				invArea := 1.0 / area
				for y := minY; y <= maxY; y++ {
					for x := minX; x <= maxX; x++ {

						// Create vec from center of pixel
						inTri, w0, w1, w2 := IsPixelInTriangle(Vec3{float32(x) + 0.5, float32(y) + 0.5, 0}, vec31, vec32, vec33, invArea)
						if inTri {

							base := ToArrayCoordsYUp(x, y, width, height, 1)
							coord := base * 4
							zCoord := base

							interpolatedZ := w0*sv1.Z + w1*sv2.Z + w2*sv3.Z
							if interpolatedZ >= zbuffer[zCoord] {

								zbuffer[zCoord] = interpolatedZ

								// inside pixel loop, interpolate and recover
								interpR := w0*r0 + w1*r1 + w2*r2
								interpG := w0*g0 + w1*g1 + w2*g2
								interpB := w0*b0 + w1*b1 + w2*b2
								interpW := w0*oneOverW0 + w1*oneOverW1 + w2*oneOverW2
								invInterpW := 1 / interpW

								// lighting
								// interpL := w0*l0OverW + w1*l1OverW + w2*l2OverW
								finalLightR := (w0*l0OverW.X + w1*l1OverW.X + w2*l2OverW.X) * invInterpW
								finalLightG := (w0*l0OverW.Y + w1*l1OverW.Y + w2*l2OverW.Y) * invInterpW
								finalLightB := (w0*l0OverW.Z + w1*l1OverW.Z + w2*l2OverW.Z) * invInterpW

								finalU := (w0*u0OverW + w1*u1OverW + w2*u2OverW) * invInterpW
								finalV := (w0*v0OverW + w1*v1OverW + w2*v2OverW) * invInterpW

								texR, texG, texB, _ := SampleTexture(textureImg, finalU, finalV)

								// 0.003921569 = 1/255
								finalR := byte(min(float32(255), float32(texR)*(float32(interpR*invInterpW)*0.003921569)*finalLightR))
								finalG := byte(min(float32(255), float32(texG)*(float32(interpG*invInterpW)*0.003921569)*finalLightG))
								finalB := byte(min(float32(255), float32(texB)*(float32(interpB*invInterpW)*0.003921569)*finalLightB))

								pixels[coord] = byte(finalR)
								pixels[coord+1] = byte(finalG)
								pixels[coord+2] = byte(finalB)
								pixels[coord+3] = 255
							}

						}
					}
				}
			}
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
		renderer.DebugText(10, 20, fmt.Sprintf("Position: %.0f, %0.f, %0.f", camera.Position.X, camera.Position.Y, camera.Position.Z))
		renderer.Present()

		// rotationY += 1.0 * deltaTime
		for i := range zbuffer {
			zbuffer[i] = float32(math.Inf(-1))
		}

		return nil
	})
}
