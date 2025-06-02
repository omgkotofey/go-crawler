package crawler

import (
	"sync"
	"time"
)

type Task struct {
	URL      string
	Depth    int64
	Timeout  time.Duration
	Cooldown time.Duration
}

type Inbox struct {
	inbox chan Task
	seen  map[string]string
	mutex *sync.RWMutex
}

func NewInbox() *Inbox {
	return &Inbox{
		inbox: make(chan Task),
		seen:  make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

func (m *Inbox) Exists(url string) bool {
	m.mutex.Lock()
	_, found := m.seen[url]
	m.mutex.Unlock()

	return found
}

func (m *Inbox) Add(url string, depth int64, timeout time.Duration) {
	m.mutex.Lock()
	m.seen[url] = url
	m.mutex.Unlock()
	m.inbox <- Task{URL: url, Depth: depth, Timeout: timeout}
}

func (m *Inbox) Next() (Task, bool) {
	url, ok := <-m.inbox

	return url, ok
}

func (m *Inbox) Close() {
	close(m.inbox)
}
