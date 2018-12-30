package logster

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/juju/errors"
)

// DefaultLogtailPath is the default path to logtail2 binary
const DefaultLogtailPath = "/usr/sbin/logtail2"

// LogtailTailer holds running parameters of logtail2
type LogtailTailer struct {
	Binary    string
	Logfile   string
	Statefile string
}

func (tailer *LogtailTailer) cmd() []string {
	return []string{"-f", tailer.Logfile, "-o", tailer.Statefile}
}

// CreateStateFile creates the state file of LogtailTailer
func (tailer *LogtailTailer) CreateStateFile() error {
	cmd := exec.Command(tailer.Binary, tailer.cmd()...)

	var stderrBuf bytes.Buffer
	var errStderr error
	stderrIn, _ := cmd.StderrPipe()
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	if err := cmd.Wait(); err != nil {
		if errStderr != nil {
			return errors.Trace(errStderr)
		}
		return errors.Annotatef(err, "create state file error: %s", string(stderrBuf.Bytes()))
	}
	return nil
}

// ReadLines reads lines via logtail and sends data to the given chan line by line
func (tailer *LogtailTailer) ReadLines(c chan string) error {
	cmd := exec.Command(tailer.Binary, tailer.cmd()...)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	defer close(c)
	return errors.Trace(cmd.Wait())
}
