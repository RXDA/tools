package tools

import (
	"fmt"
	"github.com/go-courier/geography"
	"math"
	"regexp"
	"strings"
)

//GisSliceToPostGis  gis数组转换为postgis扩展格式
// point [1,1]  -> POINT (1 1)
// line [point,point] -> [[1,1],[2,2]] -> LINESTRING (1 1, 2 2)
// polygon [line,line] -> [[[0,0],[0,1],[1,1],[1,0],[0,0]]] -> POLYGON ((0 0, 0 1, 1 1, 1 0, 0 0))
// polygon [line,line] -> [[[0,0],[0,2],[2,2],[2,0],[0,0]],[[0,0],[0,1],[1,1],[1,0],[0,0]]] -> POLYGON ((0 0, 0 2, 2 2, 2 0, 0 0), (0 0, 0 1, 1 1, 1 0, 0 0))
// multiPolygon [polygon,polygon] -> [[[[0,0],[0,1],[1,1],[1,0],[0,0]]],[[[0,0],[0,-1],[-1,-1],[-1,0],[0,0]]]]
func GisSliceToPostGis(gisStr string) (string, error) {
	gisStr = strings.ReplaceAll(gisStr, " ", "")
	if len(gisStr) < 5 {
		return "", nil
	}

	strBytes := []byte(gisStr)
	leftBracketsNum := 0
	for i := 0; i < len(strBytes); i++ {
		if strBytes[i] == '[' {
			leftBracketsNum++
		} else {
			break
		}
	}

	var geomType string
	bs := []byte(gisStr)
	switch leftBracketsNum {
	case 1:
		geomType = "POINT"
	case 2:
		geomType = "LINESTRING"
		gisStr = string(bs[1 : len(bs)-1])
	case 3:
		geomType = "POLYGON"
		gisStr = string(bs[1 : len(bs)-1])
	case 4:
		geomType = "MULTIPOLYGON"
		gisStr = string(bs[1 : len(bs)-1])
	}

	pointReg, err := regexp.Compile(`\[(\d+\.*\d{0,7})\d*\s*,\s*(\d+\.*\d{0,7})\d*]`)
	if err != nil {
		return "", err
	}

	replacedPoint := pointReg.ReplaceAllString(gisStr, "$1 $2")

	replacer := strings.NewReplacer("[", "(", "]", ")")
	replacedPoint = replacer.Replace(replacedPoint)

	return fmt.Sprintf("%s(%s)", geomType, replacedPoint), nil
}

func maxAndMin(fs ...float64) (max, min float64) {
	maxF := -math.MaxFloat64
	minF := math.MaxFloat64
	for _, v := range fs {
		if v > maxF {
			maxF = v
		}
		if v < minF {
			minF = v
		}
	}
	return maxF, minF
}

func GetPolygonBoxAndCenter(p geography.Polygon) ([][2]float64, [2]float64) {
	maxX := -math.MaxFloat64
	minX := math.MaxFloat64
	maxY := -math.MaxFloat64
	minY := math.MaxFloat64
	for _, l := range p {
		for _, p := range l {
			if p.X() > maxX {
				maxX = p.X()
			}
			if p.X() < minX {
				minX = p.X()
			}
			if p.Y() > maxY {
				maxY = p.Y()
			}
			if p.Y() < minY {
				minY = p.Y()
			}
		}
	}
	// box 从左下角开始，逆时针
	box := [][2]float64{
		{minX, minY},
		{maxX, minY},
		{maxX, maxY},
		{minX, maxY},
		{minX, minY},
	}
	center := [2]float64{(minX + maxX) / 2, (minY + maxY) / 2}
	return box, center
}

func GetMultiPolygonBoxAndCenter(mp geography.MultiPolygon) ([5][2]float64, [2]float64) {
	maxX := -math.MaxFloat64
	minX := math.MaxFloat64
	maxY := -math.MaxFloat64
	minY := math.MaxFloat64
	for _, p := range mp {
		for _, l := range p {
			for _, p := range l {
				if p.X() > maxX {
					maxX = p.X()
				}
				if p.X() < minX {
					minX = p.X()
				}
				if p.Y() > maxY {
					maxY = p.Y()
				}
				if p.Y() < minY {
					minY = p.Y()
				}
			}
		}
	}

	// box 从左下角开始，逆时针
	box := [5][2]float64{
		{minX, minY},
		{maxX, minY},
		{maxX, maxY},
		{minX, maxY},
		{minX, minY},
	}
	center := [2]float64{(minX + maxX) / 2, (minY + maxY) / 2}
	return box, center
}


func GetMultiPolygonBoxAndCenterFromGeo(mp geography.MultiPolygon) ([][2]float64, [2]float64) {
	box := mp.Bound().AsPolygon()
	center := mp.Bound().Center()
	boxResult := [][2]float64{
		(*box)[0][0],
		(*box)[0][1],
		(*box)[0][2],
		(*box)[0][3],
		(*box)[0][4],
	}
	return boxResult, center
}