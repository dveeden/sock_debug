# Overview

This is a tool to inspect the MySQL Protocol as it is transported over a UNIX domain socket.

# Building

```
go build
```

# Running

```
Usage of ./sock_debug:
  -serverSocket string
    	socket to connect to (default "/tmp/tidb.sock")
```

Start `./sock_debug` and it will start to listen on `/tmp/sock_debug.socket` and by default connect to `/tmp/tidb.sock`. Communication is printed on the console.
