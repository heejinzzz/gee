package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的一部分
	children []*node // 子结点
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true

	handlersChain []HandlerFunc // 该路由的回调函数
}

// 绑定回调函数
func (n *node) bindHandler(handlers []HandlerFunc) {
	n.handlersChain = handlers
}

// 第一个匹配成功的结点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild || part[0] == ':' || part[0] == '*' {
			return child
		}
	}
	return nil
}

// 所有匹配成功的结点，用于查找
func (n *node) matchChildren(part string) []*node {
	var nodes []*node
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 插入新的路由结点
func (n *node) insert(pattern string, parts []string, height int, handlers []HandlerFunc) {
	if len(parts) == height {
		if n.pattern != "" && n.pattern != pattern {
			panic(fmt.Sprintf("Route Conflict: \"%s\" and \"%s\"", pattern, n.pattern)) // 路由冲突
		}
		n.pattern = pattern
		n.bindHandler(handlers)
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1, handlers)
}

// 查找路由结点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
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
