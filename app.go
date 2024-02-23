package main

import (
	"context"
	"fmt"
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

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ListRoutes() []route.Route {
	ret, err := route.List()
	if err != nil {
		log.Printf("list route failed %s", err)
		return nil
	}
	return ret
}

func (a *App) AddRoute(r route.Route) error {
	return route.Add(r)
}

func (a *App) DeleteRoute(r route.Route) error {
	return route.Delete(r)
}
