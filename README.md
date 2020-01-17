# indelible

A service to create an append-only log file to which it cannot read or clobbered.

# Preamble

ext3/ext4 file systems have extended attributes (xattr) available on them.
One of value is the append-only attribute

e.g. 

```
$ touch /tmp/foo.txt
$ chattr +a /tmp/foo.txt
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

# Deployment

```
$ sudo setcap cap_linux_immutable,cap_net_raw=ep indelible
$ getcap indelible
indelible = cap_linux_immutable,cap_net_raw+ep
// eip stands for effective,inheritable,permitted
```

then you can run indelible as non-sudo/non-root and it will still have the append only and sock creation abilities.

# References

- https://blog.fpmurphy.com/2009/05/linux-security-capabilities.html
- http://www.andy-pearce.com/blog/posts/2013/Mar/file-capabilities-in-linux/
