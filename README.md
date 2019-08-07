# netfile

A Go CLI project, currently targeting Linux, that allows a file server to be deployed on any given machine which also houses sub commands for it to act as the client to connect to the server instance to make file requests.

## Current State

It's a work in progress...currently trying to get the basic stuff like `netfile server --port x`, `netfile fetch --host addr --port x`, & `netfile store --host addr --port x` working with a solid communication protocol.

The flags for specifying file requests i.e. the file name isn't set up just yet. Still working on the file transfer protocol. It's almost done. Once that is in place...will add the `--file` flag so a file/filepath can be specified.

Currently the file area the server manages is sandboxed to `/var/lib/netfile/files` on the server machine...but once the basic protocols are finished...it will be expanded so that it can track the files it has operated on and allow them to be distributed across the filesystem if the user so chooses.

## Packages

There is really only the `cmd` package...at:

### [github.com/steviesama/netfile/cmd](https://github.com/steviesama/netfile/tree/master/cmd)

...however, each command is essentially a representation of a package or functionality in [github.com/steviesama/nx/service/netfile](https://github.com/steviesama/nx/blob/master/service/netfile).

## Commands

```bash
netfile server [--port x]
```

Run this without the flags and it will start a file server listening on port 8010 by default if the port is available.

The first time it is run...the `/var/lib/netfile` folder structure will not exist yet; it will gave an error and ask that it be rerun with `sudo`. It'll create the structure and set the sudo user as the owner. Once done...rerun it normally.

```bash
netfile fetch --host addr --port x --file path
```

Run this specifying the host address and port number to connect on as well as a file path that is relative to where it is stored on the file server to initiate a download of the requested file if it exists.

Currently, the `--file` flag is not implemented.


```bash
netfile store --host addr --port x --file
```

Currently none of this command is implemented.

## Roadmap

Ultimately, this tool is going to be used to allow file backup mapping as well as deploying and managing files.

Other than obviously the server and client file storage and fetching functionality, config files...and the ability for the user to either be able to specify the config manually or to set varous config settings via a `netfile config` command and various flags would be one of the next features.

The ability to authenticate the user as well and connect via TLS via JSON Web Tokens will be an eventual feature.

May add Windows support in the future...but it is low priority as it can already be used via Windows Subsystem for Linux.