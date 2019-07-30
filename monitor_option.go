package logmonitor

import (
	"fmt"
	"github.com/lanceryou/logmonitor/filter"
	"github.com/lanceryou/logmonitor/sender"
	"github.com/lanceryou/logmonitor/watcher"
)

type MonitorObserveOptions struct {
	senders    []sender.MonitorSender
	filters    []filter.Filter
	intercepts []watcher.EventHandler
}

func (o MonitorObserveOptions) apply() error {
	if len(o.senders) == 0 {
		return fmt.Errorf("no monitor senders")
	}
	return nil
}

type MonitorObserveOption func(*MonitorObserveOptions)

func WithMonitorSenders(senders ...sender.MonitorSender) MonitorObserveOption {
	return func(o *MonitorObserveOptions) {
		o.senders = senders
	}
}

func WithFilters(fs ...filter.Filter) MonitorObserveOption {
	return func(o *MonitorObserveOptions) {
		o.filters = fs
	}
}

func WithIntercepts(intercepts ...watcher.EventHandler) MonitorObserveOption {
	return func(o *MonitorObserveOptions) {
		o.intercepts = intercepts
	}
}
