package vgl

import (
	"math"

	"github.com/vulkan-go/vulkan"

	"github.com/go-glx/glx"
)

// Api design
// -----------------------------------------------------------------------------
//  Draw2dCircle( p *Params2dCircle )
//  ^^^^  ^^^^^^
//   |  ^^   |
// const |   figure
//       2d/3d
//
// 1) all method name MUST start with Draw
// 2) next is 2d/3d API logic splitter.
//    - 2d - current API
//    - 3d - reserved for feature
// 3) figure is buildIn shader type (point, line, triangle, circle, rect, texture)
// 4) params called exactly as method, but "Params" prefix instead of "Draw"
// 5) all params struct default golang values should be some valid value (and good defaults)
// -----------------------------------------------------------------------------

// Params2dPoint is input for Draw2dPoint
type Params2dPoint struct {
	Pos       glx.Vec2  // pixel position from top,left corner of surface
	Color     glx.Color // pixel color
	NoCulling bool      // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dPoint will draw single point on current surface with current blend mode
// slow draw call, should be used only for editor/debug draw/gizmos, etc...
func (r *Render) Draw2dPoint(p *Params2dPoint) {
	localPos := r.toLocalSpace2d(p.Pos)

	if !p.NoCulling && !r.cullingPoint(localPos) {
		return
	}

	r.api.Draw(buildInShaderPoint, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeFill,
		vertexes: []shaderInputUniversal2dVertex{
			{
				pos:   localPos,
				color: p.Color.VecRGBA(),
			},
		},
	})
}

// -----------------------------------------------------------------------------

// Params2dLine is input for Draw2dLine
type Params2dLine struct {
	Pos              [2]glx.Vec2  // pixel positions from top,left corner of surface
	Color            glx.Color    // line color
	ColorGradient    [2]glx.Color // color for each vertex
	ColorUseGradient bool         // will use ColorGradient instead of Color
	Width            float32      // default=1px; max=32px; line width (1px is only guaranteed to fast GPU render).
	NoCulling        bool         // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dLine will draw line on current surface with current blend mode
func (r *Render) Draw2dLine(p *Params2dLine) {
	if p.Width < 1 {
		p.Width = 1
	}
	if p.Width > 32 {
		p.Width = 32
	}

	localPos := [2]glx.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
	}

	localColor := [2]glx.Vec4{}
	if p.ColorUseGradient {
		localColor[0] = p.ColorGradient[0].VecRGBA()
		localColor[1] = p.ColorGradient[1].VecRGBA()
	} else {
		localColor[0] = p.Color.VecRGBA()
		localColor[1] = localColor[0]
	}

	if p.Width == 1 {
		// native GPU line (faster that emulating with rect)
		if !p.NoCulling && !r.cullingLine(localPos) {
			return
		}

		r.api.Draw(buildInShaderLine, &shaderInputUniversal2d{
			mode: vulkan.PolygonModeLine,
			vertexes: []shaderInputUniversal2dVertex{
				{
					pos:   localPos[0],
					color: localColor[0],
				},
				{
					pos:   localPos[1],
					color: localColor[1],
				},
			},
		})
		return
	}

	// not all GPU support of lines with width 1px+
	// so, in case of custom width, we will emulate it with rect
	angle := localPos[0].AngleTo(localPos[1])
	offset := r.toLocalAspectRation(p.Width) / 2
	topLeft := localPos[0].PolarOffset(offset, angle+(math.Pi/2))
	bottomLeft := localPos[0].PolarOffset(offset, angle-(math.Pi/2))
	topRight := localPos[1].PolarOffset(offset, angle+(math.Pi/2))
	bottomRight := localPos[1].PolarOffset(offset, angle-(math.Pi/2))

	rectPos := [4]glx.Vec2{topLeft, topRight, bottomRight, bottomLeft}
	if !p.NoCulling && !r.cullingRect(rectPos) {
		return
	}

	r.api.Draw(buildInShaderTriangle, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeFill,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: rectPos[0], color: localColor[0]}, // tl
			{pos: rectPos[1], color: localColor[1]}, // tr
			{pos: rectPos[2], color: localColor[1]}, // br
		},
	})
	r.api.Draw(buildInShaderTriangle, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeFill,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: rectPos[2], color: localColor[1]}, // br
			{pos: rectPos[3], color: localColor[0]}, // bl
			{pos: rectPos[0], color: localColor[0]}, // tl
		},
	})
}

