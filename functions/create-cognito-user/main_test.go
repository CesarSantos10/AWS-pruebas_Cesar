package main

import (
	"context"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("success request", func(t *testing.T) {
		d := deps{}
		k, err := d.handler(context.TODO(), Event{Email: "cesar.santos@pucp.edu.pe", Password: "v8/1vZT4", Name: "Cesar Santos", Case: 5, Username: "4f4fe107-7b5c-43a9-aac2-09410b4b5443", PasswordNew: "v8/1vZT42"})
		if err != nil {
			t.Fatal("Erroraadada")
		}
		if k != "" {
			t.Fatal("Errdorbbccd")
		}
	})
}
