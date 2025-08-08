package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/server/web"
)

type CacheController struct {
	web.Controller
}

var beegoCache cache.Cache
var err error

func init() {
	beegoCache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		panic(err)
	}
	beegoCache.Put(context.Background(), "foo", "bar", 100000*time.Second)
}

func (ctrl *CacheController) GetFromCache() {
	val, err := beegoCache.Get(context.Background(), "foo")
	if err != nil {
		panic(err)
	}
	ctrl.Ctx.WriteString("Hello " + fmt.Sprintf("%v", val))
}
