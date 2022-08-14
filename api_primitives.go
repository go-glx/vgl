package vgl

import (
	"math"

	"github.com/go-glx/vgl/glm"
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
// 3) figure is buildIn shader type (point, line, triangle, circle, rect, polygon, texture)
// 4) params called exactly as method, but "Params" prefix instead of "Draw"
// 5) all params struct default golang values should be some valid value (and good defaults)
// -----------------------------------------------------------------------------

// Params2dPoint is input for Draw2dPoint
type Params2dPoint struct {
	Pos       glm.Local2D // pixel position
	Color     glm.Color   // pixel color
	NoCulling bool        // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dPoint will draw single point on current surface with current blend mode
// slow draw call, should be used only for editor/debug draw/gizmos, etc...
func (r *Render) Draw2dPoint(p *Params2dPoint) {
	pos := r.toLocalSpace2d(p.Pos)

	if !p.NoCulling && !r.cullingPoint(pos) {
		return
	}

	r.api.DrawPoint(pos, p.Color.VecRGBA())
}

// -----------------------------------------------------------------------------

// Params2dLine is input for Draw2dLine
type Params2dLine struct {
	Pos              [2]glm.Local2D // vertex positions
	Color            glm.Color      // line color
	ColorGradient    [2]glm.Color   // color for each vertex
	ColorUseGradient bool           // will use ColorGradient instead of Color
	Width            int32          // default=1px; max=32px; line width (1px is only guaranteed to fast GPU render).
	NoCulling        bool           // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dLine will draw line on current surface with current blend mode
func (r *Render) Draw2dLine(p *Params2dLine) {
	if p.Width < 1 {
		p.Width = 1
	}
	if p.Width > 32 {
		p.Width = 32
	}

	pos := [2]glm.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
	}

	color := [2]glm.Vec4{}
	if p.ColorUseGradient {
		color[0] = p.ColorGradient[0].VecRGBA()
		color[1] = p.ColorGradient[1].VecRGBA()
	} else {
		color[0] = p.Color.VecRGBA()
		color[1] = color[0]
	}

	if p.Width == 1 {
		// native GPU line (faster that emulating with rect)
		if !p.NoCulling && !r.cullingLine(pos) {
			return
		}

		r.api.DrawLine(pos, color)
		return
	}

	// not all GPU support of lines with width 1px+
	// so, in case of custom width, we will emulate it with rect
	radTo := pos[0].AngleTo(pos[1])
	offset := r.toLocalAspectRation(p.Width) / 2
	topLeft := pos[0].PolarOffset(offset, radTo+(math.Pi/2))
	bottomLeft := pos[0].PolarOffset(offset, radTo-(math.Pi/2))
	topRight := pos[1].PolarOffset(offset, radTo+(math.Pi/2))
	bottomRight := pos[1].PolarOffset(offset, radTo-(math.Pi/2))

	rectPos := [4]glm.Vec2{topLeft, topRight, bottomRight, bottomLeft}
	if !p.NoCulling && !r.cullingRect(rectPos) {
		return
	}

	r.api.DrawRect(rectPos, [4]glm.Vec4{color[0], color[1], color[1], color[0]}, true)
}

// -----------------------------------------------------------------------------

// Params2dTriangle is input for Draw2dTriangle
type Params2dTriangle struct {
	Pos              [3]glm.Local2D // vertex positions in clock-wise order
	Color            glm.Color      // color for all vertexes
	ColorGradient    [3]glm.Color   // color for each vertex
	ColorUseGradient bool           // will use ColorGradient instead of Color
	Filled           bool           // fill triangle with color/gradient
	NoCulling        bool           // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dTriangle will draw triangle on current surface with current blend mode
// Params2dTriangle.Pos must be in clock-wise order
func (r *Render) Draw2dTriangle(p *Params2dTriangle) {
	pos := [3]glm.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
		r.toLocalSpace2d(p.Pos[2]),
	}

	if !p.NoCulling && !r.cullingTriangle(pos) {
		return
	}

	color := [3]glm.Vec4{}
	if p.ColorUseGradient {
		color[0] = p.ColorGradient[0].VecRGBA()
		color[1] = p.ColorGradient[1].VecRGBA()
		color[2] = p.ColorGradient[2].VecRGBA()
	} else {
		color[0] = p.Color.VecRGBA()
		color[1] = color[0]
		color[2] = color[0]
	}

	r.api.DrawTriangle(pos, color, p.Filled)
}

// -----------------------------------------------------------------------------

// Params2dRect is input for Draw2dRect
type Params2dRect struct {
	Pos              [4]glm.Local2D // vertex positions in clock-wise order
	Color            glm.Color      // color for all vertexes
	ColorGradient    [4]glm.Color   // color for each vertex
	ColorUseGradient bool           // will use ColorGradient instead of Color
	Filled           bool           // fill rect with color/gradient
	NoCulling        bool           // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dRect will draw rect on current surface with current blend mode
// order of Params2dRect.Pos must be specified in order:
//   1) top-left
//   2) top-right
//   3) bottom-right
//   4) bottom-left
func (r *Render) Draw2dRect(p *Params2dRect) {
	pos := [4]glm.Vec2{
		r.toLocalSpace2d(p.Pos[0]),
		r.toLocalSpace2d(p.Pos[1]),
		r.toLocalSpace2d(p.Pos[2]),
		r.toLocalSpace2d(p.Pos[3]),
	}

	if !p.NoCulling && !r.cullingRect(pos) {
		return
	}

	color := [4]glm.Vec4{}
	if p.ColorUseGradient {
		color[0] = p.ColorGradient[0].VecRGBA()
		color[1] = p.ColorGradient[1].VecRGBA()
		color[2] = p.ColorGradient[2].VecRGBA()
		color[3] = p.ColorGradient[3].VecRGBA()
	} else {
		color[0] = p.Color.VecRGBA()
		color[1] = color[0]
		color[2] = color[0]
		color[3] = color[0]
	}

	r.api.DrawRect(pos, color, p.Filled)
}

// -----------------------------------------------------------------------------

// Params2dCircle is input for Draw2dCircle
type Params2dCircle struct {
	PosCenter        glm.Local2D    // pos: v1: circle center pos
	PosCenterRadius  int32          // pos: v1: radius of circle in px
	PosArea          [4]glm.Local2D // pos: v2: circle bounding box (allow drawing ellipses)
	PosUseArea       bool           // will use PosArea instead of PosCenter+PosCenterRadius
	HoleRadius       float32        // value [0 .. 1]. 0=without hole, 0.1=90% circle is visible, 1=invisible circle
	Smooth           float32        // value [-1, 0 .. 1]. -1=no smooth, 0.005=default, 1=full blur (default value will be used, if no value (0) specified)
	Color            glm.Color      // color for circle body/border
	ColorGradient    [4]glm.Color   // color for circle part (tl, tr, br, bl)
	ColorUseGradient bool           // will use ColorGradient instead of Color
	NoCulling        bool           // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dCircle will draw circle on current surface with current blend mode
func (r *Render) Draw2dCircle(p *Params2dCircle) {
	if p.HoleRadius >= 0.9999 {
		return
	}

	pos := [4]glm.Vec2{}
	if p.PosUseArea {
		pos[0] = r.toLocalSpace2d(p.PosArea[0])
		pos[1] = r.toLocalSpace2d(p.PosArea[1])
		pos[2] = r.toLocalSpace2d(p.PosArea[2])
		pos[3] = r.toLocalSpace2d(p.PosArea[3])
	} else {
		radiusX := r.toLocalAspectRationX(p.PosCenterRadius)
		radiusY := r.toLocalAspectRationY(p.PosCenterRadius)
		pos[0] = r.toLocalSpace2d(p.PosCenter).Add(-radiusX, -radiusY) // tl
		pos[1] = r.toLocalSpace2d(p.PosCenter).Add(radiusX, -radiusY)  // tr
		pos[2] = r.toLocalSpace2d(p.PosCenter).Add(radiusX, radiusY)   // br
		pos[3] = r.toLocalSpace2d(p.PosCenter).Add(-radiusX, radiusY)  // bl
	}

	if !p.NoCulling && !r.cullingRect(pos) {
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

	color := [4]glm.Vec4{}
	if p.ColorUseGradient {
		color[0] = p.ColorGradient[0].VecRGBA()
		color[1] = p.ColorGradient[1].VecRGBA()
		color[2] = p.ColorGradient[2].VecRGBA()
		color[3] = p.ColorGradient[3].VecRGBA()
	} else {
		color[0] = p.Color.VecRGBA()
		color[1] = color[0]
		color[2] = color[0]
		color[3] = color[0]
	}

	r.api.DrawCircle(
		pos,
		color,
		glm.Vec1{X: 1.0 - glm.Clamp(p.HoleRadius, 0, 1)},
		glm.Vec1{X: glm.Clamp(p.Smooth, 0, 1)},
	)
}

// -----------------------------------------------------------------------------

// Params2dPolygon is input for Draw2dPolygon
type Params2dPolygon struct {
	Pos       []glm.Local2D // vertex positions
	Color     glm.Color     // color for polygon body/border
	Filled    bool          // fill polygon with color?
	NoCulling bool          // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dPolygon will draw multi-lined polygon
// you can specify any number of edges, but for good performance
// not recommended using polygons with more than 32 edges.
//
// With vertex count < 5:
//   - this method just proxy call to another render APIs:
//     Draw2dPoint, Draw2dLine, Draw2dTriangle, Draw2dRect..
//   - With zero vertexes input, this will do nothing
//   - Params2dPolygon.NoCulling is used
//
// With vertex count >= 5:
//   - Params2dPolygon.NoCulling is ignored, it`s your responsibility
//     to not call this API, when all polygon vertexes outside of visible screen space
func (r *Render) Draw2dPolygon(p *Params2dPolygon) {
	switch len(p.Pos) {
	case 0:
		return
	case 1:
		r.Draw2dPoint(&Params2dPoint{
			Pos:       p.Pos[0],
			Color:     p.Color,
			NoCulling: p.NoCulling,
		})
		return
	case 2:
		r.Draw2dLine(&Params2dLine{
			Pos:       [2]glm.Local2D{p.Pos[0], p.Pos[1]},
			Color:     p.Color,
			NoCulling: p.NoCulling,
		})
		return
	case 3:
		r.Draw2dTriangle(&Params2dTriangle{
			Pos:       [3]glm.Local2D{p.Pos[0], p.Pos[1], p.Pos[2]},
			Color:     p.Color,
			Filled:    p.Filled,
			NoCulling: p.NoCulling,
		})
		return
	case 4:
		r.Draw2dRect(&Params2dRect{
			Pos:       [4]glm.Local2D{p.Pos[0], p.Pos[1], p.Pos[2], p.Pos[3]},
			Color:     p.Color,
			Filled:    p.Filled,
			NoCulling: p.NoCulling,
		})
		return
	default:
		// todo: draw polygon
	}
}
