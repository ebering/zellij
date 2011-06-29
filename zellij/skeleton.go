package zellij

import "strconv"
import "os"

func SkeletonMap(spec string) (*quadratic.Map,os.Error) {
	ret := TileMap(Tiles[0],0)
	currentPoint := quadratic.PointMustFromString("0,0,0,0")
	origin := quadratic.PointMustFromString("0,0,0,0")
	for _,d := range(spec) {
		translatePoint := quadratic.PointMustFromString("4,4,0,0")
		translateSaft := quadratic.PointMustFromString("0,2,0,0")
		saft := TileMap(Tiles[0],0)
		seal := TileMap(Tiles[0],0)
		heading := strconv.Atoi(d)

		translatePoint.RotatePi4(heading)
		translateSaft.RotatePi4(heading)
		saft.RotatePi4(heading)
		saft.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(translateSaft))
		saft.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(currentPoint))
	
		ret,ok = ret.Overlay(saft)
		if ok != nil {
			return nil,ok
		}

		
		currentPoint = MakeTranslation(origin,translatePoint)(currentPoint)
		seal.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(currentPoint))

		ret,ok = ret.Overlay(seal)
		if ok != nil {
			return nil,ok
		}
	}

	return ret,nil
}
