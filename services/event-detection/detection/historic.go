package detection

import (
	"strconv"
	"time"

	convtree "github.com/angrymuskrat/conv-tree"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/gonum/stat"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
)

const (
	minXLength = 0.005
	minYLength = 0.005
	maxDepth   = 20
	convNumber = 3
	gridSize   = 10
)

func HistoricGrid(postsData []data.Post, topLeft, bottomRight data.Point, maxPoints float64, tz string, gridSize float64) (convtree.ConvTree, error) {
	posts, numDays, err := splitPosts(postsData, tz, topLeft, gridSize)
	if err != nil {
		unilog.Logger().Error("unable to split posts", zap.Error(err))
		return convtree.ConvTree{}, err
	}
	averagedPosts := map[convtree.Point]float64{}
	for coordinate, postData := range posts {
		if len(postData) > 0 {
			averagedPosts[coordinate] = filterPosts(postData, numDays)
		}
	}
	tree, err := buildGrid(averagedPosts, topLeft, bottomRight, maxPoints)
	if err != nil {
		unilog.Logger().Error("unable to build historic grid", zap.Error(err))
		return convtree.ConvTree{}, err
	}
	tree.Clear()
	return tree, nil
}

func buildGrid(postsData map[convtree.Point]float64, topLeft, bottomRight data.Point, maxPoints float64) (convtree.ConvTree, error) {
	var points []convtree.Point
	for coordinate, postData := range postsData {
		numToAdd := postData
		if numToAdd < 1 {
			continue
		}
		point := convtree.Point{
			X:      coordinate.X,
			Y:      coordinate.Y,
			Weight: numToAdd,
		}
		points = append(points, point)
	}
	tl := convtree.Point{
		X:      topLeft.Lon,
		Y:      topLeft.Lat,
		Weight: 1,
	}
	br := convtree.Point{
		X:      bottomRight.Lon,
		Y:      bottomRight.Lat,
		Weight: 1,
	}
	tree, err := convtree.NewConvTree(tl, br, minXLength, minYLength, maxPoints, maxDepth, convNumber, gridSize, nil, points)
	if err != nil {
		unilog.Logger().Error("unalble to create ConvTree", zap.Error(err))
	}
	return tree, err
}

func filterPosts(posts map[string]int, numDays int) float64 {
	var postsCount []float64
	for _, v := range posts {
		postsCount = append(postsCount, float64(v))
	}
	diff := numDays - len(posts)
	for i := 0; i < diff; i++ {
		postsCount = append(postsCount, 0.0)
	}
	if len(postsCount) > 1 {
		avg := stat.Mean(postsCount, nil)
		std := stat.StdDev(postsCount, nil)
		maxValue := avg + 2*std
		var res []float64
		for _, v := range posts {
			val := float64(v)
			if val <= maxValue {
				res = append(res, val)
			}
		}
		mean := 0.0
		if len(res) > 0 {
			mean = stat.Mean(res, nil)
		}
		return mean
	} else {
		return postsCount[0]
	}
}

func splitPosts(data []data.Post, tz string, topLeft data.Point, gridSize float64) (map[convtree.Point]map[string]int, int, error) {
	posts := map[convtree.Point]map[string]int{}
	uniqueDays := map[string]string{}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		unilog.Logger().Error("unable to load timezone", zap.Error(err))
		return nil, 0, err
	}
	for _, post := range data {
		postGridLat := topLeft.Lat + float64(int((post.Lat-topLeft.Lat)/gridSize))*gridSize
		postGridLon := topLeft.Lon + float64(int((post.Lon-topLeft.Lon)/gridSize))*gridSize
		postGridPos := convtree.Point{X: postGridLon, Y: postGridLat}
		if _, ok := posts[postGridPos]; !ok {
			posts[postGridPos] = map[string]int{}
		}
		postTime := time.Unix(post.Timestamp, 0)
		postTime = postTime.In(loc)
		postYear, postMonth, postDay := postTime.Date()
		postDate := strconv.Itoa(postYear) + postMonth.String() + strconv.Itoa(postDay)
		if _, ok := posts[postGridPos][postDate]; !ok {
			posts[postGridPos][postDate] = 0
		}
		posts[postGridPos][postDate]++
		if _, ok := uniqueDays[postDate]; !ok {
			uniqueDays[postDate] = ""
		}
	}
	return posts, len(uniqueDays), nil
}
