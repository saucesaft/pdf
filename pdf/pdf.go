package pdf

import (
    "image"
    "bytes"
    "encoding/gob"
    "github.com/gen2brain/go-fitz"
	"github.com/go-gl/gl/v3.2-core/gl"
)

type Page struct {
    Contents uint32
    W, H float32
    Bytes []byte

    RemoveIndex int
    FatherList *List
}

type List struct {
    Pages []Page
}

type App struct {
    Lists []List // global pages
    Orphans []Page // orphan pages
}

func (a *App) NewList() {
    l := List {}

    a.Lists = append(a.Lists, l)
}

func getTexture(img image.Image) (uint32, float32, float32) {
    b := img.Bounds()
    w := int32(b.Max.X)
    h := int32(b.Max.Y)

    var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h,
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr( img.(*image.RGBA).Pix ) )

    return tex, float32(w), float32(h)
}

func imgsFromDoc(path string) []image.Image {
    doc, err := fitz.New(path)
	if err != nil {
		panic(err)
	}

	defer doc.Close()

    var imgs []image.Image

    for n := 0; n < doc.NumPage(); n++ {
        img, err := doc.Image(n)
        if err != nil {
            panic(err)
        }
        imgs = append(imgs, img)
    }

    return imgs
}

func (a *App) LoadPdf(path *string) {

    imgs := imgsFromDoc(*path)

    *path = ""

    l := List{}

    for i, img := range imgs {
        tex, w, h  := getTexture(img)

        p := Page{}
        b := bytes.Buffer{}
	    enc := gob.NewEncoder(&b)

        p.Contents = tex
        p.W = float32(w)
        p.H = float32(h)
        p.RemoveIndex = i
        p.FatherList = &l

        err := enc.Encode(p)
        if err != nil {
            panic(err)
        }

        p.Bytes = b.Bytes()

        l.Pages = append(l.Pages, p)

    }

    a.Lists = append(a.Lists, l)

}

func (p *Page) Remove() {}
