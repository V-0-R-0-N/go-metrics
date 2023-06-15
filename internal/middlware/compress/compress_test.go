package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"reflect"
	"testing"
)

func Test_gzipBody(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple test 1",
			args: args{
				data: []byte("123"),
			},
			want:    []byte("123"),
			wantErr: false,
		},
		{
			name: "Simple test 2",
			args: args{
				data: []byte("12"),
			},
			want:    []byte("123"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gzipBody(tt.args.data)
			if err != nil {
				t.Errorf("gzipBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			reader := bytes.NewReader(got)
			gzreader, _ := gzip.NewReader(reader)
			output, _ := io.ReadAll(gzreader)

			defer gzreader.Close()

			if !reflect.DeepEqual(output, tt.args.data) {
				if !tt.wantErr {
					t.Errorf("gzipBody() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
