DEPS="int.go point.go line.go draw.go string.go map.go poly.go overlay.go json.go"
redo-ifchange $DEPS
6g -o $3 $DEPS 1>&2
