package config_old

type videoSettings struct {
	items []VideoSetting
}

func NewVideoSettings(baseItems []VideoSetting) *videoSettings {

	vs := videoSettings{}

	if len(baseItems) > 0 {
		vs.items = baseItems
	}

	return &vs
}

// a a new videoSetting to the slice
func (vs *videoSettings) addItem(newItem *VideoSetting) {

	for i, item := range vs.items {
		if item.name == newItem.name {
			vs.items[i] = *newItem
			break
		} else {
			vs.items = append(vs.items, *newItem)
			break
		}
	}
}

// get videoSettings items
func (vs *videoSettings) Items() []VideoSetting {
	return vs.items
}

// get videoSettings items
func (vs *videoSettings) ItemByName(name string) *VideoSetting {
	for _, item := range vs.items {
		if item.name == name {
			return &item
		}
	}
	return nil
}
