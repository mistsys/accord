package db

import (
	"reflect"
	"testing"

	"github.com/mistsys/accord/id"
)

func TestLocalPSKStore_GetPSK(t *testing.T) {
	type args struct {
		key []byte
	}
	testKey, _ := id.KeyIDBytes("test", "")
	tests := []struct {
		name     string
		pskStore *LocalPSKStore
		args     args
		want     []byte
		wantErr  bool
	}{
		{
			name:     "PSK for local test works",
			pskStore: NewDummyPSKStore(),
			args: args{
				key: testKey,
			},
			want: []byte(`JpUtbRukLuIFyjeKpA4fIpjgs6MTV8eH`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.pskStore
			got, err := l.GetPSK(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalPSKStore.GetPSK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalPSKStore.GetPSK() = %v, want %v", got, tt.want)
			}
		})
	}
}
