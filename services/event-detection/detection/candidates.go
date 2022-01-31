package detection

import (
	convtree "github.com/angrymuskrat/conv-tree"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

func findCandidates(histGrid *convtree.ConvTree, posts []data.Post, maxPoints float64) (*convtree.ConvTree, bool) {
	for _, post := range posts {
		point := convtree.Point{
			X:       post.Lon,
			Y:       post.Lat,
			Content: post,
			Weight:  post.EventUtility,
		}
		histGrid.Insert(point, false)
	}
	hasAnomalies := detectCandidateTree(histGrid, maxPoints)
	if hasAnomalies {
		return histGrid, true
	}
	return nil, false
}

func detectCandidateTree(tree *convtree.ConvTree, maxPoints float64) bool {
	if tree.IsLeaf {
		if sumWeightPoints(tree.Points) >= maxPoints {
			return true
		}
		return false
	}
	res := false
	res = res || detectCandidateTree(tree.ChildBottomLeft, maxPoints)
	res = res || detectCandidateTree(tree.ChildBottomRight, maxPoints)
	res = res || detectCandidateTree(tree.ChildTopLeft, maxPoints)
	res = res || detectCandidateTree(tree.ChildTopRight, maxPoints)
	return res
}
