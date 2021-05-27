package tools

import (
	"github.com/go-courier/geography"
	"math"
)

func GetCenterAndAngle(lineString [2]geography.Point) (geography.Point, float64) {
	center := geography.Point{(lineString[0].X() + lineString[1].X()) / 2, (lineString[0].Y() + lineString[1].Y()) / 2}

	if lineString[1].Y() == lineString[0].Y() {
		if lineString[1].X() > lineString[0].X() {
			return center, 90
		} else {
			return center, 270
		}
	}
	dy := lineString[1].Y() - lineString[0].Y()
	dx := lineString[1].X() - lineString[0].X()
	k := dy / dx
	angle := 180 * math.Atan(k) / math.Pi
	switch quadrant(lineString) {
	case 2:
		angle += 180
	case 3:
		angle = angle + 180
	case 4:
		angle += 360
	}

	//fmt.Println(dx, dy, angle, quadrant(lineString))
	return center, angle
}

func quadrant(lineString [2]geography.Point) uint8 {
	if lineString[1][0]-lineString[0][0] > 0 {
		if lineString[1][1]-lineString[0][1] > 0 {
			return 1
		} else {
			return 4
		}
	} else {
		if lineString[1][1]-lineString[0][1] > 0 {
			return 2
		} else {
			return 3
		}
	}
}
