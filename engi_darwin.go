// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows,!linux,!netgo,!android

package engi

import (
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glow/gl/2.1/gl"
	"log"
	"runtime"
)

var window *glfw.Window

func run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		log.Fatal("GLFW error %v: %v\n", err, desc)
	})

	if ok := glfw.Init(); ok {
		defer glfw.Terminate()
	}

	if !config.Resizable {
		glfw.WindowHint(glfw.Resizable, 0)
	}
	glfw.WindowHint(glfw.Samples, config.Fsaa)

	width := config.Width
	height := config.Height

	monitor, err := glfw.GetPrimaryMonitor()
	if err != nil {
		log.Fatal(err)
	}
	mode, err := monitor.GetVideoMode()
	if err != nil {
		log.Fatal(err)
	}

	if config.Fullscreen {
		width = mode.Width
		height = mode.Height
		glfw.WindowHint(glfw.Decorated, 0)
	} else {
		monitor = nil
	}

	title := config.Title

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	config.Width, config.Height = window.GetSize()

	if !config.Fullscreen {
		window.SetPosition((mode.Width-width)/2, (mode.Height-height)/2)
	}

	if config.Vsync {
		glfw.SwapInterval(1)
	}

	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	GL = newgl2()

	GL.Viewport(0, 0, config.Width, config.Height)

	window.SetSizeCallback(func(window *glfw.Window, w, h int) {
		config.Width, config.Height = window.GetSize()
		responder.Resize(w, h)
	})

	window.SetCursorPositionCallback(func(window *glfw.Window, x, y float64) {
		responder.Mouse(float32(x), float32(y), MOVE)
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPosition()
		if a == glfw.Press {
			responder.Mouse(float32(x), float32(y), PRESS)
		} else {
			responder.Mouse(float32(x), float32(y), RELEASE)
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		responder.Scroll(float32(yoff))
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		if a == glfw.Press {
			responder.Key(Key(k), Modifier(m), PRESS)
		} else {
			responder.Key(Key(k), Modifier(m), RELEASE)
		}
	})

	window.SetCharacterCallback(func(window *glfw.Window, char uint) {
		responder.Type(rune(char))
	})

	responder.Preload()

	timing = NewStats(config.LogFPS)
	timing.Update()

	Files.Load(func() {})

	responder.Setup()

	for !window.ShouldClose() {
		responder.Update(float32(timing.Dt))
		GL.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.SwapBuffers()
		glfw.PollEvents()
		timing.Update()
	}
}

func exit() {
	window.SetShouldClose(true)
}

