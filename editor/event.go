package editor

type EditorEvent int

const (
	_ = iota
	evMoveUp
	evMoveDown
	evMoveLeft
	evMoveRight
)

const (
	EditorEventMoveUp    EditorEvent = evMoveUp
	EditorEventMoveDown  EditorEvent = evMoveDown
	EditorEventMoveLeft  EditorEvent = evMoveLeft
	EditorEventMoveRight EditorEvent = evMoveRight
)

type EventListener func(ed *Editor)
