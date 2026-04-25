package main

type Vertex struct {
	Pos    Vec3
	Color  *Vec4
	UV     *Vec2
	Normal *Vec3
}

type ScreenVertex struct {
	// screen position
	Pos Vec3
	// X, Y float32
	// // depth (for z-buffer)
	// Z float32
	// clip-space w
	W float32
}
