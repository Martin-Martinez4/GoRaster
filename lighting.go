package main

func CalculateSimpleLighting(vertNormals, lightValues []Vec3, lightColor, ambient Vec3, model Matrix4) {

	for i := range vertNormals {
		// normals only need rotation, not translation
		// use model matrix but ignore the translation component
		worldNormal := model.MultVec4(Vec4{vertNormals[i].X, vertNormals[i].Y, vertNormals[i].Z, 0})
		// W=0 means translation is ignored

		lightDir := Vec3{0.5, 1, 0.5}.Normalize()
		intensity := max(float32(0), Vec3{worldNormal.X, worldNormal.Y, worldNormal.Z}.Dot(lightDir))

		// 0.2 ambient
		// 0.8 lightColor
		lightValues[i].X = ambient.X + intensity*0.8*lightColor.X
		lightValues[i].Y = ambient.Y + intensity*0.8*lightColor.Y
		lightValues[i].Z = ambient.Z + intensity*0.8*lightColor.Z
	}
}

func GetLightOverW(l0, l1, l2 Vec3, sv1W, sv2W, sv3W float32) (Vec3, Vec3, Vec3) {
	return Vec3{l0.X / sv1W, l0.Y / sv1W, l0.Z / sv1W},
		Vec3{l1.X / sv2W, l1.Y / sv2W, l1.Z / sv2W},
		Vec3{l2.X / sv3W, l2.Y / sv3W, l2.Z / sv3W}
}
