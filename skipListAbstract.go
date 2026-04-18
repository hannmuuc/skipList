package main

// entry设计
type SkipListComparable[T any] interface {
	Compare(other T) int // ✅ 泛型，无断言
}

type SkipListDataNodeAbstract[K SkipListComparable[K], V any] struct {
	key1 int // 第一排序键
	key  K
	val  V
}

func NewSkipListDataNodeAbstract[K SkipListComparable[K], V any](key1 int, key K, val V) SkipListDataNodeAbstract[K, V] {
	return SkipListDataNodeAbstract[K, V]{key1: key1, key: key, val: val}
}

func (d *SkipListDataNodeAbstract[K, V]) Compare(other *SkipListDataNodeAbstract[K, V]) int {

	if d.key1 > other.key1 {
		return 1
	} else if d.key1 < other.key1 {
		return -1
	}
	return d.key.Compare(other.key)
}

func (d *SkipListDataNodeAbstract[K, V]) GetKey() K {
	return d.key
}

func (d *SkipListDataNodeAbstract[K, V]) GetVal() V {
	return d.val
}

func (d *SkipListDataNodeAbstract[K, V]) SetVal(val V) {
	d.val = val
}

// Monoid 接口
type Monoid[T any] interface {
	Identity() T
	Op(T) T
}

type IntAdd int

func (IntAdd) Identity() IntAdd     { return 0 }
func (a IntAdd) Op(b IntAdd) IntAdd { return a + b }
func (a IntAdd) Set(val IntAdd)     { a = val }

type IntMul int

func (IntMul) Identity() IntMul     { return 1 }
func (a IntMul) Op(b IntMul) IntMul { return a * b }
func (a IntMul) Set(val IntMul)     { a = val }

// 跳表节点：包含多层指针+索引计数
type SkipListNodeAbstract[K SkipListComparable[K], V any, T Monoid[T]] struct {
	data     SkipListDataNodeAbstract[K, V]   // 节点数据
	level    int                              // 当前节点的最大层级
	next     []*SkipListNodeAbstract[K, V, T] // 每一层的后继节点
	nextSize []int                            // 每一层到后继节点的元素间隔数（核心：支持TopK）
	monoid   []T                              // 每一层的元素聚合值
}

func NewSkipListNodeAbstract[K SkipListComparable[K], V any, T Monoid[T]](data SkipListDataNodeAbstract[K, V], level int) *SkipListNodeAbstract[K, V, T] {
	return &SkipListNodeAbstract[K, V, T]{
		data:     data,
		level:    level,
		next:     make([]*SkipListNodeAbstract[K, V, T], level+1),
		nextSize: make([]int, level+1),
		monoid:   make([]T, level+1),
	}
}

// 跳表核心结构体
type SkipListAbstract[K SkipListComparable[K], V any, T Monoid[T]] struct {
	head     *SkipListNodeAbstract[K, V, T] // 哨兵头节点
	nodeSize int                            // 跳表总元素个数
}

// NewSkipList 初始化跳表
func NewSkipListAbstract[K SkipListComparable[K], V any, T Monoid[T]]() *SkipListAbstract[K, V, T] {
	// 头节点为哨兵节点（key=-1，保证永远最小）
	var defaultKey K
	var defaultVal V

	head := NewSkipListNodeAbstract[K, V, T](NewSkipListDataNodeAbstract[K, V](-1, defaultKey, defaultVal), MaxLevel)
	return &SkipListAbstract[K, V, T]{
		head:     head,
		nodeSize: 0,
	}
}
