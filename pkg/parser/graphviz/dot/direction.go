package dot

type DiagramDirection string

// Sets direction of graph layout: https://graphviz.org/docs/attrs/rankdir.
const (
	DirectionTopToBottom DiagramDirection = "TB"
	DirectionBottomToTop DiagramDirection = "BT"
	DirectionLeftToRight DiagramDirection = "LR"
	DirectionRightToLeft DiagramDirection = "RL"
)

const DefaultDirection DiagramDirection = DirectionTopToBottom
