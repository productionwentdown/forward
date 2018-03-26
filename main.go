package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/judwhite/go-svc/svc"
)

type program struct {
	logFile *os.File
	wg      sync.WaitGroup
}

func main() {
	prg := &program{}

	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}
		logPath := filepath.Join(dir, "forward.log")

		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		p.logFile = f
		log.SetOutput(f)
	}

	setup()

	return nil
}

func (p *program) Start() error {
	p.wg.Add(1)
	go func() {
		serve()
		p.wg.Done()
	}()

	log.Print("started")
	return nil
}

func (p *program) Stop() error {
	log.Print("stopping...")

	ln.Close()
	p.wg.Wait()

	log.Print("stopped")
	return nil
}
