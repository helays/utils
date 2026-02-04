package tree

// Radix tree implementation below is a based on the original work by
// Armon Dadgar in https://github.com/armon/go-radix/blob/master/radix.go
// (MIT licensed). It's been heavily modified for use as a HTTP routing tree.

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"helay.net/go/utils/v3/tools"
)

type methodTyp uint

const (
	mSTUB methodTyp = 1 << iota
	mCONNECT
	mDELETE
	mGET
	mHEAD
	mOPTIONS
	mPATCH
	mPOST
	mPUT
	mTRACE
)

var mALL = mCONNECT | mDELETE | mGET | mHEAD |
	mOPTIONS | mPATCH | mPOST | mPUT | mTRACE

var methodMap = map[string]methodTyp{
	http.MethodConnect: mCONNECT,
	http.MethodDelete:  mDELETE,
	http.MethodGet:     mGET,
	http.MethodHead:    mHEAD,
	http.MethodOptions: mOPTIONS,
	http.MethodPatch:   mPATCH,
	http.MethodPost:    mPOST,
	http.MethodPut:     mPUT,
	http.MethodTrace:   mTRACE,
}

var reverseMethodMap = map[methodTyp]string{
	mCONNECT: http.MethodConnect,
	mDELETE:  http.MethodDelete,
	mGET:     http.MethodGet,
	mHEAD:    http.MethodHead,
	mOPTIONS: http.MethodOptions,
	mPATCH:   http.MethodPatch,
	mPOST:    http.MethodPost,
	mPUT:     http.MethodPut,
	mTRACE:   http.MethodTrace,
}

// RegisterMethod adds support for custom HTTP method handlers, available
// via Router#Method and Router#MethodFunc
func RegisterMethod(method string) error {
	if method == "" {
		return nil
	}
	method = strings.ToUpper(method)
	if _, ok := methodMap[method]; ok {
		return nil
	}
	n := len(methodMap)
	if n > strconv.IntSize-2 {
		return fmt.Errorf("达到最大方法数量限制 (%d)", strconv.IntSize)
	}
	mt := methodTyp(2 << n)
	methodMap[method] = mt
	mALL |= mt
	return nil
}

type nodeTyp uint8

const (
	ntStatic   nodeTyp = iota // /home
	ntRegexp                  // /{id:[0-9]+}
	ntParam                   // /{user}
	ntCatchAll                // /api/v1/*
)

type RouteTreeNode[T comparable] struct {
	rex       *regexp.Regexp // 正则表达式节点的匹配器
	endpoints Endpoints[T]   // 叶子节点上的 HTTP 处理端点
	prefix    string         // 前缀，表示我们忽略的公共前缀部分
	// 子节点应按顺序存储以便迭代，
	// 按节点类型分组存储
	children [ntCatchAll + 1]nodes[T]
	tail     byte    // 子节点前缀的第一个字节
	typ      nodeTyp // 节点类型：静态、正则表达式、参数、通配
	label    byte    // 前缀的第一个字节
}

// Endpoints is a mapping of http method constants to handlers
// for a given route.
type Endpoints[T comparable] map[methodTyp]*Endpoint[T]

type Endpoint[T comparable] struct {
	handler   T        // 端点处理器
	pattern   string   // 路由模式，表示处理器节点的路由模式
	paramKeys []string // 在处理器节点上记录的参数键
}

func (s Endpoints[T]) Value(method methodTyp) *Endpoint[T] {
	mh, ok := s[method]
	if !ok {
		mh = &Endpoint[T]{}
		s[method] = mh
	}
	return mh
}

