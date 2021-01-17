package emit

import (
	"testing"

	"github.com/seversky/gachifinder"
)

func TestEmit_Connect(t *testing.T) {
	test := struct {
		name string
		e    gachifinder.Emitter
	}{
		name: "Elasticsearch connecting test",
		e: &Elasticsearch{
			URLs: []string{"http://192.168.219.11:9200"},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		err := test.e.Connect()
		if err != nil {
			t.Error(err.Error())
		}
		test.e.Close()
	})
}

func TestElasticsearch_Write(t *testing.T) {
	tests := []struct {
		name string
		e    gachifinder.Emitter
	}{
		{
			name: "Elasticsearch writting test",
			e: &Elasticsearch{
				URLs: []string{"http://192.168.219.11:9200"},
			},
		},
	}

	esInfo := &Elasticsearch{
		URLs: []string{"http://192.168.219.11:9200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := esInfo.Connect()
			if err != nil {
				t.Error(err.Error())
			}
			esInfo.Write()
			esInfo.Close()
		})
	}
}

func TestElasticsearch_ManualDelete(t *testing.T) {

}
