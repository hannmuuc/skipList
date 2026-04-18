package main

import "fmt"

type test1 struct {
	v int
}

func (t test1) Compare(other test1) int {
	return t.v - other.v
}

func (t test1) Identity() test1 {
	return test1{v: 0}
}
func (t test1) Op(other test1) test1 {
	return test1{v: t.v + other.v}
}

func (t test1) toInt() int {
	return t.v
}

// 测试主函数
func main() {

	skipListAbstract := NewSkipListAbstract[test1, int, test1]()
	skipList := NewSkipList()

	fmt.Printf("skipListAbstract: %+v skipList: %+v", skipListAbstract, skipList)

}
