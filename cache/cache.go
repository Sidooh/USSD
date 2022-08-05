package cache

import (
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"time"
)

var (
	Instance *ttlcache.Cache[string, string]
)

func Init() {
	fmt.Println("Initializing USSD subsystem cache")

	Instance = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](15*time.Minute),
		ttlcache.WithDisableTouchOnHit[string, string](),
	)

	go Instance.Start() // starts automatic expired item deletion
}
