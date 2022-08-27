package redis

import (
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(init bool) (*redis.Client, error) {

	var (
		Network  string
		Username string
	)

	redisOpts := &redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + `:` + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	Network = os.Getenv("REDIS_NETWORK")
	if Network != "" {
		redisOpts.Network = Network
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_DB")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.DB = value
	}

	Username = os.Getenv("REDIS_USERNAME")
	if Username != "" {
		redisOpts.Username = Username
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_MAX_RETRIES")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.MaxRetries = value
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_MIN_RETRY_BACKOFF_MS")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.MinRetryBackoff = time.Duration(value) * time.Millisecond
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_MAX_RETRY_BACKOFF_MS")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.MaxRetryBackoff = time.Duration(value) * time.Millisecond
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_DIAL_TIMEOUT_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.DialTimeout = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_READ_TIMEOUT_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.ReadTimeout = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_WRITE_TIMEOUT_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.WriteTimeout = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptBoolEnv(os.Getenv("REDIS_POOL_FIFO")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.PoolFIFO = value
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_POOL_SIZE")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.PoolSize = value
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_MIN_IDLE_CONNS")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.MinIdleConns = value
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_MAX_CONN_AGE_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.MaxConnAge = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_POOL_TIMEOUT_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.PoolTimeout = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_IDLE_TIMEOUT_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.IdleTimeout = time.Duration(value) * time.Second
	}

	if value, exist, err := FetchOptIntEnv(os.Getenv("REDIS_IDLE_CHECK_FREQUENCY_SEC")); err != nil {
		return nil, err
	} else if exist == true {
		redisOpts.IdleCheckFrequency = time.Duration(value) * time.Second
	}

	return redis.NewClient(redisOpts), nil
}

func FetchOptIntEnv(envStr string) (value int, exist bool, err error) {
	if envStr != "" {
		exist = true
		if intVal, err := strconv.Atoi(envStr); err != nil {
			return value, exist, err
		} else {
			value = intVal
		}
	} else {
		exist = false
	}

	return value, exist, nil
}

func FetchOptBoolEnv(envStr string) (value bool, exist bool, err error) {
	if envStr != "" {
		exist = true
		if boolVal, err := strconv.ParseBool(envStr); err != nil {
			return value, exist, err
		} else {
			value = boolVal
		}
	} else {
		exist = false
	}

	return value, exist, nil
}
