package gui

import (
    "fmt"
    "strconv"
    "github.com/inkyblackness/imgui-go/v2"

    "pdf/pdf"
)

var (
    loadPath string = "test.pdf" //TODO testing purposes
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

func ShowFamilies(a *pdf.App) {
    // TODO check if it doesnt have pages attached to the ID
    // TODO new list button logic (check if both numbers are the same)
    it := a.Pages.Iterator()
    for it.Next() {
	page := it.Value().(pdf.Page)

	imgui.BeginV(fmt.Sprintf("List %d", page.FamilyID), &always, 0)
	imgui.BeginChild(fmt.Sprintf("##List Child %d", page.FamilyID))

	imgui.Selectable(fmt.Sprintf("Page %d", page.BaseIndex+1))
	if (imgui.BeginDragDropSource(0)) {
	    pageKey := page.GenKey()

	    imgui.SetDragDropPayload("PAGE", []byte(pageKey), 0)
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

	imgui.EndChild()

        if imgui.BeginDragDropTarget() {
            pageKeyBytes := imgui.AcceptDragDropPayload("ORPHAN", 0)
            if pageKeyBytes != nil {
// TODO have an util function to auto convert (clean code)
		orphan, ok := a.Orphans.Get(string(pageKeyBytes))
		if !ok { // hacky workaround for double dragdrop bug
		    imgui.EndDragDropTarget()
		    imgui.End()
		    continue
		}
		fmt.Println(string(pageKeyBytes))
		a.Orphans.Remove(string(pageKeyBytes))
		base := orphan.(pdf.Orphan).PageBase
		newPage := pdf.Page{}
		newPage.PageBase = base
		newPage.FamilyID = page.FamilyID

		a.Pages.Put(newPage.GenKey(), newPage)
	    }
	    imgui.EndDragDropTarget()
	}
	imgui.End()
    }
}

func DebugWin(a *pdf.App) {
    if imgui.CollapsingHeader("Orphans") {
	it := a.Orphans.Iterator()
	for it.Next() {
	    imgui.Text(it.Key().(string))
	}
    }
    if imgui.CollapsingHeader("Pages") {
	it := a.Pages.Iterator()
	for it.Next() {
	    imgui.Text(it.Key().(string))
	}
    }
}

func ShowOrphans(a *pdf.App) {

    it := a.Orphans.Iterator()
    for it.Next() {
	orphan := it.Value().(pdf.Orphan)
//	imgui.SetNextWindowSizeConstraints(imgui.Vec2{X: -1, Y: -1}, imgui.Vec2{X: -1, Y: orphan.H}
	pageNum := strconv.Itoa(int(orphan.BaseIndex)+1)
	title := fmt.Sprintf(orphan.Filename+" Page: "+pageNum+"##"+orphan.GenKey())

	imgui.Begin(title)
        imgui.BeginChildV("No Move Child", imgui.Vec2{X: -1, Y: -1}, false, imgui.WindowFlagsNoMove)
	if (imgui.BeginDragDropSource(0)) {
	    imgui.SetDragDropPayload("ORPHAN", []byte(orphan.GenKey()), 0)
	    imgui.Text("PAGE")
	    imgui.EndDragDropSource()
	}

          wsx := imgui.WindowSize().X //TODO resize the pdf page with constraints
          imgui.Image(imgui.TextureID(orphan.Contents), imgui.Vec2{X: wsx, Y: (wsx/orphan.W)*orphan.H })

        imgui.EndChild()
	imgui.End()

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
            pageKeyBytes := imgui.AcceptDragDropPayload("PAGE", 0)
            if (pageKeyBytes != nil) {

		page, _ := a.Pages.Get(string(pageKeyBytes))

		a.Pages.Remove(string(pageKeyBytes))

		base := page.(pdf.Page).PageBase

		orphan := pdf.Orphan{}
		orphan.PageBase = base
		// check if the same page is already on the orphans
		var OrphanID uint
		_, check := a.Orphans.Get(orphan.GenTestingKey(0))
		if check {
		    var checkingInt uint64 = 1
		    for {
			fmt.Println(orphan.GenTestingKey(checkingInt))
			_, internalCheck := a.Orphans.Get(orphan.GenTestingKey(checkingInt))
			if !internalCheck {
			    OrphanID = uint(checkingInt)
			    break
			} else {
			    checkingInt++
			}
		    }
		} else {
		    OrphanID = 0
		}

		orphan.OrphanID = OrphanID

		// add to orphans
		a.Orphans.Put(orphan.GenKey(), orphan)

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
