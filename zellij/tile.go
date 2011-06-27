package zellij

import "../quadratic/quadratic"

func GenerateOrbits(tile *quadratic.Map) []*quadratic.Map {
	orbits := make([]*quadratic.Map,1)
	orbits[0] = tile
	for i:=1; i <= 8; i=i*2 {
		t := tile.Copy().RotatePi4(i)
		if( t.Isomorphic(tile)) {
			for k:=1; k < i; k++ {
				orbits = append(orbits,tile.Copy().RotatePi4(k))
			}
			break
		}
	}
	return orbits
}	
