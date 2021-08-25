package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tiem/library/firstparty/config"
	framework2 "github.com/pingcap/tiem/library/framework"
	"github.com/pingcap/tiem/micro-api/route"
	"testing"
)

func TestMain(m *testing.M) {
	f := framework2.NewUtFramework(framework2.ApiService,
		InitGin)

	err := f.StartService()
	if err == nil {
		m.Run()
	} else {
		f.GetDefaultLogger().Error(err)
	}
}

func InitGin(d *framework2.UtFramework) error {
	gin.SetMode(gin.TestMode)
	g := gin.New()

	route.Route(g)

	port := config.GetClientArgs().RestPort
	if port <= 0 {
		port = config.DefaultRestPort
	}
	addr := fmt.Sprintf(":%d", port)
	if err := g.Run(addr); err != nil {
		d.GetDefaultLogger().Fatal(err)
	}

	return nil
}
