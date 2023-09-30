package diff

import (
	"context"
	"reflect"
	"testing"
)

func TestCalculate(t *testing.T) {
	type args struct {
		ctx      context.Context
		repoPath string
		commitA  string
		commitB  string
	}
	tests := []struct {
		name    string
		args    args
		want    []Diff
		wantErr bool
	}{
		{
			name: "Test different commits",
			args: args{
				ctx:      context.Background(),
				repoPath: ".",
				commitA:  "39bf3e493ec8c0bd640c7670a0f6a1a3c92cb39d",
				commitB:  "bba8505866a7172e1269f5a8cd0ceba1258c7880",
			},
			want: []Diff{
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x63, 0x65, 0x66, 0x37, 0x61, 0x38, 0x66}, Path: ".gitignore"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x31, 0x33, 0x35, 0x36, 0x36, 0x62, 0x38}, Path: ".idea/.gitignore"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x37, 0x65, 0x65, 0x30, 0x37, 0x38, 0x64}, Path: ".idea/git-diff-view.iml"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x33, 0x35, 0x65, 0x62, 0x31, 0x64, 0x64}, Path: ".idea/vcs.xml"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x62, 0x39, 0x31, 0x63, 0x31, 0x36, 0x66}, Path: "diff/diff.go"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x63, 0x63, 0x31, 0x30, 0x37, 0x62, 0x35}, Path: "diff/entities.go"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x65, 0x64, 0x39, 0x38, 0x35, 0x31, 0x32}, Path: "diff/errors.go"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x63, 0x66, 0x36, 0x30, 0x38, 0x31, 0x32}, Path: "diff/status.go"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}},
				{Status: Status{Type: StatusModify}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x37, 0x64, 0x37, 0x34, 0x35, 0x30, 0x61}, Path: "go.mod"}, Dst: State{Mode: 0x31303036, SHA1: []uint8{0x37, 0x33, 0x38, 0x34, 0x62, 0x36, 0x66}}},
				{Status: Status{Type: StatusDelete}, Src: State{Mode: 0x31303036, SHA1: []uint8{0x33, 0x64, 0x66, 0x36, 0x30, 0x32, 0x63}, Path: "go.sum"}, Dst: State{Mode: 0x30303030, SHA1: []uint8{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.args.ctx, tt.args.repoPath, tt.args.commitA, tt.args.commitB)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDiff(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name     string
		args     args
		wantDiff Diff
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDiff, err := parseDiff(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDiff, tt.wantDiff) {
				t.Errorf("parseDiff() gotDiff = %v, want %v", gotDiff, tt.wantDiff)
			}
		})
	}
}
