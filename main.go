package main

import (
	"context"
	"io"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prometheus/tsdb/fileutil"
	flag "github.com/spf13/pflag"
)

func main() {
	var lockFile string
	flag.StringVar(&lockFile, "lockfile", "", "Path to lockfile")

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	r, err := lock(lockFile, 5)
	if err != nil {
		panic(err)
	}
	if r != nil {
		defer r.Release()
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	err = attachStdin(context.Background(), cli, flag.Arg(0), os.Stdin)
	if err != nil {
		panic(err)
	}
}

func lock(f string, retryCount int) (r fileutil.Releaser, err error) {
	if f == "" {
		return
	}
	for i := 0; i < retryCount; i++ {
		r, _, err = fileutil.Flock(f)
		if err == nil {
			return
		}
		backoff := time.Duration(math.Pow(2, float64(i))) * (50 * time.Millisecond)
		time.Sleep(backoff)
	}
	return
}

func attachStdin(ctx context.Context, cli *client.Client, containerName string, r io.Reader) error {
	c, err := cli.ContainerAttach(context.Background(), containerName, types.ContainerAttachOptions{Stdin: true, Stream: true})
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = io.Copy(c.Conn, r)
	if err != nil {
		return err
	}
	return nil
}