// -----------------------------------------------------------------------------

// Params2dTriangle is input for Draw2dTriangle
type Params2dTriangle struct {
	Pos              [3]glx.Vec2  // pixel position from top,left corner of surface in clock-wise order
	Color            glx.Color    // color for all vertexes
	ColorGradient    [3]glx.Color // color for each vertex
	ColorUseGradient bool         // will use ColorGradient instead of Color
	Filled           bool         // fill triangle with color/gradient
	NoCulling        bool         // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dTriangle will draw triangle on current surface with current blend mode
// Params2dTriangle.Pos must be in clock-wise order
func (r *Render) Draw2dTriangle(p *Params2dTriangle) {
	localPos := [3]glx.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
		r.toLocalSpace2d(p.Pos[2]),
	}

	if !p.NoCulling && !r.cullingTriangle(localPos) {
		return
	}

	localColor := [3]glx.Vec4{}
	if p.ColorUseGradient {
		localColor[0] = p.ColorGradient[0].VecRGBA()
		localColor[1] = p.ColorGradient[1].VecRGBA()
		localColor[2] = p.ColorGradient[2].VecRGBA()
	} else {
		localColor[0] = p.Color.VecRGBA()
		localColor[1] = localColor[0]
		localColor[2] = localColor[0]
	}

	var mode vulkan.PolygonMode
	if p.Filled {
		mode = vulkan.PolygonModeFill
	} else {
		mode = vulkan.PolygonModeLine
	}

	r.api.Draw(buildInShaderTriangle, &shaderInputUniversal2d{
		mode: mode,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: localPos[0], color: localColor[0]},
			{pos: localPos[1], color: localColor[1]},
			{pos: localPos[2], color: localColor[2]},
		},
	})
}

// -----------------------------------------------------------------------------

// Params2dRect is input for Draw2dRect
type Params2dRect struct {
	Pos              [4]glx.Vec2  // pixel position from top,left corner of surface in clock-wise order
	Color            glx.Color    // color for all vertexes
	ColorGradient    [4]glx.Color // color for each vertex
	ColorUseGradient bool         // will use ColorGradient instead of Color
	Filled           bool         // fill rect with color/gradient
	NoCulling        bool         // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dRect will draw rect on current surface with current blend mode
// order of Params2dRect.Pos must be specified in order:
//   1) top-left
//   2) top-right
//   3) bottom-right
//   4) bottom-left
func (r *Render) Draw2dRect(p *Params2dRect) {
	localPos := [4]glx.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
		r.toLocalSpace2d(p.Pos[2]),
		r.toLocalSpace2d(p.Pos[3]),
	}

	if !p.NoCulling && !r.cullingRect(localPos) {
		return
	}

	localColor := [4]glx.Vec4{}
	if p.ColorUseGradient {
		localColor[0] = p.ColorGradient[0].VecRGBA()
		localColor[1] = p.ColorGradient[1].VecRGBA()
		localColor[2] = p.ColorGradient[2].VecRGBA()
		localColor[3] = p.ColorGradient[3].VecRGBA()
	} else {
		localColor[0] = p.Color.VecRGBA()
		localColor[1] = localColor[0]
		localColor[2] = localColor[0]
		localColor[3] = localColor[0]
	}

	if !p.Filled {
		r.api.Draw(buildInShaderRect, &shaderInputUniversal2d{
			mode: vulkan.PolygonModeLine,
			vertexes: []shaderInputUniversal2dVertex{
				{pos: localPos[0], color: localColor[0]},
				{pos: localPos[1], color: localColor[1]},
				{pos: localPos[2], color: localColor[2]},
				{pos: localPos[3], color: localColor[3]},
			},
		})

		return
	}

	// drawing two triangles faster, that outlined rect
	// with custom polygon mode
	// when rect is filled, two triangles visually looks same as rect

	r.api.Draw(buildInShaderTriangle, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeFill,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: localPos[0], color: localColor[0]}, // tl
			{pos: localPos[1], color: localColor[1]}, // tr
			{pos: localPos[2], color: localColor[2]}, // br
		},
	})
	r.api.Draw(buildInShaderTriangle, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeFill,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: localPos[2], color: localColor[2]}, // br
			{pos: localPos[3], color: localColor[3]}, // bl
			{pos: localPos[0], color: localColor[0]}, // tl
		},
	})
}

