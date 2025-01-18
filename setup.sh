#!/bin/bash

git submodule update --init --recursive
git submodule update --remote --recursive
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"

cd third_party/cloudapi

protoc -I=. -I=third_party/googleapis --go_out=. --go-grpc_out=. yandex/cloud/compute/v1/*.proto
protoc -I=. -I=third_party/googleapis --go_out=. --go-grpc_out=. yandex/cloud/operation/operation.proto

cd ../..

mkdir -p /etc/yacloud_revival

touch /etc/yacloud_revival/general.log
chmod 666 /etc/yacloud_revival/general.log
cat /dev/null > /etc/yacloud_revival/general.log
cp configs/config.yaml /etc/yacloud_revival/config.yaml
cp systemd/env /etc/yacloud_revival/env

go mod tidy
go build -o systemd/yacloud_revival ./cmd/yacloud_revival/main.go

rm -rf /usr/bin/yacloud_revival || true
cp systemd/yacloud_revival /usr/bin/yacloud_revival
chmod +x /usr/bin/yacloud_revival

cp systemd/yacloud_revival.service /etc/systemd/system/yacloud_revival.service
cp systemd/log_eraser.service /etc/systemd/system/log_eraser.service
cp systemd/log_eraser.timer /etc/systemd/system/log_eraser.timer
cp systemd/erase_log.sh /usr/bin/erase_log.sh
chmod +x /usr/bin/erase_log.sh

systemctl daemon-reload

systemctl enable yacloud_revival.service
systemctl start yacloud_revival.service

systemctl enable log_eraser.timer
systemctl start log_eraser.timer

