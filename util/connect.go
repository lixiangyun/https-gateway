package util

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func ListenAndConnect(domain string, port int) error {
	var wait sync.WaitGroup

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	token := GetToken(128)
	var success bool
	var close bool

	wait.Add(1)
	go func() {
		defer lis.Close()
		defer wait.Done()
		for {
			if success || close {
				break
			}

			conn, err := lis.Accept()
			if err != nil {
				continue
			}

			var test [1024]byte
			cnt, err := conn.Read(test[:])
			if err == nil && strings.Compare(token, string(test[:cnt])) == 0 {
				success = true
			}
			conn.Close()
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", domain, port))
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			conn.Write([]byte(token))
		}
		close = true
	}()

	wait.Wait()

	if success {
		return nil
	}

	return fmt.Errorf("connect %s:%d network is blocked", domain, port)
}

func CheckPortFree(port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	list.Close()
	return nil
}