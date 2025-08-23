package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"regexp"
)

func (d *Database) InitQuery(ctx context.Context) bun.IDB {
	tx := getTxFromContext(ctx)
	if tx != nil {
		return tx
	}
	return d.DB
}

func (d *Database) NewSelectQ(ctx context.Context, resultObject any) *Database {
	d.Query = d.InitQuery(ctx).NewSelect().Model(resultObject)
	return d
}

func (d *Database) AddMultipleORSearch(value string, columns ...string) *bun.SelectQuery {
	if columns == nil || value == "" {
		return d.Query
	}

	d.Query.WhereGroup(" AND ", func(query *bun.SelectQuery) *bun.SelectQuery {
		for _, column := range columns {
			d.Query.WhereOr(fmt.Sprintf("%s ILIKE ?", column), "%"+value+"%")
		}
		return query
	})

	return d.Query
}

func getTxFromContext(ctx context.Context) *bun.Tx {
	tx, ok := ctx.Value(txKeyData).(*bun.Tx)
	if !ok {
		return nil
	}
	return tx
}

func Censored(query string) string {
	rePassword := regexp.MustCompile(`(?i)("password" = ')[^']*(')`)
	reEncryptedPassword := regexp.MustCompile(`(?i)("encrypted_password" = ')[^']*(')`)

	censoredQuery := rePassword.ReplaceAllString(query, `$1[CENSORED]$2`)
	censoredQuery = reEncryptedPassword.ReplaceAllString(censoredQuery, `$1[CENSORED]$2`)

	return censoredQuery
}
