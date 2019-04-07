# You can get a list of syscalls via strace:
# $ strace -qcf ./team-suse

dump = """\
futex
epoll_pwait
nanosleep
read
write
openat
epoll_ctl
close
rt_sigaction
mmap
sched_yield
lstat
fstat
mprotect
rt_sigprocmask
connect
munmap
sigaltstack
set_robust_list
clone
setsockopt
socket
getsockname
gettid
getpeername
fcntl
readlinkat
getrandom
newfstatat
getsockopt
epoll_create1
brk
access
execve
arch_prctl
sched_getaffinity
getdents64
set_tid_address
prlimit64"""

whitelist = dump.split("\n")
whitelist.append("exit_group")  # I guess we alwas need to exit the program
output = ['"{}"'.format(elem) for elem in whitelist ]
output = ", ".join(output)

print(output)