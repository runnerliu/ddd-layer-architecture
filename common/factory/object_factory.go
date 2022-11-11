package factory

import "sync"

// ObjectFactory 对象工厂
type ObjectFactory struct {
	keyToObject      map[string]interface{}
	keyToDestroyFunc map[string]func(interface{}) error
	// mu 读写锁
	mu sync.RWMutex
}

// Get 获取对象
func (o *ObjectFactory) Get(key string, createFunc func(string) (interface{}, error),
	destroyFunc func(interface{}) error) (interface{}, error) {
	// 如果已创建过对象，那么直接返回
	o.mu.RLock()
	obj, found := o.keyToObject[key]
	o.mu.RUnlock()
	if found {
		return obj, nil
	}

	// 加写锁
	o.mu.Lock()
	defer o.mu.Unlock()

	// 双重检查，防止重复创建
	if obj, found := o.keyToObject[key]; found {
		return obj, nil
	}

	// 创建实例
	obj, err := createFunc(key)
	if err != nil {
		return nil, err
	}

	// 缓存实例
	o.keyToObject[key] = obj
	if destroyFunc != nil {
		o.keyToDestroyFunc[key] = destroyFunc
	}

	return obj, nil
}

// Destroy 销毁工厂及其中缓存的对象
func (o *ObjectFactory) Destroy() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.keyToObject == nil {
		return
	}

	for key, obj := range o.keyToObject {
		if obj == nil {
			continue
		}
		if destroyFunc, found := o.keyToDestroyFunc[key]; found {
			_ = destroyFunc(obj)
		}
	}
	o.keyToObject = nil
	o.keyToDestroyFunc = nil
}

var (
	defaultObjectFactoryOnce sync.Once
	defaultObjectFactory     *ObjectFactory
)

// NewObjectFactory 获取默认的 ObjectFactory 单例
func NewObjectFactory() *ObjectFactory {
	defaultObjectFactoryOnce.Do(func() {
		defaultObjectFactory = &ObjectFactory{
			keyToObject:      make(map[string]interface{}),
			keyToDestroyFunc: make(map[string]func(interface{}) error)}
	})
	return defaultObjectFactory
}
