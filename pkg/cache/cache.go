package cache

import (
	"encoding/json"
	"errors"
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

// TODO: Move to util helpers
func interfaceToString(from interface{}) string {
	record, _ := json.Marshal(from)
	return string(record)
}

func stringToInterface(from string, to interface{}) {
	_ = json.Unmarshal([]byte(from), to)
}

func Set(key string, value interface{}, time time.Duration) {
	if Instance != nil {
		stringVal := interfaceToString(value)
		fmt.Println("------------	SET CACHE:", key, stringVal)
		Instance.Set(key, stringVal, time)
	}
}

func SetString(key string, value string, time time.Duration) {
	if Instance != nil {
		Instance.Set(key, value, time)
	}
}

func Get(key string, to interface{}) error {
	value := Instance.Get(key)
	fmt.Println("------------	GET CACHE", key, value)
	if value != nil && !value.IsExpired() {
		err := json.Unmarshal([]byte(value.Value()), &to)
		return err
	}

	return errors.New("item not found")
}

func Remove(key string) {
	Instance.Delete(key)
}

func GetString(key string) (string, error) {
	value := Instance.Get(key)
	if value != nil && !value.IsExpired() {
		return value.Value(), nil
	}

	return "", errors.New("item not found")
}
