package zellij

import "../quadratic/quadratic"

var Points map[string]*quadratic.Point
var Tiles []string
var TileMaps []*quadratic.Map
var Workers int = 1

func init() {
	Points = make(map[string]*quadratic.Point)
	Points["a"], _ = quadratic.PointFromString("0,0,0,2")
	Points["b"], _ = quadratic.PointFromString("-2,0,2,0")
	Points["c"], _ = quadratic.PointFromString("2,-2,2,0")
	Points["d"], _ = quadratic.PointFromString("-2,2,2,0")
	Points["e"], _ = quadratic.PointFromString("2,0,2,0")
	Points["f"], _ = quadratic.PointFromString("0,0,4,-2")
	Points["g"], _ = quadratic.PointFromString("-2,0,-2,2")
	Points["h"], _ = quadratic.PointFromString("2,0,-2,2")
	Points["i"], _ = quadratic.PointFromString("-4,0,0,0")
	Points["j"], _ = quadratic.PointFromString("0,-2,0,0")
	Points["k"], _ = quadratic.PointFromString("-4,2,0,0")
	Points["l"], _ = quadratic.PointFromString("0,0,0,0")
	Points["m"], _ = quadratic.PointFromString("4,-2,0,0")
	Points["n"], _ = quadratic.PointFromString("0,2,0,0")
	Points["o"], _ = quadratic.PointFromString("4,0,0,0")
	Points["p"], _ = quadratic.PointFromString("-2,0,2,-2")
	Points["q"], _ = quadratic.PointFromString("0,0,-4,2")
	Points["r"], _ = quadratic.PointFromString("2,0,2,-2")
	Points["s"], _ = quadratic.PointFromString("-2,0,-2,0")
	Points["t"], _ = quadratic.PointFromString("2,-2,-2,0")
	Points["u"], _ = quadratic.PointFromString("-2,2,-2,0")
	Points["v"], _ = quadratic.PointFromString("2,0,-2,0")
	Points["w"], _ = quadratic.PointFromString("0,0,0,-2")



	Points["A"], _ = quadratic.PointFromString("-2,-2,2,0")
	Points["B"], _ = quadratic.PointFromString("2,2,2,0")
	Points["C"], _ = quadratic.PointFromString("-2,-2,-2,0")
	Points["D"], _ = quadratic.PointFromString("2,2,-2,0")

	Tiles = make([]string,9 )
	Tiles[0] = "adehnrvuwtspjgbc"
	Tiles[1] = "beovsi"
	Tiles[2] = "cdhvtj"
	Tiles[3] = "dBrh"
	Tiles[4] = "Bnhe"
	Tiles[5] = "leov"
	Tiles[6] = "ibelvs"
	Tiles[7] = "ibehnrvs"
	Tiles[8] = "bcfdehmrvuqtspkg"

	for _,t := range(Tiles) {
		base := TileMap(t,0)
		TileMaps = append(TileMaps, GenerateOrbits(base)...)
	}
	
	VertexFigures = []byte{
		leftRotate(5,0),
		leftRotate(5,1),
		leftRotate(5,2),
		leftRotate(5,3),
		leftRotate(5,4),
		leftRotate(5,5),
		leftRotate(5,6),
		leftRotate(5,7),
		leftRotate(9,0),
		leftRotate(9,1),
		leftRotate(9,2),
		leftRotate(9,3),
		leftRotate(9,4),
		leftRotate(9,5),
		leftRotate(9,6),
		leftRotate(9,7),
		leftRotate(85,0),
		leftRotate(85,1),
		leftRotate(51,0),
		leftRotate(51,1),
		leftRotate(51,2),
		leftRotate(51,3)}
}

