package quadratic

import "container/list"
import "container/vector"
import "sort"
import "os"
import "math"
import "fmt"

// Stores a vertex, which is a point that has a list of incident edges
// IMPORTANT INVARIANT: The vector of outgoing edges is a vector of things of type *Edge and it is KEPT IN SORTED ORDER BY HEADING
// DON'T COCK THAT UP.
type Vertex struct {
	*Point
	outgoingEdges *vector.Vector
	copy *Vertex //used only in the copy routine
}

func NewVertex(p *Point) (* Vertex) {
	nv := new(Vertex)
	nv.outgoingEdges = new(vector.Vector)
	nv.Point = p.Copy()
	return nv
}

func (v * Vertex) EdgeIndexCounterClockwiseOf(e *Edge) int {
	return sort.Search(v.outgoingEdges.Len(), func (i int) bool {
		return e.Less(v.outgoingEdges.At(i))
	})
}
	

// Stores a directed half edge in a planar map
type Edge struct {
	start,end *Vertex
	next,prev,twin *Edge
	face *Face
	newFace *Face
}

func (e * Edge) Coterminal(f * Edge) bool {
	return e.start.Point.Equal(f.start.Point) ||
		e.start.Point.Equal(f.end.Point) ||
		e.end.Point.Equal(f.start.Point) ||
		e.end.Point.Equal(f.end.Point)
}

func (e * Edge) Line() (*Line){
	return NewLine(e.start.Point,e.end.Point)
}

func (e * Edge) Heading() float64 {
	return math.Atan2(e.end.y.Sub(e.start.y).Float64(),e.end.x.Sub(e.start.x).Float64()) 
}

func (e * Edge) Less(v interface{}) bool {
	f := v.(*Edge)
	if e.start.Equal(f.start.Point) && e.end.Equal(f.end.Point) {
		return false
	}
	
	return e.Heading() < f.Heading()
}

func (e * Edge) Parallel(f * Edge) (bool) {
	return math.Fabs(e.Heading()-f.Heading()) < FLOAT64_EPSILON || math.Fabs(e.twin.Heading()-f.Heading()) < FLOAT64_EPSILON 
}

// Represents a face
type Face struct {
	boundary *Edge
	Value interface{}
}

// Represents a planar map
type Map struct {
	verticies, edges, faces *list.List
}

func NewMap() (* Map) {
	return &Map{new(list.List),new(list.List),new(list.List)}
}


// Given two verticies returns a twin pair of half edges between them,
// the first goes from start to end, the second is the twin. 
// Correctly sets next and prev based on the other incident edges, and adds them
// to the relevant set of pointers in the vertex.
func NewEdgePair(start,end *Vertex) (*Edge,*Edge) {
	e,f := new(Edge),new(Edge)
	e.start,e.end = start,end
	f.start,f.end = e.end,e.start

	if start.outgoingEdges.Len() == 0 {
		start.outgoingEdges.Push(e)
		start.outgoingEdges.Push(f)
		sort.Sort(start.outgoingEdges)
	} else {
		i := start.EdgeIndexCounterClockwiseOf(e)
		n := start.outgoingEdges.Len()
		e.prev = start.outgoingEdges.At(i %n ).(*Edge).twin
		f.next = start.outgoingEdges.At( (i-1+n) % n ).(*Edge)
		start.outgoingEdges.Insert(i,e)
	}

	if end.outgoingEdges.Len() == 0 {
		end.outgoingEdges.Push(e)
		end.outgoingEdges.Push(f)
		sort.Sort(end.outgoingEdges)
	} else {
		i := end.EdgeIndexCounterClockwiseOf(f)
		n := end.outgoingEdges.Len()
		f.prev = end.outgoingEdges.At(i %n).(*Edge).twin
		e.next = end.outgoingEdges.At( (i-1+n) %n ).(*Edge)
		end.outgoingEdges.Insert(i,f)
	}
	
	e.twin,f.twin = f,e

	return e,f
}

// Next returns the next edge along a face
func (e *Edge) Next() (*Edge) {
	return e.next
}

// Joins f to e as the next of e, updating next and prev pointers for both edges. 
// If this operation would not make sense (i.e. e and f are not at the same vertex)
// it does nothing.
func (e *Edge) JoinNext(f *Edge) {
	if e.end == f.start {
		e.next = f
		f.prev = e
	}
}

// Returns the previous edge along a face
func (e *Edge) Prev() (*Edge) {
	return e.prev
}

// Joins f to e prior to e instead of after it. Also does not act if this
// operation does not make sense.
func (e *Edge) JoinPrev(f *Edge) {
	if f.end == e.start {
		e.prev = f
		f.next = e
	}
}

// Returns the twin of an edge. Note that twins cannot be set.
func (e *Edge) Twin() (*Edge){
	return e.twin
}

func (m *Map) AddVertex(p * Point) (* Vertex) {
	L := m.verticies
	for l := L.Front(); l != nil; l = l.Next() {
		if( l.Value.(*Vertex).Point.Equal(p)) {
			return l.Value.(*Vertex)
		}
	}
	nv := NewVertex(p)
	_ = m.verticies.PushFront(nv)
	return nv
}	

