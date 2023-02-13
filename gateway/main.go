package main

import (
	"liujun/Time_Micro_GateWay/router"
)

func main() {

	r := router.Router()
	r.Run(":8001")
}
