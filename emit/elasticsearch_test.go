package emit

import (
	"testing"

	"github.com/seversky/gachifinder"
)

func TestEmit_Connect(t *testing.T) {
	tests := []struct {
		name string
		e    gachifinder.Emitter
	}{
		{
			name: "Elasticsearch connecting test",
			e: &Elasticsearch{
				URLs: []string{"http://192.168.219.11:9200"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Connect()
			tt.e.Close()
		})
	}
}

func TestElasticsearch_Write(t *testing.T) {
	tests := []struct {
		name string
		e    gachifinder.Emitter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Write()
		})
	}
}
