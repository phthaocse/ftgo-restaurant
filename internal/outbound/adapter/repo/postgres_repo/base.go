package postgres_repo

import (
	"context"
	"fmt"
	"ftgo-restaurant/internal/outbound/interface/logger"
	"github.com/jackc/pgconn"
	"github.com/spf13/viper"
	"time"
)

func Init(logger logger.Logger) (*pgconn.PgConn, error) {
	var err error
	dbConnUri := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		viper.GetString("POSTGRESQL_USER"),
		viper.GetString("POSTGRESQL_PASSWORD"),
		viper.GetString("POSTGRESQL_HOST"),
		viper.GetString("POSTGRESQL_PORT"),
		viper.GetString("POSTGRESQL_DB_NAME"),
	)
	connConf, err := pgconn.ParseConfig(dbConnUri)
	if err != nil {
		logger.Errorf("Postgres connection config wrong: %v", err)
	}
	connConf.ConnectTimeout = 5 * time.Second
	maxRetries := 10
	numRetry := 0
	var pgConn *pgconn.PgConn
	for numRetry < maxRetries {
		numRetry++
		logger.Infof("Connecting to Postgres DB at host %s", connConf.Host)
		pgConn, err = pgconn.ConnectConfig(context.Background(), connConf)
		if err != nil {
			logger.Errorf("Connect to Postgres DB at host %s failed %s", connConf.Host, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Infof("Connected to Postgres DB at host %s successfully", connConf.Host)
		break
	}
	return pgConn, err
}
