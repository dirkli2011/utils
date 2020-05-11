package httplib

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttplib(t *testing.T) {

	Convey("Httplib单元测试", t, func() {

		// 简单Get请求
		body, err := Get("http://www.baidu.com")
		if err != nil {
			t.Log(err.Error())
		}
		So(err, ShouldBeNil)
		So(len(body) > 0, ShouldBeTrue)

		// 简单Post请求
		body, err = Post("http://www.baidu.com", map[string]string{
			"hello":    "world",
			"username": "dirkli",
			"file1":    "@../readme.md",
		})
		if err != nil {
			t.Log(err.Error())
		}
		So(err, ShouldBeNil)
		So(len(body) > 0, ShouldBeTrue)

		// 复杂请求，自定义request
		req := Request{
			Method:      METHOD_POST,
			Url:         "http://www.baidu.com",
			ConnTimeout: 1000,
			RWTimeout:   2000,
			PostType:    POST_XML,
			PostParams: `<xml>
<hello>world</hello>
<username>dirkli</username>
</xml>`,
		}
		res := Call(req)
		So(res.Status, ShouldEqual, 200)
		So(len(res.Body) > 0, ShouldBeTrue)

		// 并发请求
		reqs := make(map[interface{}]Request)
		reqs["a"] = Request{
			Method: METHOD_GET,
			Url:    "http://www.baidu.com",
		}
		reqs["b"] = Request{
			Method: METHOD_GET,
			Url:    "http://www.baidu.com/xxx.html",
		}
		resps := make(map[interface{}]*Response)
		for idx, res := range MultiCall(reqs) {
			resps[idx] = res
		}
		So(resps["a"].Status, ShouldEqual, 200)
		So(resps["b"].Status, ShouldEqual, 200)

	})
}
