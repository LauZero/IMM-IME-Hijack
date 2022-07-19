package main

import (
	"syscall"
	"time"
)

var (
	dllMain, _         = syscall.LoadLibrary("C:\\Users\\...\\Desktop\\DLLMain.dll")
	procInstallHook, _ = syscall.GetProcAddress(dllMain, "InstallHook")
	procUnHook, _      = syscall.GetProcAddress(dllMain, "UnHook")
)

func main() {
	var nargs uintptr = 0
	syscall.Syscall(uintptr(procInstallHook), nargs, 0, 0, 0)
	// defer syscall.Syscall(uintptr(procUnHook), nargs, 0, 0, 0)
	time.Sleep(time.Minute)
}
