package main

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

// 测试1：空跳表测试
func TestSkipList_Empty(t *testing.T) {
	sl := NewSkipList()

	// 校验结构
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("空跳表校验失败: %v", err)
	}
	t.Log("✅ 空跳表测试通过")
}

// 测试2：5个节点（你原来的样例）
func TestSkipList_5Nodes(t *testing.T) {
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

	// 核心校验
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("5节点跳表校验失败: %v", err)
	}
	t.Log("✅ 5节点跳表测试通过")
}

// 测试3：20个节点（大规模测试）
func TestSkipList_20Nodes(t *testing.T) {
	sl := NewSkipList()

	data := []SkipListDataNode{
		NewSkipListDataNode(1, 1, 10),
		NewSkipListDataNode(1, 2, 20),
		NewSkipListDataNode(1, 3, 30),
		NewSkipListDataNode(1, 4, 40),
		NewSkipListDataNode(1, 5, 50),

		NewSkipListDataNode(2, 1, 60),
		NewSkipListDataNode(2, 2, 70),
		NewSkipListDataNode(2, 3, 80),
		NewSkipListDataNode(2, 4, 90),
		NewSkipListDataNode(2, 5, 100),

		NewSkipListDataNode(3, 1, 110),
		NewSkipListDataNode(3, 2, 120),
		NewSkipListDataNode(3, 3, 130),
		NewSkipListDataNode(3, 4, 140),
		NewSkipListDataNode(3, 5, 150),

		NewSkipListDataNode(4, 1, 160),
		NewSkipListDataNode(4, 2, 170),
		NewSkipListDataNode(4, 3, 180),
		NewSkipListDataNode(4, 4, 190),
		NewSkipListDataNode(4, 5, 200),
	}

	for _, v := range data {
		sl.Insert(&v)
	}

	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("20节点跳表校验失败: %v", err)
	}
	t.Log("✅ 20节点大规模跳表测试通过")
}

// TestSkipList_1000RandomNodes 随机1000个无序节点 压力测试
func TestSkipList_10000RandomNodes(t *testing.T) {
	// 初始化随机种子

	sl := NewSkipList()

	// 用 map 保证 key1+key2 不重复
	keySet := make(map[[2]int]bool)
	var data []SkipListDataNode

	// 生成 1000 个唯一随机节点
	for len(data) < 10000 {
		key1 := rand.Intn(10000)
		key2 := rand.Intn(10000)
		key := [2]int{key1, key2}

		if !keySet[key] {
			keySet[key] = true
			val := rand.Intn(100000)
			data = append(data, NewSkipListDataNode(key1, key2, val))
		}
	}

	// 无序插入跳表（完全乱序，最能测试跳表稳定性）
	for _, v := range data {
		sl.Insert(&v) // 用你正常的 Insert 方法，不是指定层级
	}

	// 核心校验：结构 + 排序 + nextSize
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("10000随机节点校验失败: %v", err)
	}

	t.Logf("✅ 10000 个随机无序节点测试通过！总节点数：%d", sl.nodeSize)
}

// TestSkipList_DeleteBasic 基础删除测试：删头、删中、删尾
func TestSkipList_DeleteBasic(t *testing.T) {
	sl := NewSkipList()

	data := []SkipListDataNode{
		NewSkipListDataNode(1, 2, 100),
		NewSkipListDataNode(3, 1, 200),
		NewSkipListDataNode(2, 5, 300),
		NewSkipListDataNode(1, 1, 400),
		NewSkipListDataNode(3, 2, 500),
	}

	// 插入
	for _, v := range data {
		sl.Insert(&v)
	}

	// 1. 删除头节点(1,1)
	sl.Delete(&data[3])
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("删除头节点失败: %v", err)
	}

	// 2. 删除中间节点(2,5)
	sl.Delete(&data[2])
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("删除中间节点失败: %v", err)
	}

	// 3. 删除尾节点(3,2)
	sl.Delete(&data[4])
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("删除尾节点失败: %v", err)
	}

	t.Log("✅ 基础删除（头/中/尾）全部通过")
}

// TestSkipList_DeleteRandomNodes 随机删除 500 个节点（压力测试）
func TestSkipList_DeleteRandomNodes(t *testing.T) {
	sl := NewSkipList()

	// 生成 1000 个唯一随机节点
	keySet := make(map[[2]int]bool)
	var data []SkipListDataNode

	for len(data) < 2000 {
		key1 := rand.Intn(10000)
		key2 := rand.Intn(10000)
		key := [2]int{key1, key2}
		if !keySet[key] {
			keySet[key] = true
			val := rand.Intn(100000)
			data = append(data, NewSkipListDataNode(key1, key2, val))
		}
	}

	// 插入
	for _, v := range data {
		sl.Insert(&v)
	}

	// 随机删除 500 个节点
	for i := 0; i < 500; i++ {
		idx := rand.Intn(len(data))
		sl.Delete(&data[idx])
		// 每删一次都校验（最严格）
		if err := sl.CheckSkipListFull(); err != nil {
			t.Fatalf("第 %d 次删除后出错: %v", i+1, err)
		}
	}

	// 最终校验
	if err := sl.CheckSkipListFull(); err != nil {
		t.Fatalf("最终校验失败: %v", err)
	}

	t.Logf("✅ 1000 节点随机删除 500 个成功！剩余节点数：%d", sl.nodeSize)
}

