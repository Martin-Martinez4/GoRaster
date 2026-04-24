package main

type Vec2 struct {
	X, Y float32
}
type Vec3 struct {
	X, Y, Z float32
}
type Vec4 struct {
	X, Y, Z, A float32
}

func (a Vec3) Dot(b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y
}

func (a Vec3) Sub(b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, 0}
}
