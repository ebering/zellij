package zellij

import "../quadratic/quadratic"

var Points map[string]*quadratic.Point
var Tiles []string
var WhiteTiles map[string]bool
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

	Points["A"], _ = quadratic.PointFromString("-2,0,2,2")
	Points["B"], _ = quadratic.PointFromString("2,0,2,2")
	Points["C"], _ = quadratic.PointFromString("2,2,2,0")
	Points["D"], _ = quadratic.PointFromString("2,2,-2,0")
	Points["E"], _ = quadratic.PointFromString("2,0,-2,-2")
	Points["F"], _ = quadratic.PointFromString("-2,0,-2,-2")
	Points["G"], _ = quadratic.PointFromString("-2,-2,-2,0")
	Points["H"], _ = quadratic.PointFromString("-2,-2,2,0")
	Points["I"], _ = quadratic.PointFromString("-2,2,-2,2")
	Points["J"], _ = quadratic.PointFromString("2,-2,-2,2")
	Points["K"], _ = quadratic.PointFromString("2,-2,2,-2")
	Points["L"], _ = quadratic.PointFromString("-2,2,2,-2")
	Points["M"], _ = quadratic.PointFromString("6,-2,-2,0")
	Points["N"], _ = quadratic.PointFromString("6,-2,2,0")
	Points["O"], _ = quadratic.PointFromString("2,-2,6,-4")
	Points["P"], _ = quadratic.PointFromString("-6,4,2,-2")
	Points["Q"], _ = quadratic.PointFromString("-2,2,-6,4")
	Points["R"], _ = quadratic.PointFromString("6,-4,-2,2")

	Tiles = make([]string, 1)
	Tiles[0] = "adehnrvuwtspjgbc"
	Tiles = append(Tiles, "Cnhe")
	Tiles = append(Tiles,"beovsi")
	Tiles = append(Tiles,"AaBeCnDvEwFsGjHb")
	Tiles = append(Tiles, "jcehnrvt")
	Tiles = append(Tiles, "jcehmrvt")
	Tiles = append(Tiles,"cdhvtj")
	Tiles = append(Tiles,"dCrh")
	//Tiles = append(Tiles, "leov")
	Tiles = append(Tiles, "ibelvs")
	Tiles = append(Tiles,"bcfdehmrvuqtspkg")
	Tiles = append(Tiles, "bel")
	//Tiles = append(Tiles,"jgkp")
	//Tiles = append(Tiles,"kfmq")
	Tiles = append(Tiles,"AaBeCrpHb")
	Tiles = append(Tiles, "pdCr")
	Tiles = append(Tiles, "jcNmrvt")
	Tiles = append(Tiles, "jcehmMt")
	Tiles = append(Tiles, "bcfderwtKpjg")
	Tiles = append(Tiles, "bcfOKpjg")
	Tiles = append(Tiles, "pOcRhQuP")

	WhiteTiles = make(map[string]bool,len(Tiles))

	for _, t := range Tiles {
		base := TileMap(t, 0)
		TileMaps = append(TileMaps, GenerateOrbits(base)...)
		WhiteTiles[t] = false
	}

	WhiteTiles["Cnhe"] = true
	WhiteTiles["dCrh"] = true
	WhiteTiles["bel"] = true
	WhiteTiles["pdCr"] = true

	VertexFigures = []byte{
		leftRotate(5, 0),
		leftRotate(5, 1),
		leftRotate(5, 2),
		leftRotate(5, 3),
		leftRotate(5, 4),
		leftRotate(5, 5),
		leftRotate(5, 6),
		leftRotate(5, 7),
		leftRotate(9, 0),
		leftRotate(9, 1),
		leftRotate(9, 2),
		leftRotate(9, 3),
		leftRotate(9, 4),
		leftRotate(9, 5),
		leftRotate(9, 6),
		leftRotate(9, 7),
		/*leftRotate(41,0),
		leftRotate(41,1),
		leftRotate(41,2),
		leftRotate(41,3),
		leftRotate(41,4),
		leftRotate(41,5),
		leftRotate(41,6),
		leftRotate(41,7),
		leftRotate(75, 0),
		leftRotate(75, 1),
		leftRotate(75, 2),
		leftRotate(75, 3),
		leftRotate(75, 4),
		leftRotate(75, 5),
		leftRotate(75, 6),
		leftRotate(75, 7),*/
		leftRotate(85, 0),
		leftRotate(85, 1),
		leftRotate(51, 0),
		leftRotate(51, 1),
		leftRotate(51, 2),
		leftRotate(51, 3)}
}
