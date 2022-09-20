package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("success request", func(t *testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "cesar.santos@pucp.edu.pe", Password: "v8/1vZT4", Name: "Cesar Santos", Case: 2})
		if err != nil {
			t.Fatal("Erroraaaa")
		}
		if k != "" {
			t.Fatal("Errorbbc")
		}
	})
}
