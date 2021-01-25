package pdf

import (
	g"github.com/AllenDang/giu"
)

type Page struct {
	Filename string
	DocId, FamId uint
	Texture *g.Texture
	W, H float32
}

func newPage (t *g.Texture, d, f uint, file string, w, h float32) Page {
	return Page {
		Texture: t,
		DocId: d,
		FamId: f,
		Filename: file,
		W: w,
		H: h,
	}
}

func (p Page) Key() Key {
	return Key {p.Filename, p.DocId, p.FamId}
}

func (p *Page) UpdateSize(w, h float32) {
	p.W = w
	p.H = h
}


type Key struct {
	Filename string
	DocId, FamId uint
}