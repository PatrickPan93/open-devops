package common

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Printf("%+v", errors.Wrap(err, "common.GetLocalIP: Error while getting local IP"))
		return ""
	}
	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	conn.Close()
	return localIP
}

func ShellCommand(shellStr string) (string, error) {
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctxt, "sh", "-c", shellStr)
	var buf bytes.Buffer
	// 标准错误重定向到标准输出
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return buf.String(), errors.Wrap(err, fmt.Sprintf("command.ShellCommand: Error while executing shell command: %s", shellStr))
	}
	if err := cmd.Wait(); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("command.ShellCommand: TimeOut while executing shell command: %s", shellStr))
	}
	return buf.String(), nil
}
