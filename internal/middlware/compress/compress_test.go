package compress

import (
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gzipBody(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("gzipBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gzipBody() got = %v, want %v", got, tt.want)
			}
		})
	}
}
