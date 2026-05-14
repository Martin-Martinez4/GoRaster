package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"

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

	tileSize := 16
	binCols := 1280 / tileSize
	binRows := 720 / tileSize
	rasterTris := make([]Bin, binCols*binRows)

	// start workers
	// limit to 4 just because for now
	numWorkers := runtime.GOMAXPROCS(0)

	var nextBin atomic.Uint32

	for by := 0; by < binRows; by++ {
		for bx := 0; bx < binCols; bx++ {

			idx := ToArrayCoordsYUp(bx, by, binCols, binRows, 1)

			rasterTris[idx].MinX = bx * tileSize
			rasterTris[idx].MinY = by * tileSize

			rasterTris[idx].MaxX = min(width-1, (bx+1)*tileSize-1)
			rasterTris[idx].MaxY = min(height-1, (by+1)*tileSize-1)
		}
	}

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

	// var lastTime uint64
	// var fps float32

	// frameCountFrequency := 10
	// var frameCount int
	// var fpsAcc float32

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

	targetFPS := 60.0

	targetFrameMS := 1000.0 / targetFPS

	sdl.RunLoop(func() error {
		frameStart := sdl.Ticks()
		// currentTime := frameStart
		// deltaTime := float32(currentTime-lastTime) / 1000.0
		// lastTime = currentTime

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

		for i := range rasterTris {
			rasterTris[i].Tris = rasterTris[i].Tris[:0]
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
				invSv1W := 1 / sv1.W
				invSv2W := 1 / sv2.W
				invSv3W := 1 / sv3.W

				r0, g0, b0, r1, g1, b1, r2, g2, b2 := ClippedColorComponents(
					trianglesAfterClipping[j].colors[0],
					trianglesAfterClipping[j].colors[1],
					trianglesAfterClipping[j].colors[2],
					invSv1W,
					invSv2W,
					invSv3W,
				)

				u0OverW, v0OverW, u1OverW, v1OverW, u2OverW, v2OverW := ClippedUVComponents(
					trianglesAfterClipping[j].textcoords[0],
					trianglesAfterClipping[j].textcoords[1],
					trianglesAfterClipping[j].textcoords[2],
					invSv1W,
					invSv2W,
					invSv3W,
				)

				// perspective correct prep for light, same as color
				l0OverW, l1OverW, l2OverW := GetLightOverW(
					trianglesAfterClipping[j].lights[0],
					trianglesAfterClipping[j].lights[1],
					trianglesAfterClipping[j].lights[2],
					sv1.W,
					sv2.W,
					sv3.W,
				)

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

				minX = max(0, minX/tileSize)
				minY = max(0, minY/tileSize)
				maxX = min((width-1)/tileSize, maxX/tileSize)
				maxY = min((height-1)/tileSize, maxY/tileSize)

				oneOverW0 := 1.0 / sv1.W
				oneOverW1 := 1.0 / sv2.W
				oneOverW2 := 1.0 / sv3.W

				// area := edge(vec31, vec32, vec33)
				// invArea := 1.0 / area

				// add tris to bins

				for y := minY; y <= maxY; y++ {
					for x := minX; x <= maxX; x++ {
						base := ToArrayCoordsYUp(x, y, binCols, binRows, 1)

						rasterTris[base].Tris = append(rasterTris[base].Tris, RasterTri{
							Points:   [3]Vec4{sv1, sv2, sv3},
							UVs:      [6]float32{u0OverW, v0OverW, u1OverW, v1OverW, u2OverW, v2OverW},
							Colors:   [9]float32{r0, g0, b0, r1, g1, b1, r2, g2, b2},
							Lights:   [9]float32{l0OverW.X, l0OverW.Y, l0OverW.Z, l1OverW.X, l1OverW.Y, l1OverW.Z, l2OverW.X, l2OverW.Y, l2OverW.Z},
							OneOverW: [3]float32{oneOverW0, oneOverW1, oneOverW2},
							InvArea:  1.0 / edge(vec31, vec32, vec33),
						})

					}
				}

				// screen space stuff

			}
		}

		var wg sync.WaitGroup

		wg.Add(numWorkers)

		for w := 0; w < numWorkers; w++ {
			go func() {
				defer wg.Done()

				for {
					idx := int(nextBin.Add(1)) - 1

					if idx >= len(rasterTris) {
						return
					}

					for _, tri := range rasterTris[idx].Tris {

						vec31 := Vec3{
							X: tri.Points[0].X,
							Y: tri.Points[0].Y,
							Z: tri.Points[0].Z,
						}

						vec32 := Vec3{
							X: tri.Points[1].X,
							Y: tri.Points[1].Y,
							Z: tri.Points[1].Z,
						}

						vec33 := Vec3{
							X: tri.Points[2].X,
							Y: tri.Points[2].Y,
							Z: tri.Points[2].Z,
						}

						for y := rasterTris[idx].MinY; y <= rasterTris[idx].MaxY; y++ {
							for x := rasterTris[idx].MinX; x <= rasterTris[idx].MaxX; x++ {

								// Create vec from center of pixel
								inTri, w0, w1, w2 := IsPixelInTriangle(Vec3{float32(x) + 0.5, float32(y) + 0.5, 0}, vec31, vec32, vec33, tri.InvArea)
								if inTri {

									base := ToArrayCoordsYUp(x, y, width, height, 1)
									coord := base * 4
									zCoord := base

									interpolatedZ := w0*tri.Points[0].Z + w1*tri.Points[1].Z + w2*tri.Points[2].Z
									if interpolatedZ >= zbuffer[zCoord] {

										zbuffer[zCoord] = interpolatedZ

										// inside pixel loop, interpolate and recover
										interpR := w0*tri.Colors[0] + w1*tri.Colors[3] + w2*tri.Colors[6]
										interpG := w0*tri.Colors[1] + w1*tri.Colors[4] + w2*tri.Colors[7]
										interpB := w0*tri.Colors[2] + w1*tri.Colors[5] + w2*tri.Colors[8]
										interpW := w0*tri.OneOverW[0] + w1*tri.OneOverW[1] + w2*tri.OneOverW[2]
										invInterpW := 1 / interpW

										// lighting
										// interpL := w0*l0OverW + w1*l1OverW + w2*l2OverW
										finalLightR := (w0*tri.Lights[0] + w1*tri.Lights[3] + w2*tri.Lights[6]) * invInterpW
										finalLightG := (w0*tri.Lights[1] + w1*tri.Lights[4] + w2*tri.Lights[7]) * invInterpW
										finalLightB := (w0*tri.Lights[2] + w1*tri.Lights[5] + w2*tri.Lights[8]) * invInterpW

										finalU := (w0*tri.UVs[0] + w1*tri.UVs[2] + w2*tri.UVs[4]) * invInterpW
										finalV := (w0*tri.UVs[1] + w1*tri.UVs[3] + w2*tri.UVs[5]) * invInterpW

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
			}()
		}

		wg.Wait()
		nextBin.Store(0)

		texture.Update(nil, pixels, int32(width*4))

		renderer.Clear()
		renderer.RenderTexture(texture, nil, nil)

		// fps = fpsAcc / float32(frameCount)

		frameTime := sdl.Ticks() - frameStart

		if frameTime < uint64(targetFrameMS) {
			sdl.Delay(uint32(targetFrameMS) - uint32(frameTime))
		}
		renderer.DebugText(10, 10, fmt.Sprintf("FPS: %.0f", 1000.0/float32(frameTime)))
		renderer.DebugText(10, 20, fmt.Sprintf("Position: %.0f, %0.f, %0.f", camera.Position.X, camera.Position.Y, camera.Position.Z))
		renderer.Present()

		// rotationY += 1.0 * deltaTime
		for i := range zbuffer {
			zbuffer[i] = float32(math.Inf(-1))
		}

		return nil
	})
}
