DEPS="zellij.go init.go generations.go vertex.go"
redo-ifchange $DEPS
6g -o $3 $DEPS 1>&2
