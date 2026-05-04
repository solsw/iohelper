package iohelper

import (
	"bytes"
	"testing"
)

func TestNilSafeWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		nw      *NilSafeWriter
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "nil NilSafeWriter",
			nw:      nil,
			args:    args{p: []byte("test")},
			want:    0,
			wantErr: true,
		},
		{
			name:    "nil writer",
			nw:      NewNilSafeWriter(nil),
			args:    args{p: []byte("test")},
			want:    4,
			wantErr: false,
		},
		{
			name:    "non-nil writer",
			nw:      NewNilSafeWriter(bytes.NewBuffer(nil)),
			args:    args{p: []byte("test")},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.nw.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("NilSafeWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NilSafeWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