func (n *RouteTreeNode[T]) InsertRoute(method methodTyp, pattern string, handler T) (*RouteTreeNode[T], error) {
	var parent *RouteTreeNode[T]
	search := pattern

	for {
		// Handle key exhaustion
		if len(search) == 0 {
			// Insert or update the node's leaf handler
			if err := n.setEndpoint(method, handler, pattern); err != nil {
				return nil, err
			}
			return n, nil
		}

		// We're going to be searching for a wild node next,
		// in this case, we need to get the tail
		var label = search[0]
		var segTail byte
		var segEndIdx int
		var segTyp nodeTyp
		var segRexpat string
		if label == '{' || label == '*' {
			var err error
			segTyp, _, segRexpat, segTail, _, segEndIdx, err = patNextSegment(search)
			if err != nil {
				return nil, err
			}
		}

		var prefix string
		if segTyp == ntRegexp {
			prefix = segRexpat
		}

		// Look for the edge to attach to
		parent = n
		n = n.getEdge(segTyp, label, segTail, prefix)

		// No edge, create one
		if n == nil {
			child := &RouteTreeNode[T]{label: label, tail: segTail, prefix: search}
			hn, err := parent.addChild(child, search)
			if err != nil {
				return nil, err
			}
			if err = hn.setEndpoint(method, handler, pattern); err != nil {
				return nil, err
			}

			return hn, nil
		}

		// Found an edge to match the pattern

		if n.typ > ntStatic {
			// We found a param node, trim the param from the search path and continue.
			// This param/wild pattern segment would already be on the tree from a previous
			// call to addChild when creating a new node.
			search = search[segEndIdx:]
			continue
		}

		// Static nodes fall below here.
		// Determine longest prefix of the search key on match.
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			// the common prefix is as long as the current node's prefix we're attempting to insert.
			// keep the search going.
			search = search[commonPrefix:]
			continue
		}

		// Split the node
		child := &RouteTreeNode[T]{
			typ:    ntStatic,
			prefix: search[:commonPrefix],
		}
		if err := parent.replaceChild(search[0], segTail, child); err != nil {
			return nil, err
		}

		// Restore the existing node
		n.label = n.prefix[commonPrefix]
		n.prefix = n.prefix[commonPrefix:]

		if _, err := child.addChild(n, n.prefix); err != nil {
			return nil, err
		}

		// If the new key is a subset, set the method/handler on this node and finish.
		search = search[commonPrefix:]
		if len(search) == 0 {
			if err := child.setEndpoint(method, handler, pattern); err != nil {
				return nil, err
			}
			return child, nil
		}

		// Create a new edge for the node
		subchild := &RouteTreeNode[T]{
			typ:    ntStatic,
			label:  search[0],
			prefix: search,
		}
		hn, err := child.addChild(subchild, search)
		if err != nil {
			return nil, err
		}
		if err = hn.setEndpoint(method, handler, pattern); err != nil {
			return nil, err
		}
		return hn, nil
	}
}

// addChild appends the new `child` node to the tree using the `pattern` as the trie key.
// For a URL router like chi's, we split the static, param, regexp and wildcard segments
// into different nodes. In addition, addChild will recursively call itself until every
// pattern segment is added to the url pattern tree as individual nodes, depending on type.
func (n *RouteTreeNode[T]) addChild(child *RouteTreeNode[T], prefix string) (*RouteTreeNode[T], error) {
	search := prefix

	// handler leaf node added to the tree is the child.
	// this may be overridden later down the flow
	hn := child
	// Parse next segment
	segTyp, _, segRexpat, segTail, segStartIdx, segEndIdx, err := patNextSegment(search)
	if err != nil {
		return nil, err
	}

	// Add child depending on next up segment
	switch segTyp {

	case ntStatic:
		// Search prefix is all static (that is, has no params in path)
		// noop

	default:
		// Search prefix contains a param, regexp or wildcard

		if segTyp == ntRegexp {
			rex, err := regexp.Compile(segRexpat)
			if err != nil {
				return nil, fmt.Errorf("路由参数中的正则表达式模式 '%s' 无效", segRexpat)
			}
			child.prefix = segRexpat
			child.rex = rex
		}

		if segStartIdx == 0 {
			// Route starts with a param
			child.typ = segTyp

			if segTyp == ntCatchAll {
				segStartIdx = -1
			} else {
				segStartIdx = segEndIdx
			}
			if segStartIdx < 0 {
				segStartIdx = len(search)
			}
			child.tail = segTail // for params, we set the tail

			if segStartIdx != len(search) {
				// add static edge for the remaining part, split the end.
				// its not possible to have adjacent param nodes, so its certainly
				// going to be a static node next.

				search = search[segStartIdx:] // advance search position

				nn := &RouteTreeNode[T]{
					typ:    ntStatic,
					label:  search[0],
					prefix: search,
				}
				hn, err = child.addChild(nn, search)
				if err != nil {
					return nil, err
				}
			}

		} else if segStartIdx > 0 {
			// Route has some param

			// starts with a static segment
			child.typ = ntStatic
			child.prefix = search[:segStartIdx]
			child.rex = nil

			// add the param edge node
			search = search[segStartIdx:]

			nn := &RouteTreeNode[T]{
				typ:   segTyp,
				label: search[0],
				tail:  segTail,
			}
			hn, err = child.addChild(nn, search)
			if err != nil {
				return nil, err
			}

		}
	}
	n.children[child.typ] = append(n.children[child.typ], child)
	n.children[child.typ].Sort()
	return hn, nil
}

