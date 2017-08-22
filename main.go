package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Editor struct {
	Text [][]rune
	X    int
	Y    int
}

func NewEditor() *Editor {
	return &Editor{
		[][]rune{{}},
		0,
		0,
	}
}

func (ed *Editor) MoveLeft() {
	if ed.X == 0 {
		return
	}
	ed.X--
}

func (ed *Editor) MoveRight() {
	if len(ed.Text[ed.Y]) == ed.X {
		return
	}
	ed.X++
}

func (ed *Editor) MoveUp() {
	if ed.Y == 0 {
		return
	}

	ed.Y--
	if len(ed.Text[ed.Y])-1 < ed.X {
		ed.X = len(ed.Text[ed.Y])
	}
}

func (ed *Editor) MoveDown() {
	if len(ed.Text)-1 == ed.Y {
		return
	}

	ed.Y++
	if len(ed.Text[ed.Y])-1 < ed.X {
		ed.X = len(ed.Text[ed.Y])
	}
}

func (ed *Editor) GoToLineStart() {
	ed.X = 0
}

func (ed *Editor) GoToLineEnd() {
	ed.X = len(ed.Text[ed.Y])
}

func (ed *Editor) AddLine() {
	currentLine := ed.Text[ed.Y]

	remainingLine := make([]rune, len(currentLine[:ed.X]))
	copy(remainingLine, currentLine[:ed.X])
	newLine := make([]rune, len(currentLine[ed.X:]))
	copy(newLine, currentLine[ed.X:])

	ed.Text[ed.Y] = remainingLine

	newText := make([][]rune, len(ed.Text)+1)
	headLines := ed.Text[:ed.Y+1]
	tailLines := ed.Text[ed.Y+1:]
	copy(newText[:ed.Y+1], headLines)
	newText[ed.Y+1] = newLine
	copy(newText[ed.Y+2:], tailLines)
	ed.Text = newText

	ed.X = 0
	ed.Y++
}

func (ed *Editor) AddRune(r rune) {

	if len(ed.Text[ed.Y]) == ed.X {
		ed.Text[ed.Y] = append(ed.Text[ed.Y], r)
		ed.X++
		return
	}

	newLine := make([]rune, len(ed.Text[ed.Y])+1)
	head := ed.Text[ed.Y][:ed.X]
	tail := ed.Text[ed.Y][ed.X:]

	copy(newLine[:ed.X], head)
	newLine[ed.X] = r
	copy(newLine[ed.X+1:], tail)

	ed.Text[ed.Y] = newLine

	ed.MoveRight()
}

func (ed *Editor) RemoveBackwardRune() {
	if ed.X == 0 && ed.Y == 0 {
		return
	}

	if ed.X == 0 {
		currentLine := ed.Text[ed.Y]
		currentLineLen := len(currentLine)

		prevLine := ed.Text[ed.Y-1]
		prevLineLen := len(prevLine)

		newLine := make([]rune, prevLineLen+currentLineLen)
		copy(newLine[:prevLineLen], prevLine)
		copy(newLine[prevLineLen:], currentLine)
		ed.Text[ed.Y-1] = newLine

		ed.X = prevLineLen
		ed.removeLine(ed.Y)
		ed.Y--

		return
	}

	ed.X--

	newLine := make([]rune, len(ed.Text[ed.Y])-1)
	head := ed.Text[ed.Y][:ed.X+1]
	tail := ed.Text[ed.Y][ed.X+1:]

	copy(newLine[:ed.X], head)
	copy(newLine[ed.X:], tail)

	ed.Text[ed.Y] = newLine
}

