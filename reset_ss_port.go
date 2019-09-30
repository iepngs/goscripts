package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	file := `/etc/shadowsocks-go/config.json`
	f, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	jsonConfig := string(f)
	fmt.Printf("origin json content: \n%s\n", jsonConfig)

	type Config struct {
		Server      string `json:"server"`
		Server_port uint32 `json:"server_port"`
		Local_port  uint32 `json:"local_port"`
		Password    string `json:"password"`
		Method      string `json:"method"`
		Timeout     uint32 `json:"timeout"`
	}
	config := Config{}
	err = json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		panic(err)
	}
	config.Server_port += 1
	newConfig, _ := json.Marshal(config)
	buf := bytes.Buffer{}
	jsonConfig = string(newConfig)
	buf.WriteString(jsonConfig)
	err = ioutil.WriteFile(file, buf.Bytes(), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Printf("new json content: \n%s\n\n", jsonConfig)
	bytes, err := exec.Command("/etc/init.d/shadowsocks-go", "restart").Output()
	if err != nil {
		panic(err)
	}
	resp := string(bytes)
	fmt.Println(resp)
}
