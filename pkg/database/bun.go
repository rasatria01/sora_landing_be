package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/pkg/config"
	"sora_landing_be/pkg/logger"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var dbInstance = &Database{}
var once = &sync.Once{}

const (
	OperationFieldName     = "operation"
	OperationTimeFieldName = "operation_time_ms"
)

func InitDB(c config.Database) {
	once.Do(func() {
		// Use DATABASE_URL if set, otherwise build from config
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
		}

		sqlDB, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to database", err)
		}

		instance := bun.NewDB(sqlDB, pgdialect.New())
		dbInstance.DB = instance
		dbInstance.RegisterModel((*domain.ArticleTag)(nil))
		dbInstance.SetMaxIdleConns(c.MaxOpenIdleConn)
		dbInstance.SetMaxOpenConns(c.MaxOpenConn)
		dbInstance.SetConnMaxIdleTime(c.MaxIdleConn)

		if err := dbInstance.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

	})
	dbInstance.AddQueryHook(NewQueryHook(logger.Log, 200*time.Millisecond))
}

func NewQueryHook(logger *logger.ZapLogger, slowDuration time.Duration) *QueryHook {
	return &QueryHook{
		logger:       logger,
		slowDuration: slowDuration,
	}
}

func (qh QueryHook) BeforeQuery(ctx context.Context, _ *bun.QueryEvent) context.Context {
	return ctx
}

func (qh QueryHook) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	queryDuration := time.Since(event.StartTime)
	fields := []zapcore.Field{
		zap.String(OperationFieldName, event.Operation()),
		zap.Int64(OperationTimeFieldName, queryDuration.Milliseconds()),
	}

	if event.Err != nil {
		fields = append(fields, zap.Error(event.Err))
		qh.logger.Error(event.Query, fields...)
		return
	} else {
		qh.logger.Info(Censored(event.Query), fields...)
	}

	if queryDuration >= qh.slowDuration {
		qh.logger.Debug(event.Query, fields...)
	}
}

func GetDB() *Database {
	return dbInstance
}

type txKey string

var txKeyData txKey = "tx"

func RunInTx(
	ctx context.Context, db *Database, opts *sql.TxOptions, fn func(ctx context.Context, tx bun.Tx) error,
) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKeyData, &tx)

	var done bool

	defer func() {
		if !done {
			_ = tx.Rollback()
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return err
	}

	done = true
	return tx.Commit()
}
