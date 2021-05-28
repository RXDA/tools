package tools

import (
	"fmt"
	"github.com/go-courier/geography"
	"math"
	"reflect"
	"testing"
)

func TestGetCenterAndAngle(t *testing.T) {
	type args struct {
		lineString [2]geography.Point
	}
	tests := []struct {
		name  string
		args  args
		want  geography.Point
		want1 float64
	}{
		{
			name: "1象限45",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {1, 1}},
			},
			want:  geography.Point{0.5, 0.5},
			want1: 45,
		},
		{
			name: "1象限30",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {math.Sqrt(3), 1}},
			},
			want:  geography.Point{math.Sqrt(3) / 2, 0.5},
			want1: 30,
		},
		{
			name: "2象限",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {-1, 1}},
			},
			want:  geography.Point{-0.5, 0.5},
			want1: 135,
		},
		{
			name: "2象限30",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {-math.Sqrt(3), 1}},
			},
			want:  geography.Point{-math.Sqrt(3) / 2, 0.5},
			want1: 150,
		},
		{
			name: "3象限",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {-1, -1}},
			},
			want:  geography.Point{-0.5, -0.5},
			want1: 225,
		},
		{
			name: "3象限30",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {-math.Sqrt(3), -1}},
			},
			want:  geography.Point{-math.Sqrt(3) / 2, 0.5},
			want1: 210,
		},
		{
			name: "4象限",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {1, -1}},
			},
			want:  geography.Point{0.5, -0.5},
			want1: 315,
		},
		{
			name: "4象限30",
			args: args{
				lineString: [2]geography.Point{{0, 0}, {math.Sqrt(3), -1}},
			},
			want:  geography.Point{math.Sqrt(3) / 2, -0.5},
			want1: 330,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetCenterAndAngle(tt.args.lineString)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCenterAndAngle() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetCenterAndAngle() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

var data = [][2]geography.Point{
	{{106.3729325, 29.9325858}, {106.3746457, 29.9318720}},
	{{106.3762637, 29.9312057}, {106.3775010, 29.9305395}},
}

func TestName(t *testing.T) {
	for _, v := range data {
		p, angle := GetCenterAndAngle(v)
		angle -= 90
		if angle < 0 {
			angle += 360
		}
		angle = 360 - angle
		fmt.Printf("{points: [%.7f,%.7f], rotate: %.3f},\n", p[0], p[1], angle)
	}
}