func (n *RouteTreeNode[T]) replaceChild(label, tail byte, child *RouteTreeNode[T]) error {
	for i := 0; i < len(n.children[child.typ]); i++ {
		if n.children[child.typ][i].label == label && n.children[child.typ][i].tail == tail {
			n.children[child.typ][i] = child
			n.children[child.typ][i].label = label
			n.children[child.typ][i].tail = tail
			return nil
		}
	}
	return fmt.Errorf("替换缺失的子节点")
}

func (n *RouteTreeNode[T]) getEdge(ntyp nodeTyp, label, tail byte, prefix string) *RouteTreeNode[T] {
	nds := n.children[ntyp]
	for i := 0; i < len(nds); i++ {
		if nds[i].label == label && nds[i].tail == tail {
			if ntyp == ntRegexp && nds[i].prefix != prefix {
				continue
			}
			return nds[i]
		}
	}
	return nil
}

func (n *RouteTreeNode[T]) setEndpoint(method methodTyp, handler T, pattern string) error {
	// Set the handler for the method type on the node
	if n.endpoints == nil {
		n.endpoints = make(Endpoints[T])
	}

	paramKeys, err := patParamKeys(pattern)
	if err != nil {
		return err
	}

	if method&mSTUB == mSTUB {
		n.endpoints.Value(mSTUB).handler = handler
	}
	if method&mALL == mALL {
		h := n.endpoints.Value(mALL)
		h.handler = handler
		h.pattern = pattern
		h.paramKeys = paramKeys
		for _, m := range methodMap {
			h := n.endpoints.Value(m)
			h.handler = handler
			h.pattern = pattern
			h.paramKeys = paramKeys
		}
	} else {
		h := n.endpoints.Value(method)
		h.handler = handler
		h.pattern = pattern
		h.paramKeys = paramKeys
	}
	return nil
}

func (n *RouteTreeNode[T]) FindRoute(rctx *Context, method methodTyp, path string) (*RouteTreeNode[T], Endpoints[T], T) {
	// Reset the context routing pattern and params
	rctx.routePattern = ""
	rctx.routeParams.Keys = rctx.routeParams.Keys[:0]
	rctx.routeParams.Values = rctx.routeParams.Values[:0]
	var zero T
	// Find the routing handlers for the path
	rn := n.findRoute(rctx, method, path)
	if rn == nil {
		return nil, nil, zero
	}

	// Record the routing params in the request lifecycle
	rctx.URLParams.Keys = append(rctx.URLParams.Keys, rctx.routeParams.Keys...)
	rctx.URLParams.Values = append(rctx.URLParams.Values, rctx.routeParams.Values...)

	// Record the routing pattern in the request lifecycle
	if rn.endpoints[method].pattern != "" {
		rctx.routePattern = rn.endpoints[method].pattern
		rctx.RoutePatterns = append(rctx.RoutePatterns, rctx.routePattern)
	}

	return rn, rn.endpoints, rn.endpoints[method].handler
}

