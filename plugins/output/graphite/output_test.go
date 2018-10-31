package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/amyangfei/go-logster/logster"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	name  string
	value float64
	ts    int64
}

func prepareData() []TestData {
	inputData := []TestData{
		{"name1", 1.1, time.Now().Unix()},
		{"name2", 2.1, time.Now().Unix()},
		{"name3", 3.1, time.Now().Unix()},
	}
	return inputData
}

var Logger zerolog.Logger

func init() {
	Logger = zerolog.New(os.Stdout)
}

func mockServer(t *testing.T, count int, finCh chan bool, host, protocol string, expected []string) {
	if protocol == "tcp" {
		l, err := net.Listen(protocol, host)
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			buf, err := ioutil.ReadAll(conn)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("gogogo: %s\n", string(buf[:]))
			close(finCh)
		}
	} else if protocol == "udp" {
		addr, err := net.ResolveUDPAddr("udp", host)
		if err != nil {
			t.Fatal(err)
		}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		buffer := make([]byte, 1024)
		result := make([]string, 0)
		for count > 0 {
			count--
			n, err := conn.Read(buffer)
			assert.Nil(t, err)
			result = append(result, string(buffer[:n]))
		}
		assert.ElementsMatch(t, expected, result)
		close(finCh)
	} else {
		t.Fatalf("protocol %s not support", protocol)
	}
}

func TestGraphiteOutput(t *testing.T) {
	output := &GraphiteOutput{}
	var err error
	err = output.Init(
		"P", "S", `{"host": "127.0.0.1:12345", "protocol": "udp"}`, false, Logger)
	assert.Nil(t, err)

	inputData := prepareData()
	metrics := []*logster.Metric{}
	expected := make([]string, 0)
	finCh := make(chan bool)
	for _, line := range inputData {
		metrics = append(metrics,
			&logster.Metric{Name: line.name, Value: line.value, Timestamp: line.ts})
		expected = append(expected, fmt.Sprintf("P.%s.S %v %d", line.name, line.value, line.ts))
	}
	go mockServer(t, len(inputData), finCh, "0.0.0.0:12345", "udp", expected)
	// FIXME: better way to wait for TCP/UDP socket readable
	time.Sleep(1 * time.Second)
	err = output.Submit(metrics)
	assert.Nil(t, err)
	<-finCh
}
