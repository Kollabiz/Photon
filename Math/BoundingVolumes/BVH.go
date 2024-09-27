package BoundingVolumes

import "Photon/Math"

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

type BVHNode struct {
	AABB     *AABoundingBox
	Children []BVHNode
}

func NewBVHNode(p1, p2 Math.Vector3) *BVHNode {
	return &BVHNode{&AABoundingBox{p1, p2}, []BVHNode{}}
}

func (node *BVHNode) AddChild(child *BVHNode) {
	node.Children = append(node.Children, *child)
}

func (node *BVHNode) AddChildren(children []BVHNode) {
	node.Children = append(node.Children, children...)
}
