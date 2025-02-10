package PhotonMapping

import (
	"Photon/Math"
	"Photon/Math/BoundingVolumes"
	"Photon/Structs"
	"math"
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
	Position  Math.Vector3
	NextPoint *CameraPoint
	I         Math.Vector3
	R         Math.Vector3
	Mat       *Structs.Material
	N         Math.Vector3
}

type KDTreeSpace struct {
	Points    []CameraPoint
	SplitAxis uint8
	Domain    BoundingVolumes.AABoundingBox
	Subspace1 *KDTreeSpace
	Subspace2 *KDTreeSpace
}

func getPointCloudBoundaries(cloud []CameraPoint) BoundingVolumes.AABoundingBox {
	// Infinitely small bounding box (in fact, it has negative size)
	box := BoundingVolumes.AABoundingBox{
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

func ConstructKDTree(pointCloud []CameraPoint, maxPointsPerDomain int) *KDTreeSpace {
	bounds := getPointCloudBoundaries(pointCloud)
	root := &KDTreeSpace{
		Points:    pointCloud,
		SplitAxis: 0,
		Domain:    bounds,
		Subspace1: nil,
		Subspace2: nil,
	}
	nodeQueue := []KDTreeSpace{*root}

	for len(nodeQueue) > 0 {
		node := nodeQueue[0]
		// If the domain is small enough, we can already skip it
		if len(node.Points) <= maxPointsPerDomain {
			nodeQueue = nodeQueue[1:]
			continue
		}
		// Domain min
		dmin := plane(node.SplitAxis, node.Domain.Point1)
		// Domain max
		dmax := plane(node.SplitAxis, node.Domain.Point2)

		// Split threshold. Everything less is on one side of the split plane, and everything that's greater is on the
		// other
		splitThreshold := (dmax-dmin)/4 + dmin

		// Subspace containing all points whose respective coordinates are LESS than splitThreshold
		ndLess := &KDTreeSpace{
			SplitAxis: (node.SplitAxis + 1) % 3, // Split planes should change every iteration
			Domain: BoundingVolumes.AABoundingBox{
				Point1: node.Domain.Point1,                                               // The MIN point of the domain won't change for that node
				Point2: planeSet(node.SplitAxis, node.Domain.Point2, dmax-(dmax-dmin)/2), // The MAX point is set
				// to the middle of the cloud domain (only along the split plane's perpendicular axis)
			},
		}
		// Subspace containing all points whose respective coordinates are GREATER than splitThreshold
		ndGreater := &KDTreeSpace{
			SplitAxis: (node.SplitAxis + 1) % 3, // Split planes should change every iteration
			Domain: BoundingVolumes.AABoundingBox{
				Point1: planeSet(node.SplitAxis, node.Domain.Point1, dmax-(dmax-dmin)/2),
				Point2: node.Domain.Point2, // The MAX point for that domain won't change
				// to the middle of the cloud domain (only along the split plane's perpendicular axis)
			},
		}
		// Assigning points
		for i := 0; i < len(node.Points); i++ {
			// checking which domain the point falls into
			if plane(node.SplitAxis, node.Points[i].Position) <= splitThreshold {
				ndLess.Points = append(ndLess.Points, node.Points[i])
			} else {
				ndGreater.Points = append(ndGreater.Points, node.Points[i])
			}
		}

		// if a node has no points, we just skip it
		if len(ndLess.Points) > 0 {
			node.Subspace1 = ndLess
		}
		if len(ndGreater.Points) > 0 {
			node.Subspace2 = ndGreater
		}
		nodeQueue = nodeQueue[1:]
		nodeQueue = append(nodeQueue, *ndLess, *ndGreater)
	}

	return root
}

func (tree *KDTreeSpace) LocateNeighborPoints(point Math.Vector3, pointRadius float64) []CameraPoint {
	currentNode := tree
	for currentNode.Subspace1 != nil && currentNode.Subspace2 != nil {
		// Domain split threshold
		splitThreshold := (plane(currentNode.SplitAxis, currentNode.Domain.Point2)-plane(currentNode.SplitAxis, currentNode.Domain.Point1))/4 + plane(currentNode.SplitAxis, currentNode.Domain.Point1)
		// Distance of the point to the split plane
		d := splitThreshold - plane(currentNode.SplitAxis, point)
		// The point overlaps both subspaces. To avoid abrupt lighting falloff at the domain edges, we should return
		// both overlapped domains
		if math.Abs(d) < pointRadius {
			return currentNode.Points
		} else {
			if d >= 0 { // The point belongs in the LESS domain
				if currentNode.Subspace1 == nil { // There is no LESS domain, just return the entire node
					return currentNode.Points
				} else {
					currentNode = currentNode.Subspace1
				}
			} else { // The point belongs in the GREATER domain
				if currentNode.Subspace2 == nil { // There is no GREATER domain, just return the entire node
					return currentNode.Points
				} else {
					currentNode = currentNode.Subspace2
				}
			}
		}
	}
	return currentNode.Points
}
