package main

import "testing"

func Test_writeConfig(t *testing.T) {
	writeConfig("测试", PROTO_GRPC, "10.10.18.19", "6000")
}
