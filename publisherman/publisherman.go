package publisherman

import (
	"sync"
	"fmt"
	"sync/atomic"
	"errors"
)

type roomInfo struct {
	NumClients *int64
}

type publisherMan struct {
	publishersPath map[string] roomInfo
	mutex sync.RWMutex
}

var instance *publisherMan = nil
var once sync.Once

func GetInstance() *publisherMan {
	once.Do(func() {
		instance = &publisherMan{
			publishersPath : make(map[string] roomInfo),
		}
	})
	return instance
}

func (pubman *publisherMan) Add(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	fmt.Printf("Adding %s to publisher map\n", publishPath)
	pubman.publishersPath[publishPath] = roomInfo{NumClients : new(int64)}
}

func (pubman *publisherMan) Remove(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	fmt.Printf("Deleteing %s from publisher map\n", publishPath)
	delete(pubman.publishersPath, publishPath)
}

func (pubman *publisherMan) IncrementViewerCount(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	_, ok := pubman.publishersPath[publishPath]
	if ok {
		atomic.AddInt64(pubman.publishersPath[publishPath].NumClients, 1)
	}
}

func (pubman *publisherMan) DeccrementViewerCount(publishPath string) {
	pubman.mutex.Lock()
	defer pubman.mutex.Unlock()

	_, ok := pubman.publishersPath[publishPath]
	if ok {
		atomic.AddInt64(pubman.publishersPath[publishPath].NumClients, -1)
	}
}

func (pubman *publisherMan) GetRoomInfo(room string) (roomInfo, error) {
	pubman.mutex.RLock()
	defer pubman.mutex.RUnlock()

	if info, ok := pubman.publishersPath[room]; !ok {
		return roomInfo{}, errors.New("Room doesn't exist'")
	} else {
		return info, nil
	}
}

func (pubman *publisherMan) GetRoomsInfo() map[string] roomInfo {

	pubman.mutex.RLock()
	defer pubman.mutex.RUnlock()

	info_copy := pubman.publishersPath
	return info_copy
}
