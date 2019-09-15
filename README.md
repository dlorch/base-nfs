<p align="center">
    <img height="250" width="250" src="https://github.com/dlorch/base-nfs/blob/master/base-nfs.png?raw=true">
</p>

[![pipeline status](https://gitlab.com/dlorch/base-nfs/badges/master/pipeline.svg)](https://gitlab.com/dlorch/base-nfs/commits/master)

A fully functional NFS (Network File System) Version 3 server
running an in-memory file system. Includes auxiliary services
like portmap and mount.

## Development

Following `make` targets are available. For some targets, [Docker]
and [Docker Compose] are necessary.

```
$ make
integration-logs               show logs from nfs-server [requires Docker Compose]
integration-setup              build docker images for integration tests [requires Docker Compose]
integration-shell              enter shell on tester [requires Docker Compose]
integration-teardown           destroy resources associated to integration tests [requires Docker Compose]
integration                    run all integration tests [requires Docker Compose]
unittests                      run all unit tests
```

### General Tips for Debugging

Retrieve logs from NFS server:

```
$ make integration-logs
```

Enter a shell. The NFS server is available as `nfs-server`:

```
$ make integration-shell
# showmount -e nfs-server  # inside the shell
```

Cleanup (removes containers and logs):

```
$ make integration-teardown
```

### Tips for Debugging Protocol Problems

Enter a shell and use `tcpdump` to record traffic, which can later
be analyzed using [Wireshark].

```
$ make integration-shell
# tpcdump -i eth0 -w /tests/dump.pcap &  # run in background
```

This will capture the network traffic to the NFS server. The
directory `/tests/` is mounted to the host system. Now run your
commands:

```
# showmount -e nfs-server
```

Get `tcpdump` back in the foreground, and close it with CTRL-C:

```
# fg
^C
```

You can now analyze `test/dump.pcap` with [Wireshark].

### Relevant RFCs for this project

* [RFC1057] RPC: Remote Procedure Call Protocol Specification Version 2
* [RFC1813] NFS Version 3 Protocol Specification
* [RFC1014] XDR: External Data Representation Standard

### Golang concepts and best practices considered

* [Standard Go Project Layout]
* [What's in a name?]
* [Twelve Go Best Practices]
* [Object Oriented Inheritance in Go]

[Docker]: https://www.docker.com/
[Docker Compose]: https://docs.docker.com/compose/
[Wireshark]: https://www.wireshark.org/
[Standard Go Project Layout]: https://github.com/golang-standards/project-layout
[What's in a name?]: https://talks.golang.org/2014/names.slide
[Twelve Go Best Practices]: https://talks.golang.org/2013/bestpractices.slide
[Object Oriented Inheritance in Go]: https://hackthology.com/object-oriented-inheritance-in-go.html
[RFC1057]: https://tools.ietf.org/html/rfc1057
[RFC1813]: https://tools.ietf.org/html/rfc1813
[RFC1014]: https://tools.ietf.org/html/rfc1014