package proc

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"github.com/lixiangyun/https-gateway/util"
	"os/exec"
	"sync"
	"time"
)

type Cmd struct {
	cmd    string
	args   []string
	stdout *bytes.Buffer
	stdin  *bytes.Buffer
	stderr *bytes.Buffer
	code   int
}

func (c *Cmd)Stdout() string {
	return c.stdout.String()
}

func (c *Cmd)Stderr() string {
	return c.stderr.String()
}

func (c *Cmd)Stdin() string {
	return c.stdin.String()
}

func (c *Cmd)Run() int  {
	return c.RunWithStdin("")
}

func (c *Cmd)RunWithStdin(in string) int {
	cmd := exec.Command(c.cmd, c.args...)
	if in != "" {
		c.stdin = bytes.NewBufferString(in)
		cmd.Stdin = c.stdin
	}
	cmd.Stderr = c.stderr
	cmd.Stdout = c.stdout
	err := cmd.Run()
	if err != nil {
		logs.Warn("exec fail", err.Error())
	}
	return cmd.ProcessState.ExitCode()
}

func NewCmd(cmd string, args...string) *Cmd {
	c := new(Cmd)
	c.cmd = cmd
	c.stderr = bytes.NewBuffer(make([]byte, 1024))
	c.stdout = bytes.NewBuffer(make([]byte, 1024))
	c.args = append(c.args, args...)
	return c
}

type Proc struct {
	sync.Mutex

	exc    *exec.Cmd
	cmd    string
	pid    int
	args   []string
	listen func(proc *Proc, err error)
	stdin  []byte
	stdout *util.LogChannel
	stderr *util.LogChannel
	alive  int

	kill   bool
	exit   chan struct{}
	stop   chan struct{}
	code   int
}

func NewProc(cmd string, args...string) *Proc {
	proc := new(Proc)
	proc.cmd = cmd
	proc.stop = make(chan struct{},1)
	proc.exit = make(chan struct{},1)
	proc.args = append(proc.args, args...)
	proc.stderr = util.NewLogChannel()
	proc.stdout = util.NewLogChannel()
	return proc
}

func (proc *Proc)Stdout() <- chan []byte {
	return proc.stdout.Listen()
}

func (proc *Proc)Stderr() <- chan []byte {
	return proc.stderr.Listen()
}

func (proc *Proc)GetPid() int {
	return proc.pid
}

func (proc *Proc)ExitCode() int {
	return proc.code
}

func (proc *Proc)retryStart() {
	go func() {
		proc.Lock()
		exc := exec.Command(proc.cmd, proc.args...)
		exc.Stderr = proc.stderr
		exc.Stdout = proc.stdout
		proc.exc = exc
		proc.alive = 0

		err := proc.exc.Start()
		if err != nil {
			logs.Error("proc start fail!", err.Error())
		}else {
			logs.Info("proc start success!", proc.cmd)
			proc.pid = proc.exc.Process.Pid
		}
		proc.Unlock()

		err = exc.Wait()
		if err != nil {
			logs.Error("proc wait fail!", err.Error())
		}else {
			logs.Info("proc exit!", proc.cmd)
		}

		proc.Lock()
		proc.exc = nil
		proc.alive = 0
		proc.code = exc.ProcessState.ExitCode()
		proc.exit <- struct{}{}
		proc.Unlock()
	}()
}


func (proc *Proc)Start() {
	proc.retryStart()

	go func() {
		tick := time.NewTicker(9*time.Second)
		for  {
			select {
			case <- proc.stop: {
				proc.shutdown()
			}
			case <- proc.exit: {
				if proc.kill {
					return
				}else {
					proc.retryStart()
				}
			}
			case <- tick.C: {
				if proc.kill && proc.exc != nil {
					proc.shutdown()
				}else {
					proc.alive++
				}
			}
			}
			time.Sleep(3*time.Second)
		}
		tick.Stop()
	}()
}

func (proc *Proc)shutdown() {
	proc.Lock()
	defer proc.Unlock()

	if proc.exc != nil {
		proc.exc.Process.Kill()
		logs.Info("proc shutdown", proc.pid)
	}
	proc.kill = true
}

func (proc *Proc)IsAlive() bool {
	if proc.alive > 2 {
		return true
	}
	return false
}

func (proc *Proc)Stop() {
	proc.stop <- struct{}{}
}

