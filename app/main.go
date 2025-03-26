package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Entrypoint(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	InitBootSystem(ctx)

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "Buddy Server",
		Prefork:               false,
		ServerHeader:          "Buddy Server",
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${latency} | ${status} | ${method} | ${path} | req = ${body} | rsp = ${resBody}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Shanghai",
	}))

	app.Get("/status", func(c *fiber.Ctx) error {
		return c.SendString(" buddy server is OK!")

	})

	SetRoute(app)
	mux := adaptor.FiberApp(app)
	addr := viper.GetString("server.webaddr")
	webEngine := nbhttp.NewEngine(nbhttp.Config{
		Name:                    "Buddy Web Server",
		Network:                 "tcp",
		IOMod:                   nbhttp.IOModMixed,
		MaxBlockingOnline:       100_000,
		Addrs:                   []string{addr},
		MaxLoad:                 4 << 20,
		ReleaseWebsocketPayload: true,
		ReadBufferSize:          4 << 10,
		Handler:                 mux,
	})

	slog.Info("web service boot", "addr", addr)
	logging.SetLevel(logging.LevelNone)
	err := webEngine.Start()
	if err != nil {
		slog.Error("boot service failed", "error", err)
		cancel()
		return
	}

	// websocket engine
	wsmux := &http.ServeMux{}
	wsmux.HandleFunc("/ws", onWebsocket)

	wsaddr := viper.GetString("server.wsaddr")
	wsEngine := nbhttp.NewEngine(nbhttp.Config{
		Name:                    "Buddy Websocket Server",
		Network:                 "tcp",
		NPoller:                 runtime.NumCPU(),
		IOMod:                   nbhttp.IOModMixed,
		MaxBlockingOnline:       10_000,
		Addrs:                   []string{wsaddr},
		MaxLoad:                 4 << 10,
		ReleaseWebsocketPayload: true,
		ReadBufferSize:          4 << 10,
		Handler:                 wsmux,
	})

	slog.Info("websocket boot", "addr", wsaddr)
	err = wsEngine.Start()
	if err != nil {
		slog.Error("boot service failed", "error", err)
		cancel()
		return
	}

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	webEngine.Shutdown(ctx)
	wsEngine.Shutdown(ctx)

	slog.Info("service shutdown")
	cancel()
	// wait clean code execution compeleted
	time.Sleep(time.Second)

	return
}
