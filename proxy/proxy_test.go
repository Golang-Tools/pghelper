package proxy

import (
	"context"
	"testing"

	log "github.com/Golang-Tools/loggerhelper"
	"github.com/go-pg/pg"
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
	err := proxy.InitFromURL(TEST_PG_URL)
	if err != nil {
		assert.FailNow(t, err.Error(), "init from url error")
	}
	defer proxy.Close()
	ctx := context.Background()
	if err != nil {
		assert.FailNow(t, err.Error(), "FlushDB error")
	}

	// 测试
	var res int
	_, err = proxy.QueryOneContext(ctx, pg.Scan(&res), "SELECT 1")
	if err != nil {
		assert.FailNow(t, err.Error(), "conn set error")
	}
	log.Info("get result", log.Dict{"res": res})
	// assert.Equal(t, "ok", res)
}

// func Test_pgProxy_reset(t *testing.T) {
// 	proxy := New()
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		t.Log("inited db")
// 		return nil
// 	})
// 	err := proxy.InitFromURL(TEST_pg_URL)
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "init from url error")
// 	}
// 	defer proxy.Close()
// 	options, err := pg.ParseURL(TEST_pg_URL)
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "reset from url error")
// 	}
// 	cli := pg.NewClient(options)
// 	err = proxy.SetConnect(cli)
// 	if err != nil {
// 		assert.Equal(t, ErrProxyAllreadySettedUniversalClient, err)
// 	} else {
// 		assert.FailNow(t, "not get error")
// 	}
// }

// func Test_pgProxy_InitFromURL_with_cb(t *testing.T) {
// 	proxy := New()
// 	// 测试cb顺序执行
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		ctx := context.Background()
// 		_, err := cli.Set(ctx, "a", "t", 0).Result()
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		ctx := context.Background()
// 		res, err := cli.Get(ctx, "a").Result()
// 		if err != nil {
// 			return err
// 		}
// 		assert.Equal(t, "t", res)
// 		return nil
// 	})
// 	err := proxy.InitFromURL(TEST_pg_URL)
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "init from url error")
// 	}
// 	defer proxy.Close()
// }

// func Test_pgProxy_regist_cb_after_InitFromURL(t *testing.T) {
// 	proxy := New()
// 	err := proxy.InitFromURL(TEST_pg_URL)
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "init from url error")
// 	}
// 	defer proxy.Close()
// 	err = proxy.Regist(func(cli pg.UniversalClient) error {
// 		ctx := context.Background()
// 		_, err := cli.Set(ctx, "a", "t", 0).Result()
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		assert.Equal(t, ErrProxyAllreadySettedUniversalClient, err)
// 	} else {
// 		assert.FailNow(t, "not get error")
// 	}
// }

// func Test_pgProxy_InitFromURL_with_parallel_cb(t *testing.T) {
// 	proxy := New()
// 	// 测试cb顺序执行
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		time.Sleep(3 * time.Second)
// 		ctx := context.Background()
// 		res, err := cli.Get(ctx, "a").Result()
// 		if err != nil {
// 			return err
// 		}
// 		assert.Equal(t, "e", res)
// 		_, err = cli.Set(ctx, "a", "t", 0).Result()
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		time.Sleep(2 * time.Second)
// 		ctx := context.Background()
// 		res, err := cli.Get(ctx, "a").Result()
// 		if err != nil {
// 			return err
// 		}
// 		assert.Equal(t, "s", res)
// 		_, err = cli.Set(ctx, "a", "e", 0).Result()
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	proxy.Regist(func(cli pg.UniversalClient) error {
// 		time.Sleep(1 * time.Second)
// 		ctx := context.Background()
// 		_, err := cli.Set(ctx, "a", "s", 0).Result()
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	err := proxy.InitFromURLParallelCallback(TEST_pg_URL)
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "init from url error")
// 	}
// 	defer proxy.Close()
// 	time.Sleep(4 * time.Second)
// 	ctx := context.Background()
// 	res, err := proxy.Get(ctx, "a").Result()
// 	if err != nil {
// 		assert.FailNow(t, err.Error(), "init from url error")
// 	}
// 	assert.Equal(t, "t", res)
// }
