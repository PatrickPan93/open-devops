#!/bin/bash
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build /Users/patrick/go/src/open-devops/src/modules/agent/agent.go
scp agent root@192.168.40.50:/root/agent
scp agent root@192.168.40.61:/root/agent
scp agent root@192.168.40.62:/root/agent
scp open-devops-agent.yaml root@192.168.40.50:/root/
scp open-devops-agent.yaml root@192.168.40.61:/root/
scp open-devops-agent.yaml root@192.168.40.62:/root/
