package config

import (
	"sync"
	"testing"
)

func TestProperties_PutDecimal(t *testing.T) {
	type fields struct {
		props         map[string]*value
		resolvedProps map[string]string
		RWMutex       sync.RWMutex
	}
	type args struct {
		k string
		v float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Properties{
				props:         tt.fields.props,
				resolvedProps: tt.fields.resolvedProps,
				RWMutex:       tt.fields.RWMutex,
			}
			got, err := p.PutDecimal(tt.args.k, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutDecimal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PutDecimal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
