package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/juju/zaputil/zapctx"
	"github.com/kkoch986/ai-skeletons-output-generation/command/internal/common"
	"github.com/kkoch986/ai-skeletons-output-generation/server"
	"github.com/kkoch986/ai-skeletons-output-generation/viseme"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var Command = &cli.Command{
	Name:  "server",
	Usage: "start the http server to handle output generation requests",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Usage: "the hostname to listen on, default to an empty string to listen on all interfaces",
			Value: "",
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "the port for the server to listen on",
			EnvVars: []string{"PORT"},
			Value:   8080,
		},
	},
	Action: func(c *cli.Context) error {
		ctx, cancel := common.InitContext()
		defer cancel()

		host := c.String("host")
		port := c.Int("port")

		ctx, flush, err := common.Logger(ctx, c.String("log-level"))
		if err != nil {
			return err
		}
		defer flush()
		logger := zapctx.Logger(ctx)

		logger.Info("starting web server")
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		r.Use(ginzap.RecoveryWithZap(logger, true))

		// The generate action
		r.POST("/generate", func(c *gin.Context) {
			r := &server.GenerateRequest{
				VisemeSequence: viseme.Sequence{},
			}
			err := c.BindJSON(r)
			if err != nil {
				c.JSON(http.StatusBadRequest, &ErrorResponse{err.Error()})
				return
			}
			resp, err := server.HandleGenerate(ctx, r)
			if err != nil {
				c.JSON(http.StatusInternalServerError, &ErrorResponse{err.Error()})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		})

		server := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: r,
		}

		go server.ListenAndServe()

		<-ctx.Done()
		logger.Info("context cancelled, shutting server down")
		err = server.Shutdown(context.Background())
		if err != nil {
			logger.Error("error shutting down server", zap.Error(err))
		}
		return err
	},
}
