package main

func cullBackFace(vs0, vs1, vs2 Vec4) bool {

	edge1 := Vec3{vs1.X - vs0.X, vs1.Y - vs0.Y, vs1.Z - vs0.Z}
	edge2 := Vec3{vs2.X - vs0.X, vs2.Y - vs0.Y, vs2.Z - vs0.Z}
	normal := edge1.Cross(edge2)

	// in view space camera is at origin so camera ray is just -vertex
	cameraRay := Vec3{-vs0.X, -vs0.Y, -vs0.Z}

	return normal.Dot(cameraRay) <= 0
}
