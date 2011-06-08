package main

import "cairo"
import "http"
import "log"
import "os"
import "math"

import "./quadratic/quadratic"
import "./zellij"

var ZellijTilings <-chan *quadratic.Map
var reset chan<- int

func main() {
	http.HandleFunc("/",MainScreen)
	http.HandleFunc("/svg",RenderSVG)
	http.HandleFunc("/start",StartTiling)
	http.HandleFunc("/tiles",RenderTiles)
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

	/*t := zellij.PathMap("tspjgbc")
	u := zellij.PathMap("tspjg").Translate(quadratic.NewVertex(zellij.Points["j"]),quadratic.NewVertex(zellij.Points["p"]))
	v,ok := t.Overlay(u)
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	v,ok = u.Overlay(t)
	if ok != nil {
		os.Stderr.WriteString(ok.String()+"\n")
	}
	t = v*/
	
	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)
	image.SetLineWidth(.1)
	image.Translate(72*2.,72*2.)
	image.Scale(4.,4.)
	t.DrawEdges(image)

	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func StartTiling(w http.ResponseWriter, req *http.Request) {
	if reset != nil {
		reset <- 1
	}

	ZellijTilings,reset = zellij.TileRegion(new(quadratic.Integer),quadratic.NewInteger(40,0),new(quadratic.Integer),quadratic.NewInteger(40,0))
	
	w.WriteHeader(http.StatusOK)
}
	

func RenderSVG(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	polys := req.Form["polys"]

	if len(polys) == 0 {
		http.ServeFile(w,req,"svg/empty.svg")
		return
	}

	e := os.Remove("svg/test-surface.svg")
	if e != nil {
		os.Stderr.WriteString(e.String())
	}	

	maps := make([]*quadratic.Map,len(polys))

	for i:=0; i < len(polys); i++ {
		m,ok := quadratic.PolygonMapFromString(polys[i])
		if ok != nil {
			http.ServeFile(w,req,"svg/empty.svg")
			os.Stderr.WriteString("couldn't parse points\n")
		}	
		maps[i] = m
	}

	for i:=1; i < len(polys); i++ {
		_,e := maps[0].Overlay(maps[i])
		if e != nil {
			os.Stderr.WriteString(e.String())
		}
	}

	image := cairo.NewSurface("svg/test-surface.svg",72*4,72*4)
	image.SetSourceRGB(0.,0.,0.)

	maps[0].DrawEdges(image)
	maps[1] = maps[0].Copy()
	/*maps[1].Translate(quadratic.NewVertex(quadratic.NewPoint(new(quadratic.Integer),new(quadratic.Integer))),
		quadratic.NewVertex(quadratic.NewPoint(new(quadratic.Integer),quadratic.NewInteger(100,0))))*/
	image.SetSourceRGB(1.,0.,0.)
	maps[1].DrawEdges(image)

	image.Finish()
	http.ServeFile(w,req,"svg/test-surface.svg")
}

func DrawCircle(ctx *cairo.Surface, x,y float64) {
	ctx.Arc(x,y,1.,0.,2.*math.Pi)
	ctx.Fill()
}
