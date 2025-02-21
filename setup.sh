#!/bin/bash

git submodule update --init --recursive
git submodule update --remote --recursive
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"

cd third_party/cloudapi

protoc -I=. -I=third_party/googleapis --go_out=. --go-grpc_out=. yandex/cloud/compute/v1/*.proto
protoc -I=. -I=third_party/googleapis --go_out=. --go-grpc_out=. yandex/cloud/operation/operation.proto

cd ../..

sudo mkdir -p /etc/yacloud_revival

sudo touch /etc/yacloud_revival/general.log
sudo chmod 666 /etc/yacloud_revival/general.log
sudo cat /dev/null > /etc/yacloud_revival/general.log
sudo cp configs/config.yaml /etc/yacloud_revival/config.yaml
sudo cp systemd/env /etc/yacloud_revival/env

go mod tidy
go build -o systemd/yacloud_revival ./cmd/yacloud_revival/main.go

sudo rm -rf /usr/bin/yacloud_revival || true
sudo cp systemd/yacloud_revival /usr/bin/yacloud_revival
sudo chmod +x /usr/bin/yacloud_revival

sudo cp systemd/yacloud_revival.service /etc/systemd/system/yacloud_revival.service
sudo cp systemd/log_eraser.service /etc/systemd/system/log_eraser.service
sudo cp systemd/log_eraser.timer /etc/systemd/system/log_eraser.timer
sudo cp systemd/erase_log.sh /usr/bin/erase_log.sh
sudo chmod +x /usr/bin/erase_log.sh

sudo systemctl daemon-reload

sudo systemctl enable log_eraser.timer
sudo systemctl stop log_eraser.timer || true
sudo systemctl start log_eraser.timer


sudo systemctl enable yacloud_revival.service
sudo systemctl stop yacloud_revival.service || true
sudo systemctl start yacloud_revival.service