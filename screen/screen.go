package screen

import (
	"github.com/bananaumai/simple-text-editor/editor"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Screen struct {
	ed      *editor.Editor
	color   termbox.Attribute
	offsetX int
	offsetY int
	cursorX int
	cursorY int
	width   int
	height  int
}

func NewScreen(ed *editor.Editor) *Screen {
	sc := &Screen{}

	ed.AddEventListener(editor.EditorEventMoveUp, sc.decreaseOffsetY)
	ed.AddEventListener(editor.EditorEventMoveDown, sc.increaseOffsetY)
	ed.AddEventListener(editor.EditorEventMoveLeft, sc.decreaseOffsetX)
	ed.AddEventListener(editor.EditorEventMoveRight, sc.increaseOffsetX)
	sc.ed = ed

	sc.color = termbox.ColorDefault
	termbox.Clear(sc.color, sc.color)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	sc.width, sc.height = termbox.Size()

	return sc
}

func (sc *Screen) increaseOffsetY(ed *editor.Editor) {
	if ed.Y-sc.offsetY >= sc.height {
		sc.offsetY = ed.Y - sc.height + 1
	}
}

func (sc *Screen) decreaseOffsetY(ed *editor.Editor) {
	if ed.Y < sc.offsetY {
		sc.offsetY--
	}
}

func (sc *Screen) increaseOffsetX(ed *editor.Editor) {
	const rightBuffer = 1
	if ed.X-sc.offsetX >= sc.width-rightBuffer {
		sc.offsetX = ed.X - sc.width + 1
	}
}

func (sc *Screen) decreaseOffsetX(ed *editor.Editor) {
	if ed.X < sc.offsetX {
		sc.offsetX--
	}
}

func (sc *Screen) Run() {

	defer termbox.Close()

mainloop:
	for {
		ev := termbox.PollEvent()

		if ev.Type == termbox.EventResize {
			sc.width, sc.height = termbox.Size()
			sc.Draw()
			continue
		}

		if ev.Type != termbox.EventKey {
			continue
		}

		switch ev.Key {
		case termbox.KeyEsc:
			break mainloop
		case termbox.KeyEnter:
			sc.ed.AddLine()
		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			sc.ed.MoveLeft()
		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			sc.ed.MoveRight()
		case termbox.KeyArrowUp, termbox.KeyCtrlP:
			sc.ed.MoveUp()
		case termbox.KeyArrowDown, termbox.KeyCtrlN:
			sc.ed.MoveDown()
		case termbox.KeyCtrlA:
			sc.ed.GoToLineStart()
		case termbox.KeyCtrlE:
			sc.ed.GoToLineEnd()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			sc.ed.RemoveBackwardRune()
		case termbox.KeyDelete, termbox.KeyCtrlD:
			sc.ed.RemoveForwardRune()
		case termbox.KeySpace:
			sc.ed.AddRune(' ')
		case termbox.KeyTab:
			sc.ed.AddRune('\t')
		default:
			if ev.Ch != 0 {
				sc.ed.AddRune(ev.Ch)
			}
		}

		sc.Draw()
	}
}

func (sc *Screen) Draw() {
	termbox.Clear(sc.color, sc.color)

	sc.cursorX = 0
	sc.cursorY = sc.ed.Y - sc.offsetY

	text := sc.ed.Text[sc.offsetY:]
	if len(text) > sc.height {
		text = text[:sc.height]
	}

	for y, line := range text {

		if len(line) < sc.offsetX {
			break
		}
		line = line[sc.offsetX:]
		if len(line) > sc.width {
			line = line[:sc.width]
		}

		x := 0
		for i, r := range line {
			termbox.SetCell(x, y, r, sc.color, sc.color)
			width := runewidth.RuneWidth(r)
			if sc.cursorY == y && sc.ed.X-sc.offsetX > i {
				sc.cursorX += width
			}
			x += width
		}
	}

	termbox.SetCursor(sc.cursorX, sc.cursorY)

	termbox.Flush()
}
