package main

import "cairo"
import "http"
import "log"
import "os"
import "strconv"
import "strings"
import "fmt"
import "math"

import "runtime"

import "./quadratic/quadratic"
import "./zellij/zellij"

import "json"

func init() {
	runtime.GOMAXPROCS(4)
	zellij.Workers = 2
}

var ZellijTilings <-chan *quadratic.Map
var reset chan int
var CurrentTiling *quadratic.Map
var MotifDatabase zellij.Database

func main() {
	http.HandleFunc("/", MainScreen)
	http.HandleFunc("/start", StartTiling)
	http.HandleFunc("/renderTiles", RenderTiles)
	http.HandleFunc("/nextTiling", NextTiling)
	http.HandleFunc("/previewSkeleton", DrawSkel)
	http.HandleFunc("/emptySvg", EmptySvg)
	http.HandleFunc("/embellish",EmbellishFrame)
	MotifDatabase = zellij.LoadDatabase("tiles")
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

func MainScreen(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "ui/skeleton.html")
}

func EmptySvg(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "svg/empty")
}

func NextTiling(w http.ResponseWriter, req *http.Request) {
	CurrentTiling = <-ZellijTilings
	RenderTiles(w, req)
}

func RenderTiles(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String() + "\n")
	}

	if CurrentTiling == nil {
		EmptySvg(w, req)
		return
	}

	style := req.FormValue("style")

	image := cairo.NewSurface("svg/test-surface.svg", 72*4, 72*4)
	image.SetSourceRGB(0., 0., 0.)
	image.SetLineWidth(.1)
	image.Translate(72*2., 72*2.)
	image.Scale(4., 4.)
	if style == "edges" {
		image.SetSourceRGBA(0., 0., 0., 1.)
		CurrentTiling.DrawEdges(image)
	} else if style == "plain" {
		CurrentTiling.ColourFaces(image, zellij.PlainBrush)
	} else {
		CurrentTiling.ColourDebugFaces(image)
		CurrentTiling.DrawDebugEdges(image)
	}

	image.Finish()
	http.ServeFile(w, req, "svg/test-surface.svg")
}

func EmbellishFrame(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/embellishment.svg")
	if e != nil {
		os.Stderr.WriteString(e.String() + "\n")
	}
	frame := quadratic.NewMap()
	jsonErr := json.NewDecoder(strings.NewReader(req.FormValue("frame"))).Decode(frame)
	if jsonErr != nil {
		os.Stderr.WriteString("json error: " + jsonErr.String() + "\n")
		return
	}
	zig,err := zellij.Embellish(frame,MotifDatabase)
	if err != nil {
		os.Stderr.WriteString(err.String() + "\n")
		return
	}
	fmt.Fprintf(os.Stderr,"embellishment has %v faces\n", zig.Faces.Len())
	bx := zig.Verticies.At(0).(*quadratic.Vertex).Point
	tx := zig.Verticies.At(0).(*quadratic.Vertex).Point
	by := zig.Verticies.At(0).(*quadratic.Vertex).Point
	ty := zig.Verticies.At(0).(*quadratic.Vertex).Point
	for i:= 1; i < zig.Verticies.Len(); i++ {
		v := zig.Verticies.At(i).(*quadratic.Vertex).Point
		if v.Y().Float64() < by.Y().Float64() {
			by = v
		} else if v.Y().Float64() > ty.Y().Float64() {
			ty = v
		}
		if v.X().Float64() < bx.X().Float64() {
			bx = v
		} else if v.X().Float64() > tx.X().Float64() {
			tx = v
		}
	}
	xwidth := math.Floor(tx.X().Float64()-bx.X().Float64()+20.0)
	ywidth := math.Floor(ty.Y().Float64()-by.Y().Float64()+20.0)

	image := cairo.NewSurface("svg/embellishment.svg", xwidth*4, ywidth*4)
	image.SetSourceRGB(0., 0., 0.)
	image.SetLineWidth(.1)
	image.Translate(-bx.X().Float64()*4.+40, -by.Y().Float64()*4.+40)
	image.Scale(4., 4.)
	zig.ColourFaces(image,zellij.CreateZellijBrush(zellij.Colour{1.,0.,0.,1.},zellij.Colour{0.,1.,0.,1.},zellij.Colour{0.,0.,1.,1.}))
	image.Finish()
	http.ServeFile(w, req, "svg/embellishment.svg")
}
	

func DrawSkel(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String() + "\n")
	}
	skeleton := req.FormValue("skeleton")
	image := cairo.NewSurface("svg/test-surface.svg", 72*4, 72*4)
	image.SetSourceRGB(0., 0., 0.)
	image.SetLineWidth(.1)
	image.Translate(72*2., 72*2.)
	image.Scale(4., 4.)
	skel, ok := zellij.SkeletonMap(skeleton)
	if ok != nil {
		os.Stderr.WriteString(ok.String() + "\n")
	}
	skel.DrawEdges(image)
	image.Finish()
	http.ServeFile(w, req, "svg/test-surface.svg")
}

func StartTiling(w http.ResponseWriter, req *http.Request) {
	if reset != nil {
		select {
			case <-reset:
				reset <- 1
			case reset <- 1:
		}
	}

	tilingType := req.FormValue("type")
	tileSymmetry := req.FormValue("symmetry")

	if tilingType == "skeleton" {
		skeleton := req.FormValue("skeleton")
		showIntermediate, ok := strconv.Atob(req.FormValue("intermediate"))
		if ok != nil || skeleton == "" {
			w.WriteHeader(http.StatusNotFound)
		}
		ZellijTilings, reset = zellij.TileSkeleton(skeleton, tileSymmetry, showIntermediate)
		w.WriteHeader(http.StatusOK)
		return
	} else if tilingType == "plane" {
		maxtiles, okm := strconv.Atoi(req.FormValue("maxtiles"))
		showIntermediate, oks := strconv.Atob(req.FormValue("intermediate"))
		if okm != nil || oks != nil || maxtiles == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		ZellijTilings, reset = zellij.TilePlane(maxtiles, tileSymmetry, showIntermediate)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
