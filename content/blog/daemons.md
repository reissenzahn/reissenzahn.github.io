+++
title = "Daemons"
date = "2024-04-06"
tags = ["Linux"]
subtitle = "Some notes from TLPI 37, excluding 37.5"
+++

## Overview
- A daemon is a long-lived process that runs in the background and has no controlling terminal.
- The lack of a controlling terminal ensures that the kernel never automatically generates any job-control or terminal-related signals for the daemon.
- Examples of daemons include: *cron*, *sshd*, *httpd*, *inetd*.
- Many standard daemons run as privileged processes with effective user ID of 0.
- Certain daemons are run as kernel threads such as *pdflush*, which periodically flushes dirty pages to disk.
- Since daemons are long-lived, it is necessary to be particularly wary of possible memory and file descriptor leaks.

## Creating a Daemon
1. Perform a `fork()`, after which the parent exits and the child becomes a child of the init process. The child is guaranteed not to be a process group leader since it inherited its process group ID from its parent and has its own unique process ID.
2. The child process calls `setsid()` to start a new session and free itself of any association with a controlling terminal.
3. If the daemon might later open a terminal device then one of two approaches can be used to prevent the device from becoming the controlling terminal:
    - Specify the `O_NOCTTY` flag on any `open()` that may apply to a terminal device.
    - Perform a second `fork()` and again have the parent exit and the child continue to ensure that the child is not the session leader (which prevents it from reacquiring a controlling terminal).
4. Clear the process *umask* to ensure that the daemon creates files and directories with the requested permissions.
5. Change the current working directory of the process and close unnecessary open file descriptors inherited from the parent to avoid preventing a file system from being unmounted.
6. After closing file descriptors 0, 1 and 2, open `/dev/null` and use `dup2()` to make all those descriptors refer to this device to prevent issues with library functions that make assumptions about these file descriptors.
```c
// become_daemon.c
#include <syslog.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>

#define BD_NO_CHDIR 01
#define BD_NO_CLOSE_FILES 02
#define BD_NO_REOPEN_STD_FDS 04
#define BD_NO_UMASK0 010
#define BD_MAX_CLOSE 8192

int become_daemon(int flags) {
  // become a background process
  switch (fork()) {
    case -1: return -1;
    case 0: break;  // child
    default: _exit(EXIT_SUCCESS);  // parent
  }

  // become leader of new session
  if (setsid() == -1) return -1;

  // ensure we are not session leader
  switch (fork()) {
    case -1: return -1;
    case 0: break;
    default: _exit(EXIT_SUCCESS);
  }
  
  // clear file mode creation mask
  if (!(flags & BD_NO_UMASK0)) umask(0);

  // change to root directory
  if (!(flags & BD_NO_CHDIR)) chdir("/");

  // close all open files
  if (!(flags & BD_NO_CLOSE_FILES)) {
    int maxfd = sysconf(_SC_OPEN_MAX);
    if (maxfd == -1) {
      maxfd = BD_MAX_CLOSE;
    }

    for (int fd = 0; fd < maxfd; fd++) {
      close(fd);
    }
  }

  // reopen standard file descriptors to /dev/null
  if (!(flags & BD_NO_REOPEN_STD_FDS)) {
    close(STDIN_FILENO);
    int fd = open("/dev/null", O_RDWR);
    if (fd != STDIN_FILENO) return -1;
    if (dup2(STDIN_FILENO, STDOUT_FILENO) != STDOUT_FILENO) return -1;
    if (dup2(STDIN_FILENO, STDERR_FILENO) != STDERR_FILENO) return -1;
  }

  return 0;
}

int main(int argc, char *argv[]) {
    become_daemon(0);
    sleep(20);
    exit(EXIT_SUCCESS);
}
```
```bash
gcc become_daemon.c -o become_daemon
./become_daemon
ps -C become_daemon -o "pid ppid pgid sid tty command"
```

## Handling SIGHUP
- The `SIGHUP` signal provides a way to prompt a daemon to reinitialize itself by rereading its configuration file and reopening any log files it may be using.
- Some daemons close all their files and restart themselves with an `exec()` on receipt of `SIGHUP`.
- A `SIGHUP` signal is generated for the controlling process on disconnection of a controlling terminal. Since a daemon has no controlling terminal, the kernel never generates this signal for a daemon.
```c
#include <sys/stat.h>
#include <signal.h>
#include <stdio.h>

#define INTERVAL 15

static volatile sig_atomic_t hup_received = 0;

static void sighup_handler(int sig) {
  hup_received = 1;
}

int main(int argc, char *argv[]) {
  struct sigaction sa;
  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_RESTART;
  sa.sa_handler = sighup_handler;
  if (sigaction(SIGHUP, &sa, NULL) == -1) return -1;
  
  if (become_daemon(0) == -1) return -1;
  
  int unslept = INTERVAL;
  for (;;) {
    unslept = sleep(unslept);

    if (hup_received) {
      // re-initialize
      hup_received = 0;
    }

    if (unslept == 0) {
      unslept = INTERVAL;
    }
  }
}
```
```bash
killall -HUP daemon_SIGHUP
killall daemon_SIGHUP
```

## Handling SIGTERM
- A daemon typically terminates only when the system shuts down.
- Many daemons are stopped by application-specific scripts executed during system shutdowns while others will simply receive a `SIGTERM` sent by the init process during system shutdown.
- By default, `SIGTERM` terminates a process. If the daemon needs to perform any cleanup before terminating, it should do so by establishing a handler for this signal.
- This handler must be designed to perform such cleanup quickly, since init follows up the `SIGTERM` signal with a `SIGKILL` signal after 5 seconds.
