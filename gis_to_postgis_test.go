package tools

import (
	"encoding/json"
	"github.com/go-courier/geography"
	"reflect"
	"strings"
	"testing"
)

func TestGisToPostGis(t *testing.T) {
	type args struct {
		gisStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "point",
			args: args{gisStr: "[1,1]"},
			want: "POINT(1 1)",
		},
		{
			name: "line",
			args: args{gisStr: "[[1,1],[2,2]]"},
			want: "LINESTRING(1 1,2 2)",
		},
		{
			name: "polygon1",
			args: args{gisStr: "[[[30,10],[40,40],[20,40],[10,20],[30,10]]]"},
			want: "POLYGON((30 10,40 40,20 40,10 20,30 10))",
		},
		{
			name: "polygon2",
			args: args{gisStr: "[[[35,10],[45,45],[15,40],[10,20],[35,10]],[[20,30],[35,35],[30,20],[20,30]]]"},
			want: "POLYGON((35 10,45 45,15 40,10 20,35 10),(20 30,35 35,30 20,20 30))",
		},
		{
			name: "multiPolygon1",
			args: args{gisStr: "[[[[30,20],[45,40],[10,40],[30,20]]],[[[15,5],[40,10],[10,20],[5,10],[15,5]]]]"},
			want: "MULTIPOLYGON(((30 20, 45 40, 10 40, 30 20)),((15 5, 40 10, 10 20, 5 10, 15 5)))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GisSliceToPostGis(tt.args.gisStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GisToPostGis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want = strings.ReplaceAll(tt.want, ", ", ",")
			if got != tt.want {
				t.Errorf("GisToPostGis() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMultiPolygonBoxAndCenter(t *testing.T) {
	mp1Str := "[[[[0,0],[5,0],[5,5],[0,5],[0,0]]],[[[6,6],[7,6],[7,7],[6,7],[6,6]]]]" // 两个不相邻的正方形
	var mp1 geography.MultiPolygon
	err := json.Unmarshal([]byte(mp1Str), &mp1)
	if err != nil {
		panic(err)
	}
	type args struct {
		mp geography.MultiPolygon
	}
	tests := []struct {
		name  string
		args  args
		want  [][2]float64
		want1 [2]float64
	}{
		{
			name: "1",
			args: args{
				mp: mp1,
			},
			want:  [][2]float64{{0, 0}, {7, 0}, {7, 7}, {0, 7}, {0, 0}},
			want1: [2]float64{3.5, 3.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetMultiPolygonBoxAndCenterFromGeo(tt.args.mp)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMultiPolygonBoxAndCenter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetMultiPolygonBoxAndCenter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
