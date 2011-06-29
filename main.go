package main

import "cairo"
import "http"
import "log"
import "os"

import "runtime"

import "./quadratic/quadratic"
import "./zellij/zellij"

func init() {
	runtime.GOMAXPROCS(3)
	ZellijTilings,reset = zellij.TileSkeleton("0246")
}

var ZellijTilings <-chan *quadratic.Map
var reset chan<- int

func main() {
	http.HandleFunc("/",MainScreen)
	http.HandleFunc("/start",StartTiling)
	http.HandleFunc("/tiles",RenderTiles)
	http.HandleFunc("/previewTiles",DrawTiles)
	http.HandleFunc("/previewSkeleton",DrawSkel)
	err := http.ListenAndServe(":8080",nil)
	if err != nil {
		log.Fatal("ListenAndServe: ",err.String())
	}
}

func MainScreen(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w,req,"ui/ui.html")
}

func RenderTiles(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String()+"\n")
	}	
	
	t := <-ZellijTilings

	/*t := zellij.TileMap(zellij.Tiles[0],0)
	u := zellij.TileMap(zellij.Tiles[5],0)
	
	v,ok := t.Overlay(u,zellij.Overlay)
	u = zellij.TileMap(zellij.Tiles[4],0)
	v,ok = v.Overlay(u,zellij.Overlay)
	u = zellij.TileMap(zellij.Tiles[3],0)
	v,ok = v.Overlay(u,zellij.Overlay)
	u = zellij.TileMap(zellij.Tiles[2],0)
	v,ok = v.Overlay(u,zellij.Overlay)
	u = zellij.TileMap(zellij.Tiles[1],0).Translate(quadratic.NewVertex(zellij.Points["i"]),quadratic.NewVertex(zellij.Points["n"]))
	v,ok = v.Overlay(u,zellij.Overlay)
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	/*v,ok = v.Overlay(t.Translate(quadratic.NewVertex(zellij.Points["s"]),quadratic.NewVertex(zellij.Points["e"])),zellij.Overlay)
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	t = v
	if(zellij.LegalVertexFigures(t)) {
		os.Stderr.WriteString("ok\n")
	}*/
	
	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)
	image.SetLineWidth(.1)
	image.Translate(72*2.,72*2.)
	image.Scale(4.,4.)
	t.ColourFaces(image)
	image.SetSourceRGBA(0.,0.,0.,1.)
	t.DrawEdges(image)

	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func DrawTiles(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String()+"\n")
	}	
	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)
	image.SetLineWidth(.1)
	image.Translate(72*2.,72*2.)
	image.Scale(4.,4.)
	for _,t := range(zellij.TileMaps) {
		t.ColourFaces(image)
	}
	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func DrawSkel(w http.ResponseWriter, req *http.Request) {
	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String()+"\n")
	}	
	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)
	image.SetLineWidth(.1)
	image.Translate(72*2.,72*2.)
	image.Scale(4.,4.)
	skel,ok := zellij.SkeletonMap("0246")
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	skel.DrawEdges(image)
	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func StartTiling(w http.ResponseWriter, req *http.Request) {
	if reset != nil {
		reset <- 1
	}

	ZellijTilings,reset = zellij.TilePlane()
	
	w.WriteHeader(http.StatusOK)
}
	
