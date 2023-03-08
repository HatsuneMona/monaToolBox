package bootstrap

import (
	"context"
	"fmt"
	"monaToolBox/global"
	"monaToolBox/routers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RunServer 启动服务器
func RunServer() {
	r := routers.GinRootRouter()
	r.Run(global.Config.Server.Listen + ":" + global.Config.Server.Port)

	srv := &http.Server{
		Addr:    global.Config.Server.Listen + ":" + global.Config.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Fatal(fmt.Sprintf("listen err: %s\n", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Fatal(fmt.Sprintf("Server Shutdown Err: %s\n", err))
	}
	global.Log.Info("Server exiting")
}
