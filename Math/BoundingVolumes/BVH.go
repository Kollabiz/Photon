package BoundingVolumes

import (
	"Photon/Math"
	"Photon/Math/Mesh"
	"Photon/Utils"
	"math"
	"math/rand"
	"strconv"
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
	AABB             *AABoundingBox
	Child1           *BVHNode
	Child2           *BVHNode
	Mesh             *Mesh.Mesh
	TriangleClusters []Mesh.TriangleCluster
}

// Joining BVH nodes

func JoinedNode(node1, node2 *BVHNode) BVHNode {
	minX := math.Min(node1.AABB.Point1.X, node2.AABB.Point1.X)
	maxX := math.Max(node1.AABB.Point1.X, node2.AABB.Point1.X)
	minY := math.Min(node1.AABB.Point1.Y, node2.AABB.Point1.Y)
	maxY := math.Max(node1.AABB.Point1.Y, node2.AABB.Point1.Y)
	minZ := math.Min(node1.AABB.Point1.Z, node2.AABB.Point1.Z)
	maxZ := math.Max(node1.AABB.Point1.Z, node2.AABB.Point1.Z)
	return BVHNode{
		AABB:             NewAABB(Math.Vector3{minX, minY, minZ}, Math.Vector3{maxX, maxY, maxZ}),
		Child1:           node1,
		Child2:           node2,
		Mesh:             nil,
		TriangleClusters: nil,
	}
}

// The hardest part, BVH node from a mesh

func BVHFromMesh(mesh *Mesh.Mesh, pointRatio float64) *BVHNode {
	Utils.Log("Creating acceleration structures for mesh " + mesh.MeshName)
	clusterCount := int(float64(len(mesh.Triangles)) * pointRatio)
	Utils.Log(strconv.Itoa(clusterCount) + " triangle clusters will be created")

	clusters := make([]Mesh.TriangleCluster, clusterCount)

	// Preparing clusters
	for i := 0; i < len(clusters); i++ {
		triIdx := rand.Intn(len(mesh.Triangles))
		tri := &mesh.Triangles[triIdx]
		clusters[i].AABB = NewAABB(
			Math.Vector3{
				X: min(tri.V1Pos.X, tri.V2Pos.X, tri.V3Pos.X),
				Y: min(tri.V1Pos.Y, tri.V2Pos.Y, tri.V3Pos.Y),
				Z: min(tri.V1Pos.Z, tri.V2Pos.Z, tri.V3Pos.Z),
			},
			Math.Vector3{
				X: max(tri.V1Pos.X, tri.V2Pos.X, tri.V3Pos.X),
				Y: max(tri.V1Pos.Y, tri.V2Pos.Y, tri.V3Pos.Y),
				Z: max(tri.V1Pos.Z, tri.V2Pos.Z, tri.V3Pos.Z),
			},
		)
	}

	Utils.Log("Assigning triangles to clusters")
	// Iterating through all the triangles and assigning them to clusters
	for j := 0; j < len(mesh.Triangles); j++ {
		tri := &mesh.Triangles[j]
		midPoint := tri.Middle()
		// Finding the closest cluster
		closestCluster := &clusters[0]
		closestDist := closestCluster.AABB.MiddlePoint().Sub(midPoint).LenSq()
		for k := 1; k < len(clusters); k++ {
			d := clusters[k].AABB.MiddlePoint().Sub(midPoint).LenSq()
			if d < closestDist {
				closestCluster = &clusters[k]
				closestDist = d
			}
		}
		closestCluster.Triangles = append(closestCluster.Triangles, *tri)
		closestCluster.AABB.Point1 = Math.Vector3{
			X: min(tri.V1Pos.X, tri.V2Pos.X, tri.V3Pos.X, closestCluster.AABB.Point1.X),
			Y: min(tri.V1Pos.Y, tri.V2Pos.Y, tri.V3Pos.Y, closestCluster.AABB.Point1.Y),
			Z: min(tri.V1Pos.Z, tri.V2Pos.Z, tri.V3Pos.Z, closestCluster.AABB.Point1.Z),
		}
		closestCluster.AABB.Point2 = Math.Vector3{
			X: max(tri.V1Pos.X, tri.V2Pos.X, tri.V3Pos.X, closestCluster.AABB.Point2.X),
			Y: max(tri.V1Pos.Y, tri.V2Pos.Y, tri.V3Pos.Y, closestCluster.AABB.Point2.Y),
			Z: max(tri.V1Pos.Z, tri.V2Pos.Z, tri.V3Pos.Z, closestCluster.AABB.Point2.Z),
		}
	}

	Utils.Log("Building mesh AABB")
	aabb := NewAABB(clusters[0].AABB.Point1, clusters[0].AABB.Point2)
	for j := 1; j < len(clusters); j++ {
		aabb.Point1 = Math.Vector3{
			X: min(aabb.Point1.X, clusters[j].AABB.Point1.X),
			Y: min(aabb.Point1.Y, clusters[j].AABB.Point1.Y),
			Z: min(aabb.Point1.Z, clusters[j].AABB.Point1.Z),
		}
		aabb.Point2 = Math.Vector3{
			X: max(aabb.Point2.X, clusters[j].AABB.Point2.X),
			Y: max(aabb.Point2.Y, clusters[j].AABB.Point2.Y),
			Z: max(aabb.Point2.Z, clusters[j].AABB.Point2.Z),
		}
	}

	Utils.Log("Assembling BVH node")
	node := &BVHNode{
		AABB:             aabb,
		Child1:           nil,
		Child2:           nil,
		Mesh:             mesh,
		TriangleClusters: clusters,
	}

	Utils.LogSuccess("Done building acceleration structures for mesh " + mesh.MeshName)
	return node
}
