package main

import (
	"context"
	"log"

	"github.com/7sunarni/route-panel/pkg/route"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListRoutes() []route.Route {
	ret, err := route.List()
	if err != nil {
		log.Printf("list route failed %s", err)
	}
	return ret
}

func (a *App) AddRoute(r route.Route) string {
	log.Printf("add route %s %s %s", r.Destination, r.Mask, r.InterfaceName)
	err := route.Add(r)
	if err != nil {
		log.Printf("add route failed %s", err)
	}
	if err == nil {
		return ""
	}
	return err.Error()
}

func (a *App) DeleteRoute(r route.Route) string {
	log.Printf("delete route %s %s %s", r.Destination, r.Mask, r.InterfaceName)
	err := route.Delete(r)
	if err != nil {
		log.Printf("delete route failed %s", err)
	}
	if err == nil {
		return ""
	}
	return err.Error()
}

func (a *App) EditRoute(before, after route.Route) string {
	log.Printf("delete route before %s %s %s, after %s %s %s",
		before.Destination, before.Mask, before.InterfaceName,
		after.Destination, after.Mask, after.InterfaceName,
	)
	if err := route.Delete(before); err != nil {
		log.Printf("delete route failed %s", err)
		return err.Error()
	}
	// TODO: rollback
	log.Printf("edit route %s %s %s", after.Destination, after.Gateway, after.InterfaceName)
	err := route.Add(after)
	if err != nil {
		log.Printf("add route failed %s", err)
	}
	if err == nil {
		return ""
	}
	return err.Error()
}
