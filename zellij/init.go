package zellij

import "../quadratic/quadratic"

var Points map[string]*quadratic.Point
var Tiles []string

func init() {
	Points = make(map[string]*quadratic.Point)
	Points["a"], _ = quadratic.PointFromString("0,0,0,2")
	Points["b"], _ = quadratic.PointFromString("-2,0,2,0")
	Points["c"], _ = quadratic.PointFromString("2,-2,2,0")
	Points["d"], _ = quadratic.PointFromString("-2,2,2,0")
	Points["e"], _ = quadratic.PointFromString("2,0,2,0")

	Points["g"], _ = quadratic.PointFromString("-2,0,-2,2")
	Points["h"], _ = quadratic.PointFromString("2,0,-2,2")
	Points["i"], _ = quadratic.PointFromString("-4,0,0,0")
	Points["j"], _ = quadratic.PointFromString("0,-2,0,0")

	Points["l"], _ = quadratic.PointFromString("0,0,0,0")

	Points["n"], _ = quadratic.PointFromString("0,2,0,0")
	Points["o"], _ = quadratic.PointFromString("4,0,0,0")
	Points["p"], _ = quadratic.PointFromString("-2,0,2,-2")

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

	Tiles = make([]string,6 )
	Tiles[0] = "adehnrvuwtspjgbc"
	Tiles[1] = "beovsi"
	Tiles[2] = "Abgj"
	Tiles[3] = "Cjps"
	Tiles[4] = "Dvrn"
	Tiles[5] = "Bnhe"
}
