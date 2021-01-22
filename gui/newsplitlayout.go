package gui

import (
	"fmt"
	g"github.com/AllenDang/giu"
)

type SplitLayoutStateNew struct {
	delta   float32
	sashPos float32
}

func (s *SplitLayoutStateNew) Dispose() {
	// Nothing to do here.
}

type SplitLayoutWidgetNew struct {
	id                  string
	direction           g.SplitDirection
	layout1             g.Widget
	layout2             g.Widget
	originItemSpacingX  float32
	originItemSpacingY  float32
	originFramePaddingX float32
	originFramePaddingY float32
	sashPos             float32
	border              bool
}

func SplitLayoutNew(id string, direction g.SplitDirection, border bool, sashPos float32, layout1, layout2 g.Widget) *SplitLayoutWidgetNew {
	return &SplitLayoutWidgetNew{
		id:        id,
		direction: direction,
		sashPos:   sashPos,
		layout1:   layout1,
		layout2:   layout2,
		border:    border,
	}
}

func (s *SplitLayoutWidgetNew) restoreItemSpacing(layout g.Widget) g.Layout {
	return g.Layout{
		g.Custom(func() {
			g.PushItemSpacing(s.originItemSpacingX, s.originItemSpacingY)
			g.PushFramePadding(s.originFramePaddingX, s.originFramePaddingY)
		}),
		layout,
		g.Custom(func() {
			g.PopStyleV(2)
		}),
	}
}

// Build Child panel. If layout is a SplitLayout, set the frame padding to zero.
func (s *SplitLayoutWidgetNew) buildChild(id string, width, height float32, layout g.Widget, flags g.WindowFlags) g.Widget {
	_, isSplitLayoutWidget := layout.(*SplitLayoutWidgetNew)

	return g.Layout{
		g.Custom(func() {
			if isSplitLayoutWidget || !s.border {
				g.PushFramePadding(0, 0)
			}
		}),
		g.Child(id).
		Border(!isSplitLayoutWidget && s.border).
		Size(width, height).
		Layout(s.restoreItemSpacing(layout)).
		Flags(flags),

		g.Custom(func() {
			if isSplitLayoutWidget || !s.border {
				g.PopStyle()
			}
		}),
	}
}

func (s *SplitLayoutWidgetNew) Build() {
	var splitLayoutStateNew *SplitLayoutStateNew
	// Register state
	stateId := fmt.Sprintf("SplitLayout_%s", s.id)
	if state := g.Context.GetState(stateId); state == nil {
		splitLayoutStateNew = &SplitLayoutStateNew{delta: 0.0, sashPos: s.sashPos}
		g.Context.SetState(stateId, splitLayoutStateNew)
	} else {
		splitLayoutStateNew = state.(*SplitLayoutStateNew)
	}

	itemSpacingX, itemSpacingY := g.GetItemInnerSpacing()
	s.originItemSpacingX, s.originItemSpacingY = itemSpacingX, itemSpacingY

	s.originFramePaddingX, s.originFramePaddingY = g.GetFramePadding()
	s.originFramePaddingX /= g.Context.GetPlatform().GetContentScale()
	s.originFramePaddingY /= g.Context.GetPlatform().GetContentScale()

	var layout g.Layout

	splitLayoutStateNew.sashPos += splitLayoutStateNew.delta

	lcf := g.WindowFlagsNoScrollbar | g.WindowFlagsNoBringToFrontOnFocus

	if s.direction == g.DirectionHorizontal {
		layout = g.Layout{
			g.Line(
				s.buildChild(fmt.Sprintf("%s_layout1", stateId), splitLayoutStateNew.sashPos, 0, s.layout1, lcf),
				g.VSplitter(fmt.Sprintf("%s_vsplitter", stateId), &(splitLayoutStateNew.delta)).Size(itemSpacingX+6, 0),
				s.buildChild(fmt.Sprintf("%s_layout2", stateId), 0, 0, s.layout2, lcf),
			),
		}
	} else {
		layout = g.Layout{
			s.buildChild(fmt.Sprintf("%s_layout1", stateId), 0, splitLayoutStateNew.sashPos, s.layout1, lcf),
			g.HSplitter(fmt.Sprintf("%s_hsplitter", stateId), &(splitLayoutStateNew.delta)).Size(0, itemSpacingY),
			s.buildChild(fmt.Sprintf("%s_layout2", stateId), 0, 0, s.layout2, lcf),
		}
	}

	g.PushItemSpacing(0, 0)
	layout.Build()
	g.PopStyle()
}