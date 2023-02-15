package redis

import (
	"context"
	"errors"
	"fmt"
	"liujun/Time_Micro_GateWay/common"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	USING       = 1
	FREE        = 2
	INITNUM     = 20
	MAXNUM      = 60
	PINGSTEP    = 10   //两次ping之间的间隔
	RETRY_TIMES = 3    //重试次数
	ALIVE_TIME  = 7200 //连接存活时间上限,2个小时,看配置的timeout来设
)

type RedisConn struct {
	Red      *redis.Client //单机和主从哨兵对象
	status   int
	pingTime int64 //最后一次ping的时间
	time     int64 //初始化时间
}

type RedisPool struct {
	sync.RWMutex
	maxConnNum  int             //最大连接数
	initConnNum int             //初始连接数
	idleConns   chan *RedisConn //空闲连接（未初始化）
	cacheConns  chan *RedisConn //缓存连接

	pushConnCount int64 //已放回的连接数量
	popConnCount  int64 //已取出的连接数量

	use_pool     bool //是否启用连接池
	close_status bool //是否关闭状态
}

var (
	RED *RedisPool
)

func init() {
	RED = OpenPoll()
}

func OpenPoll() *RedisPool {
	pool := newPoll()
	for i := 0; i < pool.maxConnNum; i++ {
		if i < pool.initConnNum { //连接数在初始连接数内，将redisconn加入缓存连接中
			conn, err := pool.addRedisConn()
			fmt.Println(conn)
			if err != nil {
				pool.addIdleConn()
				continue
			}
			pool.cacheConns <- conn
		} else { //超过初始连接数，将空的redisconn加入空闲连接中
			conn := new(RedisConn)
			pool.idleConns <- conn
		}
	}
	return pool
}

func newPoll() *RedisPool {
	return &RedisPool{
		maxConnNum:   MAXNUM,
		initConnNum:  INITNUM,
		idleConns:    make(chan *RedisConn, MAXNUM),
		cacheConns:   make(chan *RedisConn, MAXNUM),
		use_pool:     true,
		close_status: false,
	}
}

