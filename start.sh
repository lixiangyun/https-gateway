#!/bin/bash

mkdir -p $HTTPS_GATEWAY_HOME/log
mkdir -p $HTTPS_GATEWAY_HOME/etcd
mkdir -p $HTTPS_GATEWAY_HOME/console

nohup etcd --data-dir $HTTPS_GATEWAY_HOME/etcd --listen-client-urls http://localhost:18003 --listen-peer-urls http://localhost:18002 --advertise-client-urls http://localhost:18003 >$HTTPS_GATEWAY_HOME/log/etcd.log &

sleep 5

./console -log $HTTPS_GATEWAY_HOME/log -home $HTTPS_GATEWAY_HOME/console -etcds http://localhost:18003 $*