package main

import (
"fmt"
"net"
"net/http"
"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock [wwwroot]")
		return
	}

	fmt.Println("Unix HTTP server")

	root := "."
	if len(os.Args) > 2 {
		root = os.Args[2]
	}

	server := http.Server{
		Handler: http.FileServer(http.Dir(root)),
	}

	// unixListener, err := net.Listen("unix", os.Args[1])
	unixListener, err := net.Listen("tcp",  "127.0.0.1:")
	if err != nil {
		panic(err)
	}
	fmt.Println(unixListener.Addr().String())
	server.Serve(unixListener)
}

