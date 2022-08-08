package vgl

import (
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
	Width            int            // default=1px; max=32px; line width (1px is only guaranteed to fast GPU render).
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

	// todo: draw
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
func (r *Render) Draw2dRect(p *Params2dRect) {
	// todo: draw
}

// -----------------------------------------------------------------------------

// Params2dCircle is input for Draw2dCircle
type Params2dCircle struct {
	Pos              glm.Local2D  // circle center pos
	Radius           float32      // radius of circle in px
	OutlineWidth     float32      // default=radius; width of circle border, when width < radius, circle will have hole in center
	Color            glm.Color    // color for circle body/border
	ColorGradient    [2]glm.Color // 0 = color of circle center, 1 = color of circle border
	ColorUseGradient bool         // will use ColorGradient instead of Color
	Filled           bool         // fill circle with color/gradient
	NoCulling        bool         // will send render command to GPU, even if all vertexes outside of visible screen
}

// Draw2dCircle will draw circle on current surface with current blend mode
func (r *Render) Draw2dCircle(p *Params2dCircle) {
	if p.Radius <= 0 {
		return
	}

	if p.OutlineWidth == 0 {
		p.OutlineWidth = p.Radius
	}

	// todo: draw
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
