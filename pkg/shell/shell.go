package shell

import (
	"bufio"
	"os/exec"
	"sync"
	"time"
)

var shell *Shell

type Buffer struct {
	text string
	time time.Time
}

func NewBuffer(text string) *Buffer {
	return &Buffer{text: text, time: time.Now()}
}

type CommandExecute struct {
	mu          sync.Mutex
	Complete    bool   `json:"complete"`
	Command     string `json:"command"`
	args        []string
	subscribeCh []chan *Buffer
	cmd         *exec.Cmd
	StartDate   time.Time `json:"start_date"`
}

func NewCommandExecute(command string, args ...string) *CommandExecute {
	return &CommandExecute{
		Complete:    false,
		Command:     command,
		args:        args,
		subscribeCh: make([]chan *Buffer, 0),
	}
}

type Shell struct {
	commands map[int]*CommandExecute
	lastId   int
	mu       sync.Mutex
}

func NewShell() *Shell {
	if shell != nil {
		return shell
	}
	shell = &Shell{
		commands: make(map[int]*CommandExecute),
		lastId:   0,
	}
	return shell
}

func (s *Shell) List() map[int]*CommandExecute {
	return s.commands
}

func (s *Shell) Kill(id int) error {
	cmd := s.commands[id]
	return cmd.Kill()
}

func (s *Shell) Execute(command string, args ...string) (int, *CommandExecute) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastId++
	execute := NewCommandExecute(command, args...)
	s.commands[s.lastId] = execute
	return s.lastId, execute
}

func (e *CommandExecute) Subscribe(bufCh chan *Buffer) {
	e.subscribeCh = append(e.subscribeCh, bufCh)
}

func (e *CommandExecute) callSubscribe(b *Buffer) {
	for _, ch := range e.subscribeCh {
		ch <- b
	}
}

func (e *CommandExecute) Kill() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	// Kill current process Command.
	err := e.cmd.Process.Kill()
	if err != nil {
		return err
	}
	e.Complete = true
	return nil
}

func (e *CommandExecute) Run(quiteCh chan int) error {
	// Stop last Command execution if already a Command running
	if e.cmd != nil {
		if err := e.cmd.Process.Kill(); err != nil {
			return err
		}
	}
	e.cmd = exec.Command(e.Command, e.args...)
	stdout, err := e.cmd.StdoutPipe()
	e.StartDate = time.Now()

	if err != nil {
		return err
	}

	rd := bufio.NewReader(stdout)
	if err := e.cmd.Start(); err != nil {
		return err
	}

	defer func() {
		e.mu.Lock()
		_ = stdout.Close()
		_ = e.cmd.Process.Kill()
		e.Complete = true
		e.mu.Unlock()
	}()

LoopMessage:
	for {
		select {
		case <-quiteCh:
			break LoopMessage
		default:
			str, err := rd.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					break LoopMessage
				}
				return err
			}
			buffer := NewBuffer(str)
			e.callSubscribe(buffer)
		}
	}

	return nil
}
