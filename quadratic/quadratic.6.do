DEPS="qint.go qpoint.go qline.go draw.go string.go qmap.go qpoly.go"
redo-ifchange $DEPS
6g -o $3 $DEPS 1>&2
