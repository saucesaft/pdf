package gui

import (
    "fmt"
    "bytes"
    "encoding/gob"

	"github.com/inkyblackness/imgui-go/v2"

    "pdf/pdf"
)

var (
    loadPath string
    always bool = true
)

// NEW PDF
func NewPDFWindow(a *pdf.App) {
    imgui.Begin("Wanna load a pdf?")
        imgui.InputText("<- path", &loadPath)
        imgui.SameLine()
        if (imgui.Button("GO")) {
            a.LoadPdf(&loadPath)
        }
    imgui.End()
}

func ShowLists(a *pdf.App) {
    for index, list := range a.Lists {
        imgui.BeginV(fmt.Sprintf("List %d", index), &always, imgui.WindowFlagsAlwaysAutoResize)
        if imgui.BeginDragDropTarget() {
            lolByte := imgui.AcceptDragDropPayload("PAGE", 0) // TODO change to add to list
            if lolByte != nil {
                fmt.Println(lolByte)
            }
        }
        if len(list.Pages) == 0 {
            imgui.Text("Drag some pages here")
        } else {
            for i, page := range list.Pages {
                imgui.Selectable(fmt.Sprintf("Page %d", i))
                if (imgui.BeginDragDropSource(0)) {
                    imgui.SetDragDropPayload("PAGE", page.Bytes, 0)
                    imgui.Text("PAGE")
                    imgui.EndDragDropSource()
                }
                if imgui.IsItemHovered() {
                    imgui.BeginTooltip()
                    imgui.Image(
                        imgui.TextureID(page.Contents),
                        imgui.Vec2{X: page.W / 16, Y: page.H / 16},
                    )
                    imgui.EndTooltip()
                }
            }
        }
        imgui.End()
    }

}

func ShowOrphans(a *pdf.App) {
    for _, orphan := range a.Orphans {

        imgui.BeginChildV("No Move Child", imgui.Vec2{X: -1, Y: -1}, false, imgui.WindowFlagsNoMove)

          if (imgui.BeginDragDropSource(0)) {
              imgui.SetDragDropPayload("PAGE", []byte{'L','O','L'}, 0)
              imgui.Text("PAGE")
              imgui.EndDragDropSource()
          }

          wsx := imgui.WindowSize().X
          imgui.Image(imgui.TextureID(orphan.Contents), imgui.Vec2{X: wsx, Y: (wsx/orphan.W)*orphan.H })

          imgui.EndChild()

    }

}


// BACKGROUND
func BackgroundDropHandler(ww float32, wh float32, a *pdf.App) {

    imgui.SetNextWindowPos(imgui.Vec2{X: 0.0, Y: 0.0})
    imgui.SetNextWindowSize(imgui.Vec2{X: ww, Y: wh})
    imgui.BeginV("Background", &always, imgui.WindowFlagsNoBackground |
        imgui.WindowFlagsNoBringToFrontOnFocus | imgui.WindowFlagsNoDecoration |
        imgui.WindowFlagsNoMove)
        imgui.Dummy(imgui.Vec2{X: ww-15, Y: wh-15})
        if (imgui.BeginDragDropTarget()) {
            pageBytes := imgui.AcceptDragDropPayload("PAGE", 0)
            if (pageBytes != nil) {

                buf := bytes.NewBuffer(pageBytes)
                dec := gob.NewDecoder(buf)

                var p pdf.Page
                err := dec.Decode(&p)
                if err != nil {
                    panic(err)
                }

                //p.FatherList.Pages = append(p.FatherList.Pages[:p.RemoveIndex], p.FatherList.Pages[p.RemoveIndex+1:]...)

                a.Orphans = append(a.Orphans, p) // TODO, when dragging remove from list

                p.Remove() // TODO removes itself from the list
            }
            imgui.EndDragDropTarget()
        }
    imgui.End()
}

// NEW LIST
func NewListWindow(a *pdf.App) {
    imgui.Begin("New List")
    if imgui.Button("Create new list") {
        a.NewList()
    }
    imgui.End()
}
