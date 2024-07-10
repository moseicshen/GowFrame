package gow

import (
	"fmt"
	"strings"
)

type node struct {
	// 以此节点为结尾的完整路径
	pattern string
	// 此节点层对应的路径
	part string
	// 子节点
	children []*node
	// 是否为动态路由（part[0] = '*' 或 part[0] = ':'）
	isWild bool
	//中间件
	middleware []HandlerFunc
}

// matchChild 第一个找到的点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 所有找到的点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归结束
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 在此高度层的part
	part := parts[height]
	// 寻找part对应的child节点
	child := n.matchChild(part)
	// 未找到则创建
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	// 递归插入
	child.insert(pattern, parts, height+1)
}

// search 递归查找节点位置
func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		// 是否为完整路径？
		// 比如插入了 /p/go/doc, 查找 /p/go
		// 虽然能够匹配到go, 但是由于不是完整路径, 此时node.pattern为空, 无法匹配
		if n.pattern == "" {
			return nil
		} else {
			return n
		}
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}
