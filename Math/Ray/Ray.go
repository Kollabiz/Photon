package Ray

import (
	"Photon/Math"
	"Photon/Math/BoundingVolumes"
	"Photon/Math/Mesh"
	"math"
)

const epsilon = 0.00001

type Ray struct {
	Origin    Math.Vector3
	Direction Math.Vector3
}

func NewRay(origin, direction Math.Vector3) Ray {
	return Ray{origin, direction}
}

func (r *Ray) GetPoint(t float64) Math.Vector3 {
	return r.Direction.FMul(t).Add(r.Origin)
}

func (r *Ray) IntersectTriangle(tri *Mesh.Triangle) (bool, Math.Vector3, Math.Vector2) {
	var h, s, q Math.Vector3
	e1 := tri.Edge1.Vector()
	e2 := tri.Edge2.Vector()
	var a, f, u, v float64
	h = r.Direction.Cross(e2)
	a = h.Dot(e1)

	if a > -epsilon && a < epsilon {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	f = 1 / a
	s = r.Origin.Sub(tri.FirstVertPosition())
	u = f * s.Dot(h)

	if u < 0 || u > 1 {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	q = s.Cross(e1)
	v = r.Direction.Dot(q)

	if v < 0 || u+v > 1 {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	t := e2.Dot(q)

	if t <= epsilon {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	return true, r.GetPoint(t), Math.Vector2{u, v}
}

func (r *Ray) IntersectAABB(aabb *BoundingVolumes.AABoundingBox) bool {
	dirFrac := Math.Vector3{
		X: 1.0 / r.Direction.X,
		Y: 1.0 / r.Direction.Y,
		Z: 1.0 / r.Direction.Z,
	}
	t1 := (aabb.Point1.X - r.Origin.X) * dirFrac.X
	t2 := (aabb.Point2.X - r.Origin.X) * dirFrac.X
	t3 := (aabb.Point1.Y - r.Origin.Y) * dirFrac.Y
	t4 := (aabb.Point2.Y - r.Origin.Y) * dirFrac.Y
	t5 := (aabb.Point1.Z - r.Origin.Z) * dirFrac.Z
	t6 := (aabb.Point2.Z - r.Origin.Z) * dirFrac.Z

	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	if tmax < 0 || tmax < tmin {
		return false
	}

	return true
}
