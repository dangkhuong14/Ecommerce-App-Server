package component

import (
	"flag"

	sctx "github.com/viettranx/service-context"
)

type configComponent struct {
	id             string
	urlRPCCategory string
}

func NewConfigComponent(id string) *configComponent {
	return &configComponent{
		id: id,
	}
}

func (c *configComponent) GetURLRPCCategory() string {return c.urlRPCCategory}

func (c *configComponent) ID() string { return c.id }

func (c *configComponent) InitFlags() {
	flag.StringVar(
		&c.urlRPCCategory,
		"url-rpc-category",
		"http://localhost:3000/v1/category/rpc",
		"URL of category service using for RPC",
	)
}

func (c *configComponent) Activate(_ sctx.ServiceContext) error {
	return nil
}

func (c *configComponent) Stop() error {
	return nil
}
