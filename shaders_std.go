package vgl

const (
	buildInShaderPoint    = "buildIn.point"
	buildInShaderLine     = "buildIn.line"
	buildInShaderTriangle = "buildIn.triangle"
	buildInShaderCircle   = "buildIn.circle"
	buildInShaderRect     = "buildIn.rect"
)

var stdShaders = []ParamsRegisterShader{
	stdShaderPoint,
	stdShaderLine,
	stdShaderTriangle,
	stdShaderCircle,
	stdShaderRect,
}
