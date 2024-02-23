package route

import (
	"encoding/json"
	"log"
	"testing"
)

func TestListRoute(t *testing.T) {
	routes, err := List()
	if err != nil {
		t.Fatalf("ListRoute() error = %v", err)
	}

	for _, route := range routes {
		data, _ := json.Marshal(route)
		log.Printf("default route %t %s ", route.IsDefaultRoute(), string(data))
	}
}

func TestAddRoute(t *testing.T) {
	if err := Add(Route{
		InterfaceIP: "192.168.149.1",
		Destnation:  "192.168.149.125",
		Mask:        "255.255.255.255",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDelRoute(t *testing.T) {
	if err := Delete(Route{
		InterfaceIP: "192.168.149.1",
		Destnation:  "192.168.149.125",
		Mask:        "255.255.255.255",
	}); err != nil {
		t.Fatal(err)
	}
}