func (rp *RedisPool) addRedisConn() (*RedisConn, error) {
	conn := new(RedisConn)
	redis_host := common.Config.String("redis_host")
	redis_port, _ := common.Config.Int("redis_port")
	redis_addr := fmt.Sprintf("%s:%d", redis_host, redis_port)
	redis_database, _ := common.Config.Int("redis_database")
	redis_poolsize, _ := common.Config.Int("redis_poolsize")
	redis_minidlecon, _ := common.Config.Int("redis_minidlecon")
	client := redis.NewClient(&redis.Options{
		Addr:         redis_addr,
		Password:     common.Config.String("redis_password"),
		DB:           redis_database,
		PoolSize:     redis_poolsize,
		MinIdleConns: redis_minidlecon,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(pong)
	conn.Red = client
	conn.status = FREE
	conn.time = time.Now().Unix()
	conn.pingTime = time.Now().Unix()
	return conn, err
}

func (rp *RedisPool) addIdleConn() {
	conn := new(RedisConn)
	select {
	case rp.idleConns <- conn: //往idleconns通道写数据
		return
	default: //idleconns通道关闭或者写满了
		return
	}
}

func (rp *RedisPool) Close() error {
	rp.Lock()
	idleConns := rp.idleConns
	cacheConns := rp.cacheConns
	rp.close_status = true
	rp.idleConns = nil
	rp.cacheConns = nil

	if idleConns == nil || cacheConns == nil {
		return nil
	}
	close(cacheConns)
	for cacheConn := range cacheConns { //将连接池redisconn关闭
		cacheConn.Red.Close()
	}
	close(idleConns)
	return nil
}

func (rp *RedisPool) addCacheConns() error {
	conn, err := rp.addRedisConn()
	if err != nil {
		rp.addIdleConn()
		return err
	}
	select {
	case rp.cacheConns <- conn:
		return err
	default:
		return errors.New("redis缓存池满或已关闭")
	}
}

func (rp *RedisPool) pushRedisConn(conn *RedisConn) error {
	if rp.close_status {
		return errors.New("redis缓存池满或已关闭")
	}
	if conn == nil {
		return nil
	}
	select {
	case rp.cacheConns <- conn:
		return nil
	default:
		return errors.New("redis缓存池满或已关闭")
	}
}

func (rp *RedisPool) getRedisConn() (cacheConns, idleConns chan *RedisConn) {
	rp.RLock()
	cacheConns = rp.cacheConns
	idleConns = rp.idleConns
	rp.RUnlock()
	return
}

func (rp *RedisPool) getCacheConns(cacheConns chan *RedisConn) (redisConn *RedisConn) {
	for len(cacheConns) > 0 {
		redisConn = <-cacheConns
		now := time.Now().Unix()
		if redisConn != nil {
			if ALIVE_TIME < now-redisConn.time { //redis连接超过存活时间
				redisConn.Red.Close()
				err := rp.addCacheConns()
				if err != nil {
					rp.addIdleConn()
				}
				redisConn = nil
			} else if now-redisConn.pingTime > PINGSTEP {
				_, err := redisConn.Red.Ping(context.Background()).Result()
				if err != nil {
					redisConn.Red.Close()
					rp.addIdleConn()
					redisConn = nil
				} else {
					redisConn.pingTime = time.Now().Unix()
				}
			}
		}
		if redisConn != nil {
			atomic.AddInt64(&rp.popConnCount, 1)
			break
		}
	}
	return
}

func (rp *RedisPool) getIdleConns(cacheConns chan *RedisConn, idleConns chan *RedisConn) (redisConn *RedisConn, err error) {
	select {
	case redisConn = <-idleConns:
		redisConn, err = rp.addRedisConn()
		if err != nil {
			rp.addIdleConn()
			return nil, err
		}
		return redisConn, nil
	case redisConn = <-cacheConns:
		now := time.Now().Unix()
		if redisConn != nil {
			if now-redisConn.time > ALIVE_TIME {
				redisConn.Red.Close()
				err = rp.addCacheConns()
				if err != nil {
					rp.addIdleConn()
				}
				redisConn = nil
			} else if now-redisConn.pingTime > PINGSTEP {
				_, err = redisConn.Red.Ping(context.Background()).Result()
				if err != nil {
					redisConn.Red.Close()
					rp.addIdleConn()
					redisConn = nil
				} else {
					redisConn.pingTime = time.Now().Unix()
				}
			}
		}
	}
	if redisConn != nil {
		atomic.AddInt64(&rp.popConnCount, 1)
	}
	return redisConn, nil
}

// 连接池获取连接
func (rp *RedisPool) popRedisConn() (redisConn *RedisConn, err error) {
	if rp.close_status {
		return nil, errors.New("redisconn已关闭")
	}
	if !rp.use_pool {
		redisConn, err = rp.addRedisConn()
		return redisConn, err
	}
	cacheConns, idleConns := rp.getRedisConn()
	if cacheConns == nil || idleConns == nil {
		return nil, errors.New("连接池为空")
	}

	for i := 0; i < RETRY_TIMES; i++ {
		redisConn = rp.getCacheConns(cacheConns)
		if redisConn == nil {
			redisConn, err = rp.getIdleConns(cacheConns, idleConns)
			if err == nil {
				return
			}
		} else {
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
	err = errors.New("连接超时")
	return
}

func (rp *RedisPool) Hget(key, field string) (string, error) {
	redis_conn, err := rp.popRedisConn()
	defer rp.pushRedisConn(redis_conn)
	if err != nil {
		return "", err
	}
	res, err := redis_conn.Red.HGet(context.Background(), key, field).Result()
	return res, err
}

func (rp *RedisPool) Sismember(key string, member interface{}) bool {
	redis_conn, err := rp.popRedisConn()
	defer rp.pushRedisConn(redis_conn)
	if err != nil {
		return false
	}
	res, err := redis_conn.Red.SIsMember(context.Background(), key, member).Result()
	if err != nil {
		return false
	}
	return res
}
