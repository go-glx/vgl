package shader

type Shader struct {
	meta       *Meta
	moduleVert *Module
	moduleFrag *Module
}

func (s *Shader) Meta() *Meta {
	return s.meta
}

func (s *Shader) ModuleVert() *Module {
	return s.moduleVert
}

func (s *Shader) ModuleFrag() *Module {
	return s.moduleFrag
}
