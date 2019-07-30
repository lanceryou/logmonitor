package logmonitor

import (
	"github.com/lanceryou/logmonitor/watcher"
	"io/ioutil"
	"os"
	"sync"
)

type Monitor struct {
	watcher watcher.Watcher

	mu             sync.RWMutex
	monitorFiles   map[string]*MonitorFile
	monitoredNames map[string][]string
}

func (m *Monitor) Observe(path string, opts ...MonitorObserveOption) (err error) {
	var opt MonitorObserveOptions
	for _, o := range opts {
		o(&opt)
	}

	if err = opt.apply(); err != nil {
		return
	}

	m.mu.RLock()
	_, ok := m.monitoredNames[path]
	m.mu.RUnlock()
	if ok {
		return
	}

	m.mu.Lock()
	m.monitoredNames[path] = []string{}
	m.mu.Unlock()
	return m.watcher.AddWatch(path, m.eventHandle(path, &opt))
}

func (m *Monitor) eventHandle(path string, opt *MonitorObserveOptions) func(event watcher.Event) error {
	return func(event watcher.Event) (err error) {
		for _, h := range opt.intercepts {
			if err = h(event); err != nil {
				return
			}
		}

		m.mu.RLock()
		file, ok := m.monitorFiles[event.Path]
		m.mu.RUnlock()
		if !ok {
			m.mu.Lock()
			file = NewMonitorFile(event.Path, opt.senders, WithFileFilters(opt.filters...))
			m.setMonitorName(path, event.Path)
			m.monitorFiles[event.Path] = file
			m.mu.Unlock()
		}

		return file.Handle(event)
	}
}

func (m *Monitor) UnObserve(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, name := range m.monitoredNames[path] {
		m.monitorFiles[name].Close()
		delete(m.monitorFiles, name)
	}

	delete(m.monitoredNames, path)
}

func isDir(path string) bool {
	fi, e := os.Stat(path)
	if e != nil {
		return false
	}

	return fi.IsDir()
}

func getFiles(path string) (files []string, err error) {
	if !isDir(path) {
		return []string{path}, nil
	}

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	sep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		files = append(files, path+sep+fi.Name())
	}
	return
}

func (m *Monitor) setMonitorName(path, name string) {
	for _, n := range m.monitoredNames[path] {
		if n == name {
			return
		}
	}

	m.monitoredNames[path] = append(m.monitoredNames[path], name)
}

func NewMonitor() *Monitor {
	w, err := watcher.NewLogWatcher()
	if err != nil {
		panic(err)
	}
	return &Monitor{
		watcher:        w,
		monitorFiles:   make(map[string]*MonitorFile),
		monitoredNames: make(map[string][]string),
	}
}
