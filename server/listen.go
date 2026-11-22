package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// unixBindPrefix makes the bind address a UNIX socket path
// (TERMOTE_BIND=unix:/run/user/1000/termote.sock): the server is then reached
// only by who may open the socket (0600, this OS user), for example through
// an SSH forward to it, and holds no port.
const unixBindPrefix = "unix:"

// listenServer opens the server's listener: TCP on bind and port, or the UNIX
// socket bind names. A socket file left by a server that died is replaced;
// anything else at that path is refused.
func listenServer(bind, port string) (net.Listener, error) {
	path, ok := strings.CutPrefix(bind, unixBindPrefix)
	if !ok {
		return net.Listen("tcp", net.JoinHostPort(bind, port))
	}
	if path == "" || !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("TERMOTE_BIND %q: a UNIX socket needs an absolute path", bind)
	}
	if fi, err := os.Lstat(path); err == nil {
		if fi.Mode()&os.ModeSocket == 0 {
			return nil, fmt.Errorf("TERMOTE_BIND %q: %s is not a socket", bind, path)
		}
		// A live server still answers there: never take its socket.
		if c, err := net.Dial("unix", path); err == nil {
			c.Close()
			return nil, fmt.Errorf("TERMOTE_BIND %q: a server already listens there", bind)
		}
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}
	// Only this OS user may connect. Until the chmod the socket has the
	// process's umask, so its directory should be private too.
	ln, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(path, 0o600); err != nil {
		ln.Close()
		return nil, err
	}
	return ln, nil
}

// isUnixBind reports whether bind names a UNIX socket.
func isUnixBind(bind string) bool { return strings.HasPrefix(bind, unixBindPrefix) }
