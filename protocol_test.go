package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewRequest(t *testing.T) {
	Convey("Should be able to instantiate a new Request", t, func() {
		payload, err := NewRequest("test", []string{})
		So(err, ShouldEqual, nil)
		payloadStr := string(payload[:])
		So(payloadStr, ShouldEqual, "{\"command\":\"test\",\"arguments\":[]}\n")
	})
}

func TestNewResponseListCheckMethods(t *testing.T) {
	Convey("Should be able to new Response '__list_check_methods'", t, func() {
		payload := []byte(
			"{\"name\":\"__list_check_methods\",\"stdout\":[\"check_test\",\"check_second_test\"]}")
		resp, err := NewResponse(payload)
		So(err, ShouldEqual, nil)
		So(resp.Name, ShouldEqual, "__list_check_methods")
	})
}

func TestNewResponseCheckReturn(t *testing.T) {
	Convey("Should be able to new Response 'check_test'", t, func() {
		payload := []byte(
			"{\"name\":\"check_test\",\"status\":0,\"stdout\":[\"check_test output\"],\"metrics\":[{\"okay\":1}]}")
		resp, err := NewResponse(payload)
		So(err, ShouldEqual, nil)
		So(resp.Status, ShouldEqual, 0)
	})
}

/* EOF */