// -----------------------------------------------------------------------------

// Params2dCircle is input for Draw2dCircle
type Params2dCircle struct {
	Center           glx.Vec2     // position in pixels from top,left corner of surface
	Radius           float32      // radius in pixels
	HoleRadius       float32      // value [0 .. 1]. 0=without hole, 0.1=90% circle is visible, 1=invisible circle
	Smooth           float32      // value [-1, 0 .. 1]. -1=no smooth, 0.005=default, 1=full blur (default value will be used, if no value (0) specified)
	Color            glx.Color    // color for circle body/border
	ColorGradient    [4]glx.Color // color for circle part (tl, tr, br, bl)
	ColorUseGradient bool         // will use ColorGradient instead of Color
	NoCulling        bool         // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dCircle will draw circle on current surface with current blend mode
func (r *Render) Draw2dCircle(p *Params2dCircle) {
	if p.HoleRadius >= 0.9999 {
		return
	}

	localRadius := r.toLocalAspectRation(p.Radius)
	localPos := [4]glx.Vec2{
		r.toLocalSpace2d(p.Center.Add(glx.Vec2{X: -p.Radius, Y: -p.Radius})), // tl
		r.toLocalSpace2d(p.Center.Add(glx.Vec2{X: +p.Radius, Y: -p.Radius})), // tr
		r.toLocalSpace2d(p.Center.Add(glx.Vec2{X: +p.Radius, Y: +p.Radius})), // br
		r.toLocalSpace2d(p.Center.Add(glx.Vec2{X: -p.Radius, Y: +p.Radius})), // bl
	}

	if !p.NoCulling && !r.cullingRect(localPos) {
		return
	}

	if p.Smooth == 0 {
		// default value (if no specified)
		p.Smooth = 0.005
	}

	if p.Smooth == -1 {
		// specified to turn off
		p.Smooth = 0
	}

	localColor := [4]glx.Vec4{}
	if p.ColorUseGradient {
		localColor[0] = p.ColorGradient[0].VecRGBA()
		localColor[1] = p.ColorGradient[1].VecRGBA()
		localColor[2] = p.ColorGradient[2].VecRGBA()
		localColor[3] = p.ColorGradient[3].VecRGBA()
	} else {
		localColor[0] = p.Color.VecRGBA()
		localColor[1] = localColor[0]
		localColor[2] = localColor[0]
		localColor[3] = localColor[0]
	}

	r.api.Draw(buildInShaderCircle, &shaderInputCircle2d{
		vertexes: []shaderInputCircle2dVertex{
			{pos: localPos[0], color: localColor[0]},
			{pos: localPos[1], color: localColor[1]},
			{pos: localPos[2], color: localColor[2]},
			{pos: localPos[3], color: localColor[3]},
		},
		center: glx.Vec2{
			X: (localPos[0].X + localPos[1].X + localPos[2].X + localPos[3].X) / 4,
			Y: (localPos[0].Y + localPos[1].Y + localPos[2].Y + localPos[3].Y) / 4,
		},
		radius:    glx.Vec1{X: localRadius},
		thickness: glx.Vec1{X: glx.Clamp(p.HoleRadius, 0, 1)},
		smooth:    glx.Vec1{X: glx.Clamp(p.Smooth, 0, 1)},
	})

	// todo: remove
	r.api.Draw(buildInShaderRect, &shaderInputUniversal2d{
		mode: vulkan.PolygonModeLine,
		vertexes: []shaderInputUniversal2dVertex{
			{pos: localPos[0], color: localColor[0]},
			{pos: localPos[1], color: localColor[1]},
			{pos: localPos[2], color: localColor[2]},
			{pos: localPos[3], color: localColor[3]},
		},
	})
}
