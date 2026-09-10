package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/ent"
	"github.com/gcc798/microservice-kit/kratos/application/sys-rpc/internal/conf"
	"github.com/gcc798/microservice-kit/kratos/pkg/configx"
	"github.com/redis/go-redis/v9"
)

type Resources struct {
	Ent     *ent.Client
	SQLDB   *sql.DB
	Redis   *redis.Client
	Storage *StorageManager
	WeChat  *WeChatManager
	JWT     *JWTManager
	Auth    AuthPolicy
	stopCh  chan struct{}
}

type AuthPolicy struct {
	AllowConcurrent bool
	ShareToken      bool
}

var idCounter atomic.Int64

func init() {
	idCounter.Store(time.Now().UnixNano())
}

func NewResources(dataCfg *conf.Data, authCfg *conf.Auth, jwtCfg *conf.JWT, obsCfg *conf.Observability) (*Resources, error) {
	if dataCfg.GetDatabase().GetDriver() == "" || dataCfg.GetDatabase().GetDsn() == "" {
		return nil, errors.New("database config is required")
	}
	var (
		client *ent.Client
		sqlDB  *sql.DB
		err    error
	)
	switch dataCfg.GetDatabase().GetDriver() {
	case "postgres":
		sqlDB, err = openObservedPostgresDB(dataCfg.GetDatabase().GetDsn(), dbObservability{
			slowThreshold: time.Duration(obsCfg.GetDbSlowThresholdMs()) * time.Millisecond,
		})
		if err != nil {
			return nil, err
		}
		client = ent.NewClient(ent.Driver(entsql.OpenDB(dataCfg.GetDatabase().GetDriver(), sqlDB)))
	default:
		client, err = ent.Open(dataCfg.GetDatabase().GetDriver(), dataCfg.GetDatabase().GetDsn())
		if err != nil {
			return nil, err
		}
	}
	if dataCfg.GetRedis().GetAddr() == "" {
		if client != nil {
			_ = client.Close()
		}
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		return nil, errors.New("redis config is required")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     dataCfg.GetRedis().GetAddr(),
		Password: dataCfg.GetRedis().GetPassword(),
		DB:       int(dataCfg.GetRedis().GetDb()),
	})
	redisClient.AddHook(newRedisMetricsHook(redisObservability{
		slowThreshold: time.Duration(obsCfg.GetRedisSlowThresholdMs()) * time.Millisecond,
	}))
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		return nil, err
	}
	res := &Resources{
		Ent:     client,
		SQLDB:   sqlDB,
		Redis:   redisClient,
		Storage: NewStorageManager(client),
		WeChat:  NewWeChatManager(dataCfg.GetWechat()),
		JWT:     NewJWTManager(jwtCfg.GetSecret(), jwtCfg.GetExpire(), jwtCfg.GetIssuer()),
		Auth: AuthPolicy{
			AllowConcurrent: authCfg.GetAllowConcurrent(),
			ShareToken:      authCfg.GetShareToken(),
		},
		stopCh: make(chan struct{}),
	}
	if sqlDB != nil {
		startDBPoolMetricsSampler(sqlDB, configx.ParseDurationOrDefault(obsCfg.GetDbPoolSampleInterval(), 15*time.Second), res.stopCh)
	}
	return res, nil
}

func (r *Resources) Close() error {
	if r == nil {
		return nil
	}
	if r.stopCh != nil {
		close(r.stopCh)
		r.stopCh = nil
	}
	if r.Redis != nil {
		_ = r.Redis.Close()
	}
	if r.Ent == nil {
		if r.SQLDB != nil {
			return r.SQLDB.Close()
		}
		return nil
	}
	return r.Ent.Close()
}

func (r *Resources) withTx(ctx context.Context, fn func(tx *ent.Tx) error) (err error) {
	tx, err := r.Ent.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err = fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("%w: rollback failed: %v", err, rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

func nextID() int64 {
	return idCounter.Add(1)
}
