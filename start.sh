#!/bin/bash
mkdir -p /home/data/etcd
mkdir -p /home/log/etcd

nohup etcd --data-dir /home/data/etcd 2>&1 1>/home/log/etcd/etcd.log etcd --listen-client-urls http://localhost:18003 --listen-peer-urls http://localhost:18002 --advertise-client-urls http://localhost:18003 &
sleep 5

mkdir -p /home/data/console
mkdir -p /home/log/console

./console -log /home/log/console -etcds http://localhost:18003 $*