// TestSkipList_DeleteAll 删光所有节点测试
func TestSkipList_DeleteAll(t *testing.T) {
	sl := NewSkipList()

	data := []SkipListDataNode{
		NewSkipListDataNode(1, 1, 10),
		NewSkipListDataNode(1, 2, 20),
		NewSkipListDataNode(2, 3, 30),
		NewSkipListDataNode(3, 4, 40),
	}

	for _, v := range data {
		sl.Insert(&v)
	}

	// 全部删除
	for _, v := range data {
		sl.Delete(&v)
		if err := sl.CheckSkipListFull(); err != nil {
			t.Fatalf("删除过程中出错: %v", err)
		}
	}

	// 最终节点数必须为 0
	if sl.nodeSize != 0 {
		t.Fatalf("删光后节点数应为 0，实际为 %d", sl.nodeSize)
	}

	t.Log("✅ 全部删除测试通过")
}

// TestSkipList_LowerBound_Key1Zero 仅测试 key1=0 的下界查询
func TestSkipList_FindTopK_Key1Zero(t *testing.T) {
	sl := NewSkipList()

	// 所有数据 key1=0
	data := []SkipListDataNode{
		NewSkipListDataNode(0, 10, 100),
		NewSkipListDataNode(0, 30, 200),
		NewSkipListDataNode(0, 20, 300),
		NewSkipListDataNode(0, 50, 400),
		NewSkipListDataNode(0, 40, 500),
	}

	for _, v := range data {
		sl.Insert(&v)
	}

	// 期望顺序（key2 升序）
	expectKey2 := []int{10, 20, 30, 40, 50}

	// 测试 Top1~Top5
	for k := 1; k <= 5; k++ {
		node, ok := sl.FindTopK(k)
		if !ok {
			t.Fatalf("Top%d 查找失败", k)
		}
		if node.key1 != 0 {
			t.Fatalf("Top%d key1 必须=0，实际=%d", k, node.key1)
		}
		if node.key2 != expectKey2[k-1] {
			t.Fatalf("Top%d key2 期望=%d 实际=%d",
				k, expectKey2[k-1], node.key2)
		}
	}

	// 越界
	if _, ok := sl.FindTopK(0); ok {
		t.Fatal("FindTopK(0) 应该失败")
	}
	if _, ok := sl.FindTopK(6); ok {
		t.Fatal("FindTopK(6) 应该失败")
	}

	t.Log("✅ FindTopK key1=0 测试通过")
}

// TestSkipList_FindTopK_LargeScale 大规模测试 FindTopK (1w 节点, key1=0)
func TestSkipList_FindTopK_LargeScale(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	sl := NewSkipList()
	const count = 10000

	var key2List []int
	keySet := make(map[int]bool)
	for len(key2List) < count {
		k2 := rand.Intn(1000000)
		if !keySet[k2] {
			keySet[k2] = true
			key2List = append(key2List, k2)
		}
	}

	// ✅ 修复：取地址 &
	for _, k2 := range key2List {
		node := NewSkipListDataNode(0, k2, rand.Intn(100000))
		sl.Insert(&node)
	}

	sort.Ints(key2List)

	testCnt := 200
	for i := 0; i < testCnt; i++ {
		var k int
		switch {
		case i == 0:
			k = 1
		case i == 1:
			k = count
		case i == 2:
			k = count / 2
		default:
			k = 1 + rand.Intn(count)
		}

		node, ok := sl.FindTopK(k)
		if !ok {
			t.Fatalf("FindTopK(%d) 失败", k)
		}
		if node.key1 != 0 {
			t.Fatal("key1 必须等于 0")
		}
		if node.key2 != key2List[k-1] {
			t.Fatalf("FindTopK(%d) 错误: 期望 key2=%d 实际=%d", k, key2List[k-1], node.key2)
		}
	}

	t.Log("✅ FindTopK 10000 节点大规模测试通过")
}

// TestSkipList_LowerBound_LargeScale 大规模测试 LowerBound (1w 节点, key1=0)
func TestSkipList_LowerBound_LargeScale(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	sl := NewSkipList()
	const count = 10000

	var key2List []int
	keySet := make(map[int]bool)
	for len(key2List) < count {
		k2 := rand.Intn(1000000)
		if !keySet[k2] {
			keySet[k2] = true
			key2List = append(key2List, k2)
		}
	}

	// ✅ 修复：取地址 &
	for _, k2 := range key2List {
		node := NewSkipListDataNode(0, k2, rand.Intn(100000))
		sl.Insert(&node)
	}

	sort.Ints(key2List)

	testCnt := 500
	for i := 0; i < testCnt; i++ {
		targetKey2 := rand.Intn(1000000 + 200000)
		target := NewSkipListDataNode(0, targetKey2, 0)

		node, rank, hasNode := sl.LowerBound(&target)

		expRank := 0
		expKey2 := 0
		for j := len(key2List) - 1; j >= 0; j-- {
			if key2List[j] <= targetKey2 {
				expRank = j + 1
				expKey2 = key2List[j]
				break
			}
		}

		expHasNode := expRank != 0

		if hasNode != expHasNode {
			t.Fatalf("targetKey2=%d hasNode 错误", targetKey2)
		}
		if rank != expRank {
			t.Fatalf("targetKey2=%d rank 错误: 期望=%d 实际=%d", targetKey2, expRank, rank)
		}
		if hasNode {
			if node.key1 != 0 {
				t.Fatal("key1 必须等于 0")
			}
			if node.key2 != expKey2 {
				t.Fatalf("targetKey2=%d node 错误: 期望=%d 实际=%d", targetKey2, expKey2, node.key2)
			}
		}
	}

	t.Log("✅ LowerBound 10000 节点大规模测试通过")
}
