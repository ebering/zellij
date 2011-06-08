package quadratic

import "container/vector"
import "container/heap"
import "os"
import "fmt"
import "sort"

type sweepEvent struct {
	point *Point
	coincidentEdge *Edge
}

func (e *sweepEvent) Less(f interface{}) bool {
	return e.point.Less(f.(*sweepEvent).point)
}

type sweepStatus struct {
	segments *vector.Vector
	sweepLocation *Point
}

func (s *sweepStatus) Less(i,j int) bool {
	iL := s.segments.At(i).(*Edge).Line()
	jL := s.segments.At(j).(*Edge).Line()
	return iL.LessAt(jL,s.sweepLocation)
}

func (s *sweepStatus) Len() int {
	return s.segments.Len()
}

func (s *sweepStatus) Swap(i,j int) {
	s.segments.Swap(i,j)
}
	

func (m *Map) Overlay(n * Map) (*Map,os.Error) {
	o := NewMap()
	m.Edges.Do(func(l interface{}) {
		e := l.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			f,g := NewEdgePair(NewVertex(e.start.Copy()),NewVertex(e.end.Copy()))
			o.Edges.Push(f)
			o.Edges.Push(g)
		}
	})
	n.Edges.Do(func(l interface{}) {
		e := l.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			f,g := NewEdgePair(NewVertex(e.start.Copy()),NewVertex(e.end.Copy()))
			o.Edges.Push(f)
			o.Edges.Push(g)
		}
	})

	Q := new(vector.Vector)
	T := new(sweepStatus)
	T.segments = new(vector.Vector)

	o.Edges.Do(func(l interface{}) {
		e,_ := l.(*Edge)
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
		evt,_ := heap.Pop(Q).(*sweepEvent)
		//fmt.Fprintf(os.Stderr,"event: %v\n",evt.point)
		L := new(vector.Vector)
		Lswp := new(sweepStatus)
		Lswp.segments = L
		Lswp.sweepLocation = evt.point
		if evt.coincidentEdge != nil {
			L.Push(evt.coincidentEdge)
		}
		for Q.Len() > 0 && evt.point.Equal(Q.At(0).(*sweepEvent).point) {
			evt,_ := heap.Pop(Q).(*sweepEvent)
			if evt.coincidentEdge != nil {
				L.Push(evt.coincidentEdge)
			}
		}
		sort.Sort(Lswp)
		for i := 0; i < L.Len()-1;  {
			if L.At(i).(*Edge).Equal(L.At(i+1).(*Edge)) {
				L.Delete(i+1)
			} else {
				i++
			}
		}
		//fmt.Fprintf(os.Stderr,"L: %v\n",L)
		R := new(vector.Vector)
		for i:=0; i < T.segments.Len(); {
			e := T.segments.At(i).(*Edge)
			if e.end.Point.Equal(evt.point) {
				R.Push(e)
				T.segments.Delete(i)
			} else if e.Line().On(evt.point) {
				return nil,os.NewError("intersection not at an endpoint")
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
			return nil,os.NewError("event point with no edges terminal at it "+evt.point.String()+fmt.Sprintf("current status: %v",T.segments))
		} else if L.Len() == 0 {
			above := sort.Search(T.Len(),func (i int) bool {
				return T.segments.At(i).(*Edge).Line().Below(evt.point)
			})
			//fmt.Fprintf(os.Stderr,"Testing status point, no new edge. above: %v, Len: %v\n",above,T.Len())
			if 0 < above && above < T.Len() {
				sa := T.segments.At(above).(*Edge)
				sb := T.segments.At(above-1).(*Edge)

				cross,_,_ := sb.Line().IntersectAsFloats(sa.Line())
				if cross && !sa.Coterminal(sb) {
					return nil,os.NewError("intersection not at an endpoint")
				}
			}
		} else {
			aboveL := sort.Search(T.Len(),func (i int) bool {
				return L.Last().(*Edge).Line().LessAt(T.segments.At(i).(*Edge).Line(),evt.point)
			})
			belowL := aboveL - 1
			//fmt.Fprintf(os.Stderr,"Testing status point, new edges. above: %v, Len: %v\n",aboveL,T.Len())
			if 0 <= belowL && belowL < T.Len() {
				sa := L.At(0).(*Edge)
				sb := T.segments.At(belowL).(*Edge)

				cross,_,_ := sa.Line().IntersectAsFloats(sb.Line())
				if cross && !sa.Coterminal(sb) {
					return nil,os.NewError("intersection not at an endpoint")
				}
			}
			if aboveL < T.Len() {
				sa := T.segments.At(aboveL).(*Edge)
				sb := L.Last().(*Edge)

				cross,_,_ := sa.Line().IntersectAsFloats(sb.Line())
				if cross && !sa.Coterminal(sb) {
					return nil,os.NewError("intersection not at an endpoint")
				}
			}
		}
		L.Do(func(l interface{}) {
			T.segments.Push(l)
		})

		// This is the barrier between preparing the new vertex (below) and determining if the new vertex is good and updating the sweep (above)

		R.Do(func(r interface{}) {
			L.Push(r.(*Edge).twin)
		})
		nv := NewVertex(evt.point.Copy())
		nv.outgoingEdges = L
		L.Do(func(l interface{}) {
			l.(*Edge).start = nv
			l.(*Edge).twin.end = nv
		})
		sort.Sort(nv.outgoingEdges)
		o.Verticies.Push(nv)
	}

	return o,nil
}
