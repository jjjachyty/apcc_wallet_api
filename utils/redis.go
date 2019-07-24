package utils

import (
	"errors"
	"time"

	"pccqcpa.com.cn/app/fpm/api/utils"

	"github.com/go-redis/redis"
)

type cacheRedis struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
}

var cache cacheRedis
var flag int

func init() {
	var cacheCfg = utils.GetCacheCfg()
	if 0 < len(cacheCfg.Cluster.Addrs) { //集群
		flag = 2
	} else if "" != cacheCfg.Single.Server { //单机
		flag = 1
	}

	switch flag {
	case 1: //单机模式
		cache.client = redis.NewClient(&redis.Options{
			Addr:     cacheCfg.Single.Server,
			Password: cacheCfg.Single.PassWord, // no password set
			DB:       cacheCfg.Single.DB,       // use default DB
		})
		_, err := cache.client.Ping().Result()
		if nil != err {
			// utils.SysLog.Errorf("单机版缓存连接失败,请检测toml配置文件%v", err)
			utils.SysLog.Panic("单机版缓存连接失败,请检查config/app.toml配置文件\n" + err.Error())
		} else {
			utils.SysLog.Infoln("启用单机版缓存")
			cache.client.FlushAll()
		}

	case 2: //集群模式
		cache.clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: cacheCfg.Cluster.Addrs,
		})
		_, err := cache.clusterClient.Ping().Result()
		if nil != err {
			utils.SysLog.Error("集群版缓存连接失败,请检测toml配置文件", err)
			panic("集群版缓存连接失败,请检查config/app.toml配置文件\n" + err.Error())
		} else {
			utils.SysLog.Infoln("启用集群版缓存")
			cache.clusterClient.FlushAll()
		}

	default:
		utils.SysLog.Infoln("Redis缓存未开启,无法使用")
	}

}

//Set 设置缓存 key 键 value []byte值 如果是实体或者指针需要实现MarshalBinary方法 建议使用json.Marshal做处理 expiration过期时间 0为不过期
func Set(key string, value interface{}, expiration time.Duration) error {
	var statusCmd *redis.StatusCmd
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			statusCmd = cache.clusterClient.Set(key, value, expiration)
		} else {
			statusCmd = cache.client.Set(key, value, expiration)
		}

	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Get Value失败", err)
		return err
	}
	return statusCmd.Err()
}

//HMSet 设置Hash Table缓存 key 键 maps 字段-值
func HMSet(key string, maps map[string]interface{}) error {
	var statusCmd *redis.StatusCmd
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			statusCmd = cache.clusterClient.HMSet(key, maps)
		} else {
			statusCmd = cache.client.HMSet(key, maps)
		}

	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Hash SET 失败", err)
		return err
	}
	return statusCmd.Err()
}

//SAdd sadd(key, member) 设置Set集合类型 key 键 sets slice
func SAdd(key string, sets ...interface{}) error {
	var intCmd *redis.IntCmd
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			intCmd = cache.clusterClient.SAdd(key, sets...)
		} else {
			intCmd = cache.client.SAdd(key, sets...)
		}

	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Set集合类型 SAdd 失败", err)
		return err
	}
	return intCmd.Err()
}

//SMembers 根据key获取Set 集合值
func SMembers(key string) ([]string, error) {
	var ssCmd *redis.StringSliceCmd
	if 0 < flag {
		if 2 == flag { //集群
			ssCmd = cache.clusterClient.SMembers(key)
		} else {

			ssCmd = cache.client.SMembers(key)

		}

	} else {
		err := errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("SMembers 获取Set集合 失败", err)
		return nil, err
	}
	return ssCmd.Result()
}

//HMGet 根据key,字段 获取字段对应的值
func HMGet(key string, fields ...string) ([]interface{}, error) {
	var ssCmd *redis.SliceCmd
	if 0 < flag {
		if 2 == flag { //集群
			ssCmd = cache.clusterClient.HMGet(key, fields...)
		} else {

			ssCmd = cache.client.HMGet(key, fields...)

		}

	} else {
		err := errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("HMGet 失败", err)
		return nil, err
	}
	return ssCmd.Result()
}

//HGetAll 根据key 获取所有的Hash Table缓存值 key
func HGetAll(key string) (map[string]string, error) {
	var ssmCmd *redis.StringStringMapCmd
	if 0 < flag {
		if 2 == flag { //集群
			//取key下面所有的值
			ssmCmd = cache.clusterClient.HGetAll(key)
		} else {
			//取key下面所有的值
			ssmCmd = cache.client.HGetAll(key)
		}

	} else {
		err := errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("HGetAll 失败", err)
		return nil, err
	}
	return ssmCmd.Result()
}

//HDel 根据key fields 删除HashTable 中的值
func HDel(key string, fields ...string) error {
	var intCmd *redis.IntCmd
	if 0 < flag {
		if 2 == flag { //集群

			intCmd = cache.clusterClient.HDel(key, fields...)
		} else {

			intCmd = cache.client.HDel(key, fields...)
		}

	} else {
		err := errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("HDel 失败", err)
		return err
	}
	return intCmd.Err()
}

//Get 设置缓存 key 键 value
func Get(key string) (string, error) {
	var stringCmd *redis.StringCmd
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			stringCmd = cache.clusterClient.Get(key)
		} else {
			stringCmd = cache.client.Get(key)

		}

	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Get Value失败", err)
		return "", err

	}
	return stringCmd.Result()
}

//Keys 模糊查询 *：通配任意多个字符 ?：通配单个字符[]：通配括号内的某一个字符
//PS：据说会导致cpu很高 不建议使用模糊搜索
func Keys(pattern string) ([]string, error) {
	var ssCmd *redis.StringSliceCmd
	if 0 < flag {
		if 2 == flag { //集群
			ssCmd = cache.clusterClient.Keys(pattern)
		} else {
			ssCmd = cache.client.Keys(pattern)
		}

	} else {
		err := errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Del Key失败", err)
		return nil, err
	}
	return ssCmd.Result()
}

//Del 删除缓存 key 键 value
func Del(key string) error {
	var intCmd *redis.IntCmd
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			intCmd = cache.clusterClient.Del(key)
		} else {
			intCmd = cache.client.Del(key)
		}

	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Del Key失败", err)
		return err

	}
	return intCmd.Err()
}

//Subscribe 订阅redis的广播
func Subscribe(channel string) (*redis.PubSub, error) {
	var pubsub *redis.PubSub
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			pubsub = cache.clusterClient.Subscribe(channel)
		} else {
			pubsub = cache.client.Subscribe(channel)
		}

		// Wait for subscription to be created before publishing message.
		_, err = pubsub.ReceiveTimeout(time.Second)
		if err != nil {
			utils.SysLog.Errorf("订阅通道[%s]失败%v", err)
		}
	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Subscribe消息失败", err)
	}
	return pubsub, err
}

//Publish 发送redis的广播
func Publish(channel string, message interface{}) error {
	var err error
	if 0 < flag {
		if 2 == flag { //集群
			err = cache.clusterClient.Publish(channel, message).Err()
		} else {
			err = cache.client.Publish(channel, message).Err()
		}
		if err != nil {
			utils.SysLog.Errorf("向通道[%s]发布失败%v", err)

		}
	} else {
		err = errors.New("redis 缓存配置未开启,请先开启Redis 缓存配置")
		utils.SysLog.Error("Publish消息失败", err)

	}
	return err
}
