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

Go Concepts:
* [Object Oriented Inheritance in Go]

[Object Oriented Inheritance in Go]: https://hackthology.com/object-oriented-inheritance-in-go.html