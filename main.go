package main

import (
	"github.com/bananaumai/simple-text-editor/editor"
	"github.com/bananaumai/simple-text-editor/screen"
)

func main() {
	ed := editor.NewEditor()
	sc := screen.NewScreen(ed)
	sc.Run()
}
