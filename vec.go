package main

type Vec2 struct {
	X, Y float32
}
type Vec3 struct {
	X, Y, Z float32
}
type Vec4 struct {
	X, Y, Z, W float32
}

func (a Vec3) Dot(b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vec3) Sub(b Vec3) Vec3 {
	return Vec3{a.X - b.X, a.Y - b.Y, 0}
}

func (v Vec4) MultMatrix4(m Matrix4) Vec4 {
	return Vec4{
		X: v.X*m.At(0, 0) + v.Y*m.At(1, 0) + v.Z*m.At(2, 0) + v.W*m.At(3, 0),
		Y: v.X*m.At(0, 1) + v.Y*m.At(1, 1) + v.Z*m.At(2, 1) + v.W*m.At(3, 1),
		Z: v.X*m.At(0, 2) + v.Y*m.At(1, 2) + v.Z*m.At(2, 2) + v.W*m.At(3, 2),
		W: v.X*m.At(0, 3) + v.Y*m.At(1, 3) + v.Z*m.At(2, 3) + v.W*m.At(3, 3),
	}
}

func vec3Clone(v *Vec3) Vec3 {
	result := Vec3{v.X, v.Y, v.Z}
	return result
}

func tex2Clone(t *Tex2) Tex2 {
	result := Tex2{t.U, t.V}
	return result
}

func vec3FromVec4(v *Vec4) Vec3 {
	return Vec3{v.X, v.X, v.Z}
}
