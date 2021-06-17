package marker

import "errors"

type Collector struct {
	*Registry
}

func NewCollector(registry *Registry) *Collector {
	return &Collector{
		registry,
	}
}

func (collector *Collector) Collect(pkg *Package) error {

	if pkg == nil {
		return errors.New("pkg(package) cannot be nil")
	}

	return nil
}
