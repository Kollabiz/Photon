package BoundingVolumes

import (
	"Photon/Math"
	"Photon/Math/Mesh"
	"math"
)

type AABoundingBox struct {
	Point1 Math.Vector3
	Point2 Math.Vector3
}

func NewAABB(p1, p2 Math.Vector3) *AABoundingBox {
	return &AABoundingBox{
		Point1: p1,
		Point2: p2,
	}
}

func (aabb *AABoundingBox) MiddlePoint() Math.Vector3 {
	return aabb.Point1.Add(aabb.Point2).FDiv(2)
}

type BVHNode struct {
	AABB     *AABoundingBox
	Children []BVHNode
	Mesh     *Mesh.Mesh
}

func NewBVHNodeFromMesh(mesh *Mesh.Mesh) *BVHNode {
	node := &BVHNode{}
	// Constructing AABB
	// We need minimal and maximal point
	minX := math.Inf(1)
	minY := math.Inf(1)
	minZ := math.Inf(1)
	maxX := math.Inf(-1)
	maxY := math.Inf(-1)
	maxZ := math.Inf(-1)
	for ti := 0; ti < len(mesh.Triangles); ti++ {
		tri := mesh.Triangles[ti]
		p1 := tri.FirstVertPosition()
		p2 := tri.SecondVertPosition()
		p3 := tri.ThirdVertPosition()
		lmx := math.Min(math.Min(p1.X, p2.X), p3.X)
		lmy := math.Min(math.Min(p1.Y, p2.Y), p3.Y)
		lmz := math.Min(math.Min(p1.Z, p2.Z), p3.Z)
		lmax := math.Max(math.Max(p1.X, p2.X), p3.X)
		lmay := math.Max(math.Max(p1.Y, p2.Y), p3.Y)
		lmaz := math.Max(math.Max(p1.Z, p2.Z), p3.Z)
		if lmx < minX {
			minX = lmx
		}
		if lmy < minY {
			minY = lmy
		}
		if lmz < minZ {
			minZ = lmz
		}
		if lmax > maxX {
			maxX = lmax
		}
		if lmay > maxY {
			maxY = lmay
		}
		if lmaz > maxZ {
			maxZ = lmaz
		}
	}
	aabb := NewAABB(Math.Vector3{minX, minY, minZ}, Math.Vector3{maxX, maxY, maxZ})
	node.AABB = aabb
	node.Mesh = mesh
	return node
}

func NewBVHNode(p1, p2 Math.Vector3) *BVHNode {
	return &BVHNode{&AABoundingBox{p1, p2}, []BVHNode{}, nil}
}

func JoinBVHNodes(n1, n2 *BVHNode) BVHNode {
	jNode := BVHNode{}
	jNode.Children = []BVHNode{*n1, *n2}
	jNode.AABB = &AABoundingBox{
		Point1: Math.Vector3{
			X: math.Min(math.Min(n1.AABB.Point1.X, n1.AABB.Point2.X), math.Min(n2.AABB.Point1.X, n2.AABB.Point2.X)),
			Y: math.Min(math.Min(n1.AABB.Point1.Y, n1.AABB.Point2.Y), math.Min(n2.AABB.Point1.Y, n2.AABB.Point2.Y)),
			Z: math.Min(math.Min(n1.AABB.Point1.Z, n1.AABB.Point2.Z), math.Min(n2.AABB.Point1.Z, n2.AABB.Point2.Z)),
		},
		Point2: Math.Vector3{
			X: math.Max(math.Max(n1.AABB.Point1.X, n1.AABB.Point2.X), math.Max(n2.AABB.Point1.X, n2.AABB.Point2.X)),
			Y: math.Max(math.Max(n1.AABB.Point1.Y, n1.AABB.Point2.Y), math.Max(n2.AABB.Point1.Y, n2.AABB.Point2.Y)),
			Z: math.Max(math.Max(n1.AABB.Point1.Z, n1.AABB.Point2.Z), math.Max(n2.AABB.Point1.Z, n2.AABB.Point2.Z)),
		},
	}
	return jNode
}

func (node *BVHNode) AddChild(child *BVHNode) {
	node.Children = append(node.Children, *child)
}

func (node *BVHNode) AddChildren(children []BVHNode) {
	node.Children = append(node.Children, children...)
}

func (node *BVHNode) Copy() *BVHNode {
	return &BVHNode{
		AABB: NewAABB(
			node.AABB.Point1,
			node.AABB.Point1,
		),
		Children: nil,
	}
}
