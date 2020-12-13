package config_old

type videoExtensions struct {
	items []videoExtension
}

func NewVideoExtensions(baseItems []videoExtension) *videoExtensions {
	vs := videoExtensions{}
	if len(baseItems) > 0 {
		vs.items = baseItems
	}
	return &vs
}

func (vs *videoExtensions) addItem(newItem *videoExtension) {

	for i, item := range vs.items {
		if item.extension == newItem.extension {
			vs.items[i] = *newItem
			return
		}
	}
	vs.items = append(vs.items, *newItem)
}

func (vs *videoExtensions) Items() []videoExtension {
	return vs.items
}

//
func (vs *videoExtensions) GetNamesByExt(ext string) []string {

	for _, item := range vs.items {
		if item.extension == ext {
			return item.videoSettingName
		}
	}
	return nil
}
