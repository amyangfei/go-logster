package logster

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStateFile(t *testing.T) {
	var lf, sf *os.File
	var err error

	lf, err = ioutil.TempFile("", "fake_tailer.log")
	assert.Nil(t, err)
	defer os.Remove(lf.Name())

	sf, err = ioutil.TempFile("", "fake_tailer.state")
	assert.Nil(t, err)
	defer os.Remove(sf.Name())

	tailer := &LogtailTailer{
		Binary:    DefaultLogtailPath,
		Statefile: sf.Name(),
	}
	err = tailer.CreateStateFile()
	assert.NotNil(t, err)

	tailer.Logfile = lf.Name()
	err = tailer.CreateStateFile()
	assert.Nil(t, err)
}

func TestReadlines(t *testing.T) {

	lf, err := ioutil.TempFile("", "fake_tailer.log")
	assert.Nil(t, err)
	defer os.Remove(lf.Name())
	sf, err := ioutil.TempFile("", "fake_tailer.state")
	assert.Nil(t, err)
	defer os.Remove(sf.Name())

	testCases := []struct {
		line  string
		count int
	}{
		{"test one\nline two\nthis is line three\n", 3},
		{"test two\n", 1},
		{"test two\nline file\nthis is line six with emoji: 😊\n", 3},
	}

	f, err := os.OpenFile(lf.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	assert.Nil(t, err)
	defer f.Close()

	tailer := &LogtailTailer{
		Binary:    DefaultLogtailPath,
		Logfile:   lf.Name(),
		Statefile: sf.Name(),
	}

	var wg sync.WaitGroup
	for _, tc := range testCases {
		c := make(chan string)
		wg.Add(1)
		go func(expected string, count int) {
			defer wg.Done()
			lines := []string{}
			for line := range c {
				lines = append(lines, line)
			}
			assert.Equal(t, count, len(lines))
			assert.Equal(t, expected, strings.Join(lines, "\n")+"\n")
		}(tc.line, tc.count)

		_, err := f.WriteString(tc.line)
		assert.Nil(t, err)
		err = tailer.ReadLines(c)
		assert.Nil(t, err)
		wg.Wait()
	}
}
