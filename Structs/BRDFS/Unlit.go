package BRDFS

import "Photon/Math"

type UnlitBRDF struct {
}

func (u UnlitBRDF) Sample(view, indescent, normal, lightColor, albedo Math.Vector3, lightIntensity, roughness, metallic, ior float64) Math.Vector3 {
	return albedo
}