// Recursive edge traversal by checking all nodeTyp groups along the way.
// It's like searching through a multi-dimensional radix trie.
func (n *RouteTreeNode[T]) findRoute(rctx *Context, method methodTyp, path string) *RouteTreeNode[T] {
	nn := n
	search := path

	for t, nds := range nn.children {
		ntyp := nodeTyp(t)
		if len(nds) == 0 {
			continue
		}

		var xn *RouteTreeNode[T]
		xsearch := search

		var label byte
		if search != "" {
			label = search[0]
		}

		switch ntyp {
		case ntStatic:
			xn = nds.findEdge(label)
			if xn == nil || !strings.HasPrefix(xsearch, xn.prefix) {
				continue
			}
			xsearch = xsearch[len(xn.prefix):]

		case ntParam, ntRegexp:
			// short-circuit and return no matching route for empty param values
			if xsearch == "" {
				continue
			}

			// serially loop through each node grouped by the tail delimiter
			for idx := 0; idx < len(nds); idx++ {
				xn = nds[idx]

				// label for param nodes is the delimiter byte
				p := strings.IndexByte(xsearch, xn.tail)

				if p < 0 {
					if xn.tail == '/' {
						p = len(xsearch)
					} else {
						continue
					}
				} else if ntyp == ntRegexp && p == 0 {
					continue
				}

				if ntyp == ntRegexp && xn.rex != nil {
					if !xn.rex.MatchString(xsearch[:p]) {
						continue
					}
				} else if strings.IndexByte(xsearch[:p], '/') != -1 {
					// avoid a match across path segments
					continue
				}

				prevlen := len(rctx.routeParams.Values)
				rctx.routeParams.Values = append(rctx.routeParams.Values, xsearch[:p])
				xsearch = xsearch[p:]

				if len(xsearch) == 0 {
					if xn.isLeaf() {
						h := xn.endpoints[method]
						if h != nil && !tools.IsZero(h.handler) {
							rctx.routeParams.Keys = append(rctx.routeParams.Keys, h.paramKeys...)
							return xn
						}

						for endpoints := range xn.endpoints {
							if endpoints == mALL || endpoints == mSTUB {
								continue
							}
							rctx.methodsAllowed = append(rctx.methodsAllowed, endpoints)
						}

						// flag that the routing context found a route, but not a corresponding
						// supported method
						rctx.methodNotAllowed = true
					}
				}

				// recursively find the next node on this branch
				fin := xn.findRoute(rctx, method, xsearch)
				if fin != nil {
					return fin
				}

				// not found on this branch, reset vars
				rctx.routeParams.Values = rctx.routeParams.Values[:prevlen]
				xsearch = search
			}

			rctx.routeParams.Values = append(rctx.routeParams.Values, "")

		default:
			// catch-all nodes
			rctx.routeParams.Values = append(rctx.routeParams.Values, search)
			xn = nds[0]
			xsearch = ""
		}

		if xn == nil {
			continue
		}

		// did we find it yet?
		if len(xsearch) == 0 {
			if xn.isLeaf() {
				h := xn.endpoints[method]
				if h != nil && !tools.IsZero(h.handler) {
					rctx.routeParams.Keys = append(rctx.routeParams.Keys, h.paramKeys...)
					return xn
				}

				for endpoints := range xn.endpoints {
					if endpoints == mALL || endpoints == mSTUB {
						continue
					}
					rctx.methodsAllowed = append(rctx.methodsAllowed, endpoints)
				}

				// flag that the routing context found a route, but not a corresponding
				// supported method
				rctx.methodNotAllowed = true
			}
		}

		// recursively find the next node..
		fin := xn.findRoute(rctx, method, xsearch)
		if fin != nil {
			return fin
		}

		// Did not find final handler, let's remove the param here if it was set
		if xn.typ > ntStatic {
			if len(rctx.routeParams.Values) > 0 {
				rctx.routeParams.Values = rctx.routeParams.Values[:len(rctx.routeParams.Values)-1]
			}
		}

	}

	return nil
}

func (n *RouteTreeNode[T]) findEdge(ntyp nodeTyp, label byte) *RouteTreeNode[T] {
	nds := n.children[ntyp]
	num := len(nds)
	idx := 0

	switch ntyp {
	case ntStatic, ntParam, ntRegexp:
		i, j := 0, num-1
		for i <= j {
			idx = i + (j-i)/2
			if label > nds[idx].label {
				i = idx + 1
			} else if label < nds[idx].label {
				j = idx - 1
			} else {
				i = num // breaks cond
			}
		}
		if nds[idx].label != label {
			return nil
		}
		return (*RouteTreeNode[T])(nds[idx])

	default: // catch all
		return (*RouteTreeNode[T])(nds[idx])
	}
}

func (n *RouteTreeNode[T]) isLeaf() bool {
	return n.endpoints != nil
}

func (n *RouteTreeNode[T]) findPattern(pattern string) bool {
	nn := n
	for _, nds := range nn.children {
		if len(nds) == 0 {
			continue
		}

		n = nn.findEdge(nds[0].typ, pattern[0])
		if n == nil {
			continue
		}

		var idx int
		var xpattern string

		switch n.typ {
		case ntStatic:
			idx = longestPrefix(pattern, n.prefix)
			if idx < len(n.prefix) {
				continue
			}

		case ntParam, ntRegexp:
			idx = strings.IndexByte(pattern, '}') + 1

		case ntCatchAll:
			idx = longestPrefix(pattern, "*")

		default:
			return false
		}

		xpattern = pattern[idx:]
		if len(xpattern) == 0 {
			return true
		}

		return n.findPattern(xpattern)
	}
	return false
}

