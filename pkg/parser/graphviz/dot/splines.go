package dot

type DiagramSpline string

// Controls how, and if, edges are represented: https://graphviz.org/docs/attrs/splines.
const (
	SplineCurved   DiagramSpline = "curved"
	SplineLine     DiagramSpline = "line"
	SplineNone     DiagramSpline = ""
	SplineOrtho    DiagramSpline = "ortho"
	SplinePolyline DiagramSpline = "polyline"
	SplineSpline   DiagramSpline = "spline"
)

const DefaultSpline DiagramSpline = SplineNone
