# chatserver

`chatserver` is a simple TCP chat server written in Go that can easily be used
with telnet as a client.

# Installation

Simply clone and build the go source in your preferred way. It has no external
dependencies.

# Server usage

`chatserver [-config]`

In the absence of a `-config` path the server uses a default config:

- host `localhost`
- port `3001`
- logLocation `/tmp/chatserver.log`

You can change the config using a JSON config file located at `-config`. Sample
JSON:

```
{
    "host": "localhost",
    "port": 3001,
    "logLocation": "/tmp/chatserver.log"
}
```

All messages are logged in the file specified at `logLocation`.

# Client usage

Using telnet as a client, a sample telnet session:

```
$ telnet localhost 3001
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
Please enter your name: bobby jones
Welcome bobby jones! You are alone here. We hope that changes soon.
(4:20PM) kate has logged in!
(4:20PM) kate: hey bobby!
hey kate! How are you?
(4:20PM) kate: I'm doing great, thanks for asking!
(4:21PM) kate has logged out
```
