package physical

type (
	Families struct {
		GraphicsFamilyId uint32
		PresentFamilyId  uint32
		supportGraphics  bool
		supportPresent   bool
	}
)

func (f *Families) UniqueIDs() []uint32 {
	uniqueFamilies := make(map[uint32]any)

	if f.supportGraphics {
		uniqueFamilies[f.GraphicsFamilyId] = struct{}{}
	}

	if f.supportPresent {
		uniqueFamilies[f.PresentFamilyId] = struct{}{}
	}

	ids := make([]uint32, 0, len(uniqueFamilies))
	for id := range uniqueFamilies {
		ids = append(ids, id)
	}

	return ids
}
