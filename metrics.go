package logContext

import (
	"fmt"
	"github.com/DataDog/datadog-go/statsd"
	"os"
	"sync"
	"time"
)

var MetricsNamespace = ""

func LogErrors() bool {
	log := os.Getenv("METRICS_LOG_ERRORS")
	return log == "" || log == "y" || log == "true" || log == "yes"
}

type Metrics struct {
	client    *statsd.Client
	tags      []string
	waitGroup *sync.WaitGroup

	lock *sync.RWMutex
}

func NewMetrics() (Metrics, error) {
	var metrics Metrics
	var lock sync.RWMutex
	metrics.lock = &lock
	var waitGroup sync.WaitGroup
	metrics.waitGroup = &waitGroup

	err := metrics.Connect()
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (m Metrics) Fork(tags ...string) Metrics {
	var lock sync.RWMutex
	var waitGroup sync.WaitGroup
	copy := Metrics{
		tags:      append(m.tags, tags...),
		lock:      &lock,
		waitGroup: &waitGroup,
	}

	return copy
}

func (m *Metrics) AppendTags(tag ...string) {
	m.tags = append(m.tags, tag...)
}

func GetDatadogDns() string {
	dns := os.Getenv("DATADOG_DNS")
	if dns == "" {
		dns = "metrics-datadog.default.svc.cluster.local"
	}

	return dns
}

func (m *Metrics) Connect() error {
	if m.lock == nil {
		var lock sync.RWMutex
		m.lock = &lock
	}

	if m.waitGroup == nil {
		var waitGroup sync.WaitGroup
		m.waitGroup = &waitGroup
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	dns := GetDatadogDns()
	if dns == "" {
		return fmt.Errorf("Datadog dns is empty")
	}

	addr := LookupIp(dns)
	client, err := statsd.New(addr + ":8125")
	if err != nil {
		return fmt.Errorf("new statsd: %s", err)
	}

	m.client = client
	m.client.Namespace = MetricsNamespace
	m.client.Tags = m.tags

	return nil
}

func (m *Metrics) Conn() (*statsd.Client, error) {
	client := m.Client()
	if client == nil {
		err := m.Connect()
		if err != nil {
			return nil, fmt.Errorf("retrieve connection: %s", err)
		}
	}

	return m.Client(), nil
}

func (m *Metrics) Client() *statsd.Client {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.client
}

/// Forwarding methods to the inner client.

func (m *Metrics) Gauge(ctx Context, name string, value float64, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("gauge: %s", err.Error()))
			}
			return
		}

		err = client.Gauge(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("gauge: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Count(ctx Context, name string, value int64, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("count: %s", err.Error()))
			}
			return
		}

		err = client.Count(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("count: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Histogram(ctx Context, name string, value float64, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("histogram: %s", err.Error()))
			}
			return
		}

		err = client.Histogram(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("histogram: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Distribution(ctx Context, name string, value float64, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("distribution: %s", err.Error()))
			}
			return
		}

		err = client.Distribution(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("distribution: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Incr(ctx Context, name string, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("incr: %s", err.Error()))
			}
			return
		}

		err = client.Incr(name, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("incr: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Decr(ctx Context, name string, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("decr: %s", err.Error()))
			}
			return
		}

		err = client.Decr(name, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("decr: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Set(ctx Context, name string, value string, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("set: %s", err.Error()))
			}
			return
		}

		err = client.Set(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("set: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Timing(ctx Context, name string, value time.Duration, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("timing: %s", err.Error()))
			}
			return
		}

		err = client.Timing(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("timing: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) TimeInMilliseconds(ctx Context, name string, value float64, tags []string, rate float64) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("timing milliseconds: %s", err.Error()))
			}
			return
		}

		err = client.TimeInMilliseconds(name, value, tags, rate)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("timing milliseconds: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Event(ctx Context, event *statsd.Event) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("event: %s", err.Error()))
			}
			return
		}

		err = client.Event(event)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("event: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) SimpleEvent(ctx Context, title, text string) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("simple event: %s", err.Error()))
			}
			return
		}

		err = client.SimpleEvent(title, text)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("simple event: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) ServiceCheck(ctx Context, check *statsd.ServiceCheck) {
	m.waitGroup.Add(1)
	go func() {
		defer m.waitGroup.Done()

		client, err := m.Conn()
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("service check: %s", err.Error()))
			}
			return
		}

		err = client.ServiceCheck(check)
		if err != nil {
			if !LogErrors() {
				ctx.Elog(fmt.Sprintf("service check: %s", err.Error()))
			}
			return
		}
	}()
}

func (m *Metrics) Wait() {
	m.waitGroup.Wait()
}

func (m *Metrics) Close() error {
	m.Wait()

	if m.client != nil {
		return m.client.Close()
	}

	return nil
}
