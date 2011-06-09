package quadratic

import "container/vector"
import "sort"
import "math"

// Stores a vertex, which is a point that has a list of incident edges
// IMPORTANT INVARIANT: The vector of outgoing edges is a vector of things of type *Edge and it is KEPT IN SORTED ORDER BY HEADING
// DON'T COCK THAT UP.
type Vertex struct {
	*Point
	outgoingEdges *vector.Vector
	copy *Vertex //used only in the copy routine
	inFace *Face
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

func (v *Vertex) Less(u interface{}) bool {
	return v.Point.Less(u.(*Vertex).Point)
}

// Stores a directed half edge in a planar map
type Edge struct {
	start,end *Vertex
	next,prev,twin *Edge
	face *Face
	newFace *Face
	visited bool
}

func (e * Edge) Coterminal(f * Edge) bool {
	return e.start.Point.Equal(f.start.Point) ||
		e.start.Point.Equal(f.end.Point) ||
		e.end.Point.Equal(f.start.Point) ||
		e.end.Point.Equal(f.end.Point)
}

func (e * Edge) Equal( f * Edge) bool {
	return e.start.Point.Equal(f.start.Point) && e.end.Point.Equal(f.end.Point)
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
	Verticies, Edges, Faces *vector.Vector
}

func NewMap() (* Map) {
	return &Map{new(vector.Vector),new(vector.Vector),new(vector.Vector)}
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
	V := m.Verticies
	for i :=0; i < V.Len(); i++ {
		if( V.At(i).(*Vertex).Point.Equal(p)) {
			return V.At(i).(*Vertex)
		}
	}
	nv := NewVertex(p)
	m.Verticies.Push(nv)
	return nv
}	

func (m *Map) JoinVerticies(u * Vertex, v * Vertex) (* Edge) {
	e,eTwin := NewEdgePair(u,v)
	m.Edges.Push(e)
	m.Edges.Push(eTwin)
	return e
}

func (m *Map) Copy() (* Map) {
	c := NewMap()
	m.Verticies.Do(func (v interface{}) {
		cv := NewVertex(v.(*Vertex).Copy())
		v.(*Vertex).copy = cv
		c.Verticies.Push(cv)
	})
	m.Edges.Do(func (f interface{}) {
		e := f.(*Edge)
		if e.start.Point.Less(e.end.Point) {
			c.JoinVerticies(e.start.copy,e.end.copy)
		}
	})
	return c
}

func (m *Map) Translate(from,to *Vertex) (*Map) {
	T := MakeTranslation(from.Point,to.Point)
	m.Verticies.Do(func (v interface{}) {
		v.(*Vertex).Point = T(v.(*Vertex).Point)
	})
	return m
}

func (m *Map) AdjacencyMatrix() [][]bool {
	mat := make([][]bool,m.Verticies.Len())
	for i,_ := range(mat) {
		mat[i] = make([]bool,m.Verticies.Len())
	}
	sort.Sort(m.Verticies)
	for i:=0; i < m.Verticies.Len(); i++ {
		u := m.Verticies.At(i).(*Vertex)
		for j:=0; j < m.Verticies.Len(); j++ {
			if i == j { continue }
			v := m.Verticies.At(j).(*Vertex)
			u.outgoingEdges.Do(func (e interface{}) {
				mat[i][j] = mat[i][j] || v.Equal(e.(*Edge).end.Point)
			})
		}
	}
	return mat
}

func (m *Map) Isomorphic(n *Map) bool {
	if m.Verticies.Len() != n.Verticies.Len() || m.Edges.Len() != n.Edges.Len() { return false }
	mA := m.AdjacencyMatrix()
	nA := n.AdjacencyMatrix()
	for i,r := range(mA) {
		for j,v := range(r) {
			if v != nA[i][j] { return false }
		}
	}
	return true
}
