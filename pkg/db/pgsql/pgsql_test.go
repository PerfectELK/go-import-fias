package pgsql

import (
	"github.com/PerfectELK/go-import-fias/internal/config"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.InitConfig(false)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestConnect(t *testing.T) {
	p := Processor{}

	err := p.Connect()
	if err != nil {
		t.Error(err)
	}
}
