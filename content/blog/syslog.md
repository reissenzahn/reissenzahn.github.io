+++
title = "Syslog"
date = "2024-04-06"
tags = ["Linux"]
+++

*Some notes from TLPI 37.5*

### Overview
- The syslog facility provides a single, centralized logging facility that can be used to log messages by all applications on the system.
- It has two principal components: the *syslogd* daemon and the *syslog(3)* library function.
- *syslogd* accepts log messages from a UNIX domain socket `/dev/log` which holds locally produced messages and (if enabled) an Internet domain socket UDP port 514 which holds messages sent across a TCP/IP network.
- Each message has a number of attributes including a facility which specifies the type of program generating the message and a severity level.
- The daemon examines the facility and level of each message and then passes it along to any of several possible destinations as configured in `/etc/syslog.conf`.
- Possible destinations include a terminal, disk file, FIFO, one or more logged-in users or a process on another system connected via a TCP/IP network. A single message may be sent to multiple destinations (or none at all).
- The *syslog(3)* library function uses its supplied arguments to construct a message in a standard format that is then placed on the `/dev/log` socket.
- The *klogd* daemon collects kernel log messages (produced by the kernel using `printk()`) using the `/proc/kmsg` file or the *syslog(2)* system call and places them on `/dev/log` using *syslog(3)*.
- The *logger(1)* shell command can be used to add entries to the system log.
```c
openlog(argv[0], LOG_PID | LOG_CONS | LOG_NOWAIT, LOG_LOCALO);

syslog(LOG_ERROR, "Bad argument: %s", argv[1]);
syslog(LOG_USER | LOG_INFO, "Exiting");

setlogmask(LOG_MASK(LOG_EMERG) | LOG_MASK(LOG_ALERT) | LOG_MASK(LOG_CRIT) | LOG_MASK(LOG_ERR));

setlogmask(LOG_UPTO(LOG_ERR));

closelog();
```

### openlog()
```c
#include <syslog.h>

void openlog(const char *ident, int log_options, int facility);
```
- `openlog()` establishes a connection to the system log facility and sets defaults that apply to subsequent `syslog()` calls. The use of `openlog()` is optional.
- `ident` is a pointer to a string (typically the program name) that is included in each message written by `syslog()`. As long as it continues to call `syslog()`, the application should ensure that the referenced string is not later changed.
- `log_options` is a bit mask created by ORing together options, including:
   - `LOG_NDELAY`: Open the connection to `/dev/log` immediately rather than when the first message is logged.
   - `LOG_PID`: Log the process ID of the caller with each message.
   - `LOG_PERROR`: Write messages to standard error as well as to the system log.
- `facility` specifies the default facility value to be used in subsequent calls to `syslog()`: `LOG_USER` (default), `LOG_AUTH`, `LOG_KERN`, etc.

### syslog()
```c
#include <syslog.h>

void syslog(int priority, const char *format, ...);
```
- `syslog()` writes a log message.
- `priority` is created by ORing together a facility value and a level value.
- Available level values are: `LOG_EMERG`, `LOG_ALERT`, `LOG_CRIT`, `LOG_ERR`, `LOG_WARNING`, `LOG_NOTICE`, `LOG_INFO`, `LOG_DEBUG`.
- `format` is a format string that is used with the following arguments in the manner of `printf()`.
- The format string does not need to include a terminating newline and may include the sequence `%m` which is replaced by the equivalent of `strerror(errno)`.

### setlogmask()
`setlogmask()` sets a mask that filters the messages written by `syslog()`; any message whose level is not included in the current mask setting is discarded. Returns previous log priority mask.
```c
#include <syslog.h>
int setlogmask(int mask_priority);
```
- `mask_priority` is a bit mask obtained by ORing together the results of applying `LOG_MASK()` to level values.
- `LOG_UPTO()` creates a bit mask filtering all messages of a certain level and above.

### closelog()
```c
#include <syslog.h>

void closelog(void);
```
- `closelog()` deallocates the file descriptor used for the `/dev/log` socket.

### /etc/syslog.conf
- `/etc/syslog.conf` controls the operation of `syslogd` by specifying rules of the form *facility.level action*.
- Together, the facility and level are referred to as the selector as they select the messages to which the rule applies.
- The action specifies where to send the messages matching this selector.
- A `SIGHUP` signal can be sent to *syslogd* to ask it to reload the `/etc/syslog.conf` file.
```txt
# messages from all facilities with level of LOG_ERR or higher should be sent to the /dev/tty10 console device
*.err /dev/tty10

# LOG_AUTH messages with level of LOG_NOTICE or higher should be sent to any consoles or terminals where root is logged in
auth.notice root

# all messages except those for the LOG_MAIL and LOG_NEWS facilities should be sent to the file /var/log/messages (the - indicates that a sync to disk does not occur on each write to the file)
*.debug;mail.none;news.none -/var/log/messages
```
```bash
killall -HUP syslogd
```
