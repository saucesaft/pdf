package pdf

import (
    "image"
    "strconv"
    "github.com/gen2brain/go-fitz"
    "github.com/go-gl/gl/v3.2-core/gl"
    "github.com/emirpasic/gods/maps/treemap"
)

type PageBase struct {
    Contents uint32
    W, H float32
//    Bytes []byte // representation of a page message used for serialization
    BaseIndex uint
    Filename string
  //  ShouldDelete bool
}

type Orphan struct { // lone page that lurks around
    PageBase
    OrphanID uint
}
func (o Orphan) GenTestingKey(argKey uint64) string {
    oid := strconv.FormatUint(argKey, 10)
    pid := strconv.FormatUint(uint64(o.BaseIndex), 10)

    return o.Filename+"-"+oid+"-"+pid
}
func (o Orphan) GenKey() string {
    oid := strconv.FormatUint(uint64(o.OrphanID), 10)
    pid := strconv.FormatUint(uint64(o.BaseIndex), 10)

    return o.Filename+"-"+oid+"-"+pid
}

type Page struct { // lone page that has a connection with another page
    PageBase
    FamilyID uint
}

func (p Page) GenKey() string {
    fid := strconv.FormatUint(uint64(p.FamilyID), 10)
    pid := strconv.FormatUint(uint64(p.BaseIndex), 10)

    return p.Filename+"-"+fid+"-"+pid
}
type App struct {
    Pages *treemap.Map
    Orphans *treemap.Map
    LastFamilyID uint
}

func (a *App) Init() {
    a.Pages = treemap.NewWithStringComparator()
    a.Orphans = treemap.NewWithStringComparator()
}

func (a *App) NewList() {
    // TODO will create an identifier for the list
    // set it also on lastfamilyid
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

    inPath := *path // assign to another variable for internal uses
    *path = "" // clear the imGui input

    var nextFamID uint

    if a.LastFamilyID == 0 {
	nextFamID = 1
	a.LastFamilyID = nextFamID
    } else {
	nextFamID = a.LastFamilyID + 1
	a.LastFamilyID = nextFamID
    }

    for i, img := range imgs {
	// prepare structs and textures
        tex, w, h  := getTexture(img)
        p := Page{}

	// setup struct
        p.Contents = tex
        p.W = float32(w)
        p.H = float32(h)
	p.FamilyID = nextFamID
	p.BaseIndex = uint(i)
	p.Filename = inPath

	// generate key for treemap
	pk := p.GenKey()

	// insert page to treemap with corresponding key
	a.Pages.Put(pk, p)

    }


}
