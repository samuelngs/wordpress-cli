package app

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

var files = []*file{
	&file{
		filename:  "plugins",
		directory: true,
	},
	&file{
		filename:  "theme",
		directory: true,
	},
	&file{
		filename: "theme/theme.resource",
		content:  make([]byte, 0),
	},
}

// app - the app controller
type app struct {
	name       string
	port       string
	dir        string
	conf       []byte
	proj       project.APIProject
	ctx        context.Context
	cancelFunc context.CancelFunc
	cancelled  bool
}

func (v *app) Up() error {

	for _, f := range files {
		if err := f.createIfNotExist(v.dir); err != nil {
			return err
		}
	}

	switch err := v.proj.Create(v.ctx, options.Create{
		ForceRecreate: true,
	}); {
	case err != nil && !strings.Contains(err.Error(), "driver has changed"):
		return fmt.Errorf("The project cannnot be created. Details: %s", err)
	}

	switch err := v.proj.Start(v.ctx); {
	case err != nil && !strings.Contains(err.Error(), "already exists"):
		return fmt.Errorf("The project cannnot be started. Details: %s", err)
	}

	return nil
}

func (v *app) Down() error {

	switch err := v.proj.Stop(v.ctx, 0); {
	case err != nil:
		return fmt.Errorf("The project cannnot be stopped. Details: %s", err)
	}

	return nil
}

func (v *app) Remove() error {

	switch err := v.proj.Down(v.ctx, options.Down{
		RemoveVolume:  true,
		RemoveImages:  "local",
		RemoveOrphans: true,
	}); {
	case err != nil:
		return fmt.Errorf("The project cannnot be removed. Details: %s", err)
	}

	return nil
}

func (v *app) Log(follow bool) error {

	signalCh := make(chan os.Signal, 1)
	doneCh := make(chan error)
	errCh := make(chan error)

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errCh <- v.proj.Log(v.ctx, follow)
	}()

	go func() {
		select {
		case <-signalCh:
			v.cancel()
			doneCh <- nil
		case err := <-errCh:
			doneCh <- err
		}
	}()

	return <-doneCh
}

func (v *app) Port() string {
	return v.port
}

func (v *app) cancel() {
	if v.cancelled {
		return
	}
	v.cancelFunc()
	v.cancelled = true
}

func hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))[:10]
}

func port(s string) string {
	hasher := fnv.New32a()
	hasher.Write([]byte(s))
	o := strconv.FormatUint(uint64(hasher.Sum32()), 10)
	if l := len(o); l >= 4 {
		o = o[:4]
	} else {
		o = o[:l]
	}
	n, _ := strconv.ParseUint(o, 0, 64)
	n += uint64(start)
	return strconv.FormatUint(n, 10)
}

func newApp(dir string) (*app, error) {

	// create project hash name
	name := hash(dir)

	// find unique port
	port := port(name)

	// compile docker compose template
	conf, err := compile(buildConfig(name, dir, port))
	if err != nil {
		return nil, err
	}

	// import docker compose project
	proj, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{conf},
			ProjectName:  name,
		},
	}, nil)
	if err != nil {
		return nil, err
	}

	// cancelable context instance
	ctx, cancelFunc := context.WithCancel(context.Background())

	v := &app{
		name:       name,
		port:       port,
		dir:        dir,
		conf:       conf,
		proj:       proj,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	return v, nil
}
