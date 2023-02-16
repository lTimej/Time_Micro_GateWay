package main

import (
	"context"
	"errors"
	"log"
	"strconv"
)

func RefreshRbacMethodsWithDb() error {
	methods := []Method{}
	DB.Find(&methods)
	if len(methods) <= 0 {
		return errors.New("methods为空")
	}
	m := make(map[string]interface{})
	for _, method := range methods {
		m[KEY_RBAC_METHOD_PREFIX+method.Name] = method.Id
	}
	_, err := RED.HMSet(context.Background(), KEY_RBAC_METHOD, m).Result()
	return err
}

func RefreshRbacRoleMethodsWithDb() error {
	role_methods := []RoleMethod{}
	DB.Find(&role_methods)
	if len(role_methods) <= 0 {
		return errors.New("role_methods为空")
	}
	m := make(map[int][]string)
	for _, role_method := range role_methods {
		if v, ok := m[role_method.RoleId]; ok {
			m[role_method.RoleId] = append(v, strconv.Itoa(role_method.MethodId))
		} else {
			m[role_method.RoleId] = []string{strconv.Itoa(role_method.MethodId)}
		}
	}
	for role_id, method := range m {
		_, err := RED.SAdd(context.Background(), KEY_RBAC_ROLE_PREFIX+strconv.Itoa(role_id), method).Result()
		if err != nil {
			return errors.New("rbac_role_ Err:" + err.Error())
		}
	}
	return nil
}

func RefreshRbacHandler() error {
	err := RefreshRbacMethodsWithDb()
	if err != nil {
		log.Println("method_err:", err)
		return err
	}
	err = RefreshRbacRoleMethodsWithDb()
	if err != nil {
		log.Println("role_method_err:", err)
		return err
	}
	return nil
}
