package id

import "testing"

/*
The numbers here were cross validated using the following simple python function:
import hashlib, struct
def key_id(s, salt):
    m = hashlib.md5()
    m.update(salt)
    m.update(s)
    b = m.digest()[:4]
    return struct.unpack('>I', b)[0]
*/

func TestKeyID(t *testing.T) {
	type args struct {
		s    string
		salt string
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			name: "test id with empty salt",
			args: args{
				s:    "test",
				salt: "",
			},
			want:    uint32(160394189),
			wantErr: false,
		},
		{
			name: "test with a known source path and AWS account - no salt",
			args: args{
				s:    "staging_ec2_660610034966_us-east-1",
				salt: "",
			},
			want:    uint32(1160306280),
			wantErr: false,
		},
		{
			name: "test with a known source path and AWS account - with salt",
			args: args{
				s:    "staging_ec2_660610034966_us-east-1",
				salt: "hUYh5x4N2DOnTIce",
			},
			want:    uint32(3299138274),
			wantErr: false,
		},
		{
			name: "empty strings return error",
			args: args{
				s:    "",
				salt: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KeyID(tt.args.s, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("KeyID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("KeyID() = %v, want %v", got, tt.want)
			}
		})
	}
}
