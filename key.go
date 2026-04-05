package tui

// KeyType identifies special keys.
type KeyType int

const (
	KeyRune KeyType = iota // Regular character (check Rune field)
	KeyEnter
	KeyTab
	KeyBacktab
	KeyBackspace
	KeyDelete
	KeyInsert
	KeyEscape
	KeySpace
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
	KeyCtrlD
	KeyCtrlE
	KeyCtrlF
	KeyCtrlG
	KeyCtrlH
	KeyCtrlK
	KeyCtrlL
	KeyCtrlN
	KeyCtrlO
	KeyCtrlP
	KeyCtrlQ
	KeyCtrlR
	KeyCtrlS
	KeyCtrlT
	KeyCtrlU
	KeyCtrlV
	KeyCtrlW
	KeyCtrlX
	KeyCtrlY
	KeyCtrlZ
)

// keyNames maps KeyType to human-readable names.
var keyNames = map[KeyType]string{
	KeyRune:      "rune",
	KeyEnter:     "enter",
	KeyTab:       "tab",
	KeyBacktab:   "shift+tab",
	KeyBackspace: "backspace",
	KeyDelete:    "delete",
	KeyInsert:    "insert",
	KeyEscape:    "escape",
	KeySpace:     "space",
	KeyUp:        "up",
	KeyDown:      "down",
	KeyLeft:      "left",
	KeyRight:     "right",
	KeyHome:      "home",
	KeyEnd:       "end",
	KeyPageUp:    "pgup",
	KeyPageDown:  "pgdn",
	KeyF1:        "f1",
	KeyF2:        "f2",
	KeyF3:        "f3",
	KeyF4:        "f4",
	KeyF5:        "f5",
	KeyF6:        "f6",
	KeyF7:        "f7",
	KeyF8:        "f8",
	KeyF9:        "f9",
	KeyF10:       "f10",
	KeyF11:       "f11",
	KeyF12:       "f12",
	KeyCtrlA:     "ctrl+a",
	KeyCtrlB:     "ctrl+b",
	KeyCtrlC:     "ctrl+c",
	KeyCtrlD:     "ctrl+d",
	KeyCtrlE:     "ctrl+e",
	KeyCtrlF:     "ctrl+f",
	KeyCtrlG:     "ctrl+g",
	KeyCtrlH:     "ctrl+h",
	KeyCtrlK:     "ctrl+k",
	KeyCtrlL:     "ctrl+l",
	KeyCtrlN:     "ctrl+n",
	KeyCtrlO:     "ctrl+o",
	KeyCtrlP:     "ctrl+p",
	KeyCtrlQ:     "ctrl+q",
	KeyCtrlR:     "ctrl+r",
	KeyCtrlS:     "ctrl+s",
	KeyCtrlT:     "ctrl+t",
	KeyCtrlU:     "ctrl+u",
	KeyCtrlV:     "ctrl+v",
	KeyCtrlW:     "ctrl+w",
	KeyCtrlX:     "ctrl+x",
	KeyCtrlY:     "ctrl+y",
	KeyCtrlZ:     "ctrl+z",
}

// String returns the human-readable name of a key type.
func (k KeyType) String() string {
	if name, ok := keyNames[k]; ok {
		return name
	}
	return "unknown"
}
