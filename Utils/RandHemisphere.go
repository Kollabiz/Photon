package Utils

import (
	"Photon/Math"
	"math"
	"math/rand"
	"time"
)

func randFloat(gen *rand.Rand) float64 {
	f := gen.NormFloat64() / math.MaxFloat64
	gen.Seed(time.Now().UnixMilli())
	return f
}

func RandomPointOnSphere(gen *rand.Rand) Math.Vector3 {
	phi := randFloat(gen) * math.Pi
	theta := randFloat(gen) * math.Pi
	phiCos := math.Cos(phi)
	phiSin := math.Sin(phi)
	thetaCos := math.Cos(theta)
	thetaSin := math.Sin(theta)
	return Math.Vector3{
		X: thetaCos * phiSin,
		Y: thetaSin * phiSin,
		Z: phiCos,
	}.Normalized()
}

func RandomPointOnHemisphere(gen *rand.Rand) Math.Vector3 {
	phi := gen.Float64() * math.Pi / 2
	gen.Seed(time.Now().UnixMilli())
	theta := randFloat(gen) * math.Pi
	phiCos := math.Cos(phi)
	phiSin := math.Sin(phi)
	thetaCos := math.Cos(theta)
	thetaSin := math.Sin(theta)
	return Math.Vector3{
		X: thetaCos * phiSin,
		Y: thetaSin * phiSin,
		Z: phiCos,
	}.Normalized()
}

func RandomPointOnHemisphereConstrained(cone float64, gen *rand.Rand) Math.Vector3 {
	phi := gen.Float64() * cone * math.Pi / 2
	gen.Seed(time.Now().UnixMilli())
	theta := randFloat(gen) * math.Pi
	phiCos := math.Cos(phi)
	phiSin := math.Sin(phi)
	thetaCos := math.Cos(theta)
	thetaSin := math.Sin(theta)
	return Math.Vector3{
		X: thetaCos * phiSin,
		Y: thetaSin * phiSin,
		Z: phiCos,
	}.Normalized()
}