// patNextSegment returns the next segment details from a pattern:
// node type, param key, regexp string, param tail byte, param starting index, param ending index
func patNextSegment(pattern string) (nodeTyp, string, string, byte, int, int, error) {
	ps := strings.Index(pattern, "{")
	ws := strings.Index(pattern, "*")

	if ps < 0 && ws < 0 {
		return ntStatic, "", "", 0, 0, len(pattern), nil // we return the entire thing
	}

	// Sanity check
	if ps >= 0 && ws >= 0 && ws < ps {
		return 0, "", "", 0, 0, 0, fmt.Errorf("通配符 '*' 必须是路由中的最后一个模式，否则请使用 '{param}'")
	}

	var tail byte = '/' // Default endpoint tail to / byte

	if ps >= 0 {
		// Param/Regexp pattern is next
		nt := ntParam

		// Read to closing } taking into account opens and closes in curl count (cc)
		cc := 0
		pe := ps
		for i, c := range pattern[ps:] {
			if c == '{' {
				cc++
			} else if c == '}' {
				cc--
				if cc == 0 {
					pe = ps + i
					break
				}
			}
		}
		if pe == ps {
			return 0, "", "", 0, 0, 0, fmt.Errorf("路由参数的结束分隔符 '}' 缺失")
		}

		key := pattern[ps+1 : pe]
		pe++ // set end to next position

		if pe < len(pattern) {
			tail = pattern[pe]
		}

		key, rexpat, isRegexp := strings.Cut(key, ":")
		if isRegexp {
			nt = ntRegexp
		}

		if len(rexpat) > 0 {
			if rexpat[0] != '^' {
				rexpat = "^" + rexpat
			}
			if rexpat[len(rexpat)-1] != '$' {
				rexpat += "$"
			}
		}

		return nt, key, rexpat, tail, ps, pe, nil
	}

	// Wildcard pattern as finale
	if ws < len(pattern)-1 {
		return 0, "", "", 0, 0, 0, fmt.Errorf("通配符 '*' 必须是路由中的最后一个值。请删除尾随文本或改用 '{param}'")
	}
	return ntCatchAll, "*", "", 0, ws, len(pattern), nil
}

func patParamKeys(pattern string) ([]string, error) {
	pat := pattern
	paramKeys := []string{}
	for {
		ptyp, paramKey, _, _, _, e, err := patNextSegment(pat)
		if err != nil {
			return nil, err
		}
		if ptyp == ntStatic {
			return paramKeys, nil
		}
		for i := 0; i < len(paramKeys); i++ {
			if paramKeys[i] == paramKey {
				return nil, fmt.Errorf("chi: 路由模式 '%s' 包含重复的参数键 '%s'", pattern, paramKey)
			}
		}
		paramKeys = append(paramKeys, paramKey)
		pat = pat[e:]
	}
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	m := tools.Min(len(k1), len(k2))
	var i int
	for i = 0; i < m; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

type nodes[T comparable] []*RouteTreeNode[T]

// Sort the list of nodes by label
func (ns nodes[T]) Sort()              { sort.Sort(ns); ns.tailSort() }
func (ns nodes[T]) Len() int           { return len(ns) }
func (ns nodes[T]) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes[T]) Less(i, j int) bool { return ns[i].label < ns[j].label }

// tailSort pushes nodes with '/' as the tail to the end of the list for param nodes.
// The list order determines the traversal order.
func (ns nodes[T]) tailSort() {
	for i := len(ns) - 1; i >= 0; i-- {
		if ns[i].typ > ntStatic && ns[i].tail == '/' {
			ns.Swap(i, len(ns)-1)
			return
		}
	}
}

func (ns nodes[T]) findEdge(label byte) *RouteTreeNode[T] {
	num := len(ns)
	idx := 0
	i, j := 0, num-1
	for i <= j {
		idx = i + (j-i)/2
		if label > ns[idx].label {
			i = idx + 1
		} else if label < ns[idx].label {
			j = idx - 1
		} else {
			i = num // breaks cond
		}
	}
	if ns[idx].label != label {
		return nil
	}
	return ns[idx]
}
