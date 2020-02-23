package detection

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
)

func findCandidates(histGrid *convtree.ConvTree, posts []data.Post, maxPoints int) (*convtree.ConvTree, bool) {
	for _, post := range posts {
		point := convtree.Point{
			X:       post.Lon,
			Y:       post.Lat,
			Content: post,
			Weight:  1,
		}
		histGrid.Insert(point, false)
	}
	hasAnomalies := detectCandTree(histGrid, maxPoints)
	if hasAnomalies {
		return histGrid, true
	}
	return nil, false
}

func detectCandTree(tree *convtree.ConvTree, maxPoints int) bool {
	if tree.IsLeaf {
		if len(tree.Points) >= maxPoints {
			return true
		}
		return false
	}
	res := false
	res = res || detectCandTree(tree.ChildBottomLeft, maxPoints)
	res = res || detectCandTree(tree.ChildBottomRight, maxPoints)
	res = res || detectCandTree(tree.ChildTopLeft, maxPoints)
	res = res || detectCandTree(tree.ChildTopRight, maxPoints)
	return res
}
