package schema

import (
	"fmt"
	"sort"
)

type PrioritizedFactoryLoader struct {
	loaders prioritizedLoaders
}

type prioritizedLoaders []prioritizedLoader

func (pl prioritizedLoaders) Len() int {
	return len(pl)
}

func (pl prioritizedLoaders) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

func (pl prioritizedLoaders) Less(i, j int) bool {
	return pl[i].priority < pl[j].priority
}

type prioritizedLoader struct {
	priority int
	loader   Loader
}

var _ Loader = &PrioritizedFactoryLoader{}

func (f *PrioritizedFactoryLoader) Add(loader Loader, priority int) {
	f.loaders = append(f.loaders, prioritizedLoader{
		priority: priority,
		loader:   loader,
	})

	sort.Sort(f.loaders)
}

func (f *PrioritizedFactoryLoader) IsSupported(uri string) bool {
	for _, loader := range f.loaders {
		if loader.loader.IsSupported(uri) {
			return true
		}
	}

	return false
}

func (f *PrioritizedFactoryLoader) Load(uri string) ([]byte, error) {
	for _, loader := range f.loaders {
		if loader.loader.IsSupported(uri) {
			return loader.loader.Load(uri)
		}
	}

	return nil, fmt.Errorf("no supported loader found: %s", uri)
}
