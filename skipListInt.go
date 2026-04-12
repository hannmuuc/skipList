package main

type SkipListInt struct {
	list *SkipList
}

func newSkipListInt() *SkipListInt {
	return &SkipListInt{
		list: NewSkipList(),
	}
}

func (s *SkipListInt) Size() int {
	return s.list.Size()
}

func (s *SkipListInt) Insert(key int, val int) {
	data := NewSkipListDataNode(0, key, val)
	s.list.Insert(&data)
}

// Delete 删除节点
func (s *SkipListInt) Delete(key int) {
	data := NewSkipListDataNode(0, key, 0)
	s.list.Delete(&data)
}

// key val hasNode
func (s *SkipListInt) FindTopK(k int) (int, int, bool) {
	node, hasNode := s.list.FindTopK(k)
	return node.key2, node.val, hasNode
}

// 小于等于目标数据的最大节点 + 其排名 + 是否存在
// key rank hasNode
func (s *SkipListInt) LowerBound(key int) (int, int, bool) {
	data := NewSkipListDataNode(0, key, 0)
	node, rank, hasNode := s.list.LowerBound(&data)
	return node.key2, rank, hasNode
}
