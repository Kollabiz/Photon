package Structs

import (
	"Photon/Math"
	"math"
)

type Material struct {
	// Albedo
	albedoTexture     *TextureRGB
	albedoTextureUsed bool
	albedoColor       Math.Vector3
	// Roughness
	roughnessTexture     *TextureGrayscale
	roughnessTextureUsed bool
	roughness            float64
	// Metallic
	metallicTexture     *TextureGrayscale
	metallicTextureUsed bool
	metallic            float64
	// IOR
	ior float64
}

func (material *Material) sampleTextures(uv Math.Vector2) (Math.Vector3, float64, float64) {
	var albedo Math.Vector3
	if material.albedoTextureUsed {
		albedo = material.albedoTexture.At(uv)
	} else {
		albedo = material.albedoColor
	}
	var roughness float64
	if material.roughnessTextureUsed {
		roughness = material.roughnessTexture.At(uv)
	} else {
		roughness = material.roughness
	}
	var metallic float64
	if material.metallicTextureUsed {
		metallic = material.metallicTexture.At(uv)
	} else {
		metallic = material.metallic
	}
	return albedo, roughness, metallic
}

func (material *Material) SampleLight(uv Math.Vector2, v, l, n Math.Vector3, li float64, lc Math.Vector3) Math.Vector3 {
	albedo, roughness, metallic := material.sampleTextures(uv)

	// Cook-Torrance BRDF

	h := v.Add(l).FDiv(2)

	// Geometric attenuation
	g := math.Min(math.Min(1, 2*h.Dot(n)*v.Dot(n)/v.Dot(h)), 2*h.Dot(n)*l.Dot(n)/v.Dot(h))

	// Beckmann distribution
	alpha := math.Acos(n.Dot(h))
	cosAlpha := math.Pow(math.Cos(alpha), 2)
	tanAlpha := (cosAlpha - 1) / cosAlpha
	roughness2 := roughness * roughness
	denominator := math.Pi * roughness2 * cosAlpha * cosAlpha
	d := math.Exp(-tanAlpha*tanAlpha/roughness2) / denominator

	// Fresnel
	r0 := math.Pow((1-material.ior)/(1+material.ior), 2)
	f := r0 + (1-r0)*(1-n.Dot(v))

	specularDistribution := d * f * g / (4 * v.Dot(n) * n.Dot(l))
	specular := lc.FMul(li * specularDistribution)

	metal := specular.Mul(albedo)
	glossy := specular.Add(albedo).FDiv(2)
	return Math.InterpolateVector3(glossy, metal, metallic)
}

func (material *Material) SampleSimplifiedLight(uv Math.Vector2, l, n Math.Vector3, li float64, lc Math.Vector3) Math.Vector3 {
	lnDot := ((l.Dot(n)+1)/2 + 0.2) / 1.2
	albedo, _, _ := material.sampleTextures(uv)
	return albedo.Mul(lc.FMul(lnDot * li))
}
