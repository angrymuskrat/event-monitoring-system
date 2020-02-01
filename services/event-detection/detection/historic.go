package detection

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/gonum/stat"
	convtree "github.com/visheratin/conv-tree"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GenerateGrid(data []data.Post, topLeft, bottomRight data.Point, maxPoints int, tz string, gridSize float64) (convtree.ConvTree, error) {
	posts, numDays, err := splitPosts(data, tz, topLeft, gridSize)
	if err != nil {
		unilog.Logger().Error("unable to split posts", zap.Error(err))
		return convtree.ConvTree{}, err
	}
	averagedPosts := map[convtree.Point]float64{}
	for coord, data := range posts {
		if len(data) > 0 {
			averagedPosts[coord] = filterPosts(data, numDays)
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

func buildGrid(postData map[convtree.Point]float64, topLeft, bottomRight data.Point, maxPoints int) (convtree.ConvTree, error) {
	points := []convtree.Point{}
	for coord, data := range postData {
		numToAdd := int(data)
		if numToAdd < 1 {
			continue
		}
		point := convtree.Point{
			X:      coord.X,
			Y:      coord.Y,
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
	tree, err := convtree.NewConvTree(tl, br, 0.001, 0.001, maxPoints, 50, 3, 10, nil, points)
	if err != nil {
		unilog.Logger().Error("unalble to create ConvTree", zap.Error(err))
	}
	return tree, err
}

func filterPosts(posts map[string]int, numDays int) float64 {
	data := []float64{}
	for _, v := range posts {
		data = append(data, float64(v))
	}
	diff := numDays - len(posts)
	for i := 0; i < diff; i++ {
		data = append(data, 0.0)
	}
	if len(data) > 1 {
		avg := stat.Mean(data, nil)
		std := stat.StdDev(data, nil)
		maxValue := avg + 2*std
		res := []float64{}
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
		return data[0]
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
		postGridPos := convtree.Point{X: postGridLat, Y: postGridLon}
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
