```
$ sudo showmount -e 172.17.0.2
Exports list on 127.0.0.1:
/volume1/Public                     *
$ sudo mount -o nolock,nfsvers=3 172.17.0.2:/volume1/Public /mnt
$ mount | grep Public
172.17.0.2:/volume1/Public on /mnt type nfs (rw,relatime,vers=3,rsize=131072,wsize=131072,namlen=255,hard,nolock,proto=tcp,timeo=600,retrans=2,sec=sys,mountaddr=172.17.0.2,mountvers=3,mountport=892,mountproto=tcp,local_lock=all,addr=172.17.0.2)
$ ls /mnt
gopher.go
```

Golang concepts and best practices considered:
* [Standard Go Project Layout]
* [What's in a name?]
* [Twelve Go Best Practices]
* [Object Oriented Inheritance in Go]

[Standard Go Project Layout]: https://github.com/golang-standards/project-layout
[What's in a name?]: https://talks.golang.org/2014/names.slide
[Twelve Go Best Practices]: https://talks.golang.org/2013/bestpractices.slide
[Object Oriented Inheritance in Go]: https://hackthology.com/object-oriented-inheritance-in-go.html