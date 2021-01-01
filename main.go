package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s container_name\n", os.Args[0])
		os.Exit(1)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	r, err := cli.ContainerAttach(context.Background(), os.Args[1], types.ContainerAttachOptions{Stdin: true, Stream: true})
	if err != nil {
		panic(err)
	}
	defer r.Close()

	_, err = io.Copy(r.Conn, os.Stdin)
	if err != nil {
		panic(err)
	}
}
