package data

import (
	"context"
	"fmt"
	"github.com/hypercoze/hyper-admin/app/system/service/internal/conf"
	"github.com/hypercoze/hyper-admin/app/system/service/internal/data/ent"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewEntClient)

// Data .
type Data struct {
	db *ent.Client
}

// NewEntClient 初始化 Ent 客户端
func NewEntClient(c *conf.Data, logger log.Logger) *ent.Client {
	l := log.NewHelper(logger)

	client, err := ent.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		l.Fatalf("failed opening connection to postgres: %v", err)
	}

	// 开发环境下开启自动迁移 (Auto Migration)
	// 生产环境建议通过 deploy/sql 下的脚本手动管理
	if err := client.Schema.Create(context.Background()); err != nil {
		l.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, entClient *ent.Client) (*Data, func(), error) {
	l := log.NewHelper(logger)
	d := &Data{
		db: entClient,
	}
	return d, func() {
		l.Info("message", "closing the data resources")
		if err := d.db.Close(); err != nil {
			l.Error(err)
		}
	}, nil
}

// InTx 事务包装函数
// InTx 事务包装函数 - 修复版
func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// 1. 开启事务
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return err
	}

	// 2. 执行业务逻辑
	// 注意：这里需要考虑逻辑执行中的 panic 恢复
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(ctx); err != nil {
		// 3. 发生错误时回滚
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}

	// 4. 成功后提交
	return tx.Commit()
}
