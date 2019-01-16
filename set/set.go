// 简单集合，内部元素排重、无序
package set

type Set interface {
	// 添加元素
	Add(i interface{}) bool

	// 返回元素数量
	Count() int

	// 清空
	Clear()

	// 复制，深拷贝
	Clone() Set

	// 判存
	Contains(i ...interface{}) bool

	// 判断元素是否一致
	Equal(other Set) bool

	// 并集
	Union(other Set) Set

	// 交集
	Intersect(other Set) Set

	// 对称差集，并集 - 交集
	SymmetricDifference(other Set) Set

	// 差集，返回原集合存在，other中不存在的元素组成的子集
	Difference(other Set) Set

	// 判断当前集合是否为other的子集，不能相等
	IsProperSubset(other Set) bool

	// 判断当前集合是否为other的超集，不能相等
	IsProperSuperset(other Set) bool

	// 判断当前集合是否为other的子集，可以相等
	IsSubset(other Set) bool

	// 判断当前集合是否为other的超集，可以相等
	IsSuperset(other Set) bool

	// 为每个元素执行一个func，当返回true时，停止迭代
	Each(func(interface{}) bool)

	// 返回一个用于range的channel
	Iter() <-chan interface{}

	// 返回用于range的迭代器
	Iterator() *Iterator

	// 删除一个元素
	Remove(i interface{})

	// 集合的字符串表示，用于打印
	String() string

	// 删除并返回一个元素
	Pop() interface{}

	// 返回所有子集
	PowerSet() Set

	// 返回本集合与other的笛卡尔乘积
	CartesianProduct(other Set) Set

	// 将集合类型转换为slice
	ToSlice() []interface{}
}

// 创建一个线程安全的新集合，返回集合地址
func NewSet(s ...interface{}) Set {
	set := newThreadSafeSet()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// 创建一个线程安全的新集合，返回集合地址
func NewSetWith(elts ...interface{}) Set {
	return NewSetFromSlice(elts)
}

// 使用slice创建一个线程安全的新集合，返回集合地址
func NewSetFromSlice(s []interface{}) Set {
	a := NewSet(s...)
	return a
}

// 创建一个非线程安全的新集合，返回集合地址
func NewThreadUnsafeSet() Set {
	set := newThreadUnsafeSet()
	return &set
}

// 使用slice创建一个非线程安全的新集合，返回集合地址
func NewThreadUnsafeSetFromSlice(s []interface{}) Set {
	a := NewThreadUnsafeSet()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
