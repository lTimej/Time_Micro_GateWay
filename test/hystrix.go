package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
)

func HytrixTest(n int) {
	body := "{\"username\":\"liujun\",\"password\":\"admin\",\"captcha\":\"1234\"}"
	//post 方法参数，第一个参数为请求url,第二个参数 是contentType, 第三个参数为请求体[]byte格式
	w := sync.WaitGroup{}
	for i := 1; i <= n; i++ {
		w.Add(1)
		go func() {
			responce, err := http.Post("http://localhost:8001/login", "application/json", bytes.NewBuffer([]byte(body)))
			if err != nil {
				fmt.Println("net http post method err,", err)
			}
			data := make([]byte, 1024)
			n, _ := responce.Body.Read(data)
			fmt.Println(string(data[:n]))
			defer responce.Body.Close()
			w.Done()
		}()
	}
	w.Wait()
}
