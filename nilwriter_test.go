package iohelper

import (
	"bytes"
	"io"
	"testing"
)

func TestNilWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		w       io.Writer
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "nil writer",
			w:       nil,
			args:    args{p: []byte("test")},
			want:    4,
			wantErr: false,
		},
		{
			name:    "non-nil writer",
			w:       bytes.NewBuffer(nil),
			args:    args{p: []byte("test")},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNilWriter(tt.w).Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("NilWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.w != nil {
				got = tt.w.(*bytes.Buffer).Len()
			}
			if got != tt.want {
				t.Errorf("NilWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
