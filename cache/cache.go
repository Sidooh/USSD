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
		ttlcache.WithTTL[string, string](14 * time.Minute),
	)

	go Instance.Start() // starts automatic expired item deletion

	//go func() {
	//	for {
	//		time.Sleep(1 * time.Minute)
	//		fmt.Println("==== RUNNING cache delete")
	//		Instance.Delete("token")
	//	}
	//}()
}
