// A wrapper of AutoIt (AutoItX) for the Go Programming Language

// A wrapper of AutoIt (AutoItX) for the Go Programming Language.
//
// Dependencies
//     AutoIt (with AutoItX) from http://www.autoitscript.com/site/autoit/downloads/
//
// Example
//     package main
//
//     import "github.com/brunoqc/go-autoit"
//
//     func main() {
//         success, pid := autoit.Run("notepad.exe", "", autoit.SW_NORMAL)
//         if !success {
//         	log.Panic("can't run process")
//         } else {
//         	log.Println("pid", pid)
//         }
//     }
//
// Build
//     set CGO_CFLAGS=-Ic:/AutoIt3/AutoItX/StandardDLL/VC6
//     set CGO_LDFLAGS=-lAutoItX3
//     set CGO_LDFLAGS=-lAutoItX3_x64 # for 64-bit
//     go build
package autoit

/*
#include <Windows.h>
#include <AutoIt3.h>
*/
import "C"

import (
	"encoding/binary"
	"syscall"
)

const (
	SwHide     = C.SW_HIDE     // Hidden window
	SwMinimize = C.SW_MINIMIZE // Minimized window
	SwMaximize = C.SW_MAXIMIZE // Maximized window
	SwNormal   = 4
)

const (
	EnableUserInput  = 0
	DisabelUserInput = 1
)

const (
	StateExists    = 1
	StateVisible   = 2
	StateEnabled   = 4
	StateActive    = 8
	StateMinimized = 16
	StateMaximized = 32
)

// Run a program and don't wait
// Possibles flags are SW_HIDE, SW_MINIMIZE, SW_MAXIMIZE and SW_NORMAL
// returns true on success with the pid
func Run(filename, workingdir string, flag int) (bool, int) {
	pid := C.AU3_Run((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(filename)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(workingdir)), C.long(flag))
	return C.AU3_error() == 0, int(pid)
}

// BlockInput blocks the keyboard and mouse
func BlockInput(flag int) {
	C.AU3_BlockInput(C.long(flag))
}

// WinClose closes a window
func WinClose(title, text string) {
	C.AU3_WinClose((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)))
}

// WinGetState returns a window's state
func WinGetState(title, text string) (bool, int) {
	result := C.AU3_WinGetState((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)))
	return C.AU3_error() == 0, int(result)
}

// WinActive returns true if the window is active
func WinActive(title, text string) bool {
	return int(C.AU3_WinActive((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)))) == 1
}

// WinExists returns true if the window exist
func WinExists(title, text string) bool {
	return int(C.AU3_WinExists((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)))) == 1
}

// WinGetText returns the text contained in a window
func WinGetText(title, text string, bufSize int) (result string) {
	// TODO: test if bufSize is not greater than 64KB
	if bufSize < 1 {
		panic("bufSize must be greater than 0")
	}

	data := make([]uint16, bufSize)

	C.AU3_WinGetText((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)), (*_Ctype_WCHAR)(&data[0]), (C.int)(bufSize))

	for _, char := range data {
		if char == 0x0 {
			break
		}

		buf := make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, char)

		// FIXME: shoudln't have to only use the first byte
		result += string(buf[0])
	}

	return
}

// WinActivate set the focus on a window
func WinActivate(title, text string) {
	C.AU3_WinActivate((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)))
}

// Send simulates input on the keyboard
// flag: 0: normal, 1: raw
func Send(keys string, flag int) {
	C.AU3_Send((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(keys)), C.long(flag))
}

// PixelGetColor returns the color of the pixel at the specified location
// return -1 if the location is invalid
func PixelGetColor(x, y int) int {
	return int(C.AU3_PixelGetColor(C.long(x), C.long(y)))
}

// Opt is used to set/get a property
func Opt(option string, param int) int {
	return int(C.AU3_Opt((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(option)), C.long(param)))
}

// ControlClick clicks on a control without using the mouse pointer
// TODO: x, y should be center by defaut
func ControlClick(title, text, controlID, button string, clicks, x, y int) int {
	return int(C.AU3_ControlClick((*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(title)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(text)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(controlID)), (*_Ctype_WCHAR)(syscall.StringToUTF16Ptr(button)), C.long(clicks), C.long(x), C.long(y)))
}

// PixelChecksum returns a checksum of the pixel in a region
func PixelChecksum(left, top, right, bottom, step int) int64 {
	return int64(C.AU3_PixelChecksum(C.long(left), C.long(top), C.long(right), C.long(bottom), C.long(step)))
}

// MouseMove moves the mouse's pointer to a specific location
func MouseMove(x, y, speed int) {
	C.AU3_MouseMove(C.long(x), C.long(y), C.long(speed))
}
