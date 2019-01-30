package file

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// 获取当前路径
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// 获取当前目录
func SelfDir() string {
	return Dir(SelfPath())
}

// 判断文件或目录是否存在
func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 判断是否为文件
func IsFile(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// 判断是否为目录
func IsDir(path string) bool {
	f, e := os.Stat(path)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// 判断文件是否可写
func IsWritable(path string) bool {
	err := unix.Access(path, unix.O_RDWR)
	if err == nil {
		return true
	}
	return false
}

// 获取文件名后缀,不包括.
func Ext(file string) string {
	f := filepath.Ext(file)
	if f == "" {
		return f
	}
	return f[1:]
}

// 返回路径的文件名
func Basename(path string) string {
	return filepath.Base(path)
}

// 返回路径
func Dir(path string) string {
	return filepath.Dir(path)
}

// 创建文件夹
func MkdirAll(path string) error {
	if !Exist(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	if !IsWritable(path) {
		return errors.New("path [" + path + "] is not writable!")
	}
	return nil
}

// 写文件
func PutContent(file string, content []byte) (int, error) {
	err := MkdirAll(filepath.Dir(file))
	if err != nil {
		return 0, err
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}
	n, err := f.Write(content)
	f.Close()
	return n, err
}

// 读文件
func GetContent(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// 文件拷贝
func Copy(src, dst string) error {
	// 软链文件，指向原始文件
	linfo, err := os.Readlink(src)
	if err == nil || len(linfo) > 0 {
		return os.Symlink(linfo, dst)
	}

	// 普通文件
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// 返回指定目录下的文件列表
func ListFiles(path string, abs bool) []string {
	if !Exist(path) || !IsDir(path) {
		return nil
	}

	p := ""
	if abs {
		p, _ = filepath.Abs(path)
		p += "/"
	}

	var items []string
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		items = append(items, p+f.Name())

	}
	return items
}

// 返回指定目录下的所有文件列表，包括子目录下的文件
func ListAll(path string, abs bool) []string {
	if !Exist(path) || !IsDir(path) {
		return nil
	}

	p := ""
	if abs {
		p, _ = filepath.Abs(path)
		p += "/"
	}

	var items []string
	filepath.Walk(path, func(filename string, f os.FileInfo, err error) error {
		// 排除根目录
		if filename == path {
			return nil
		}
		if f.IsDir() {
			if abs {
				p = p + f.Name() + "/"
			}
			return nil
		}
		items = append(items, p+f.Name())
		return nil
	})

	return items
}

// 删除文件,有子目录时会报错
func Remove(path string) bool {
	if !Exist(path) {
		return true
	}
	err := os.Remove(path)
	if err == nil {
		return true
	}
	return false
}

// 删除全部
func RemoveAll(path string) bool {
	if !Exist(path) {
		return true
	}
	err := os.RemoveAll(path)
	if err == nil {
		return true
	}
	return false
}

// 在指定路径下查找文件
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); Exist(fullpath) {
			return
		}
	}
	fullpath = ""
	err = errors.New(filename + " not fould in paths")
	return
}
