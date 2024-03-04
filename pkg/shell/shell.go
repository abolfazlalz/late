package shell

import (
	"bufio"
	"log"
	"os/exec"
)

type Buffer struct {
	text string
}

type Shell struct {
	bufferCh chan *Buffer
}

func NewShell(ch chan *Buffer) *Shell {
	return &Shell{
		bufferCh: ch,
	}
}

func (s Shell) Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		log.Fatal("Buffer Error:", err)
	}

	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
			return err
		}
		s.bufferCh <- &Buffer{text: str}
	}

	return nil
}