func (ed *Editor) RemoveForwardRune() {
	if ed.X == len(ed.Text[ed.Y]) && ed.Y == len(ed.Text)-1 {
		return
	}

	if ed.X == len(ed.Text[ed.Y]) {
		currentLine := ed.Text[ed.Y]
		currentLineLen := len(currentLine)

		nextLine := ed.Text[ed.Y+1]
		nextLineLen := len(nextLine)

		newLine := make([]rune, nextLineLen+currentLineLen)
		copy(newLine[:currentLineLen], currentLine)
		copy(newLine[currentLineLen:], nextLine)
		ed.Text[ed.Y] = newLine

		ed.X = currentLineLen
		ed.removeLine(ed.Y + 1)

		return
	}

	newLine := make([]rune, len(ed.Text[ed.Y])-1)
	head := ed.Text[ed.Y][:ed.X]
	tail := ed.Text[ed.Y][ed.X+1:]

	copy(newLine[:ed.X], head)
	copy(newLine[ed.X:], tail)

	ed.Text[ed.Y] = newLine
}

func (ed *Editor) removeLine(lineOffset int) {

	newText := make([][]rune, len(ed.Text)-1)
	head := ed.Text[:lineOffset]
	tail := ed.Text[lineOffset+1:]

	copy(newText[:lineOffset], head)
	copy(newText[lineOffset:], tail)

	ed.Text = newText
}

type Screen struct {
	ed          *Editor
	prevEditorX int
	prevEditorY int
	offsetX     int
	offsetY     int
	cursorX     int
	cursorY     int
}

func NewScreen(ed *Editor) *Screen {
	return &Screen{
		ed: ed,
	}
}

func (sc *Screen) Draw() {
	const color = termbox.ColorDefault

	termbox.Clear(color, color)

	windowWidth, windowHeight := termbox.Size()

	sc.updateOffsetX(windowWidth)
	sc.updateOffsetY(windowHeight)

	sc.cursorX = 0
	sc.cursorY = sc.ed.Y - sc.offsetY

	text := sc.ed.Text[sc.offsetY:]
	if len(text) > windowHeight {
		text = text[:windowHeight]
	}

	for y, line := range text {

		if len(line) <= sc.offsetX {
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

	sc.prevEditorX = sc.ed.X
	sc.prevEditorY = sc.ed.Y

	termbox.SetCursor(sc.cursorX, sc.cursorY)

	termbox.Flush()
}

func (sc *Screen) updateOffsetX(width int) {
	if sc.offsetX <= sc.ed.X && sc.ed.X <= width {
		return
	}
	sc.offsetX = sc.ed.X - width
}

func (sc *Screen) updateOffsetY(height int) {
	if sc.offsetY <= sc.ed.Y && sc.ed.Y < height {
		return
	}
	if sc.ed.Y > sc.prevEditorY && sc.ed.Y >= height {
		sc.offsetY = sc.ed.Y + 1 - height
	}
	if sc.ed.Y < sc.prevEditorY && sc.ed.Y < sc.offsetY {
		sc.offsetY--
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	const color = termbox.ColorDefault
	termbox.Clear(color, color)

	ed := NewEditor()
	sc := NewScreen(ed)

mainloop:
	for {
		ev := termbox.PollEvent()

		if ev.Type != termbox.EventKey {
			continue
		}

		switch ev.Key {
		case termbox.KeyEsc:
			break mainloop
		case termbox.KeyEnter:
			ed.AddLine()
		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			ed.MoveLeft()
		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			ed.MoveRight()
		case termbox.KeyArrowUp, termbox.KeyCtrlP:
			ed.MoveUp()
		case termbox.KeyArrowDown, termbox.KeyCtrlN:
			ed.MoveDown()
		case termbox.KeyCtrlA:
			ed.GoToLineStart()
		case termbox.KeyCtrlE:
			ed.GoToLineEnd()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			ed.RemoveBackwardRune()
		case termbox.KeyDelete, termbox.KeyCtrlD:
			ed.RemoveForwardRune()
		case termbox.KeySpace:
			ed.AddRune(' ')
		case termbox.KeyTab:
			ed.AddRune('\t')
		default:
			if ev.Ch != 0 {
				ed.AddRune(ev.Ch)
			}
		}
		sc.Draw()
	}
}
