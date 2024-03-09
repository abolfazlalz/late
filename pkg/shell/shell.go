package shell

import (
	"bufio"
	"os/exec"
	"sync"
	"time"
)

type Buffer struct {
	text string
	time time.Time
}

func NewBuffer(text string) *Buffer {
	return &Buffer{text: text, time: time.Now()}
}

type CommandExecute struct {
	bufferCh chan *Buffer
	isDone   bool
	command  string
	args     []string
}

func NewCommandExecute(command string, args ...string) *CommandExecute {
	return &CommandExecute{
		bufferCh: make(chan *Buffer),
		isDone:   false,
		command:  command,
		args:     args,
	}
}

type Shell struct {
	commands map[int]*CommandExecute
	lastId   int
	mu       sync.Mutex
}

func NewShell() *Shell {
	return &Shell{
		commands: make(map[int]*CommandExecute),
		lastId:   0,
	}
}

func (s *Shell) Execute(command string, args ...string) (int, *CommandExecute) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastId++
	execute := NewCommandExecute(command, args...)
	s.commands[s.lastId] = execute
	return s.lastId, execute
}

func (cmdExec *CommandExecute) Run(quiteCh chan int) error {
	cmd := exec.Command(cmdExec.command, cmdExec.args...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	rd := bufio.NewReader(stdout)
	if err := cmd.Start(); err != nil {
		return err
	}

	for {
		select {
		case <-quiteCh:
			return nil
		default:
			str, err := rd.ReadString('\n')
			if err != nil {
				return err
			}
			cmdExec.bufferCh <- NewBuffer(str)
		}
	}
}
