package main

// 测试主函数
func main() {
	sl := NewSkipList()

	data := []SkipListDataNode{
		NewSkipListDataNode(1, 2, 100),
		NewSkipListDataNode(3, 1, 200),
		NewSkipListDataNode(2, 5, 300),
		NewSkipListDataNode(1, 1, 400),
		NewSkipListDataNode(3, 2, 500),
	}
	levelList := []int{0, 0, 3, 3, 3}

	for i, v := range data {
		sl.InsertWithLevel(&v, levelList[i])
	}
	sl.Delete(&data[0])

	sl.FindTopK(1)

}
