package scrape

import (
	"testing"
	"fmt"

	"github.com/seversky/gachifinder"
)

func TestPortalNews_Do(t *testing.T) {
	tests := []struct {
		name	string
		p		*PortalNews
		s		gachifinder.Scraper
	}{
		{
			name: 	"Scrape Portal news",
			p: 	&PortalNews {
					Scrape {
						VisitDomains: []string {
							"https://news.naver.com/",
							// "https://news.daum.net/",
						},
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

			done := make(chan bool)
			cd := make(chan gachifinder.GachiData)

			go tt.s.Do(tt.p.ParsingHandler, cd, done)

			emitData := make([]gachifinder.GachiData, 0, 20)
			for c := true; c;{
				select {
				case data := <-cd:
					emitData = append(emitData, data)
				case <-done:
					c = false
				}
			}

			length := len(emitData)
			if length == 0 {
				t.Error("There is not any collected data")
			}

			fmt.Println(length)
			for _, data := range emitData {
				fmt.Println(data.Timestamp)
				fmt.Println(data.Title)
				fmt.Println(data.Creator)
			}
		})
	}
}
