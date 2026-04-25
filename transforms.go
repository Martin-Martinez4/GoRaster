package main

import "math"

func RotationAlongY(radians float32) *Matrix4 {

	rads64 := float64(radians)

	return &Matrix4{
		float32(math.Cos(rads64)), 0, float32(math.Sin(rads64)), 0,
		0, 1, 0, 0,
		float32(-math.Sin(rads64)), 0, float32(math.Cos(rads64)), 0,
		0, 0, 0, 1,
	}
}

func Translate(tx, ty, tz float32) *Matrix4 {
	m := IdentityMatrix()

	m[12] = tx
	m[13] = ty
	m[14] = tz

	return &m
}
