package logmonitor

import (
	"bufio"
	"github.com/lanceryou/logmonitor/filter"
	"github.com/lanceryou/logmonitor/sender"
	"github.com/lanceryou/logmonitor/watcher"
	"os"
)

type MonitorFile struct {
	offset  int64
	path    string
	file    *os.File
	senders []sender.MonitorSender
	opt     MonitorFileOptions
}

func (mf *MonitorFile) Handle(e watcher.Event) (err error) {
	if e.Op == watcher.Delete {
		return mf.Close()
	}

	if mf.file == nil {
		file, err := os.OpenFile(mf.path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}

		ret, err := file.Seek(mf.offset, 2)
		if err != nil {
			return err
		}
		mf.offset = ret
		mf.file = file

	}
	scan := bufio.NewScanner(mf.file)
	for scan.Scan() {
		for _, s := range mf.senders {
			text := filter.Filters(scan.Text(), mf.opt.filters)
			if text != "" {
				s.Send(text)
			}
		}
	}

	return
}

func (mf *MonitorFile) Close() error {
	if mf.file == nil {
		return nil
	}

	return mf.file.Close()
}

type MonitorFileOptions struct {
	filters []filter.Filter
}

type MonitorFileOption func(*MonitorFileOptions)

func WithFileFilters(fs ...filter.Filter) MonitorFileOption {
	return func(o *MonitorFileOptions) {
		o.filters = fs
	}
}

func NewMonitorFile(path string, senders []sender.MonitorSender, opts ...MonitorFileOption) (mf *MonitorFile) {
	var op MonitorFileOptions
	for _, o := range opts {
		o(&op)
	}

	mf = &MonitorFile{
		path:    path,
		senders: senders,
		opt:     op,
	}
	return
}
