package form

import (
	"errors"
	"sort"
)

type PrioritizedFieldFactory struct {
	factories prioritizedFieldFactories
}

type prioritizedFieldFactories []prioritizedFieldFactory

func (pl prioritizedFieldFactories) Len() int {
	return len(pl)
}

func (pl prioritizedFieldFactories) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl prioritizedFieldFactories) Less(i, j int) bool {
	return pl[i].priority < pl[j].priority
}

type prioritizedFieldFactory struct {
	priority int
	factory  FieldFactory
}

var _ FieldFactory = &PrioritizedFieldFactory{}

func (f *PrioritizedFieldFactory) Add(factory FieldFactory, priority int) {
	f.factories = append(f.factories, prioritizedFieldFactory{
		priority: priority,
		factory:  factory,
	})

	sort.Sort(f.factories)
}

func (f *PrioritizedFieldFactory) IsSupported(uri string) bool {
	for _, factory := range f.factories {
		if factory.factory.IsSupported(uri) {
			return true
		}
	}

	return false
}

func (f *PrioritizedFieldFactory) Create(uri, path string, options FieldOptions) (Field, error) {
	for _, factory := range f.factories {
		if factory.factory.IsSupported(uri) {
			return factory.factory.Create(uri, path, options)
		}
	}

	return nil, errors.New("no supported factory found")
}
