package godutch_test

import (
	. "github.com/otaviof/godutch"
	gocache "github.com/patrickmn/go-cache"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func populatedCache() *gocache.Cache {
	var metrics []map[string]int
	var cache *gocache.Cache
	var resp *Response

	metrics = append(metrics, map[string]int{"okay": 1})
	cache = gocache.New(time.Minute, 20*time.Second)
	resp = &Response{
		Name:    "check_test",
		Status:  0,
		Stdout:  []string{"Mocked"},
		Metrics: metrics,
	}
	cache.Set("check_test", resp, gocache.DefaultExpiration)

	return cache
}

func TestNewCarbonService(t *testing.T) {
	var err error
	var cfg *Config = mockNewConfig(t)
	var carbonService *CarbonService
	var cache *gocache.Cache = populatedCache()

	carbonService = NewCarbonService(cfg.Service["carbonrelay"], cache)

	Convey("Should be able to send metrics into Carbon", t, func() {
		err = carbonService.Send()
		So(err, ShouldEqual, nil)
	})
}

/* EOF */
