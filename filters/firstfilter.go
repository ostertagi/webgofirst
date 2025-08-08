package filters

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/server/web/context"
)

var LogManager = func(ctx *context.Context) {
	fmt.Println("IP :: " + ctx.Request.RemoteAddr + ",Time :: " + time.Now().Format(time.RFC850))
}

func MyFilterFunc(ctx *context.Context) {
	fmt.Println("do something here")
}
