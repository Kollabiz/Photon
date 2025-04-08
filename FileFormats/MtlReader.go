package FileFormats

import (
	"Photon/Math"
	"Photon/Structs"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func tryParseColor(line []string) (r, g, b float64, err error) {
	r, err = strconv.ParseFloat(line[1], 64)
	g, err = strconv.ParseFloat(line[2], 64)
	b, err = strconv.ParseFloat(line[3], 64)
	return
}

func checkLen(ln []string, desired int, component string) {
	if len(ln) < desired {
		panic(component + " has not enough values")
	}
}

type MTLParser struct {
	albedoTextures    map[string]*Structs.TextureRGB
	roughnessTextures map[string]*Structs.TextureGrayscale
	metallicTextures  map[string]*Structs.TextureGrayscale
	Brdf              Structs.IBRDF
}

func (parser *MTLParser) DropTables() {
	parser.albedoTextures = make(map[string]*Structs.TextureRGB)
	parser.roughnessTextures = make(map[string]*Structs.TextureGrayscale)
	parser.metallicTextures = make(map[string]*Structs.TextureGrayscale)
}

func (parser *MTLParser) lookupOrOpenAlbedoTexture(tex string) *Structs.TextureRGB {
	if parser.albedoTextures[tex] != nil {
		return parser.albedoTextures[tex]
	}
	otex := Structs.ReadTextureRGB(tex)
	parser.albedoTextures[tex] = otex
	return otex
}

func (parser *MTLParser) lookupOrOpenMetallicTexture(tex string) *Structs.TextureGrayscale {
	if parser.albedoTextures[tex] != nil {
		return parser.metallicTextures[tex]
	}
	otex := Structs.ReadTextureGrayscale(tex)
	parser.metallicTextures[tex] = otex
	return otex
}

func (parser *MTLParser) lookupOrOpenRoughnessTexture(tex string) *Structs.TextureGrayscale {
	if parser.albedoTextures[tex] != nil {
		return parser.roughnessTextures[tex]
	}
	otex := Structs.ReadTextureGrayscale(tex)
	parser.roughnessTextures[tex] = otex
	return otex
}

func (parser *MTLParser) Parse(mtlFile string) map[string]*Structs.Material {
	var parsedMaterials map[string]*Structs.Material = make(map[string]*Structs.Material)
	var currentMaterial *Structs.Material
	var currentMaterialName string

	fl, err := os.Open(mtlFile)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fl)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		switch strings.ToLower(line[0]) {
		case "newmtl": // new material
			if currentMaterial != nil {
				parsedMaterials[currentMaterialName] = currentMaterial
			}
			currentMaterial = &Structs.Material{BRDF: parser.Brdf}
			currentMaterialName = line[1]
			break

		case "kd": // Albedo color
			checkLen(line, 4, "kd")
			r, g, b, err := tryParseColor(line)
			if err != nil {
				panic(err)
			}
			currentMaterial.SetAlbedo(Math.Vector3{r, g, b})
			break
		case "map_kd":
			checkLen(line, 2, "map_kd")
			mkd := parser.lookupOrOpenAlbedoTexture(line[1])
			currentMaterial.SetAlbedoTexture(mkd)
			break

		case "ns": // Roughness
			checkLen(line, 2, "ns")
			ks, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				panic(err)
			}
			currentMaterial.SetRoughness(1 - ks/1000)
			break
		case "map_ks":
			checkLen(line, 2, "map_ks")
			mks := parser.lookupOrOpenRoughnessTexture(line[1])
			currentMaterial.SetRoughnessTexture(mks)
			break

		case "ka": // Metallic. MTL is old so Blender uses ambient light (Ka) for metallic
			checkLen(line, 2, "ka")
			ka, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				panic(err)
			}
			currentMaterial.SetMetallic(ka)
			break
		case "map_ka":
			checkLen(line, 2, "map_ka")
			mka := parser.lookupOrOpenMetallicTexture(line[1])
			currentMaterial.SetMetallicTexture(mka)
			break

		case "Ni": // IOR
			checkLen(line, 2, "ni")
			ni, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				panic(err)
			}
			currentMaterial.SetIOR(ni)
			break
		}
	}

	parsedMaterials[currentMaterialName] = currentMaterial
	return parsedMaterials
}
