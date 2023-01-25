package cache

import (
	"github.com/jellydator/ttlcache/v3"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	Init()

	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	item := Instance.Set("key", "value", ttlcache.DefaultTTL)

	if item == nil {
		t.Errorf("item = nil; want item")
	}

	if item.Key() != "key" {
		t.Errorf("item.Key() = %s; want 'key'", item.Key())
	}

	if item.Value() != "value" {
		t.Errorf("item.Value() = %s; want 'value'", item.Value())
	}

	if item.IsExpired() {
		t.Errorf("item.IsExpired() = %v; want 'false'", true)
	}

	if item.TTL().String() != (15 * time.Minute).String() {
		t.Errorf("item.TTL() is %s; want 15m0s", item.TTL())
	}

	item = Instance.Set("key", "value", 5*time.Second)

	if item.TTL().String() != (5 * time.Second).String() {
		t.Errorf("item.TTL() is %s; want 5s", item.TTL())
	}
}

func TestGet(t *testing.T) {
	nonExistentItem := Instance.Get("key1")

	if nonExistentItem != nil {
		t.Errorf("nonExistentItem is %v; want nil", nonExistentItem)
	}
}

func TestExpiry(t *testing.T) {
	item := Instance.Set("key", "value", 1*time.Second)

	if item == nil {
		t.Errorf("item = nil; want item")
	}

	exists := Instance.Get("key")

	if exists == nil {
		t.Errorf("exists = %v; want item", exists)
	}

	time.Sleep(1 * time.Second)

	nonExistentItem := Instance.Get("key")

	if nonExistentItem != nil {
		t.Errorf("nonExistentItem = %v; want nil", nonExistentItem)
	}
}
