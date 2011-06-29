package zellij

import "strconv"
import "os"
import "../quadratic/quadratic"
import "strings"

func SkeletonMap(spec string) (*quadratic.Map,os.Error) {
	var ret *quadratic.Map
	currentPoint := quadratic.PointMustFromString("0,0,0,0")
	origin := quadratic.PointMustFromString("0,0,0,0")
	for _,d := range(strings.Split(spec,"",-1)) {
		translatePoint := quadratic.PointMustFromString("8,4,0,0")
		translateSaft := quadratic.PointMustFromString("4,2,0,0")
		saft := TileMap(Tiles[1],0)
		seal := TileMap(Tiles[0],0)
		heading,_ := strconv.Atoi(d)

		seal.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(currentPoint))
		if ret == nil {
			ret = seal
		} else {
			r,ok := ret.Overlay(seal,Overlay)
			if ok != nil {
				return nil,ok
			}
			ret = r
		}

		translatePoint.RotatePi4(heading)
		translateSaft.RotatePi4(heading)
		saft.RotatePi4(heading)
		saft.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(translateSaft))
		saft.Translate(quadratic.NewVertex(origin),quadratic.NewVertex(currentPoint))
	
		r,ok := ret.Overlay(saft,Overlay)
		ret = r
		if ok != nil {
			return nil,ok
		}

		currentPoint = quadratic.MakeTranslation(origin,translatePoint)(currentPoint)
	}

	ret.Faces.Do(func (f interface{}) {
		if f.(*quadratic.Face).Value.(string) == "outer" && f.(*quadratic.Face).Inner() {
			f.(*quadratic.Face).Value = "active"
		}
	})

	return ret,nil
}
