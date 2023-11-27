package commands

import (
	"sort"
	"strings"
	"sync"
)

type Description struct {
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	Parameters  ParamsDesc `json:"parameters"`
}

var (
	makeDescriptionsOnce = sync.Once{}
	descriptionsList     []Description
	// descriptionsMap      map[string]Description
)

func makeDescriptionsList() {
	makeDescriptionsOnce.Do(func() {
		descriptionsList = make([]Description, 0, len(byName))
		for _, def := range byName {
			descriptionsList = append(descriptionsList, def.description)
		}
		sort.Slice(descriptionsList, func(i, j int) bool {
			return strings.Compare(descriptionsList[i].Name, descriptionsList[j].Name) <= 0
		})
	})
}

func List() []Description {
	makeDescriptionsList()
	return descriptionsList
}

func GetDescription(name string) Description {
	return byName[name].description
}

// func Map() map[string]Description {
// 	makeDescriptionsList()
// 	return descriptionsMap
// }
