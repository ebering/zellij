package zellij

import "../quadratic/quadratic"
import "os"
import "fmt"

func Embellish(frame *quadratic.Map,motifs Database) (*quadratic.Map,os.Error) {
	embellishments := make([]*quadratic.Map,frame.Faces.Len()-1)
	fmt.Fprintf(os.Stderr,"embellishing %v frames\n",frame.Faces.Len()-1)

	i := 0;
	abort := false
	frame.Faces.Do(func (f interface{}) {
		F := f.(*quadratic.Face)
		Fr := F.Frame()
		e,iso,ok :=  motifs.Matching(Fr)

		if ok != nil && F.Inner() {
			abort = true
			return
		}

		if F.Inner() {
			os.Stderr.WriteString("got an embellishment\n")
			embellishments[i] = iso(e.Variations[0].Copy()).Translate(iso(e.Frame.Copy()).Verticies.At(0).(*quadratic.Vertex),
								       Fr.Verticies.At(0).(*quadratic.Vertex))
			i++
		}

	})

	if abort {
		return nil,os.NewError("cannot match a frame")
	}

	ret := embellishments[0];

	for i :=1 ; i < len(embellishments); i++ {
		var ok os.Error
		ret,ok = ret.Overlay(embellishments[i],EmbellishmentOverlay)
		if ok != nil {
			os.Stderr.WriteString("embellish overlay failure:" +ok.String()+"\n")
		}
	}

	return ret,nil
}
