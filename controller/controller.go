package controller

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhost/hasura-auth/connector"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Error string `json:"error"`
}

type JWTSecret struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type Controller struct {
	proxySchema string
	proxyHost   string
	hasuraConn  *connector.HasuraPostgresConnector
	logger      logrus.FieldLogger
	jwtSecret   JWTSecret
}

func New(proxySchema, proxyHost string, jwtSecret JWTSecret, logger logrus.FieldLogger) *Controller {
	return &Controller{
		proxySchema: proxySchema,
		proxyHost:   proxyHost,
		logger:      logger,
		jwtSecret:   jwtSecret,
		hasuraConn:  connector.NewHasuraPostgresConnector("http://graphql-engine:8080/v1/query", "hello123"),
	}
}

func (ctrl *Controller) SetupRouter(
	trustedProxies []string, apiRootPrefix string, middleware ...gin.HandlerFunc,
) (*gin.Engine, error) {
	router := gin.New()
	if err := router.SetTrustedProxies(trustedProxies); err != nil {
		return nil, fmt.Errorf("problem setting trusted proxies: %w", err)
	}

	router.Use(gin.Recovery())

	for _, mw := range middleware {
		router.Use(mw)
	}

	router.POST("/token", ctrl.RefreshToken)
	router.Any("/:path", ctrl.Proxy)
	router.Any("/:path/*path", ctrl.Proxy)

	return router, nil
}

func (ctrl *Controller) Proxy(ctx *gin.Context) {
	req := ctx.Request

	req.URL.Scheme = ctrl.proxySchema
	req.URL.Host = ctrl.proxyHost

	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(req)
	if err != nil {
		ctx.AbortWithStatusJSON(
			500,
			Error{fmt.Sprintf("error in roundtrip: %v", err)},
		)
		return
	}

	// step 3: return real server response to upstream.
	for k, vv := range resp.Header {
		for _, v := range vv {
			ctx.Header(k, v)
		}
	}
	defer resp.Body.Close()

	if _, err := bufio.NewReader(resp.Body).WriteTo(ctx.Writer); err != nil {
		_ = ctx.Error(fmt.Errorf("error writing response: %w", err))
		return
	}
}

func (ctrl *Controller) Healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"healthz": "ok",
	})
}
