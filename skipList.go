package main

import (
	"fmt"
	"math/rand"
)

const (
	MaxLevel = 40  // 跳表最大层级
	Prob     = 0.5 // 层级晋升概率
)

// 跳表核心结构体
type SkipList struct {
	head     *SkipListNode // 哨兵头节点
	nodeSize int           // 跳表总元素个数
}

// 数据节点：复合key(key1+key2)排序，val存储值
type SkipListDataNode struct {
	key1 int // 第一排序键
	key2 int // 第二排序键
	val  int // 存储的值
}

// 跳表节点：包含多层指针+索引计数
type SkipListNode struct {
	data     SkipListDataNode // 节点数据
	level    int              // 当前节点的最大层级
	next     []*SkipListNode  // 每一层的后继节点
	nextSize []int            // 每一层到后继节点的元素间隔数（核心：支持TopK）
}

// NewSkipList 初始化跳表
func NewSkipList() *SkipList {
	// 头节点为哨兵节点（key=-1，保证永远最小）
	head := NewSkipListNode(NewSkipListDataNode(-1, 0, 0), MaxLevel)
	return &SkipList{
		head:     head,
		nodeSize: 0,
	}
}

// NewSkipListDataNode 创建数据节点
func NewSkipListDataNode(key1 int, key2 int, val int) SkipListDataNode {
	return SkipListDataNode{key1: key1, key2: key2, val: val}
}

// Compare 复合键比较：key1优先，key2次之
// 返回值：1(大于)、-1(小于)、0(等于)
func (d SkipListDataNode) Compare(d2 *SkipListDataNode) int {
	if d.key1 > d2.key1 {
		return 1
	} else if d.key1 < d2.key1 {
		return -1
	}
	if d.key2 > d2.key2 {
		return 1
	} else if d.key2 < d2.key2 {
		return -1
	}
	return 0
}

// GetVal 获取值
func (d *SkipListDataNode) GetVal() int {
	return d.val
}

// SetVal 修改值
func (d *SkipListDataNode) SetVal(val int) {
	d.val = val
}

// NewSkipListNode 创建跳表节点
func NewSkipListNode(data SkipListDataNode, level int) *SkipListNode {
	return &SkipListNode{
		data:     data,
		level:    level,
		next:     make([]*SkipListNode, level+1), // 层级0~level
		nextSize: make([]int, level+1),
	}
}

// Size 获取跳表元素总数
func (s *SkipList) Size() int {
	return s.nodeSize
}

// randomLevel 随机生成节点层级（保证跳表平衡性）
func (s *SkipList) randomLevel() int {
	level := 0
	// 满足概率则晋升，最大不超过MaxLevel-1
	for rand.Float32() < Prob && level < MaxLevel {
		level++
	}
	return level
}

// FindTopK 查找第K个元素（排名从1开始）
func (s *SkipList) FindTopK(k int) (*SkipListDataNode, bool) {
	if k <= 0 || k > s.nodeSize {
		return nil, false
	}
	node := s.head
	remain := k // 剩余需要跳过的元素数

	// 从最高层向下遍历
	for i := MaxLevel; i >= 0; i-- {
		// 跳过当前层的连续节点，直到剩余数不足
		for node.next[i] != nil && remain > node.nextSize[i] {
			remain -= node.nextSize[i]
			node = node.next[i]
		}
	}

	// 到达0层的目标节点
	node = node.next[0]
	return &node.data, true
}

// LowerBound 下界查找：返回<=目标数据的最大节点 + 其排名 + 是否存在
func (s *SkipList) LowerBound(d *SkipListDataNode) (*SkipListDataNode, int, bool) {
	if d == nil {
		return nil, 0, false
	}
	node := s.head
	rank := 0 // 最终节点的排名

	// 从最高层向下遍历
	for i := MaxLevel; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].data.Compare(d) <= 0 {
			rank += node.nextSize[i]
			node = node.next[i]
		}
	}

	// 头节点无数据
	if node == s.head {
		return nil, rank, false
	}
	return &node.data, rank, true
}

func (s *SkipList) Insert(d *SkipListDataNode) bool {
	newLevel := s.randomLevel()
	return s.InsertWithLevel(d, newLevel)

}

func (s *SkipList) Delete(d *SkipListDataNode) bool {
	if d == nil {
		return false
	}

	preNodes := make([]*SkipListNode, MaxLevel+1)
	current := s.head

	for i := MaxLevel; i >= 0; i-- {

		for current.next[i] != nil && current.next[i].data.Compare(d) < 0 {
			current = current.next[i]
		}
		preNodes[i] = current
	}

	if current.next[0] == nil || current.next[0].data.Compare(d) != 0 {
		return false
	}

	deleteLevel := current.next[0].level
	for i := 0; i <= deleteLevel; i++ {
		if preNodes[i].next[i] != nil {
			preNodes[i].nextSize[i] += preNodes[i].next[i].nextSize[i]
			preNodes[i].next[i] = preNodes[i].next[i].next[i]
		}
		preNodes[i].nextSize[i]--
	}

	for i := deleteLevel + 1; i <= MaxLevel; i++ {
		preNodes[i].nextSize[i]--
	}

	s.nodeSize--

	return true
}

// Insert 插入/更新数据（key相同则更新val，key不同则插入）
func (s *SkipList) InsertWithLevel(d *SkipListDataNode, level int) bool {
	if d == nil {
		return false
	}

	// 记录每一层的前驱节点
	preNodes := make([]*SkipListNode, MaxLevel+1)
	// 记录每一层跳过的元素数
	skipSize := make([]int, MaxLevel+1)
	current := s.head

	// 1. 查找每一层的前驱节点
	for i := MaxLevel; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].data.Compare(d) < 0 {
			skipSize[i+1] += current.nextSize[i]
			current = current.next[i]
		}
		preNodes[i] = current
	}

	// 2. 键已存在：直接更新值
	if current.next[0] != nil && current.next[0].data.Compare(d) == 0 {
		current.next[0].data.SetVal(d.GetVal())
		return true
	}

	// 3. 键不存在：生成新节点层级
	newLevel := level
	newNode := NewSkipListNode(*d, newLevel)

	// 4. 插入新节点，更新指针和nextSize
	for i := 0; i <= newLevel; i++ {
		newNode.next[i] = preNodes[i].next[i]
		// 新节点的间隔数 = 原前驱的间隔数 - 已跳过的数
		newNode.nextSize[i] = preNodes[i].nextSize[i] - skipSize[i]
		// 前驱节点指向新节点，间隔数=跳过数+1
		preNodes[i].next[i] = newNode
		preNodes[i].nextSize[i] = skipSize[i] + 1
		skipSize[i+1] += skipSize[i]
	}

	// 5. 新节点层级之上的所有前驱，间隔数+1
	for i := newLevel + 1; i <= MaxLevel; i++ {
		preNodes[i].nextSize[i]++
	}

	// 总元素数+1
	s.nodeSize++
	return true
}

// Print 打印跳表（0层，所有元素，按排序顺序）
func (s *SkipList) Print(printLevel int) {
	fmt.Printf("跳表总元素数: %d\n", s.nodeSize)

	for i := 0; i <= printLevel; i++ {
		fmt.Printf("--------level:%d-------\n", i)
		node := s.head
		for node != nil {
			fmt.Printf("key1=%d, key2=%d, val=%d nextSize=%d\n", node.data.key1, node.data.key2, node.data.val, node.nextSize[i])
			node = node.next[i]
		}

	}

	fmt.Println("------------------------")
}
