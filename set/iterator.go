// 集合迭代器示例
/*
	// 迭代处理每一个元素
	for elem := range s1.Iter() {
		fmt.Println(elem)
	}

	// 使用迭代器处理每一个元素，可以随时Stop
	it = s2.Iterator()
	for elem := range it.C {
		if (elem == 'xxx') {
			it.Stop()
		}
	}

	// 使用函数处理每一个元素，函数返回true时停止迭代
	s3.Each(func(elem interface{}) bool {
		fmt.Println(elem, "====4")
		return elem.(string) <= "2"
	})
*/

package set

// 迭代器结构，C用来遍历集合元素
type Iterator struct {
	C    <-chan interface{}
	stop chan struct{}
}

// 停止迭代操作，同时关闭C的channel
func (i *Iterator) Stop() {
	// 允许多次调用Stop函数，恢复多次关闭channel导致的panic
	defer func() {
		recover()
	}()

	close(i.stop)

	// 释放C中的元素
	for range i.C {
	}
}

// 返回一个迭代器实例
func newIterator() (*Iterator, chan<- interface{}, <-chan struct{}) {
	itemChan := make(chan interface{})
	stopChan := make(chan struct{})
	return &Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
