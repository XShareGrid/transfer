package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

type Protocol int

const (
	PROTO_HTTP Protocol = iota
	PROTO_GRPC
)

var fileName = "恩捷跨厂区服务统计.xlsx"
var genPath = "../"

func main() {
	// 读取excel
	f, err := xlsx.OpenFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var s *xlsx.Sheet
	var ok bool
	if s, ok = f.Sheet["工作表1"]; !ok {
		log.Fatal("无法找到工作表1")
	}
	for i, r := range s.Rows {
		if i == 0 {
			continue
		}
		if r.Cells[0].String() == "" {
			break
		}
		if r.Cells[2].String() == "待定" {
			continue
		}
		confName := r.Cells[0].String() + "_" + r.Cells[1].String()
		p := strings.ToLower(r.Cells[4].String())
		proto := PROTO_HTTP
		if p == "grpc" {
			proto = PROTO_GRPC
		}
		src := fmt.Sprintf("%s:%s", r.Cells[2].String(), r.Cells[3].String())
		dstPort := r.Cells[9].String()
		err = writeConfig(confName, proto, src, dstPort)
		if err != nil {
			log.Fatalf("第%d行写文件失败，错误："+err.Error(), i+1)
		}
		log.Printf("第%d配置%s生成OK\n", i+1, confName)
	}

}

func writeConfig(confName string, proto Protocol, src string, dstPort string) error {
	tpl := `server {
    listen %s;                                  # 代理端口
    server_name %s;                      # 代理服务名
    location / {
        %s%s;     # 目标地址
		client_max_body_size 200m;
    }
}

`
	var cfg string
	switch proto {
	case PROTO_HTTP:
		cfg = fmt.Sprintf(tpl, dstPort, confName, "proxy_pass http://", src)
	case PROTO_GRPC:
		tpl = `server {
    listen %s http2;                                  # 代理端口
    server_name %s;                      # 代理服务名
    location / {
        %s%s;     # 目标地址
    }
}
`
		cfg = fmt.Sprintf(tpl, dstPort, confName, "grpc_pass grpc://", src)
	}

	return os.WriteFile(genPath+confName+".conf", []byte(cfg), 0644)

}
