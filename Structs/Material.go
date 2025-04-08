package Structs

import (
	"Photon/Math"
)

// BRDF function interface
// Kind of crutch-y, but fast (at least I hope it is)

type IBRDF interface {
	Sample(view, indescent, normal, lightColor, albedo Math.Vector3, lightIntensity, roughness, metallic, ior float64) Math.Vector3
	// View, Light, Normal, light color, light intensity, albedo, roughness, metallic, ior -> color
}

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
	// BRDF function
	BRDF IBRDF
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

func (material *Material) SampleAlbedo(uv Math.Vector2) Math.Vector3 {
	if material.albedoTextureUsed {
		return material.albedoTexture.At(uv)
	} else {
		return material.albedoColor
	}
}

func (material *Material) SampleLight(uv Math.Vector2, v, l, n Math.Vector3, li float64, lc Math.Vector3) Math.Vector3 {
	albedo, roughness, metallic := material.sampleTextures(uv)
	return material.BRDF.Sample(v, l, n, lc, albedo, li, roughness, metallic, material.ior)
}

func (material *Material) SampleSimplifiedLight(uv Math.Vector2, l, n Math.Vector3, li float64, lc Math.Vector3) Math.Vector3 {
	lnDot := ((l.Dot(n)+1)/2 + 0.2) / 1.2
	albedo, _, _ := material.sampleTextures(uv)
	return albedo.Mul(lc.FMul(lnDot * li))
}

// Setters for MTL reader

func (material *Material) SetMetallic(ka float64) {
	material.metallic = ka
	material.metallicTextureUsed = false
}

func (material *Material) SetMetallicTexture(mapKa *TextureGrayscale) {
	material.metallicTexture = mapKa
	material.metallicTextureUsed = true
}

func (material *Material) SetAlbedo(kd Math.Vector3) {
	material.albedoColor = kd
	material.albedoTextureUsed = false
}

func (material *Material) SetAlbedoTexture(mapKd *TextureRGB) {
	material.albedoTexture = mapKd
	material.albedoTextureUsed = true
}

func (material *Material) SetRoughness(ks float64) {
	material.roughness = ks
	material.roughnessTextureUsed = false
}

func (material *Material) SetRoughnessTexture(mapKs *TextureGrayscale) {
	material.roughnessTexture = mapKs
	material.roughnessTextureUsed = true
}

func (material *Material) SetIOR(ior float64) {
	material.ior = ior
}

func (material *Material) GetRoughness(uv Math.Vector2) float64 {
	if material.roughnessTextureUsed {
		return material.roughnessTexture.At(uv)
	} else {
		return material.roughness
	}
}
