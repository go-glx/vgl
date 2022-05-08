package vulkan

func (f *vkPhysicalDeviceFamilies) uniqueIDs() []uint32 {
	uniqueFamilies := make(map[uint32]struct{})

	if f.supportGraphics {
		uniqueFamilies[f.graphicsFamilyId] = struct{}{}
	}

	if f.supportPresent {
		uniqueFamilies[f.presentFamilyId] = struct{}{}
	}

	ids := make([]uint32, 0, len(uniqueFamilies))
	for id := range uniqueFamilies {
		ids = append(ids, id)
	}

	return ids
}
