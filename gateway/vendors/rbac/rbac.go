package rbac

import (
	"liujun/Time_Micro_GateWay/utils"
	myredis "liujun/Time_Micro_GateWay/vendors/redis"

	"github.com/go-redis/redis/v8"
)

func RbacFilter(role_id int, method_name string) bool {
	method_id, err := myredis.RED.Hget("rbac_method", "rbac_method:"+method_name)
	if err != redis.Nil {
		return false
	}
	return myredis.RED.Sismember("rbac_role:"+utils.GetString(role_id), method_id)
}
