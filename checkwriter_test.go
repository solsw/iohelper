package iohelper

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestNewCheckWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	noop := func(*io.Writer, []byte) error { return nil }
	tests := []struct {
		name    string
		w       io.Writer
		check   func(*io.Writer, []byte) error
		wantErr bool
	}{
		{
			name:    "nil writer",
			w:       nil,
			check:   noop,
			wantErr: true,
		},
		{
			name:    "nil check",
			w:       buf,
			check:   nil,
			wantErr: true,
		},
		{
			name:    "valid",
			w:       buf,
			check:   noop,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCheckWriter(tt.w, tt.check)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCheckWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewCheckWriter() returned nil, want non-nil")
			}
		})
	}
}

func TestCheckWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		cw      *CheckWriter
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "nil CheckWriter",
			cw:      nil,
			args:    args{p: []byte("data")},
			want:    0,
			wantErr: true,
		},
		{
			name: "check returns error",
			cw: func() *CheckWriter {
				cw, _ := NewCheckWriter(&bytes.Buffer{}, func(*io.Writer, []byte) error {
					return errors.New("check failed")
				})
				return cw
			}(),
			args:    args{p: []byte("data")},
			want:    0,
			wantErr: true,
		},
		{
			name: "write succeeds",
			cw: func() *CheckWriter {
				cw, _ := NewCheckWriter(&bytes.Buffer{}, func(*io.Writer, []byte) error { return nil })
				return cw
			}(),
			args:    args{p: []byte("hello")},
			want:    5,
			wantErr: false,
		},
		{
			name: "check replaces writer",
			cw: func() *CheckWriter {
				replacement := bytes.NewBufferString("")
				cw, _ := NewCheckWriter(&bytes.Buffer{}, func(w *io.Writer, _ []byte) error {
					*w = replacement
					return nil
				})
				return cw
			}(),
			args:    args{p: []byte("replaced")},
			want:    8,
			wantErr: false,
		},
		{
			name: "empty payload",
			cw: func() *CheckWriter {
				cw, _ := NewCheckWriter(&bytes.Buffer{}, func(*io.Writer, []byte) error { return nil })
				return cw
			}(),
			args:    args{p: []byte{}},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cw.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
