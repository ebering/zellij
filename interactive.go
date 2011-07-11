package main

import "cairo"
import "http"
import "log"
import "os"
import "strconv"

import "runtime"

import "./quadratic/quadratic"
import "./zellij/zellij"

func init() {
	runtime.GOMAXPROCS(3)
	zellij.Workers = 1
}

var ZellijTilings <-chan *quadratic.Map
var reset chan int
var CurrentTiling *quadratic.Map

func main() {
	http.HandleFunc("/", MainScreen)
	http.HandleFunc("/start", StartTiling)
	http.HandleFunc("/renderTiles", RenderTiles)
	http.HandleFunc("/nextTiling", NextTiling)
	http.HandleFunc("/previewSkeleton", DrawSkel)
	http.HandleFunc("/emptySvg", EmptySvg)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.String())
	}
}

func MainScreen(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "ui/ui.html")
}

func EmptySvg(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "svg/empty.svg")
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
		<-reset
		reset <- zellij.Workers + 1
	}

	tilingType := req.FormValue("type")

	if tilingType == "skeleton" {
		skeleton := req.FormValue("skeleton")
		showIntermediate, ok := strconv.Atob(req.FormValue("intermediate"))
		if ok != nil || skeleton == "" {
			w.WriteHeader(http.StatusNotFound)
		}
		ZellijTilings, reset = zellij.TileSkeleton(skeleton, showIntermediate)
		w.WriteHeader(http.StatusOK)
		return
	} else if tilingType == "plane" {
		maxtiles, okm := strconv.Atoi(req.FormValue("maxtiles"))
		showIntermediate, oks := strconv.Atob(req.FormValue("intermediate"))
		if okm != nil || oks != nil || maxtiles == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		ZellijTilings, reset = zellij.TilePlane(maxtiles, showIntermediate)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
