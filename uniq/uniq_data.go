package uniq

import (
	"github.com/zhanglp92/go-utils/xmath"
)

// Option ...
type Option struct {
	Cap    int                       // 容量大小
	Filter collectionNoderNodeFilter // 数据过滤放过
	Order  bool                      // 保证进入的顺序
}

// NewUniqCollectionInt ...
func NewUniqCollectionInt(o Option) *CollectionInt {
	return &CollectionInt{base: newUniqCollection(o)}
}

// NewUniqCollectionInt64 ...
func NewUniqCollectionInt64(o Option) *CollectionInt64 {
	return &CollectionInt64{base: newUniqCollection(o)}
}

// NewUniqCollectionFloat32 ...
func NewUniqCollectionFloat32(o Option) *CollectionFloat32 {
	return &CollectionFloat32{base: newUniqCollection(o)}
}

// NewUniqCollectionFloat64 ...
func NewUniqCollectionFloat64(o Option) *CollectionFloat64 {
	return &CollectionFloat64{base: newUniqCollection(o)}
}

// CollectionInt int集合去重 ------- int ------------------
type CollectionInt struct {
	base *collection
}

// Add ...
func (m *CollectionInt) Add(v int) {
	m.base.add(v)
}

// Collection ...
func (m *CollectionInt) Collection() []int {
	response := make([]int, 0, len(m.base.dup))
	m.base.xrange(func(v CollectionNoder) {
		response = append(response, v.(int))
	})
	return response
}

// CollectionInt64 int64集合去重 ------- int64 ------------------
type CollectionInt64 struct {
	base *collection
}

// Add ...
func (m *CollectionInt64) Add(v int64) {
	m.base.add(v)
}

// Collection ...
func (m *CollectionInt64) Collection() []int64 {
	response := make([]int64, 0, len(m.base.dup))
	m.base.xrange(func(v CollectionNoder) {
		response = append(response, v.(int64))
	})
	return response
}

// CollectionFloat32 int集合去重 ------- float32 ------------------
type CollectionFloat32 struct {
	base *collection
}

// Add ...
func (m *CollectionFloat32) Add(v float32) {
	m.base.add(v)
}

// Collection ...
func (m *CollectionFloat32) Collection() []float32 {
	response := make([]float32, 0, len(m.base.dup))
	m.base.xrange(func(v CollectionNoder) {
		response = append(response, v.(float32))
	})
	return response
}

// CollectionFloat64 int集合去重 ------- float64 ------------------
type CollectionFloat64 struct {
	base *collection
}

// Add ...
func (m *CollectionFloat64) Add(v float64) {
	m.base.add(v)
}

// Collection ...
func (m *CollectionFloat64) Collection() []float64 {
	response := make([]float64, 0, len(m.base.dup))
	m.base.xrange(func(v CollectionNoder) {
		response = append(response, v.(float64))
	})
	return response
}

// ------------------------------- base ---------------------------------------------

var collectionVaule = struct{}{}

// CollectionNoder 去重集合数据节点接口
type CollectionNoder interface{}

// collectionNoderNodeFilter 集合阶段过滤函数签名
type collectionNoderNodeFilter func(v interface{}) bool

// collection 去重集合
type collection struct {
	dup map[CollectionNoder]struct{}

	queue  []CollectionNoder // 保留顺序时启用
	filter collectionNoderNodeFilter

	add    func(CollectionNoder)
	xrange func(h func(CollectionNoder))
}

// newUniqCollection ...
func newUniqCollection(o Option) (m *collection) {
	defer func() {
		if o.Order {
			m.queue = make([]CollectionNoder, 0, o.Cap)
			m.add = m.addWithOrder
			m.xrange = m.xrangeWithOrder
		} else {
			m.add = m.addWithoutOrder
			m.xrange = m.xrangeWithoutOrder
		}
	}()

	xmath.MaximinInt(&o.Cap, 1, 0, xmath.MinBorderLine)
	return &collection{
		dup:    make(map[CollectionNoder]struct{}, o.Cap),
		filter: o.Filter,
	}
}

// Add 向集合添加元素
func (m *collection) addWithOrder(v CollectionNoder) {
	if m.filter != nil && m.filter(v) {
		return
	}
	if _, ok := m.dup[v]; ok {
		return
	}
	m.queue = append(m.queue, v)
	m.dup[v] = collectionVaule
}

func (m *collection) xrangeWithOrder(h func(CollectionNoder)) {
	for k := range m.queue {
		h(k)
	}
}

// Add 向集合添加元素
func (m *collection) addWithoutOrder(v CollectionNoder) {
	if m.filter == nil || !m.filter(v) {
		m.dup[v] = collectionVaule
	}
}

func (m *collection) xrangeWithoutOrder(h func(CollectionNoder)) {
	for k := range m.dup {
		h(k)
	}
}
