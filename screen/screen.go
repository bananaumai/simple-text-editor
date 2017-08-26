package screen

import (
	"github.com/bananaumai/simple-text-editor/editor"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Screen struct {
	prevEditorX int
	prevEditorY int
	offsetX     int
	offsetY     int
	cursorX     int
	cursorY     int
}

func NewScreen() *Screen {
	return &Screen{}
}

func (sc *Screen) Draw(ed *editor.Editor) {
	const color = termbox.ColorDefault

	termbox.Clear(color, color)

	windowWidth, windowHeight := termbox.Size()

	sc.updateOffsetX(windowWidth, ed)
	sc.updateOffsetY(windowHeight, ed)

	sc.cursorX = 0
	sc.cursorY = ed.Y - sc.offsetY

	text := ed.Text[sc.offsetY:]
	if len(text) > windowHeight {
		text = text[:windowHeight]
	}

	for y, line := range text {

		if len(line) < sc.offsetX {
			break
		}
		line = line[sc.offsetX:]
		if len(line) > windowWidth {
			line = line[:windowWidth]
		}

		x := 0
		for _, r := range line {
			termbox.SetCell(x, y, r, color, color)
			width := runewidth.RuneWidth(r)
			x += width
			if sc.cursorY == y {
				sc.cursorX += width
			}
		}
	}

	sc.prevEditorX = ed.X
	sc.prevEditorY = ed.Y

	termbox.SetCursor(sc.cursorX, sc.cursorY)

	termbox.Flush()
}

func (sc *Screen) updateOffsetX(width int, ed *editor.Editor) {
	if sc.offsetX <= ed.X && ed.X <= width {
		return
	}
	sc.offsetX = ed.X - width
}

func (sc *Screen) updateOffsetY(height int, ed *editor.Editor) {
	if sc.offsetY <= ed.Y && ed.Y < height {
		return
	}
	if ed.Y > sc.prevEditorY && ed.Y >= height {
		sc.offsetY = ed.Y + 1 - height
	}
	if ed.Y < sc.prevEditorY && ed.Y < sc.offsetY {
		sc.offsetY--
	}
}
