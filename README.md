# indelible

A side-car service that can run with Linux capability `CAP_LINUX_IMMUTABLE` so that it can set your service's log files to append-only but not read, write to them, nor does it have the capability to remove the append-only setting once set.

This is intended for secure environments where you do not want an attacker to be able to erase log traces post-compromise of your service.

# Preamble

ext3/ext4 file systems have extended attributes (xattr) available on them.
One cool one is the append-only attribute

e.g. 

```
$ touch /tmp/foo.txt
$ chattr +a /tmp/foo.txt
$ setpriv --inh-caps=-linux_immutable --bounding-set=-linux_immutable bash
$ echo "blah" >> /tmp/foo.txt
```

you'll notice this works well but you cannot delete the text once written, only append to it

# Dependencies 

- [Protobufs v3](https://github.com/google/protobuf/releases)
- [Insalling includes](https://gist.github.com/sofyanhadia/37787e5ed098c97919b8c593f0ec44d8)
- [protoc](https://github.com/protocolbuffers/protobuf/releases)
- [installing protoc](https://askubuntu.com/a/1072684)

# Running Tests

```
// Build a test binary
$ go test -c 
// Run a test binary with priviliges as 
// Service needs to run as "CAP_LINUX_IMMUTABLE"
$ sudo ./indelible.test
```

# Running in container

Here we build the service and the example client, we also drop the priviliges in the container

```
$ ./docker_run.sh golang:1.13.11-buster
root@25a0b3ae394a:/go# capsh --print | grep immutable # we have the capability here
root@25a0b3ae394a:/go# cd src/indelible/
root@25a0b3ae394a:/go/src/indelible# go build
go: downloadin
  ...

root@25a0b3ae394a:/go/src/indelible# nohup ./indelible &
root@25a0b3ae394a:/go/src/indelible# cd exampleclient/
root@25a0b3ae394a:/go/src/indelible# setpriv --no-new-privs --inh-caps=-linux_immutable --bounding-set=-linux_immutable bash
root@25a0b3ae394a:/go# capsh --print | grep immutable # we do not have the capability here
root@25a0b3ae394a:/go/src/indelible/exampleclient# go build
go: downloading
  ...
root@25a0b3ae394a:/go/src/indelible/exampleclient# ./exampleclient
Creating log file at /var/log/immutable.log
Requesting log file (/var/log/immutable.log) be marked append-only...
success
root@25a0b3ae394a:/go/src/indelible/exampleclient# echo "test line" > /var/log/immutable.log
bash: /var/log/immutable.log: Operation not permitted
root@25a0b3ae394a:/go/src/indelible/exampleclient# echo "test line" >> /var/log/immutable.log
root@25a0b3ae394a:/go/src/indelible/exampleclient# cat /var/log/immutable.log
test line
```


# References

- https://blog.fpmurphy.com/2009/05/linux-security-capabilities.html
- http://www.andy-pearce.com/blog/posts/2013/Mar/file-capabilities-in-linux/
