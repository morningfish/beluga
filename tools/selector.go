package tools

import (
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func GetLabelSelector(filters map[string]string) (*labels.Selector, error) {
	labelSelector := labels.NewSelector()
	for key, value := range filters {
		requirement, err := labels.NewRequirement(key, selection.DoubleEquals, []string{value})
		if err != nil {
			return nil, err
		}
		labelSelector = labelSelector.Add(*requirement)
	}
	return &labelSelector, nil
}

func GetFieldSelector(filters map[string]string) (*fields.Selector, error) {
	var sets = make(fields.Set)
	for k, v := range filters {
		sets[k] = v
	}
	selector := fields.SelectorFromSet(sets)
	return &selector, nil
}
