package config

import (
	"testing"
)

func TestCofParse(t *testing.T) {
	type args struct {
		file string
		in   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"asong",args{file: "/Users/songsun/go/src/asong.cloud/Golang_Dream/Gin/gin_jwt_swagger/config.yaml",in: &Server{}},false},// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CofParse(tt.args.file, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("CofParse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}