func (m *Map) JoinVerticies(u * Vertex, v * Vertex) (* Edge) {
	e,eTwin := NewEdgePair(u,v)
	_ = m.edges.PushFront(e)
	_ = m.edges.PushFront(eTwin)
	return e
}

// Adds an edge e to a planar map m
// Does not alter e
// Does not handle intersections or coincident edges that do not share endpoints at present
// Does not update face information
func (m *Map) AddEdge(e *Edge) (os.Error) {
	if m.edges.Front() == nil {
		_ = m.edges.PushFront(e)
		_ = m.edges.PushFront(e.twin)
		_ = m.verticies.PushFront(e.start)
		_ = m.verticies.PushFront(e.end)
		return nil
	}
		
	L := m.edges
	
	// Check for edge overlap, coincidence, and intersection
	for l := L.Front(); l != nil; l = l.Next() {
		f,ok := l.Value.(*Edge)
		if !ok {
			panic("Something is terribly wrong with the planar map")
		}
		if e.start.Point.Equal(f.start.Point) && e.end.Point.Equal(f.end.Point) {
			f.newFace = e.face
			f.twin.newFace = e.twin.face
			return nil
		}
		if e.start.Point.Equal(f.end.Point) && e.end.Point.Equal(f.start.Point) {
			f.newFace = e.twin.face
			f.twin.newFace = e.face
			return nil
		}
		fLine := f.Line()
		if (fLine.On(e.start.Point) || fLine.On(e.end.Point)) && e.Parallel(f) && !e.Coterminal(f) {
			return os.NewError(fmt.Sprintf("Coincident non-coterminal edges %v, %v not supported",e,f))
		}
		cross,_,_ := fLine.IntersectAsFloats(e.Line()) 
		if cross && !( e.start.Point.Equal(f.start.Point) || e.end.Point.Equal(f.start.Point) || e.start.Point.Equal(f.end.Point) || e.end.Point.Equal(f.end.Point)) {
			return os.NewError("Edges intersecting not at endpoints not supported")
		}
	}

	for l := L.Front(); l != nil; l = l.Next() {
		f,ok := l.Value.(*Edge)
		if !ok {
			return os.NewError("Something is terribly wrong with your planar map")
		}
		if e.start.Point.Equal(f.start.Point) {
			nv := m.AddVertex(e.end.Point)
			a := m.JoinVerticies(f.start,nv)
			a.newFace = e.face
			a.twin.newFace = e.twin.face
			return nil
		}
		if e.end.Point.Equal(f.start.Point) {
			nv := m.AddVertex(e.start.Point)
			a := m.JoinVerticies(f.start,nv)
			a.newFace = e.face
			a.twin.newFace = e.twin.face
			return nil
		}
	}

	u := m.AddVertex(e.start.Point)
	v := m.AddVertex(e.end.Point)
	f := m.JoinVerticies(u,v)
	f.newFace = e.face
	f.twin.newFace = e.twin.face
	return nil
}

// Merge n into m
func (m *Map) Merge(n * Map) (os.Error) {
	newEdges := n.edges

	for l := newEdges.Front(); l != nil; l = l.Next() {
		e,_ := l.Value.(*Edge)
		if e.start.Less(e.end.Point) {
			ok := m.AddEdge(e)
			if ok != nil {
				return ok
			}
		}
	}

	for l := m.edges.Front(); l != nil; l = l.Next() {
		e,_ := l.Value.(*Edge)
		if e.newFace != nil {
			if e.face == nil {
				e.face = e.newFace
			} else if e.face.Value.(string) == e.newFace.Value.(string) && e.face.Value.(string) == "inner" {
				return os.NewError("Overlapping inner faces not allowed")
			} else if e.newFace.Value.(string) == "inner" {
				e.face = e.newFace
			}
			e.newFace = nil
		}
	}

	return nil
}

func (m *Map) Copy() (* Map) {
	c := NewMap()
	for l := m.verticies.Front(); l != nil; l = l.Next() {
		cv := NewVertex(l.Value.(*Vertex).Copy())
		l.Value.(*Vertex).copy = cv
		c.verticies.PushBack(cv)
	}
	for l := m.edges.Front(); l != nil; l = l.Next() {
		e := l.Value.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			c.JoinVerticies(e.start.copy,e.end.copy)
		}
	}
	return c
}

func (m *Map) DoVerticies(f func (*Vertex)) {
	L := m.verticies
	for l := L.Front(); l != nil; l = l.Next() {
		f(l.Value.(*Vertex))
	}
}

func (m *Map) DoEdges(f func (*Edge)) {
	L := m.edges
	for l := L.Front(); l != nil; l = l.Next() {
		f(l.Value.(*Edge))
	}
}

func (m *Map) DoFaces(f func (*Face)) {
	L := m.faces
	for l := L.Front(); l != nil; l = l.Next() {
		f(l.Value.(*Face))
	}
}

func (m *Map) Translate(from,to *Vertex) (*Map) {
	T := MakeTranslation(from.Point,to.Point)
	m.DoVerticies(func (v *Vertex) {
		v.Point = T(v.Point)
	})
	return m
}
