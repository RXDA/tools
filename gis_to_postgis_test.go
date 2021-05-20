package tools

import (
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
			tt.want = strings.ReplaceAll(tt.want,", ",",")
			if got != tt.want {
				t.Errorf("GisToPostGis() got = %v, want %v", got, tt.want)
			}
		})
	}
}
