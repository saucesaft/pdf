package gui

import (
	im"github.com/AllenDang/giu/imgui"
	"github.com/JoshVarga/svgparser"
	"strconv"
	"strings"
)

/*func SetTheme() {
	im.CurrentStyle().SetColor(im.StyleColorWindowBg, im.Vec4{0.42, 0.42, 0.42, 1})
}*/

func floatColor(hex string) im.Vec4 {
	code := strings.ReplaceAll(hex, "#", "")

	values, err := strconv.ParseUint(code, 16, 32)

	if err != nil {
		panic(err)
	}

	r := float32(values >> 16)
	g := float32((values >> 8) & 0xFF)
	b := float32(values & 0xFF)

	return im.Vec4{r/255.0, g/255.0, b/255.0, 1.0}
}

func SetTheme() {
	blob := `<svg width="96px" height="64px" xmlns="http://www.w3.org/2000/svg" baseProfile="full" version="1.1">
  <rect width='96' height='64' id='background' fill='#a1a1a1'></rect>
  <!-- Foreground -->
  <circle cx='24' cy='24' r='8' id='f_high' fill='#222222'></circle>
  <circle cx='40' cy='24' r='8' id='f_med' fill='#e00b30'></circle>
  <circle cx='56' cy='24' r='8' id='f_low' fill='#888888'></circle>
  <circle cx='72' cy='24' r='8' id='f_inv' fill='#ffffff'></circle>
  <!-- Background -->
  <circle cx='24' cy='40' r='8' id='b_high' fill='#555555'></circle>
  <circle cx='40' cy='40' r='8' id='b_med' fill='#fbba2d'></circle>
  <circle cx='56' cy='40' r='8' id='b_low' fill='#b3b3b3'></circle>
  <circle cx='72' cy='40' r='8' id='b_inv' fill='#0e7242'></circle>
</svg>`

	ids := [9]string{"background", "f_high", "f_med", "f_low", "f_inv", "b_high", "b_med", "b_low", "b_inv"}

	reader := strings.NewReader(blob)

	element, _ := svgparser.Parse(reader, false)

	theme := make(map[string]im.Vec4)

	for i := 0; i < len(element.Children); i++ {

		id := element.Children[i].Attributes["id"]
		for _, x := range ids {
			if x == id {
				theme[id] = floatColor(element.Children[i].Attributes["fill"])
			}
		}
	}

	im.CurrentStyle().SetColor(im.StyleColorWindowBg, theme["background"])
	im.CurrentStyle().SetColor(im.StyleColorChildBg, theme["b_high"])
	im.CurrentStyle().SetColor(im.StyleColorMenuBarBg, theme["b_high"])

	im.CurrentStyle().SetColor(im.StyleColorText, theme["f_high"])

	im.CurrentStyle().SetColor(im.StyleColorButtonHovered, theme["f_med"])
	im.CurrentStyle().SetColor(im.StyleColorTabHovered, theme["f_med"])
	im.CurrentStyle().SetColor(im.StyleColorHeaderHovered, theme["f_med"])
	im.CurrentStyle().SetColor(im.StyleColorSeparatorHovered, theme["f_med"])
	im.CurrentStyle().SetColor(im.StyleColorScrollbarGrabActive, theme["f_med"])

	im.CurrentStyle().SetColor(im.StyleColorSeparatorActive, theme["f_low"])

	im.CurrentStyle().SetColor(im.StyleColorScrollbarGrab, theme["f_low"])

	//im.CurrentStyle().SetColor(im.StyleColorBorder, theme[""])
}