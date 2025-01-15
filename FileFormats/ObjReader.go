package FileFormats

import (
	"Photon/Math"
	"Photon/Math/Mesh"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadOBJFile(file string) []Mesh.Mesh {
	objFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	objFileBuff := bufio.NewScanner(objFile)

	// Mesh data
	var currentMesh *Mesh.Mesh
	var meshes []Mesh.Mesh
	var vertices []Math.Vector3
	var textureCoords []Math.Vector2
	var normals []Math.Vector3
	var x, y, z float64

	for objFileBuff.Scan() {
		line := strings.Split(objFileBuff.Text(), " ")
		lineType := line[0]
		switch lineType {
		case "o":
			if currentMesh != nil {
				meshes = append(meshes, *currentMesh)
			}
			currentMesh = &Mesh.Mesh{MeshName: line[1], Transform: Math.NewTransform(Math.Vector3{0, 0, 0}, Math.Vector3{0, 0, 0}, Math.Vector3{1, 1, 1})}
			break
		case "v": // Vertices
			x, err = strconv.ParseFloat(line[1], 64)
			y, err = strconv.ParseFloat(line[2], 64)
			z, err = strconv.ParseFloat(line[3], 64)
			if err != nil {
				panic(err)
			}
			vertices = append(vertices, Math.Vector3{x, y, z})
			break
		case "vt": // Texture coords
			x, err = strconv.ParseFloat(line[1], 64)
			if len(line) > 2 {
				y, err = strconv.ParseFloat(line[2], 64)
			} else {
				y = 0
			}
			if err != nil {
				panic(err)
			}
			textureCoords = append(textureCoords, Math.Vector2{x, y})
			break
		case "vn":
			x, err = strconv.ParseFloat(line[1], 64)
			y, err = strconv.ParseFloat(line[2], 64)
			z, err = strconv.ParseFloat(line[3], 64)
			if err != nil {
				panic(err)
			}
			normals = append(normals, Math.Vector3{x, y, z}.Normalized())
			break
		case "f":
			if len(line) >= 4 {
				p1 := strings.Split(line[2], "/")
				p2 := strings.Split(line[3], "/")
				p3 := strings.Split(line[4], "/")
				x, err = strconv.ParseFloat(p1[0], 64)
				po1 := vertices[int(x)]
				y, err = strconv.ParseFloat(p2[0], 64)
				po2 := vertices[int(y)]
				z, err = strconv.ParseFloat(p3[0], 64)
				po3 := vertices[int(z)]
				var tx1, tx2, tx3 Math.Vector2
				var n1, n2, n3 Math.Vector3
				if len(p1) > 1 && p1[1] != "" {
					// First vertex texture coordinates
					x, err = strconv.ParseFloat(p1[1], 64)
					tx1 = textureCoords[int(x)]
				}
				if len(p2) > 1 && p2[1] != "" {
					// Second vertex texture coordinates
					y, err = strconv.ParseFloat(p2[1], 64)
					tx2 = textureCoords[int(y)]
				}
				if len(p3) > 1 && p3[1] != "" {
					// Third vertex texture coordinates
					z, err = strconv.ParseFloat(p3[1], 64)
					tx3 = textureCoords[int(z)]
				}
				if len(p1) > 2 && p1[2] != "" {
					// First vertex normal
					x, err = strconv.ParseFloat(p1[2], 64)
					n1 = normals[int(x)]
				}
				if len(p2) > 2 && p2[2] != "" {
					// Second vertex normal
					y, err = strconv.ParseFloat(p2[2], 64)
					n2 = normals[int(y)]
				}
				if len(p3) > 2 && p3[2] != "" {
					// Third vertex normal
					z, err = strconv.ParseFloat(p3[2], 64)
					n3 = normals[int(z)]
				}
				if err != nil {
					panic(err)
				}
				face := Mesh.Triangle{
					V1Pos:    po1,
					V2Pos:    po2,
					V3Pos:    po3,
					V1Tex:    tx1,
					V2Tex:    tx2,
					V3Tex:    tx3,
					V1Normal: n1,
					V2Normal: n2,
					V3Normal: n3,
					Material: nil,
				}
				face.RecalcNormal()
				currentMesh.Triangles = append(currentMesh.Triangles, face)
				break
			}
		}
	}
	return meshes
}
