#!/bin/bash

GOOS=linux GOARCH=amd64 go build cmd/handler.go
func azure functionapp publish webreg-bot-app