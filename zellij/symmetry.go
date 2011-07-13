package zellij

import "../quadratic/quadratic"
import "fmt"

func GenerateOrbit(tile *quadratic.Map) []*quadratic.Map {
	orbit := make([]*quadratic.Map, 1)
	orbit[0] = tile
	for i := 1; i <= 8; i = i * 2 {
		t := tile.Copy().RotatePi4(i)
		if t.Isomorphic(tile) {
			for k := 1; k < i; k++ {
				orbits = append(orbits, tile.Copy().RotatePi4(k))
			}
			break
		}
	}
	return orbits
}

func DetectSymmetryGroup(shape *quadratic.Map) string {
	centroid := shape.Centroid()
	shape.Translate(quadratic.NewVertex(centroid),quadratic.NewVertex(new quadratic.Point(quadratic.Zero,quadratic.Zero)))
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
		return fmt.Sprintf("d%v",8/i)
	} else {
		return fmt.Sptintf("c%v",8/i)
	}

	return "e"
}
