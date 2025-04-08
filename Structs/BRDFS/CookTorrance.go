package BRDFS

import (
	"Photon/Math"
	"math"
)

type CookTorranceBRDF struct {
	// DEBUG STUFF

	SpecularScale float64
}

func NewCookTorranceBRDF() *CookTorranceBRDF {
	return &CookTorranceBRDF{1}
}

func safeDivide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func clampDot(a, b Math.Vector3) float64 {
	return math.Max(math.Min(a.Dot(b), 1), 0)
}

func clampLight(light Math.Vector3) Math.Vector3 {
	if math.IsNaN(light.X) {
		light.X = 0
	}
	if math.IsNaN(light.Y) {
		light.Y = 0
	}
	if math.IsNaN(light.Z) {
		light.Z = 0
	}
	return Math.Vector3{
		X: math.Sqrt(math.Max(math.Min(light.X, 1), 0)),
		Y: math.Sqrt(math.Max(math.Min(light.Y, 1), 0)),
		Z: math.Sqrt(math.Max(math.Min(light.Z, 1), 0)),
	}
}

func (brdf CookTorranceBRDF) Sample(view, indescent, normal, lightColor, albedo Math.Vector3, lightIntensity, roughness, metallic, ior float64) Math.Vector3 {
	h := indescent.Add(normal).Normalized()
	lColor := clampLight(lightColor.FMul(lightIntensity))

	// Geometric attenuation
	vDotH := clampDot(h, view)
	hDotN := clampDot(h, normal)
	vDotN := clampDot(view, normal)
	iDotN := clampDot(indescent, normal)
	G := math.Min(
		math.Min(
			1,
			2*hDotN*vDotN/vDotH,
		),
		2*hDotN*math.Abs(iDotN)/vDotH,
	)

	// Beckmann distribution
	alpha := math.Acos(hDotN)
	cosAlpha := math.Pow(math.Cos(alpha), 2)
	tanAlpha := safeDivide(cosAlpha-1, cosAlpha)
	roughness2 := roughness * roughness
	denominator := math.Pi * roughness2 * cosAlpha * cosAlpha
	D := safeDivide(math.Exp(safeDivide(-(tanAlpha*tanAlpha), roughness2)), denominator)

	// Fresnel
	r0 := math.Pow(safeDivide(1-ior, 1+ior), 2)
	F := r0 + (1-r0)*(1-vDotN)

	specular := safeDivide(D*F*G, (4*vDotN*iDotN)*brdf.SpecularScale)

	diffuse := albedo.FMul(F)
	metal := albedo.Mul(lColor).FMul(F)
	glossy := diffuse.Add(lColor.FMul(specular)).FDiv(2)

	return Math.InterpolateVector3(glossy, metal, metallic)
}
