package main

import "math"

type Tex2 struct {
	U float32
	V float32
}

type Triangle struct {
	points     [3]Vec4
	textcoords [3]Tex2
	colors     [3]Vec4
	lights     [3]Vec3
	avgDepth   float32
}

func edge(p, a, b Vec3) float32 {
	return (p.X-a.X)*(b.Y-a.Y) - (p.Y-a.Y)*(b.X-a.X)
}

func IsPixelInTriangle(p Vec3, v1, v2, v3 Vec3) (bool, float32, float32, float32) {

	area := edge(v1, v2, v3)

	w0 := edge(p, v2, v3) / area
	w1 := edge(p, v3, v1) / area
	w2 := edge(p, v1, v2) / area

	return w0 >= 0 && w1 >= 0 && w2 >= 0, w0, w1, w2

}

func ColorFromWeights(w0, w1, w2 float32, c0, c1, c2 Vec4) (byte, byte, byte, byte) {
	r := byte(w0*c0.X + w1*c1.X + w2*c2.X)
	g := byte(w0*c0.Y + w1*c1.Y + w2*c2.Y)
	b := byte(w0*c0.Z + w1*c1.Z + w2*c2.Z)
	a := byte(w0*c0.W + w1*c1.W + w2*c2.W)

	return r, g, b, a
}

func GetRectBounds(a, b, c Vec3) (minX, maxX, minY, maxY int) {

	minX = int(math.Floor(math.Min(math.Min(float64(a.X), float64(b.X)), float64(c.X))))
	maxX = int(math.Ceil(math.Max(math.Max(float64(a.X), float64(b.X)), float64(c.X))))

	minY = int(math.Floor(math.Min(math.Min(float64(a.Y), float64(b.Y)), float64(c.Y))))
	maxY = int(math.Ceil(math.Max(math.Max(float64(a.Y), float64(b.Y)), float64(c.Y))))

	return minX, maxX, minY, maxY

}
