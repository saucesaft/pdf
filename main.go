// +build glfw

package main

import (
	"fmt"
	"os"

	"pdf/backend"
	"pdf/gui"
	"pdf/pdf"

	"github.com/inkyblackness/imgui-go/v2"
)

type Platform interface {
	ShouldStop() bool
	ProcessEvents()
	DisplaySize() [2]float32
	FramebufferSize() [2]float32
	NewFrame()
	PostRender()
	ClipboardText() (string, error)
	SetClipboardText(text string)
}

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

type Renderer interface {
	PreRender(clearColor [3]float32)
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)
}

func run(p Platform, r Renderer, a pdf.App) {
	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	always := true
	clearColor := [3]float32{0.0, 0.0, 0.0}

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		// DEMO
		imgui.ShowDemoWindow(&always)

        gui.NewListWindow(&a)

        gui.ShowLists(&a)

		gui.NewPDFWindow(&a)

        gui.ShowOrphans(&a)

		// PDF (for each page in the lone pages in app)
		/*    imgui.Begin("PDF")

		      imgui.BeginChildV("No Move Child", imgui.Vec2{-1, -1}, false, imgui.WindowFlagsNoMove)

		      if (imgui.BeginDragDropSource(0)) {
		          imgui.SetDragDropPayload("PAGE", []byte{'L','O','L'}, 0)
		          imgui.Text("PAGE")
		          imgui.EndDragDropSource()
		      }

		      wsx := imgui.WindowSize().X
		      imgui.Image(imgui.TextureID(t), imgui.Vec2{wsx, (wsx/w)*h })

		      imgui.EndChild()

		      imgui.End()*/

		ww := p.FramebufferSize()[0]
		wh := p.FramebufferSize()[1]

		gui.BackgroundDropHandler(ww, wh, &a)

		// Rendering
		imgui.Render()

		r.PreRender(clearColor)

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()

	}
}

func main() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := backend.NewGLFW(io, backend.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := backend.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	app := pdf.App{}

	run(platform, renderer, app)
}
