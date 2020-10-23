package publisherman

import (
	"sync"
	"fmt"
)

type publisherMan struct {
	publishersPath map[string] struct{}
	mutex sync.RWMutex
}

var instance *publisherMan = nil
var once sync.Once

func GetInstance() *publisherMan {
	once.Do(func() {
		instance = &publisherMan{
			publishersPath : make(map[string] struct{}),
		}
	})
	return instance
}

func (pubman *publisherMan) Add(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	fmt.Printf("Adding %s to publisher map\n", publishPath)
	pubman.publishersPath[publishPath] = struct{}{}
}

func (pubman *publisherMan) Remove(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	fmt.Printf("Deleteing %s from publisher map\n", publishPath)
	delete(pubman.publishersPath, publishPath)
}

func (pubman *publisherMan) GetAllPublshers() []string {
	publishers := []string {}

	pubman.mutex.RLock()
	defer pubman.mutex.RUnlock()
	for path, _ := range pubman.publishersPath {
		publishers = append(publishers, path)
	}
	return publishers
}
