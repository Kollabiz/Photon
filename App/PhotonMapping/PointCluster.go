package PhotonMapping

import (
	"Photon/Math"
	"Photon/Structs"
	"Photon/Utils"
	"math"
	"strconv"
)

// First pass point cluster
// Contains a bunch of points which are nodes of ScreenWidth*ScreenHeight linked lists
// Each node contains:
// - Position
// - Next node
// - Incoming light vector
// - Reflected light vector
// - Material reference
// - Normal
// All point are arranged in a K-D tree. During rendering, the photon point is represented as a "sphere", so that it can
// Cover multiple nodes at a time

type CameraPoint struct {
	Position           Math.Vector3
	NextPoint          *CameraPoint
	I                  Math.Vector3
	R                  Math.Vector3
	Triangle           *Structs.Triangle
	Bary               Math.Vector2
	Color              Math.Vector3
	AccumulatedPhotons int
}

type KDTreeSpace struct {
	Points    []*CameraPoint
	SplitAxis uint8
	Domain    Structs.AABoundingBox
	Subspace1 *KDTreeSpace
	Subspace2 *KDTreeSpace
}

func getPointCloudBoundaries(cloud []*CameraPoint) Structs.AABoundingBox {
	// Infinitely small bounding box (in fact, it has negative size)
	box := Structs.AABoundingBox{
		Point1: Math.InfiniteVector3(),
		Point2: Math.NegativeInfiniteVector3(),
	}
	for i := 0; i < len(cloud); i++ {
		point := cloud[i]
		box.Point1 = Math.Vector3{
			X: math.Min(point.Position.X, box.Point1.X),
			Y: math.Min(point.Position.Y, box.Point1.Y),
			Z: math.Min(point.Position.Z, box.Point1.Z),
		}
		box.Point2 = Math.Vector3{
			X: math.Max(point.Position.X, box.Point2.X),
			Y: math.Max(point.Position.Y, box.Point2.Y),
			Z: math.Max(point.Position.Z, box.Point2.Z),
		}
	}
	return box
}

func plane(idx uint8, vec Math.Vector3) float64 {
	switch idx {
	case 0: // YZ
		return vec.X
	case 1: // XZ
		return vec.Y
	case 2: // XY
		return vec.Z
	}
	panic("Invalid split plane index in KD-tree")
}

func planeSet(idx uint8, vec Math.Vector3, val float64) Math.Vector3 {
	switch idx {
	case 0: // YZ
		return Math.Vector3{
			X: val,
			Y: vec.Y,
			Z: vec.Z,
		}
	case 1: // XZ
		return Math.Vector3{
			X: vec.X,
			Y: val,
			Z: vec.Z,
		}
	case 2: // XY
		return Math.Vector3{
			X: vec.X,
			Y: vec.Y,
			Z: val,
		}
	}
	panic("Invalid split plane index in KD-tree")
}

func splitAABB(splitPlane uint8, splitThreshold float64, aabb Structs.AABoundingBox) (Structs.AABoundingBox, Structs.AABoundingBox) {
	return Structs.AABoundingBox{
			Point1: aabb.Point1,
			Point2: planeSet(splitPlane, aabb.Point2, splitThreshold),
		}, Structs.AABoundingBox{
			Point1: planeSet(splitPlane, aabb.Point1, splitThreshold),
			Point2: aabb.Point2,
		}
}

func ConstructKDTree(pointCloud []*CameraPoint, maxPointsPerDomain int) *KDTreeSpace {
	Utils.Log("Creating K-D tree for the point cloud")
	bounds := getPointCloudBoundaries(pointCloud)
	Utils.LogSuccess("Cloud bounds found. diagonal size: " + strconv.FormatFloat(bounds.Point2.Sub(bounds.Point1).Len(), 'f', 3, 64))
	root := &KDTreeSpace{
		Points:    pointCloud,
		SplitAxis: 0,
		Domain:    bounds,
		Subspace1: nil,
		Subspace2: nil,
	}
	nodeQueue := []*KDTreeSpace{root}

	Utils.Log("Creating the tree...")

	for len(nodeQueue) > 0 {
		node := nodeQueue[0]
		if len(node.Points) <= maxPointsPerDomain {
			nodeQueue = nodeQueue[1:]
			continue
		}
		splitAxis := node.SplitAxis
		splitThreshold := plane(splitAxis, node.Domain.MiddlePoint())
		lessAABB, moreAABB := splitAABB(splitAxis, splitThreshold, node.Domain)
		lessNode, moreNode := &KDTreeSpace{
			Points:    []*CameraPoint{},
			SplitAxis: (splitAxis + 1) % 3,
			Domain:    lessAABB,
			Subspace1: nil,
			Subspace2: nil,
		}, &KDTreeSpace{
			Points:    []*CameraPoint{},
			SplitAxis: (splitAxis + 1) % 3,
			Domain:    moreAABB,
			Subspace1: nil,
			Subspace2: nil,
		}
		for i := 0; i < len(node.Points); i++ {
			point := node.Points[i]
			if plane(splitAxis, point.Position) <= splitThreshold {
				lessNode.Points = append(lessNode.Points, point)
			} else {
				moreNode.Points = append(moreNode.Points, point)
			}
		}
		node.Subspace1 = lessNode
		node.Subspace2 = moreNode
		nodeQueue = nodeQueue[1:]
		nodeQueue = append(nodeQueue, lessNode, moreNode)
	}

	Utils.LogSuccess("Done building the K-D tree for the point cloud")

	return root
}

func (tree *KDTreeSpace) LocateNeighborPoints(point Math.Vector3, phR float64) *KDTreeSpace {
	currentNode := tree
	for currentNode.Subspace1 != nil {
		if currentNode.Subspace1 == nil {
			return currentNode
		}
		splitThreshold := plane(currentNode.SplitAxis, currentNode.Domain.MiddlePoint())
		pPos := plane(currentNode.SplitAxis, point)
		d := splitThreshold - pPos
		if math.Abs(d) < phR {
			return currentNode
		}
		if d < 0 {
			currentNode = currentNode.Subspace2
		} else {
			currentNode = currentNode.Subspace1
		}
	}
	return currentNode
}
