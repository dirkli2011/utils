package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 单元测试
func TestFile(t *testing.T) {
	path := "/tmp/test/"
	filename := "test.txt"
	filecontent := "hello world"

	Convey("test", t, func() {
		So(RemoveFile(path+filename+".copy"), ShouldEqual, true)
		_, err := FilePutContent(path+filename, []byte(filecontent))
		So(MkdirAll(path), ShouldEqual, nil)
		So(err, ShouldEqual, nil)
		So(IsDir(path), ShouldEqual, true)
		So(IsExist(path), ShouldEqual, true)
		So(IsExist(path+filename), ShouldEqual, true)
		So(IsExist(path+"notexist.txt"), ShouldEqual, false)
		So(IsFile(path+filename), ShouldEqual, true)
		So(GetExt(path+filename), ShouldEqual, "txt")
		So(GetBasename(path+filename), ShouldEqual, filename)
		So(IsWritable(path), ShouldEqual, true)
		So(len(ListFiles(path, false)), ShouldEqual, 1)
		So(len(ListAllFiles(path, true)), ShouldEqual, 1)
		str, _ := FileGetContent(path + filename)
		So(string(str), ShouldEqual, filecontent)
		So(CopyFile(path+filename, path+filename+".copy"), ShouldEqual, nil)
	})

}
