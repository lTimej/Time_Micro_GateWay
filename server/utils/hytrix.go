package utils

import (
	"github.com/astaxie/beego/config"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
)

// 初始化追踪器
func GetTracer(serviceName string, host string, config config.Configer) opentracing.Tracer {

	//set up a span reporter
	url := config.String("zipkin_url")
	zipkinReporter := zipkinhttp.NewReporter(url)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, host)
	log.Println("GetTracer:", serviceName, host, url)

	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(zipkinReporter, zipkin.WithLocalEndpoint(endpoint))

	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)

	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.InitGlobalTracer(tracer)

	return tracer

}
