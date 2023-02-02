package visca

import (
	"log"
	"testing"
)

func TestDiscovery(t *testing.T) {
	dev := Discover()
	log.Println(dev)
}
