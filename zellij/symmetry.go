package zellij

import "../quadratic/quadratic"
import "fmt"
import "strconv"
import "strings"

// GenerateOrbits is the old way we dealt with tile symmetries.
func GenerateOrbits(tile *quadratic.Map) []*quadratic.Map {
	orbit := make([]*quadratic.Map, 1)
	orbit[0] = tile
	for i := 1; i <= 8; i = i * 2 {
		t := tile.Copy().RotatePi4(i)
		if t.Isomorphic(tile) {
			for k := 1; k < i; k++ {
				orbit = append(orbit, tile.Copy().RotatePi4(k))
			}
			break
		}
	}
	return orbit
}

func GenerateOrbit(shape *quadratic.Map, symmetryGroup string) []*quadratic.Map {
	orbit := make([]*quadratic.Map, 1)
	orbit[0] = shape
	groupSymbol := strings.Split(symmetryGroup, "")
	if symmetryGroup == "e" {
		return orbit
	}
	rotOrder, _ := strconv.Atoi(groupSymbol[1])
	for i := 0; i < rotOrder; i++ {
		s := shape.Copy().RotatePi4(i * (8 / rotOrder))
		if !duplicateShape(orbit, s) {
			orbit = append(orbit, s)
		}
		if groupSymbol[0] == "d" {
			s = s.Copy().ReflectXAxis()
			if !duplicateShape(orbit, s) {
				orbit = append(orbit, s)
			}
		}
	}
	return orbit
}

func duplicateShape(shapes []*quadratic.Map, s *quadratic.Map) bool {
	ret := false
	for _, t := range shapes {
		ret = ret || t.Equal(s)
	}
	return ret
}

func DetectSymmetryGroup(shape *quadratic.Map) string {
	centroid := shape.Centroid()
	shape.Translate(quadratic.NewVertex(centroid), quadratic.NewVertex(quadratic.NewPoint(quadratic.Zero, quadratic.Zero)))
	var i int
	for i = 1; i <= 8; i = i * 2 {
		s := shape.Copy().RotatePi4(i)
		if s.Isomorphic(shape) {
			break
		}
	}
	if i == 8 {
		// TODO: Check if we're d1
		return "e"
	}
	s := shape.Copy().ReflectXAxis()
	if s.Isomorphic(shape) {
		return fmt.Sprintf("d%v", 8/i)
	} else {
		return fmt.Sprintf("c%v", 8/i)
	}

	return "e"
}
