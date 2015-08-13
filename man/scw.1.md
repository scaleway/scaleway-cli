% SCW(1) Scaleway-cli User Manuals
% Scaleway Team
% August 2015
# NAME
scw \- Scaleway command line interface

# SYNOPSIS
**scw** [OPTIONS] COMMAND [arg...]

# DESCRIPTION
**scw** is a client to the Scaleway API, through the CLI.

The Scaleway CLI has near 30 commands. The commands are listed below and each
has its own man page which explain usage and arguments.

To see the man page for a command run **man scw <command>**.

# OPTIONS
**-h**, **--help**
  Print usage statement

**-D**, **--debug**=*true*|*false*
  Enable debug mode. Default is false.

**-v**, **--version**=*true*|*false*
  Print version information and quit. Default is false.

**--sensitive**=*true*|*false*
  Show sensitive data in outputs, i.e. API Token/Organization. Default is false.

# COMMANDS
**attach**
  Attach to a server serial console
  See **scw-attach(1)** for full documentation on the **attach** command.

**commit**
  Create a new snapshot from a server's volume
  See **scw-commit(1)** for full documentation on the **commit** command.

**cp**
  Copy files/folders from a PATH on the server to a HOSTDIR on the host
  See **scw-cp(1)** for full documentation on the **cp** command.

**create**
  Create a new server but do not start it
  See **scw-create(1)** for full documentation on the **create** command.

**events**
  Get real time events from the server
  See **scw-events(1)** for full documentation on the **events** command.

**exec**
  Run a command in a running server
  See **scw-exec(1)** for full documentation on the **exec** command.

**history**
  Show the history of an image
  See **scw-history(1)** for full documentation on the **history** command.

**images**
  List images
  See **scw-images(1)** for full documentation on the **images** command.

**info**
  Display system-wide information
  See **scw-info(1)** for full documentation on the **info** command.

**inspect**
  Return low-level information on a server, image, snapshot, volume or bootscript
  See **scw-inspect(1)** for full documentation on the **inspect** command.

**kill**
  Kill a running server
  See **scw-kill(1)** for full documentation on the **kill** command.

**login**
  Log in to Scaleway API
  See **scw-login(1)** for full documentation on the **login** command.

**logout**
  Log out from the Scaleway API
  See **scw-logout(1)** for full documentation on the **logout** command.

**logs**
  Fetch the logs of a server
  See **scw-logs(1)** for full documentation on the **logs** command.

**port**
  Lookup the public-facing port which is NAT-ed to PRIVATE_PORT
  See **scw-port(1)** for full documentation on the **port** command.

**ps**
  List servers
  See **scw-ps(1)** for full documentation on the **ps** command.

**rename**
  Rename a server
  See **scw-ps(1)** for full documentation on the **ps** command.

**restart**
  Restart a running server
  See **scw-restart(1)** for full documentation on the **restart** command.

**rm**
  Remove one or more servers
  See **scw-rm(1)** for full documentation on the **rm** command.

**rmi**
  Remove one or more images
  See **scw-rmi(1)** for full documentation on the **rmi** command.

**run**
  Run a command in a new server
  See **scw-run(1)** for full documentation on the **run** command.

**search**
  Search the Scaleway Hub for images
  See **scw-search(1)** for full documentation on the **search** command.

**start**
  Start a stopped server
  See **scw-start(1)** for full documentation on the **start** command.

**stop**
  Stop a running server
  See **scw-stop(1)** for full documentation on the **stop** command.

**tag**
  Tag a snapshot into an image
  See **scw-tag(1)** for full documentation on the **tag** command.

**top**
  Lookup the running processes of a server
  See **scw-top(1)** for full documentation on the **top** command.

**version**
  Show the Scaleway CLI version information
  See **scw-version(1)** for full documentation on the **version** command.

**wait**
  Block until a server stops
  See **scw-wait(1)** for full documentation on the **wait** command.


#### Client
For specific client examples please see the man page for the specific Scw
command. For example:

    man scw-run

# HISTORY
August 2015, Originally compiled by Scaleway Team based on scaleway.com source material and internal work.
