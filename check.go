package main

import (
	"errors"
	"fmt"
)

// 完美匹配你跳表的终极校验函数
func (sl *SkipList) CheckSkipListFull() error {
	if sl == nil {
		return errors.New("跳表不能为nil")
	}
	head := sl.head
	if head == nil {
		return errors.New("头节点不能为nil")
	}
	maxLevel := head.level

	// 统计0层总节点数
	realSize := 0
	cur := head
	for cur.next[0] != nil {
		realSize++
		cur = cur.next[0]
	}
	if realSize != sl.nodeSize {
		return fmt.Errorf("节点数错误：记录=%d, 实际=%d", sl.nodeSize, realSize)
	}

	// 逐层校验
	for level := 0; level <= maxLevel; level++ {
		visited := make(map[*SkipListNode]bool)
		curNode := head

		for {
			if visited[curNode] {
				return fmt.Errorf("level %d 发现环结构", level)
			}
			visited[curNode] = true

			// 节点自身合法性
			if curNode.level < 0 || len(curNode.next) != curNode.level+1 || len(curNode.nextSize) != curNode.level+1 {
				return fmt.Errorf("level %d 节点结构非法", level)
			}

			nextNode := curNode.next[level]
			realSpan := 0

			// ==============================
			// 【完全按你的规则计算】
			// ==============================
			tmp := curNode.next[0]
			for tmp != nil && tmp != nextNode {
				realSpan++
				tmp = tmp.next[0]
			}
			if tmp != nil {
				realSpan++
			}

			// 对比 nextSize
			storeSpan := curNode.nextSize[level]
			if realSpan != storeSpan {
				return fmt.Errorf("LEVEL %d NEXT SIZE 错误：(%d,%d) 期望=%d 实际=%d",
					level, curNode.data.key1, curNode.data.key2, realSpan, storeSpan)
			}

			// 排序校验
			if nextNode != nil && !isValidOrder(curNode.data, nextNode.data) {
				return fmt.Errorf("level %d 排序错误", level)
			}

			if nextNode == nil {
				break
			}
			curNode = nextNode
		}
	}

	return nil
}

func isValidOrder(a, b SkipListDataNode) bool {
	if a.key1 < b.key1 {
		return true
	}
	if a.key1 == b.key1 && a.key2 < b.key2 {
		return true
	}
	return false
}
