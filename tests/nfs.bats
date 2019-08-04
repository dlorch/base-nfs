#!/usr/bin/env bats

@test "show remote NFS mounts" {
  run showmount -e nfs-server
  [ $status -eq 0 ]
}

@test "mount directory" {
  run mount -o nolock,nfsvers=3 nfs-server:/volume1/Public /mnt
  [ $status -eq 0 ]
}

@test "list directory" {
  run ls -al /mnt
  [ "${lines[0]}" == 'total 20' ]
  [ "${lines[1]}" == 'drwxrwxrwx    4 root     root          4096 Sep 16  2018 .' ]
  [ "${lines[3]}" == '-rw-rw-rw-    1 1027     users          292 Jan 15  2014 gopher.go' ]
}

@test "cat file" {
  run cat /mnt/gopher.go
  [ $status -eq 0 ]
}

@test "create new file" {
  run touch /mnt/hello.txt
  [ $status -eq 0 ]
}

@test "write to file" {
  run echo "Hello, NFS" > /mnt/hello.txt
  [ $status -eq 0 ]
}

@test "append to file" {
  run echo "Another line" >> /mnt/hello.txt
  [ $status -eq 0 ]
}

@test "delete file" {
  run rm -f /mnt/hello.txt
  [ $status -eq 0 ]
}

@test "create new directory" {
  run mkdir /mnt/new_directory/
  [ $status -eq 0 ]
}

@test "delete directory" {
  run rmdir /mnt/new_directory/
  [ $status -eq 0 ]
}