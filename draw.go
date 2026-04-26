package main

import "math"

func drawTriangle(v0, v1, v2 Vec3, width, height int, color Vec4, pixelsBuf []byte, zBuf []float32,
) {
	drawLineZ(v0, v1, width, height, color, pixelsBuf, zBuf)
	drawLineZ(v1, v2, width, height, color, pixelsBuf, zBuf)
	drawLineZ(v2, v0, width, height, color, pixelsBuf, zBuf)
}

func drawLine(x0, y0, x1, y1 int, width, height int, color Vec4, pixelsBuf []byte) {
	deltaX := x1 - x0
	deltaY := y1 - y0

	absDeltaX := math.Abs(float64(deltaX))
	absDeltaY := math.Abs(float64(deltaY))

	longestSide := int(absDeltaY)
	if absDeltaX >= absDeltaY {
		longestSide = int(absDeltaX)
	}

	xInc := float32(deltaX) / float32(longestSide)
	yInc := float32(deltaY) / float32(longestSide)

	currentX := float32(x0)
	currentY := float32(y0)

	for i := 0; i <= longestSide; i++ {
		drawPixel(int(math.Round(float64(currentX))), int(math.Round(float64(currentY))), width, height, color, pixelsBuf)
		currentX += xInc
		currentY += yInc
	}

}

func drawLineZ(
	v0, v1 Vec3,
	width, height int,
	color Vec4,
	pixelsBuf []byte,
	zBuf []float32,
) {
	dx := v1.X - v0.X
	dy := v1.Y - v0.Y

	steps := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))
	if steps == 0 {
		return
	}

	xInc := dx / float32(steps)
	yInc := dy / float32(steps)
	zInc := (v1.Z - v0.Z) / float32(steps)

	x := v0.X
	y := v0.Y
	z := v0.Z

	for i := 0; i <= steps; i++ {
		xi := int(math.Round(float64(x)))
		yi := int(math.Round(float64(y)))

		if xi >= 0 && xi < width && yi >= 0 && yi < height {
			idx := yi*width + xi

			// depth test
			if z <= zBuf[idx] {
				zBuf[idx] = z

				coord := idx * 4
				pixelsBuf[coord] = byte(color.X)
				pixelsBuf[coord+1] = byte(color.Y)
				pixelsBuf[coord+2] = byte(color.Z)
				pixelsBuf[coord+3] = byte(color.W)
			}
		}

		x += xInc
		y += yInc
		z += zInc
	}
}

func drawPixel(x, y, width, height int, color Vec4, pixelsBuf []byte) {

	coord := ToArrayCoordsYUp(x, y, width, height, 4)

	pixelsBuf[coord] = byte(color.X)
	pixelsBuf[coord+1] = byte(color.Y)
	pixelsBuf[coord+2] = byte(color.Z)
	pixelsBuf[coord+3] = byte(color.W)

}
