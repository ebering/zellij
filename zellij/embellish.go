package zellij

import "../quadratic/quadratic"
import "os"

func Embellish(frame *quadratic.Map,motifs *Database) (*quadratic.Map,os.Error) {
	embellishments := make([]*quadratic.Map,frame.Faces.Len())

	i := 0;
	frame.Faces.Do(func (f interface{}) {
		F := f.(*quadratic.Face)
		Fr := F.Frame()
		e,iso,ok :=  motifs.Matching(Fr)
		if ok != nil && !F.Inner() {
			return nil,os.NewError("no motif matching a frame")
		}
		embellishments[i] = iso(e.Motifs[0].Copy()).Translate(iso(e.Frame).Verticies.At(0).(*quadratic.Vertex).Point,
						Fr.Verticies.At(0).(*quadratic.Vertex).Point)

	}

	ret := embellishments[0];

	for i :=1 ; i < len(embellishments); i++ {
		ret,_ = ret.Overlay(embellishments[i])
	}

	return ret,nil
}