const (
	Dash         = Key(glfw.KeyMinus)
	Apostrophe   = Key(glfw.KeyApostrophe)
	Semicolon    = Key(glfw.KeySemicolon)
	Equals       = Key(glfw.KeyEqual)
	Comma        = Key(glfw.KeyComma)
	Period       = Key(glfw.KeyPeriod)
	Slash        = Key(glfw.KeySlash)
	Backslash    = Key(glfw.KeyBackslash)
	Backspace    = Key(glfw.KeyBackspace)
	Tab          = Key(glfw.KeyTab)
	CapsLock     = Key(glfw.KeyCapsLock)
	Space        = Key(glfw.KeySpace)
	Enter        = Key(glfw.KeyEnter)
	Escape       = Key(glfw.KeyEscape)
	Insert       = Key(glfw.KeyInsert)
	PrintScreen  = Key(glfw.KeyPrintScreen)
	Delete       = Key(glfw.KeyDelete)
	PageUp       = Key(glfw.KeyPageUp)
	PageDown     = Key(glfw.KeyPageDown)
	Home         = Key(glfw.KeyHome)
	End          = Key(glfw.KeyEnd)
	Pause        = Key(glfw.KeyPause)
	ScrollLock   = Key(glfw.KeyScrollLock)
	ArrowLeft    = Key(glfw.KeyLeft)
	ArrowRight   = Key(glfw.KeyRight)
	ArrowDown    = Key(glfw.KeyDown)
	ArrowUp      = Key(glfw.KeyUp)
	LeftBracket  = Key(glfw.KeyLeftBracket)
	LeftShift    = Key(glfw.KeyLeftShift)
	LeftControl  = Key(glfw.KeyLeftControl)
	LeftSuper    = Key(glfw.KeyLeftSuper)
	LeftAlt      = Key(glfw.KeyLeftAlt)
	RightBracket = Key(glfw.KeyRightBracket)
	RightShift   = Key(glfw.KeyRightShift)
	RightControl = Key(glfw.KeyRightControl)
	RightSuper   = Key(glfw.KeyRightSuper)
	RightAlt     = Key(glfw.KeyRightAlt)
	Zero         = Key(glfw.Key0)
	One          = Key(glfw.Key1)
	Two          = Key(glfw.Key2)
	Three        = Key(glfw.Key3)
	Four         = Key(glfw.Key4)
	Five         = Key(glfw.Key5)
	Six          = Key(glfw.Key6)
	Seven        = Key(glfw.Key7)
	Eight        = Key(glfw.Key8)
	Nine         = Key(glfw.Key9)
	F1           = Key(glfw.KeyF1)
	F2           = Key(glfw.KeyF2)
	F3           = Key(glfw.KeyF3)
	F4           = Key(glfw.KeyF4)
	F5           = Key(glfw.KeyF5)
	F6           = Key(glfw.KeyF6)
	F7           = Key(glfw.KeyF7)
	F8           = Key(glfw.KeyF8)
	F9           = Key(glfw.KeyF9)
	F10          = Key(glfw.KeyF10)
	F11          = Key(glfw.KeyF11)
	F12          = Key(glfw.KeyF12)
	A            = Key(glfw.KeyA)
	B            = Key(glfw.KeyB)
	C            = Key(glfw.KeyC)
	D            = Key(glfw.KeyD)
	E            = Key(glfw.KeyE)
	F            = Key(glfw.KeyF)
	G            = Key(glfw.KeyG)
	H            = Key(glfw.KeyH)
	I            = Key(glfw.KeyI)
	J            = Key(glfw.KeyJ)
	K            = Key(glfw.KeyK)
	L            = Key(glfw.KeyL)
	M            = Key(glfw.KeyM)
	N            = Key(glfw.KeyN)
	O            = Key(glfw.KeyO)
	P            = Key(glfw.KeyP)
	Q            = Key(glfw.KeyQ)
	R            = Key(glfw.KeyR)
	S            = Key(glfw.KeyS)
	T            = Key(glfw.KeyT)
	U            = Key(glfw.KeyU)
	V            = Key(glfw.KeyV)
	W            = Key(glfw.KeyW)
	X            = Key(glfw.KeyX)
	Y            = Key(glfw.KeyY)
	Z            = Key(glfw.KeyZ)
	NumLock      = Key(glfw.KeyNumLock)
	NumMultiply  = Key(glfw.KeyKpMultiply)
	NumDivide    = Key(glfw.KeyKpDivide)
	NumAdd       = Key(glfw.KeyKpAdd)
	NumSubtract  = Key(glfw.KeyKpSubtract)
	NumZero      = Key(glfw.KeyKp0)
	NumOne       = Key(glfw.KeyKp1)
	NumTwo       = Key(glfw.KeyKp2)
	NumThree     = Key(glfw.KeyKp3)
	NumFour      = Key(glfw.KeyKp4)
	NumFive      = Key(glfw.KeyKp5)
	NumSix       = Key(glfw.KeyKp6)
	NumSeven     = Key(glfw.KeyKp7)
	NumEight     = Key(glfw.KeyKp8)
	NumNine      = Key(glfw.KeyKp9)
	NumDecimal   = Key(glfw.KeyKpDecimal)
	NumEnter     = Key(glfw.KeyKpEnter)
)