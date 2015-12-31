package godutch_test

import (
	. "github.com/otaviof/godutch"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewRequest(t *testing.T) {
	var err error
	var req []byte
	var reqStr string

	Convey("Should be able to instantiate a new Request", t, func() {
		req, err = NewRequest("test", []string{})
		So(err, ShouldEqual, nil)
		reqStr = string(req[:])
		So(reqStr, ShouldEqual, "{\"command\":\"test\",\"arguments\":[]}\n")
	})
}

func TestNewResponseListCheckMethods(t *testing.T) {
	var err error
	var req []byte = []byte(
		"{\"name\":\"__list_check_methods\",\"stdout\":[\"check_test\",\"check_second_test\"]}")
	var resp *Response

	Convey("Should be able to new Response '__list_check_methods'", t, func() {
		resp, err = NewResponse(req)
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
