package logster

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const DefaultLogtailPath = "/usr/sbin/logtail2"

type LogtailTailer struct {
	Binary    string
	Logfile   string
	Statefile string
}

func (tailer *LogtailTailer) cmd() []string {
	return []string{"-f", tailer.Logfile, "-o", tailer.Statefile}
}

func (tailer *LogtailTailer) CreateStateFile() error {
	cmd := exec.Command(tailer.Binary, tailer.cmd()...)

	var stderrBuf bytes.Buffer
	var errStderr error
	stderrIn, _ := cmd.StderrPipe()
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	if err := cmd.Wait(); err != nil {
		if errStderr != nil {
			return err
		}
		return fmt.Errorf("%s: %s", err, string(stderrBuf.Bytes()))
	}
	return nil
}

func (tailer *LogtailTailer) ReadLines(c chan string) error {
	cmd := exec.Command(tailer.Binary, tailer.cmd()...)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	defer close(c)
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
