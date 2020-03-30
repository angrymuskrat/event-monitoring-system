package detection

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
)

func findCandidates(histGrid *convtree.ConvTree, posts []data.Post, maxPoints int) (*convtree.ConvTree, bool) {
	grid := copyTree(histGrid)
	for _, post := range posts {
		point := convtree.Point{
			X:       post.Lon,
			Y:       post.Lat,
			Content: post,
			Weight:  1,
		}
		grid.Insert(point, false)
	}
	hasAnomalies := detectCandTree(grid, maxPoints)
	if hasAnomalies {
		return grid, true
	}
	return nil, false
}

func copyTree(oldTree *convtree.ConvTree) *convtree.ConvTree {
	if oldTree == nil {
		return nil
	}
	newTree := *oldTree
	newTree.ChildBottomLeft = copyTree(oldTree.ChildBottomLeft)
	newTree.ChildBottomRight = copyTree(oldTree.ChildBottomRight)
	newTree.ChildTopLeft = copyTree(oldTree.ChildTopLeft)
	newTree.ChildTopRight = copyTree(oldTree.ChildTopRight)
	return &newTree
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
