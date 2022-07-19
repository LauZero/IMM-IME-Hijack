package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	imm32                   = windows.NewLazySystemDLL("imm32.dll")
	procImmGetContext       = imm32.NewProc("ImmGetContext")
	hHook                   HHOOK
	hHinstance              HINSTANCE = NULL
)

const (
	WH_KEYBOARD        = 2
	WH_GETMESSAGE      = 3
	WH_KEYBOARD_LL     = 13
	WM_KEYDOWN         = 256
	WM_KEYFIRST        = 256
	WM_KEYUP           = 257
	WM_CHAR            = 258
	WM_SYSKEYDOWN      = 260
	WM_SYSKEYUP        = 261
	WM_KEYLAST         = 264
	WM_IME_COMPOSITION = 271
	PM_NOREMOVE        = 0x000
	PM_REMOVE          = 0x001
	PM_NOYIELD         = 0x002
	WM_LBUTTONDOWN     = 513
	WM_RBUTTONDOWN     = 516
	HC_ACTION          = 0
	NULL               = 0
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/dd162805.aspx
type POINT struct {
	X, Y int32
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644958.aspx
type MSG struct {
	Hwnd     HWND
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       POINT
	LPrivate uint32
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func MessageProc(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
	fmt.Println(nCode)
	if nCode >= HC_ACTION { // && wparam == WM_KEYDOWN { // 有键按下
		// TODO: lparam -> MSG
		msg := (*MSG)(unsafe.Pointer(lparam))
		fmt.Println(msg.Message, msg.LParam, msg.WParam, msg.Pt, msg.Time, msg.LPrivate)
		fmt.Println(msg.Hwnd)
		hIMC, _, _ := procImmGetContext.Call(uintptr(msg.Hwnd))
		fmt.Println(">>", hIMC)

		kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		fmt.Println(kbdstruct.VkCode, kbdstruct.Flags, kbdstruct.ScanCode)
		fmt.Println("==")
		switch msg.Message { // TODO: ?? lparam.MESSAGE
		case WM_IME_COMPOSITION:
			fmt.Println("IMM")
		case WM_CHAR:
			fmt.Println("CHAR")
		}
	}
	return CallNextHookEx(hHook, nCode, wparam, lparam) // 将消息还给钩子链，不要影响别人
}

func Start() {
	// defer user32.Release()
	hHook = SetWindowsHookEx(WH_GETMESSAGE, (HOOKPROC)(MessageProc), hHinstance, 0)
	// var msg MSG
	// for GetMessage(&msg, 0, 0, 0) != 0 {
	// 	fmt.Println("MSG", msg.Message)
	// }

	// UnhookWindowsHookEx(hHook)
	// hHook = 0
}

func main() {
	Start()
	time.Sleep(time.Minute)
}
