package sources

import (
	"fmt"
	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/zeroturnaround/configo/flatmap"
)

type CompositeSource struct {
	Sources []map[string]interface{} `json:"sources"`
	UppercaseKeys bool
}

func (compositeSource *CompositeSource) Get() (map[string]interface{}, error) {

	resultEnv := make(map[string]interface{})

	_, err := From(compositeSource.Sources).
		Select(func(rawSource T) (T, error) {
			loader, err := GetSource(rawSource.(map[string]interface{}))

			if err != nil {
				return nil, fmt.Errorf("Failed to parse source: %s", err)
			}

			return loader, nil
		}).
		// Resolve in parallel because some sources might use IO and will take some time
		AsParallel().AsOrdered().
		Select(func(loader T) (T, error) {
			result, err := loader.(Source).Get()

			if err != nil {
				return nil, fmt.Errorf("Failed to resolve source: %s", err)
			}

			return result, nil
		}).
		AsSequential().
		CountBy(func(partialConfig T) (bool, error) {
			for key, value := range flatmap.Flatten(partialConfig.(map[string]interface{})) {
				resultEnv[key] = value
			}
			return true, nil
		})

	return resultEnv, err
}
