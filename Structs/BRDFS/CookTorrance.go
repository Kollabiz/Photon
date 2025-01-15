package BRDFS

import (
	"Photon/Math"
	"math"
)

type CookTorranceBRDF struct {
	// DEBUG STUFF

	SpecularScale float64
}

func (brdf CookTorranceBRDF) Sample(view, indescent, normal, lightColor, albedo Math.Vector3, lightIntensity, roughness, metallic, ior float64) Math.Vector3 {
	h := indescent.Add(normal).FDiv(2)

	// Geometric attenuation
	G := min(1, 2*h.Dot(normal)*view.Dot(normal)/h.Dot(view), 2*h.Dot(normal)*indescent.Dot(normal)/h.Dot(view))

	// Beckmann distribution
	alpha := math.Acos(normal.Dot(h))
	cosAlpha := math.Pow(math.Cos(alpha), 2)
	tanAlpha := (cosAlpha - 1) / cosAlpha
	roughness2 := roughness * roughness
	denominator := math.Pi * roughness2 * cosAlpha * cosAlpha
	D := math.Exp(-(tanAlpha*tanAlpha)/roughness2) / denominator

	// Fresnel
	r0 := math.Pow((1-ior)/(1+ior), 2)
	F := r0 + (1-r0)*(1-normal.Dot(view))

	specular := D * F * G / (4 * view.Dot(normal) * indescent.Dot(normal)) * brdf.SpecularScale

	metal := albedo.Mul(lightColor.FMul(lightIntensity).FMul(specular))
	glossy := albedo.Add(lightColor.FMul(lightIntensity).FMul(specular)).FDiv(2)

	return Math.InterpolateVector3(glossy, metal, metallic)
}
