package vgl

import (
	"github.com/go-glx/vgl/glm"
)

// Api design
//  Draw2dCircleExt( vertexData[, transform, params], color, filled )
//  ^^^^  ^^^^^^
//   |  ^^   |  ^^^
// const | figure |
//      2d/3d    opt (full params)
//
// 1) all method name MUST start with Draw
// 2) next is 2d/3d API logic splitter.
//    - 2d - current API
//    - 3d - reserved for feature
// 3) figure is buildIn shader type

// Draw2dPoint will draw one pixel on canvas
// slow draw calls, should be used
// only for editor/debug draw/gizmos, etc...
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dPoint(pos glm.Local2D, color glm.Color) {
	// todo: draw
}

// Draw2dLine will draw one line on canvas
// vert is X,Y vec with Start and End position on screen
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dLine(pos [2]glm.Local2D, color [2]glm.Color) {
	// todo: draw
}

// Draw2dTriangle will draw one triangle on canvas
// vert elements MUST BE in clock-wise order
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dTriangle(pos [3]glm.Local2D, color glm.Color, filled bool) {
	r.Draw2dTriangleExt(pos, [3]glm.Color{color, color, color}, filled)
}

// Draw2dTriangleExt will draw one triangle on canvas
// vert elements MUST BE in clock-wise order
// you can specify color for each vertex, GPU will automatically blend colors
// in nice gradient.
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dTriangleExt(pos [3]glm.Local2D, color [3]glm.Color, filled bool) {
	v1, v2, v3 := r.toLocalSpace2d(pos[0]),
		r.toLocalSpace2d(pos[1]),
		r.toLocalSpace2d(pos[2])

	if !r.cullingTriangle([3]glm.Vec2{v1, v2, v3}) {
		return
	}

	r.api.DrawTriangle(
		[3]glm.Vec2{v1, v2, v3},
		[3]glm.Vec4{
			color[0].VecRGBA(),
			color[1].VecRGBA(),
			color[2].VecRGBA(),
		},
		filled,
	)
}

// Draw2dRect will draw one rect on canvas
// vert elements MUST BE in clock-wise order
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dRect(pos [4]glm.Local2D, color glm.Color) {
	// todo: draw
}

// Draw2dRectExt will draw one rect on canvas
// vert elements MUST BE in clock-wise order
// you can specify color for each vertex, GPU will automatically blend colors
// in nice gradient.
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dRectExt(pos [4]glm.Local2D, color [4]glm.Color, filled bool) {
	// todo: draw
}

// Draw2dCircle will draw one circle on canvas
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dCircle(pos glm.Local2D, radius float32, color glm.Color) {
	// todo: draw
}

// Draw2dCircleExt will draw one circle on canvas
//  outlineWidth - width of circle border, when width < radius, circle will have hole in center
//  color[0] - color of circle center
//  color[1] - color of circle border
//  GPU will blend colors in nice gradient
//
// This API has auto culling, and will discard drawing
// if all geometry outside of visible surface
func (r *Render) Draw2dCircleExt(pos glm.Local2D, radius float32, outlineWidth float32, color [2]glm.Color) {
	// todo: draw
}

// Draw2dPolygon will draw multi-lined polygon
// you can specify any number of edges, but for good performance
// not recommended using polygons with more than 32 edges.
//
// With vertex count < 5:
//  this method just proxy call to another render APIs
//  Draw2dPoint, Draw2dLine, Draw2dTriangle, Draw2dRect..
// With zero vertexes input, this will do nothing
//
// Polygon API not have automatic culling for >5 vertexes, it`s your responsibility
// to not call this, when all polygon vertexes outside of visible screen space
func (r *Render) Draw2dPolygon(pos []glm.Local2D, color glm.Color, filled bool) {
	switch len(pos) {
	case 0:
		return
	case 1:
		r.Draw2dPoint(pos[0], color)
		return
	case 2:
		r.Draw2dLine([2]glm.Local2D{pos[0], pos[1]}, [2]glm.Color{color, color})
		return
	case 3:
		r.Draw2dTriangleExt(
			[3]glm.Local2D{pos[0], pos[1], pos[2]},
			[3]glm.Color{color, color, color},
			filled,
		)
		return
	case 4:
		r.Draw2dRectExt(
			[4]glm.Local2D{pos[0], pos[1], pos[2], pos[3]},
			[4]glm.Color{color, color, color, color},
			filled,
		)
		return
	default:
		// todo: draw polygon
	}
}
