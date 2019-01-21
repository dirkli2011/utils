package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	path := "/tmp/test/"
	filename := "test.txt"
	filecontent := "hello world"

	Convey("文件单元测试", t, func() {
		So(Remove(path+filename+".copy"), ShouldEqual, true)
		_, err := PutContent(path+filename, []byte(filecontent))
		So(MkdirAll(path), ShouldEqual, nil)
		So(err, ShouldEqual, nil)
		So(IsDir(path), ShouldEqual, true)
		So(Exist(path), ShouldEqual, true)
		So(Exist(path+filename), ShouldEqual, true)
		So(Exist(path+"notexist.txt"), ShouldEqual, false)
		So(IsFile(path+filename), ShouldEqual, true)
		So(Ext(path+filename), ShouldEqual, "txt")
		So(Basename(path+filename), ShouldEqual, filename)
		So(IsWritable(path), ShouldEqual, true)
		So(len(ListFiles(path, false)), ShouldEqual, 1)
		So(len(ListAll(path, true)), ShouldEqual, 1)
		str, _ := GetContent(path + filename)
		So(string(str), ShouldEqual, filecontent)
		So(Copy(path+filename, path+filename+".copy"), ShouldEqual, nil)
	})

}
