% SCW(1) Scw User Manuals
% Scaleway Community
% JULY 2015
# NAME
scw-attach - Attach to a running server

# SYNOPSIS
**scw attach**
[**--help**]/
[**--no-stdin**[=*false*]]
SERVER

# DESCRIPTION
The **scw attach** command allows you to attach to a running server using
the server's ID or name, either to view its ongoing output or to control it
interactively.  You can't attach to the same server multiple times
simultaneously.

You can detach from the server (and leave it running) with `CTRL-q`
(for a quiet exit).

# OPTIONS
**--help**
  Print usage statement

**--no-stdin**=*true*|*false*
   Do not attach STDIN. The default is *false*.

# EXAMPLES

## Attaching to a server

In this example the top command is run inside a server, from an image called
fedora, in detached mode. The ID from the server is passed into the **scw
attach** command:

    # ID=$(scw run -d fedora)
    # scw attach $ID
    Booting Linux on physical CPU 0
    Initializing cgroup subsys cpuset
    Initializing cgroup subsys cpu
    Linux version 3.2.34-30 (build-bot@cloud.online.net) (gcc version 4.9.1 (Ubuntu/Linaro 4.9.1-10ubuntu2) ) #17 SMP Mon Apr 13 15:53:45 UTC 2015
    CPU: Marvell PJ4Bv7 Processor [562f5842] revision 2 (ARMv7), cr=10c53c7d
    CPU: PIPT / VIPT nonaliasing data cache, PIPT instruction cache
    Machine: Online Labs C1
    Using UBoot passing parameters structure
    [...]

# HISTORY
July 2015, Originally compiled by Scaleway Team
