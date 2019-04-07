package main

import (
	"fmt"
	"syscall"

	libseccomp "github.com/seccomp/libseccomp-golang"
)

func init() {
	var syscalls = []string{"futex", "epoll_pwait", "nanosleep", "read",
		"write", "openat", "epoll_ctl", "close", "rt_sigaction", "mmap",
		"sched_yield", "lstat", "fstat", "mprotect", "rt_sigprocmask",
		"connect", "munmap", "sigaltstack", "set_robust_list", "clone",
		"setsockopt", "socket", "getsockname", "gettid", "getpeername",
		"fcntl", "readlinkat", "getrandom", "newfstatat", "getsockopt",
		"epoll_create1", "brk", "access", "execve", "arch_prctl",
		"sched_getaffinity", "getdents64", "set_tid_address", "prlimit64",
		"exit_group"}
	whiteList(syscalls)

}

// Load the seccomp whitelist.
func whiteList(syscalls []string) {

	filter, err := libseccomp.NewFilter(
		libseccomp.ActErrno.SetReturnCode(int16(syscall.EPERM)))
	if err != nil {
		fmt.Printf("Error creating filter: %s\n", err)
	}
	for _, element := range syscalls {
		// fmt.Printf("[+] Whitelisting: %s\n", element)
		syscallID, err := libseccomp.GetSyscallFromName(element)
		if err != nil {
			panic(err)
		}
		filter.AddRule(syscallID, libseccomp.ActAllow)
	}
	filter.Load()
}
