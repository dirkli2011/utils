package set

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSet(t *testing.T) {
	s1 := NewSet("a", "b", "c", "f")
	s2 := NewSetWith("a", "b", "c", "d")
	s3 := NewSetFromSlice([]interface{}{"a", "b", "c"})
	Convey("集合单元测试", t, func() {
		So(s1.Count(), ShouldEqual, 4)
		So(s2.Equal(s1), ShouldBeFalse)
		So(s3.Contains("b"), ShouldBeTrue)

		// 并集
		So(s1.Union(s2).Union(s3).Count() == 5, ShouldBeTrue)
		// 交集
		So(s1.Intersect(s2).Equal(s3), ShouldBeTrue)
		// 对称差集
		So(s1.SymmetricDifference(s2).Equal(NewSet("f", "d")), ShouldBeTrue)
		// 差集
		So(s1.Difference(s2).Equal(NewSet("f")), ShouldBeTrue)
		// 判子集，不能相等
		So(s3.IsProperSubset(s1), ShouldBeTrue)
		// 判超集，不能相等
		So(s1.IsProperSuperset(s3), ShouldBeTrue)
		// 判子集，可相等
		So(s3.IsSubset(s1), ShouldBeTrue)
		// 判超集，可相等
		So(s1.IsSuperset(s3), ShouldBeTrue)

		s11 := s1.Clone()
		So(s1.IsProperSubset(s11), ShouldBeFalse)
		So(s1.IsProperSuperset(s11), ShouldBeFalse)
		So(s1.IsSubset(s11), ShouldBeTrue)
		So(s1.IsSuperset(s11), ShouldBeTrue)

		ss := s11.ToSlice()
		So(len(ss) == 4, ShouldBeTrue)
		So(len(s11.String()), ShouldEqual, len("Set{a, b, c, f}"))

		find := make(map[string]bool)
		find["f"] = false
		find["d"] = false
		find["c"] = false
		// 迭代方式1
		for elem := range s1.Iter() {
			v := elem.(string)
			if _, ok := find[v]; ok {
				find[v] = true
			}
		}
		So(find["f"], ShouldBeTrue)
		So(find["d"], ShouldBeFalse)
		So(find["c"], ShouldBeTrue)

		find["f"] = false
		find["d"] = false
		find["c"] = false
		// 迭代方式2, 调用it.Stop停止迭代
		it := s2.Iterator()
		for elem := range it.C {
			v := elem.(string)
			if _, ok := find[v]; ok {
				find[v] = true
			}
			if find["d"] && find["c"] {
				it.Stop()
			}
		}
		So(find["f"], ShouldBeFalse)
		So(find["d"], ShouldBeTrue)
		So(find["c"], ShouldBeTrue)

		find["f"] = false
		find["d"] = false
		find["c"] = false
		// 迭代方式3, 返回true时停止迭代
		s3.Each(func(elem interface{}) bool {
			v := elem.(string)
			if _, ok := find[v]; ok {
				find[v] = true
			}
			return find["c"]
		})
		So(find["f"], ShouldBeFalse)
		So(find["d"], ShouldBeFalse)
		So(find["c"], ShouldBeTrue)
	})

}
