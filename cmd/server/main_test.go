package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_getMetrics(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		metrics map[string]map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]map[string]interface{}
	}{
		{
			name: "test #1",
			fields: fields{
				Mutex:   sync.Mutex{},
				metrics: map[string]map[string]interface{}{},
			},
			want: map[string]map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Mutex:   tt.fields.Mutex,
				metrics: tt.fields.metrics,
			}
			if got := m.getMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
