package gui

import (
	//"fmt"
	"pdf/pdf"
	g"github.com/AllenDang/giu"
	im"github.com/AllenDang/giu/imgui"

)

func MenuBar(a *pdf.App) g.Layout {
	return g.Layout {
		g.MenuBar().Layout(g.Layout{
			g.Menu("~").Layout(g.Layout{
				g.MenuItem("Theme"),
				g.MenuItem("Toggle Mode"),
			}),
			g.Menu("File").Layout(g.Layout{
				g.MenuItem("Open").OnClick(func(){go a.Open()}),
			}),
		}),
	}
}

func Gui(a *pdf.App) g.Layout {
	return g.Layout{
		g.Custom(func() {
			//im.PushStyleVarFloat(im.StyleVarScrollbarSize, 0.0)

			SplitLayoutNew("Split", g.DirectionHorizontal, true, 100,
				left(a),
				right(a),
			).Build()

			//im.PopStyleVar()
		}),
	}

}

func left(a *pdf.App) g.Layout {
	it := a.Pages.Iterator()

	return 	g.Layout{
		g.MenuItem("Open").OnClick(func(){go a.Open()}),
		g.Custom(func(){
			wsx := im.GetItemRectSize().X

//			g.Label(fmt.Sprintf("%v", wsx)).Build()
			for it.Next() {
				page := it.Value().(pdf.Page)
				im.PushStyleColor(im.StyleColorButtonHovered, im.Vec4{0.1, 0.1, 0.1, 1.0})
				g.ImageButton(page.Texture).Size(wsx-12, ((wsx/page.W)*page.H)-12).Build()
				im.PopStyleColor()
			}
		}),
	}
}

func right(a *pdf.App) g.Layout {
	it := a.Pages.Iterator()

	return 	g.Layout{
		g.Custom(func(){
			wsx := im.GetItemRectSize().X

			for it.Next() {
				page := it.Value().(pdf.Page)

				g.Image(page.Texture).Size(wsx-16, ((wsx/page.W)*page.H)-16).Build()
			}
		}),
	}
}
