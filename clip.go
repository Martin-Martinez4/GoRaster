package main

import (
	"math"
)

const MAX_NUM_POLY_VERTICES = 10
const MAX_NUM_POLY_TRIANGLES = 10

const NUM_PLANES = 6

type PlaneDir int

const (
	LEFT_FRUSTUM_PLANE PlaneDir = iota
	RIGHT_FRUSTUM_PLANE
	TOP_FRUSTUM_PLANE
	BOTTOM_FRUSTUM_PLANE
	NEAR_FRUSTUM_PLANE
	FAR_FRUSTUM_PLANE
)

var frustumPlanes [NUM_PLANES]plane

type plane struct {
	Point  Vec3
	Normal Vec3
}

type Polygon struct {
	Vertices    [MAX_NUM_POLY_VERTICES]Vec4
	TextCoords  [MAX_NUM_POLY_VERTICES]Tex2
	Colors      [MAX_NUM_POLY_VERTICES]Vec4
	NumVertices int
}

func float_lerp(a, b, t float32) float32 {
	return a + t*(b-a)
}

func initFrustumPlanes(fovy, fovx, zNear, zFar float32) {
	cosHalfovX := math.Cos(float64(fovx / 2))
	sinHalfovX := math.Sin(float64(fovx / 2))

	cosHalfovY := math.Cos(float64(fovy / 2))
	sinHalfovY := math.Sin(float64(fovy / 2))

	frustumPlanes[LEFT_FRUSTUM_PLANE].Point = Vec3{0, 0, 0}
	frustumPlanes[LEFT_FRUSTUM_PLANE].Normal = Vec3{
		float32(cosHalfovX), 0, float32(-sinHalfovX),
	}

	frustumPlanes[RIGHT_FRUSTUM_PLANE].Point = Vec3{0, 0, 0}
	frustumPlanes[RIGHT_FRUSTUM_PLANE].Normal = Vec3{
		float32(-cosHalfovX), 0, float32(-sinHalfovX),
	}

	frustumPlanes[TOP_FRUSTUM_PLANE].Point = Vec3{0, 0, 0}
	frustumPlanes[TOP_FRUSTUM_PLANE].Normal = Vec3{
		0, float32(-cosHalfovY), float32(-sinHalfovY),
	}

	frustumPlanes[BOTTOM_FRUSTUM_PLANE].Point = Vec3{0, 0, 0}
	frustumPlanes[BOTTOM_FRUSTUM_PLANE].Normal = Vec3{
		0, float32(cosHalfovY), float32(-sinHalfovY),
	}

	// Near: point at -zNear, normal pointing in -Z
	frustumPlanes[NEAR_FRUSTUM_PLANE].Point = Vec3{0, 0, -zNear}
	frustumPlanes[NEAR_FRUSTUM_PLANE].Normal = Vec3{0, 0, -1}

	// Far: point at -zFar, normal pointing in +Z
	frustumPlanes[FAR_FRUSTUM_PLANE].Point = Vec3{0, 0, -zFar}
	frustumPlanes[FAR_FRUSTUM_PLANE].Normal = Vec3{0, 0, 1}
}

func TriangleFromPolygon(polygon *Polygon, triangles []Triangle, numTriangles *int) {

	for i := 0; i < polygon.NumVertices-2; i++ {
		index0 := 0
		index1 := i + 1
		index2 := i + 2

		triangles[i].points[0] = Vec4{
			polygon.Vertices[index0].X,
			polygon.Vertices[index0].Y,
			polygon.Vertices[index0].Z,
			polygon.Vertices[index0].W, // carry W through
		}
		triangles[i].points[1] = Vec4{
			polygon.Vertices[index1].X,
			polygon.Vertices[index1].Y,
			polygon.Vertices[index1].Z,
			polygon.Vertices[index1].W,
		}
		triangles[i].points[2] = Vec4{
			polygon.Vertices[index2].X,
			polygon.Vertices[index2].Y,
			polygon.Vertices[index2].Z,
			polygon.Vertices[index2].W,
		}

		triangles[i].colors[0] = polygon.Colors[index0]
		triangles[i].colors[1] = polygon.Colors[index1]
		triangles[i].colors[2] = polygon.Colors[index2]

	}
	*numTriangles = polygon.NumVertices - 2

}

