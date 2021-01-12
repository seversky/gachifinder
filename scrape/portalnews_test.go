package scrape

import (
	"testing"
	"github.com/seversky/gachifinder"
)

func TestPortalNews_Do(t *testing.T) {
	tests := []struct {
		name	string
		p		*PortalNews
		s		gachifinder.Scraper
	}{
		{
			name: 	"Scrape naver news",
			p: 	&PortalNews {
					Scrape {
						AllowedDomains: []string {
							"news.naver.com",
							"news.naver.com/main",
							"news.daum.net",
							"news.v.daum.net/v",
						},
					},
				},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s = tt.p
			tt.s.Do(tt.p.ParsingHandler)
		})
	}
}
