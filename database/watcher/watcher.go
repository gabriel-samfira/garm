package watcher

import (
	"context"
	"sync"

	"github.com/cloudbase/garm/database/common"
)

var databaseWatcher common.Watcher

func InitWatcher(ctx context.Context) {
	if databaseWatcher != nil {
		return
	}
	w := &watcher{
		producers: make(map[string]*producer),
		consumers: make(map[string]*consumer),
		quit:      make(chan struct{}),
		ctx:       ctx,
	}

	go w.loop()
	databaseWatcher = w
}

func RegisterProducer(ID string) (common.Producer, error) {
	if databaseWatcher == nil {
		return nil, common.ErrWatcherNotInitialized
	}
	return databaseWatcher.RegisterProducer(ID)
}

func RegisterConsumer(ID string, filters ...common.PayloadFilterFunc) (common.Consumer, error) {
	if databaseWatcher == nil {
		return nil, common.ErrWatcherNotInitialized
	}
	return databaseWatcher.RegisterConsumer(ID, filters...)
}

type watcher struct {
	producers map[string]*producer
	consumers map[string]*consumer

	mux    sync.Mutex
	closed bool
	quit   chan struct{}
	ctx    context.Context
}

func (w *watcher) RegisterProducer(ID string) (common.Producer, error) {
	if _, ok := w.producers[ID]; ok {
		return nil, common.ErrProducerAlreadyRegistered
	}
	p := &producer{
		id:       ID,
		messages: make(chan common.ChangePayload, 1),
		quit:     make(chan struct{}),
	}
	w.producers[ID] = p
	go w.serviceProducer(p)
	return p, nil
}

func (w *watcher) serviceProducer(prod *producer) {
	defer func() {
		w.mux.Lock()
		defer w.mux.Unlock()
		prod.Close()
		delete(w.producers, prod.id)
	}()
	for {
		select {
		case <-w.quit:
			return
		case <-w.ctx.Done():
			return
		case payload := <-prod.messages:
			for _, c := range w.consumers {
				go c.Send(payload)
			}
		}
	}
}

func (w *watcher) RegisterConsumer(ID string, filters ...common.PayloadFilterFunc) (common.Consumer, error) {
	if _, ok := w.consumers[ID]; ok {
		return nil, common.ErrConsumerAlreadyRegistered
	}
	c := &consumer{
		messages: make(chan common.ChangePayload, 1),
		filters:  filters,
		quit:     make(chan struct{}),
		id:       ID,
	}
	w.consumers[ID] = c
	go w.serviceConsumer(c)
	return c, nil
}

func (w *watcher) serviceConsumer(consumer *consumer) {
	defer func() {
		w.mux.Lock()
		defer w.mux.Unlock()
		consumer.Close()
		delete(w.consumers, consumer.id)
	}()
	for {
		select {
		case <-consumer.quit:
			return
		case <-w.quit:
			return
		case <-w.ctx.Done():
			return
		}
	}

}

func (w *watcher) Close() {
	w.mux.Lock()
	defer w.mux.Unlock()
	if w.closed {
		return
	}

	close(w.quit)
	w.closed = true

	for _, p := range w.producers {
		p.Close()
	}
	w.producers = nil

	for _, c := range w.consumers {
		c.Close()
	}
	w.consumers = nil
}

func (w *watcher) loop() {
	defer func() {
		w.Close()
	}()
	for {
		select {
		case <-w.quit:
			return
		case <-w.ctx.Done():
			return
		}
	}
}
