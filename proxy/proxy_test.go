package proxy

import (
	"context"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
)

//TEST_pg_URL 测试用的pg地址
const TEST_PG_URL = "postgres://postgres:postgres@localhost:5432/test?sslmode=disable"

//TEST_pg_WRONG_URL 错误的pg url
const TEST_PG_WRONG_URL = "redis://localhost:6379"

func Test_pgProxy_InitFromWrongURL(t *testing.T) {
	proxy := New()
	err := proxy.InitFromURL(TEST_PG_WRONG_URL)
	assert.NotNil(t, err)
}

func Test_pgProxy_InitFromURL(t *testing.T) {
	//准备
	proxy := New()
	ctx := context.Background()
	err := proxy.InitFromURL(TEST_PG_URL)
	if err != nil {
		assert.FailNow(t, err.Error(), "init from url error")
	}
	defer proxy.Close()
	// 测试
	var res int
	_, err = proxy.QueryOneContext(ctx, pg.Scan(&res), "SELECT 1")
	if err != nil {
		assert.FailNow(t, err.Error(), "conn set error")
	}
	// log.Info("get result", log.Dict{"res": res})
	assert.Equal(t, 1, res)
}

func Test_pgProxy_InitFromURL_with_cb(t *testing.T) {
	proxy := New()
	ctx := context.Background()
	// 测试cb顺序执行
	proxy.Regist(func(cli *pg.DB) error {
		var res int
		_, err := proxy.QueryOneContext(ctx, pg.Scan(&res), "SELECT 1")
		if err != nil {
			assert.FailNow(t, err.Error(), "conn set error")
		}
		// log.Info("get result", log.Dict{"res": res})
		assert.Equal(t, 1, res)
		return nil
	})
	proxy.Regist(func(cli *pg.DB) error {
		var res int
		_, err := proxy.QueryOneContext(ctx, pg.Scan(&res), "SELECT 10")
		if err != nil {
			assert.FailNow(t, err.Error(), "conn set error")
		}
		// log.Info("get result", log.Dict{"res": res})
		assert.Equal(t, 10, res)
		return nil
	})
	err := proxy.InitFromURL(TEST_PG_URL)
	if err != nil {
		assert.FailNow(t, err.Error(), "init from url error")
	}
	defer proxy.Close()
}
