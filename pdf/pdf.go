package pdf

import (
    "sync"
    "image"

    "pdf/structures"

    g"github.com/AllenDang/giu"
    "github.com/gen2brain/go-fitz"
)

type App struct {
    Test string
    EditMode bool
    LastFamID uint
    TextureMutex sync.Mutex
    HashmapMutex sync.Mutex

    ShouldScroll bool
    ScrollOffset float32
    LastWidth float32


    Pages *structures.Map
}

func (a *App) Init () {
    a.Test = "Hello World"
    a.EditMode = false
    a.ShouldScroll = false

    a.Pages = structures.NewLinkedHashMap()

}

func (a *App) ToggleEdit () {
    a.EditMode = !a.EditMode
}

func (a *App) Open() {
    filename := "test2.pdf"

    doc, err := fitz.New(filename)
    if err != nil {
        panic(err)
    }

    defer doc.Close()


    var nextFamID uint

    if a.LastFamID == 0 {
        nextFamID = 1
        a.LastFamID = nextFamID
    } else {
        nextFamID = a.LastFamID + 1
        a.LastFamID = nextFamID
    }

    for n := 0; n < doc.NumPage(); n++ {

        a.TextureMutex.Lock()
        page, key := preparePage(n, filename, nextFamID, doc)
        a.TextureMutex.Unlock()

        a.HashmapMutex.Lock()
        a.Pages.Put(key, page)
        a.HashmapMutex.Unlock()

    }

}

func (a *App) ViewerScroll(s float32) {
    a.ShouldScroll = true
    a.ScrollOffset = s
}

func (a *App) FinishScroll() {
    a.ShouldScroll = false
}

func (a *App) UpdateWidth(w float32) {
    a.LastWidth = w
}

func preparePage(n int, filename string, nfid uint, doc *fitz.Document) (Page, Key){

    img, err := doc.Image(n)
    if err != nil {
        panic(err)
    }

    tex, _ := g.NewTextureFromRgba(img.(*image.RGBA))  //TODO dont ignore error

    b := img.Bounds()

    page := newPage(tex, uint(n+1), nfid, filename, float32(b.Max.X), float32(b.Max.Y))

    key := page.Key()

    return page, key
}