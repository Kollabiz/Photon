package Structs

import (
	"Photon/Math"
	"math"
)

type Camera struct {
	transform   *Math.Transform
	basis       Math.Mat3
	focalLength float64
	lensSize    Math.Vector2
	resolution  Math.Vector2
}

func lensSizeFromResolution(resolution Math.Vector2) Math.Vector2 {
	if resolution.U > resolution.V {
		aspectRatio := resolution.U / resolution.V
		return Math.Vector2{
			U: 1,
			V: aspectRatio,
		}
	} else {
		aspectRatio := resolution.V / resolution.U
		return Math.Vector2{
			U: aspectRatio,
			V: 1,
		}
	}
}

func focalLengthFromFOV(fov float64) float64 {
	return math.Tan(fov / 180 * math.Pi)
}

func NewCamera(position, rotation Math.Vector3, resolution Math.Vector2, fov float64) *Camera {
	c := &Camera{}
	c.transform = Math.NewTransform(position, rotation, Math.Vector3{1, 1, 1})
	c.basis = Math.FromBasisMat3(c.transform.GetRotationMatrix().VecMul(Math.Vector3{Z: 1}))
	c.lensSize = lensSizeFromResolution(resolution)
	c.focalLength = focalLengthFromFOV(fov)
	c.resolution = resolution
	return c
}

func (c *Camera) recalculateBasis() {
	c.basis = Math.FromBasisMat3(c.transform.GetRotationMatrix().VecMul(Math.Vector3{Z: 1}))
}

func (c *Camera) GetCameraGrid(uv Math.Vector2) (pointPos Math.Vector3, direction Math.Vector3) {
	focalPoint := Math.Vector3{Z: -c.focalLength}
	pX := c.lensSize.U*(uv.U/c.resolution.U) - c.lensSize.U/2
	pY := c.lensSize.V*(uv.V/c.resolution.V) - c.lensSize.V/2
	point := Math.Vector3{pX, pY, 0}
	d := focalPoint.Sub(point).Normalized()
	return c.basis.VecMul(point).Add(c.transform.GetPosition()), c.basis.VecMul(d)
}

func (c *Camera) MoveTo(position Math.Vector3) {
	c.transform.SetPosition(position)
}

func (c *Camera) SetRotation(rotation Math.Vector3) {
	c.transform.SetRotation(rotation)
	c.recalculateBasis()
}

func (c *Camera) Move(offset Math.Vector3) {
	c.transform.Move(offset)
}

func (c *Camera) Rotate(rotation Math.Vector3) {
	c.transform.Rotate(rotation)
	c.recalculateBasis()
}
