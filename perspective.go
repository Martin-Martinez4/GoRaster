package main

import "math"

func Perspective(near, far, fovY, aspect float32) Matrix4 {

	// top := float32(math.Tan(float64(fovY)/2.0)) * near
	// right := top * aspect

	invTan := 1 / math.Tan(float64(fovY/2))

	return Matrix4{
		// col 0
		aspect * float32(invTan), 0, 0, 0,
		// col 1
		0, float32(invTan), 0, 0,
		// col 2
		0, 0, -(far + near) / (far - near), -1,
		// col 3
		0, 0, -(2 * far * near) / (far - near), 0,
	}
}

func PerspectiveDivide(v Vec4) Vec4 {
	return Vec4{X: v.X / v.W, Y: v.Y / v.W, Z: v.Z / v.W}
}
