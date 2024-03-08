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

func (a *App) AddRoute(r route.Route) error {
	log.Printf("Add route %s %s %s", r.Destination, r.Mask, r.InterfaceIP)
	err := route.Add(r)
	if err != nil {
		log.Printf("add route failed %s", err)
	}
	return err
}

func (a *App) DeleteRoute(r route.Route) error {
	err := route.Delete(r)
	if err != nil {
		log.Printf("delete route failed %s", err)
	}
	return err
}

func (a *App) EditRoute(before, after route.Route) error {
	if err := route.Delete(before); err != nil {
		log.Printf("delete route failed %s", err)
		return err
	}
	log.Printf("edit route %s %s %s", after.Destination, after.Gateway, after.InterfaceIP)
	err := route.Add(after)
	if err != nil {
		log.Printf("add route failed %s", err)
	}
	return err
}
