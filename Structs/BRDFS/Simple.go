package BRDFS

import (
	"Photon/Math"
	"math"
)

type SimpleBRDF struct {
}

func (s SimpleBRDF) Sample(view, indescent, normal, lightColor, albedo Math.Vector3, lightIntensity, roughness, metallic, ior float64) Math.Vector3 {
	// Lambertian diffuse
	diffuse := math.Pow(math.Max(normal.Dot(indescent), 0), 2)
	specRough := roughness*roughness*100 + 1
	specular := math.Pow(math.Max(view.Dot(indescent.Reflect(normal)), 0), specRough)
	met := albedo.FMul(diffuse).Mul(lightColor.FMul(specular * lightIntensity))
	nonMet := albedo.FMul(diffuse).Add(lightColor.FMul(specular * lightIntensity)).FDiv(2)
	return met.FMul(metallic).Add(nonMet.FMul(1 - metallic))
}
