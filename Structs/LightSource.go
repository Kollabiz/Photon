package Structs

import (
	"Photon/Math"
	"Photon/Utils"
	"math"
	"math/rand"
)

const AnArbitrarilyBigNumber = 65536

type LightSource interface {
	GetLightDirectionTo(point Math.Vector3) Math.Vector3
	GetLightIntensityTo(point Math.Vector3) float64
	GetLightIntensityInDirection(dir Math.Vector3) float64
	GetLightColor() Math.Vector3
	GetPosition() Math.Vector3
	GetID() int
	GetRandomPoint(gen *rand.Rand) Math.Vector3
}

// Point Light

type PointLight struct {
	Position  Math.Vector3
	Intensity float64
	Color     Math.Vector3
	id        int
}

func NewPointLight(position Math.Vector3, intensity float64, color Math.Vector3) *PointLight {
	l := &PointLight{
		Position:  position,
		Intensity: intensity,
		Color:     color,
		id:        rand.Int(),
	}
	return l
}

func (p *PointLight) GetRandomPoint(gen *rand.Rand) Math.Vector3 {
	return Utils.RandomPointOnSphere(gen)
}

func (p *PointLight) GetPosition() Math.Vector3 {
	return p.Position
}

func (p *PointLight) GetLightDirectionTo(point Math.Vector3) Math.Vector3 {
	return p.Position.Sub(point).Normalized()
}

func (p *PointLight) GetLightIntensityTo(point Math.Vector3) float64 {
	d := p.Position.Sub(point).LenSq()
	return p.Intensity / d
}

func (p *PointLight) GetLightIntensityInDirection(dir Math.Vector3) float64 {
	return p.Intensity
}

func (p *PointLight) GetLightColor() Math.Vector3 {
	return p.Color
}

func (p *PointLight) GetID() int {
	return p.id
}

// Sun Light

type SunLight struct {
	Direction Math.Vector3
	Intensity float64
	Color     Math.Vector3
	id        int
}

func NewSunLight(direction Math.Vector3, intensity float64, color Math.Vector3) *SunLight {
	s := &SunLight{
		Direction: direction,
		Intensity: intensity,
		Color:     color,
		id:        rand.Int(),
	}
	return s
}

func (s *SunLight) GetRandomPoint(gen *rand.Rand) Math.Vector3 {
	return s.Direction
}

func (s *SunLight) GetPosition() Math.Vector3 {
	return s.Direction.Inverse().FMul(AnArbitrarilyBigNumber)
}

func (s *SunLight) GetLightDirectionTo(point Math.Vector3) Math.Vector3 {
	return s.Direction.Normalized()
}

func (s *SunLight) GetLightIntensityTo(point Math.Vector3) float64 {
	return s.Intensity
}

func (s *SunLight) GetLightIntensityInDirection(dir Math.Vector3) float64 {
	return s.Intensity
}

func (s *SunLight) GetLightColor() Math.Vector3 {
	return s.Color
}

func (s *SunLight) GetID() int {
	return s.id
}

// Cone Light

type ConeLight struct {
	Position  Math.Vector3
	Direction Math.Vector3
	Intensity float64
	Falloff   float64
	Color     Math.Vector3
	id        int
}

func NewConeLight(position, direction Math.Vector3, intensity, falloff float64, color Math.Vector3) *ConeLight {
	s := &ConeLight{
		Position:  position,
		Direction: direction,
		Intensity: intensity,
		Falloff:   falloff,
		Color:     color,
		id:        rand.Int(),
	}
	return s
}

func (c *ConeLight) GetRandomPoint(gen *rand.Rand) Math.Vector3 {
	p := Utils.RandomPointOnHemisphereConstrained(c.Falloff, gen).FromSingleVectorBasis(c.Direction)
	return p
}

func (c *ConeLight) GetPosition() Math.Vector3 {
	return c.Position
}

func (c *ConeLight) GetLightDirectionTo(point Math.Vector3) Math.Vector3 {
	return c.Position.Sub(point).Normalized()
}

func (c *ConeLight) GetLightIntensityTo(point Math.Vector3) float64 {
	d := c.Position.Sub(point).LenSq()
	falloff := math.Max((c.Direction.Dot(c.GetLightDirectionTo(point))-c.Falloff)*(1-c.Falloff), 0)
	return d * falloff
}

func (c *ConeLight) GetLightIntensityInDirection(dir Math.Vector3) float64 {
	falloff := math.Max((c.Direction.Dot(dir)-c.Falloff)*(1-c.Falloff), 0)
	return falloff * c.Intensity
}

func (c *ConeLight) GetLightColor() Math.Vector3 {
	return c.Color
}

func (c *ConeLight) GetID() int {
	return c.id
}
