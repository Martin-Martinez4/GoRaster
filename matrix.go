package main

type Matrix4 [16]float32

func (m Matrix4) At(r, c int) float32 {
	return m[r+c*4]
}
func (m Matrix4) MultVec4(v Vec4) Vec4 {
	return Vec4{
		m[0]*v.X + m[4]*v.Y + m[8]*v.Z + m[12]*v.W,
		m[1]*v.X + m[5]*v.Y + m[9]*v.Z + m[13]*v.W,
		m[2]*v.X + m[6]*v.Y + m[10]*v.Z + m[14]*v.W,
		m[3]*v.X + m[7]*v.Y + m[11]*v.Z + m[15]*v.W,
	}
}

func MulMatrix4(out *Matrix4, a, b *Matrix4) {
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			out[row+col*4] =
				a[row+0]*b[0+col*4] +
					a[row+4]*b[1+col*4] +
					a[row+8]*b[2+col*4] +
					a[row+12]*b[3+col*4]
		}
	}
}

func ViewMatrix(camPos Vec3) Matrix4 {
	m := IdentityMatrix()

	m[12] = -camPos.X
	m[13] = -camPos.Y
	m[14] = -camPos.Z

	return m
}

func IdentityMatrix() Matrix4 {
	return Matrix4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}
