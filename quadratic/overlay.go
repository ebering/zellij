package quadratic

import "container/vector"
import "container/heap"
import "os"
import "fmt"
import "sort"

type sweepEvent struct {
	point          *Point
	coincidentEdge *Edge
}

func (e *sweepEvent) Less(f interface{}) bool {
	return e.point.Less(f.(*sweepEvent).point)
}

type sweepStatus struct {
	segments      *vector.Vector
	sweepLocation *Point
}

func (s *sweepStatus) Less(i, j int) bool {
	iL := s.segments.At(i).(*Edge).Line()
	jL := s.segments.At(j).(*Edge).Line()
	return iL.LessAt(jL, s.sweepLocation)
}

func (s *sweepStatus) Len() int {
	return s.segments.Len()
}

func (s *sweepStatus) Swap(i, j int) {
	s.segments.Swap(i, j)
}


func (m *Map) Overlay(n *Map, mergeFaces func(interface{}, interface{}) (interface{}, os.Error)) (*Map, os.Error) {
	o := NewMap()
	CopiedEdges := new(vector.Vector)
	m.Edges.Do(func(l interface{}) {
		e := l.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			f, g := NewEdgePair(NewVertex(e.start.Copy()), NewVertex(e.end.Copy()))
			f.face = e.face
			g.face = e.twin.face
			f.Generation = e.Generation
			g.Generation = e.twin.Generation
			CopiedEdges.Push(f)
			CopiedEdges.Push(g)
		}
	})
	n.Edges.Do(func(l interface{}) {
		e := l.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			f, g := NewEdgePair(NewVertex(e.start.Copy()), NewVertex(e.end.Copy()))
			f.face = e.face
			g.face = e.twin.face
			f.Generation = e.Generation
			g.Generation = e.twin.Generation
			CopiedEdges.Push(f)
			CopiedEdges.Push(g)
		}
	})

	Q := new(vector.Vector)
	T := new(sweepStatus)
	T.segments = new(vector.Vector)

	CopiedEdges.Do(func(l interface{}) {
		e, _ := l.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			evt := new(sweepEvent)
			evt.point = e.start.Point
			evt.coincidentEdge = e
			Q.Push(evt)
			evt = new(sweepEvent)
			evt.point = e.end.Point
			evt.coincidentEdge = nil
			Q.Push(evt)
		}
	})

	heap.Init(Q)

	for Q.Len() > 0 {
		evt, _ := heap.Pop(Q).(*sweepEvent)
		//fmt.Fprintf(os.Stderr,"event: %v\n",evt.point)
		L := new(vector.Vector)
		Lswp := new(sweepStatus)
		Lswp.segments = L
		Lswp.sweepLocation = evt.point
		if evt.coincidentEdge != nil {
			L.Push(evt.coincidentEdge)
		}
		for Q.Len() > 0 && evt.point.Equal(Q.At(0).(*sweepEvent).point) {
			evt, _ := heap.Pop(Q).(*sweepEvent)
			if evt.coincidentEdge != nil {
				L.Push(evt.coincidentEdge)
			}
		}
		sort.Sort(Lswp)
		for i := 0; i < L.Len()-1; {
			if L.At(i).(*Edge).Equal(L.At(i + 1).(*Edge)) {
				L.At(i).(*Edge).newFace = L.At(i + 1).(*Edge).face
				L.At(i).(*Edge).twin.newFace = L.At(i + 1).(*Edge).twin.face
				L.Delete(i + 1)
			} else {
				i++
			}
		}
		//fmt.Fprintf(os.Stderr,"L: %v\n",L)
		R := new(vector.Vector)
		for i := 0; i < T.segments.Len(); {
			e := T.segments.At(i).(*Edge)
			if e.end.Point.Equal(evt.point) {
				R.Push(e)
				T.segments.Delete(i)
			} else if e.Line().On(evt.point) {
				return nil, os.NewError("intersection not at an endpoint")
			} else {
				i++
			}
		}

		// Fill in handle event. You won't need the whole damn thing because
		// Most of the time you just abort with non-terminal intersection
		T.sweepLocation = evt.point
		sort.Sort(T)
		//fmt.Fprintf(os.Stderr,"status: %v\n",T.segments)
		if L.Len() == 0 && R.Len() == 0 {
			return nil, os.NewError("event point with no edges terminal at it " + evt.point.String() + fmt.Sprintf("current status: %v", T.segments))
		} else if L.Len() == 0 {
			above := sort.Search(T.Len(), func(i int) bool {
				return T.segments.At(i).(*Edge).Line().Below(evt.point)
			})
			//fmt.Fprintf(os.Stderr,"Testing status point, no new edge. above: %v, Len: %v\n",above,T.Len())
			if 0 < above && above < T.Len() {
				sa := T.segments.At(above).(*Edge)
				sb := T.segments.At(above - 1).(*Edge)

				cross, _, _ := sb.Line().IntersectAsFloats(sa.Line())
				if cross && !sa.Coterminal(sb) {
					return nil, os.NewError("intersection not at an endpoint")
				}
			}
		} else {
			aboveL := sort.Search(T.Len(), func(i int) bool {
				return L.Last().(*Edge).Line().LessAt(T.segments.At(i).(*Edge).Line(), evt.point)
			})
			belowL := aboveL - 1
			//fmt.Fprintf(os.Stderr,"Testing status point, new edges. above: %v, Len: %v\n",aboveL,T.Len())
			if 0 <= belowL && belowL < T.Len() {
				sa := L.At(0).(*Edge)
				sb := T.segments.At(belowL).(*Edge)

				cross, _, _ := sa.Line().IntersectAsFloats(sb.Line())
				if cross && !sa.Coterminal(sb) {
					return nil, os.NewError("intersection not at an endpoint")
				}
			}
			if aboveL < T.Len() {
				sa := T.segments.At(aboveL).(*Edge)
				sb := L.Last().(*Edge)

				cross, _, _ := sa.Line().IntersectAsFloats(sb.Line())
				if cross && !sa.Coterminal(sb) {
					return nil, os.NewError("intersection not at an endpoint")
				}
			}
		}

		// This is the barrier between preparing the new vertex (below) and determining if the new vertex is good 

		// Setting up edges
		nv := NewVertex(evt.point.Copy())
		R.Do(func(r interface{}) {
			nv.OutgoingEdges.Push(r.(*Edge).twin)
			r.(*Edge).end = nv
			r.(*Edge).twin.start = nv
			o.Edges.Push(r)
			o.Edges.Push(r.(*Edge).twin)
		})
		L.Do(func(l interface{}) {
			l.(*Edge).start = nv
			l.(*Edge).twin.end = nv
			nv.OutgoingEdges.Push(l)
		})
		sort.Sort(nv.OutgoingEdges)

		for i := 0; i < nv.OutgoingEdges.Len(); i++ {
			e := nv.OutgoingEdges.At(i).(*Edge)
			f := nv.OutgoingEdges.At((i + 1) % nv.OutgoingEdges.Len()).(*Edge)
			e.prev = f.twin
			f.twin.next = e
		}

		// Setting up nv's inFace
		// Vertical lines make this shit go tits up. 
		above := sort.Search(T.Len(), func(i int) bool {
			return T.segments.At(i).(*Edge).Line().Below(evt.point)
		})
		if 0 < above && above < T.Len() {
			//fmt.Fprintf(os.Stderr,"Testing status point, looking for vertex in face. above: %v, Len: %v\n",above,T.Len())
			sa := T.segments.At(above).(*Edge)
			sb := T.segments.At(above - 1).(*Edge)
			if sa.twin.face == sb.face {
				onface := false
				nv.OutgoingEdges.Do(func(e interface{}) {
					onface = onface || e.(*Edge).face == sb.face || e.(*Edge).newFace == sb.face
				})
				onface = onface || sa.end.Equal(evt.point) || sb.end.Equal(evt.point)
				if !onface {
					//fmt.Fprintf(os.Stderr,"Vertex %v in face %v from map %p\n",nv,sb.face.Value,)
					nv.inFace = sb.face
				}
			}
		}

		o.Verticies.Push(nv)

		// Vertex done, add any new edges to the sweep line.

		L.Do(func(l interface{}) {
			T.segments.Push(l)
		})
	}

	var leFuck string
	for i := 0; i < o.Edges.Len(); i++ {
		e, _ := o.Edges.At(i).(*Edge)
		if e.visited {
			continue
		}
		//fmt.Fprintf(os.Stderr,"found a face containing: %v,",e.start)

		F := new(Face)
		F.boundary = e
		F.fromMap = o
		e.visited = true
		oldFaces := make(map[*Face]int)
		//fmt.Fprintf(os.Stderr,"%v: ",e.start)
		if e.face != nil {
			//fmt.Fprintf(os.Stderr,"f %v ",e.face.Value)
			oldFaces[e.face]++
		}
		if e.newFace != nil {
			//fmt.Fprintf(os.Stderr,"nf %v ",e.newFace.Value)
			oldFaces[e.newFace]++
		}
		if e.start.inFace != nil {
			//fmt.Fprintf(os.Stderr,"if %v ",e.start.inFace.Value)
			oldFaces[e.start.inFace]++
		}
		if e.face == nil && e.newFace == nil {
			panic("the edge without a face\n")
		}
		e.face = F

		for f := e.Next(); f != e; f = f.Next() {
			//fmt.Fprintf(os.Stderr,"%v: ",f.start)
			f.visited = true
			if f.face != nil {
				//fmt.Fprintf(os.Stderr,"f %v ",f.face.Value)
				oldFaces[f.face]++
			}
			if f.newFace != nil {
				//fmt.Fprintf(os.Stderr,"nf %v ",f.newFace.Value)
				oldFaces[f.newFace]++
			}
			if f.start.inFace != nil {
				//fmt.Fprintf(os.Stderr,"if %v ",f.start.inFace.Value)
				oldFaces[f.start.inFace]++
			}
			//fmt.Fprintf(os.Stderr,",")
			f.face = F
		}
		//os.Stderr.WriteString("\n")

		//fmt.Fprintf(os.Stderr,"%v old faces\n",len(oldFaces))
		if len(oldFaces) > 2 {
			leFuck += fmt.Sprintf("%v faces overlapping a new face, input must have been malformed, maps m: %p n: %p\n", len(oldFaces), m, n)
			for f, _ := range oldFaces {
				leFuck = leFuck + fmt.Sprintf("face %p from: %p containing: %v\n", f, f.fromMap, f.Value)
			}
			//os.Stderr.WriteString(leFuck)
		} else if len(oldFaces) == 0 {
			panic(fmt.Sprintf("No old faces. e: %v, e.face: %+v, maps: m: %p n: %p o: %p\n", e, e.face, m, n, o))
		}

		var mFace, nFace *Face
		for f, _ := range oldFaces {
			if f.fromMap == m {
				mFace = f
			} else {
				nFace = f
			}
		}
		if mFace != nil && nFace != nil {
			v, ok := mergeFaces(mFace.Value, nFace.Value)
			if ok != nil {
				return nil, ok
			}
			if mFace.Type != "" {
				F.Type = mFace.Type
			} else if nFace.Type != "" {
				F.Type = nFace.Type
			}
			F.Value = v
		} else if mFace != nil {
			F.Value = mFace.Value
			F.Type = mFace.Type
		} else if nFace != nil {
			F.Value = nFace.Value
			F.Type = nFace.Type
		} else {
			panic(fmt.Sprintf("face didn't come from an mFace or an nFace, pointers m: %v n: %v o: %v face: %v", m, n, o, e.face))
		}

		o.Faces.Push(F)
	}

	if leFuck != "" {
		return o, os.NewError(leFuck)
	}

	o.Init()
	return o, nil
}
