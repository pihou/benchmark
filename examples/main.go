package main

import "github.com/myzhan/boomer"
import "github.com/valyala/fasthttp"

//import "time"

func foo() {
	//"http://k8s-node1.shoupihou.site:1008/app/benchmark/"
	//"https://api.longban.site/app/login/?version=2018004&phone=%2b8613249629530&password=123456&type=phone&platform=android"

	start := boomer.Now()
	url := "http://k8s-node1.shoupihou.site:1008/app/benchmark/"

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)

	elapsed := boomer.Now() - start

	// Report your test result as a success, if you write it in python, it will looks like this
	// events.request_success.fire(request_type="http", name="foo", response_time=100, response_length=10)
	if err != nil {
		boomer.Events.Publish("request_failure", "http", "foo", elapsed, err.Error())
	} else {
		boomer.Events.Publish("request_success", "http", "foo", elapsed, int64(0))
	}
}

func main() {

	task1 := &boomer.Task{
		Name:   "foo",
		Weight: 100,
		Fn:     foo,
	}
	boomer.Run(task1)
}