func CreatePolygonFromTriangle(v0, v1, v2 Vec4, t0, t1, t2 Tex2, c0, c1, c2 Vec4) Polygon {
	return Polygon{
		Vertices:    [10]Vec4{v0, v1, v2},
		TextCoords:  [10]Tex2{t0, t1, t2},
		Colors:      [10]Vec4{c0, c1, c2},
		NumVertices: 3,
	}
}

func ClipPolygonAgainstPlane(polygon *Polygon, plane PlaneDir) {

	var insideVertices [MAX_NUM_POLY_VERTICES]Vec4
	var insideTexCoords [MAX_NUM_POLY_VERTICES]Tex2
	var insideColors [MAX_NUM_POLY_VERTICES]Vec4
	numInsideVertices := 0

	for current := 0; current < polygon.NumVertices; current++ {
		currentVertex := polygon.Vertices[current]
		currentTexCoord := polygon.TextCoords[current]
		currentColor := polygon.Colors[current]

		previousVertex := polygon.Vertices[(current+polygon.NumVertices-1)%polygon.NumVertices]
		previousTexCoord := polygon.TextCoords[(current+polygon.NumVertices-1)%polygon.NumVertices]
		previousColor := polygon.Colors[(current+polygon.NumVertices-1)%polygon.NumVertices]

		// replaces the manual dot product against frustum plane normal
		currentDot := dotClipPlane(currentVertex, plane)
		previousDot := dotClipPlane(previousVertex, plane)

		if currentDot*previousDot < 0 {
			t := previousDot / (previousDot - currentDot)

			insideVertices[numInsideVertices] = Vec4{
				float_lerp(previousVertex.X, currentVertex.X, t),
				float_lerp(previousVertex.Y, currentVertex.Y, t),
				float_lerp(previousVertex.Z, currentVertex.Z, t),
				float_lerp(previousVertex.W, currentVertex.W, t), // W interpolated correctly
			}
			insideTexCoords[numInsideVertices] = Tex2{
				float_lerp(previousTexCoord.U, currentTexCoord.U, t),
				float_lerp(previousTexCoord.V, currentTexCoord.V, t),
			}
			insideColors[numInsideVertices] = Vec4{
				float_lerp(previousColor.X, currentColor.X, t),
				float_lerp(previousColor.Y, currentColor.Y, t),
				float_lerp(previousColor.Z, currentColor.Z, t),
				float_lerp(previousColor.W, currentColor.W, t),
			}
			numInsideVertices++
		}

		if currentDot > 0 {
			insideVertices[numInsideVertices] = currentVertex
			insideTexCoords[numInsideVertices] = currentTexCoord
			insideColors[numInsideVertices] = currentColor
			numInsideVertices++
		}
	}

	for i := 0; i < numInsideVertices; i++ {
		polygon.Vertices[i] = insideVertices[i]
		polygon.TextCoords[i] = insideTexCoords[i]
		polygon.Colors[i] = insideColors[i]
	}
	polygon.NumVertices = numInsideVertices
}

func ClipPolygon(polygon *Polygon) {

	ClipPolygonAgainstPlane(polygon, LEFT_FRUSTUM_PLANE)
	ClipPolygonAgainstPlane(polygon, RIGHT_FRUSTUM_PLANE)
	ClipPolygonAgainstPlane(polygon, TOP_FRUSTUM_PLANE)
	ClipPolygonAgainstPlane(polygon, BOTTOM_FRUSTUM_PLANE)
	ClipPolygonAgainstPlane(polygon, NEAR_FRUSTUM_PLANE)
	ClipPolygonAgainstPlane(polygon, FAR_FRUSTUM_PLANE)
}

func dotClipPlane(v Vec4, plane PlaneDir) float32 {
	switch plane {
	case LEFT_FRUSTUM_PLANE:
		return v.X + v.W
	case RIGHT_FRUSTUM_PLANE:
		return v.W - v.X
	case BOTTOM_FRUSTUM_PLANE:
		return v.Y + v.W
	case TOP_FRUSTUM_PLANE:
		return v.W - v.Y
	case NEAR_FRUSTUM_PLANE:
		return v.Z + v.W
	case FAR_FRUSTUM_PLANE:
		return v.W - v.Z
	}
	return 0
}
