package main

import (
	"pdf/gui"
	"pdf/pdf"

	g"github.com/AllenDang/giu"
	im"github.com/AllenDang/giu/imgui"

)

func loop(a *pdf.App) {
	always := true

	g.SingleWindowWithMenuBar("pdf").Layout(g.Layout{
		gui.MenuBar(a),
		gui.Gui(a),
	})

	im.ShowDemoWindow(&always)

}

func main() {
	app := pdf.App{}
	app.Init()

	wnd := g.NewMasterWindow("pdf", 700, 550, 0, nil)
    wnd.Run(func () {loop(&app)})
}