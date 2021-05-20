package tools

import (
	"fmt"
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
	gisStr = strings.ReplaceAll(gisStr," ","")
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

	pointReg, err := regexp.Compile(`\[(\d+\.*\d*)\s*,\s*(\d+\.*\d*)]`)
	if err != nil {
		return "", err
	}

	replacedPoint := pointReg.ReplaceAllString(gisStr, "$1 $2")

	replacer := strings.NewReplacer("[", "(", "]", ")")
	replacedPoint = replacer.Replace(replacedPoint)

	return fmt.Sprintf("%s(%s)", geomType, replacedPoint), nil
}
