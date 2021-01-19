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
			name: 	"Scrape Test for Portal News",
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
			cd, done := tt.s.Do(tt.p.ParsingHandler)
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

			fmt.Println("The number of the collected data:", length)
			for _, data := range emitData {
				fmt.Println(data.Timestamp)
				fmt.Println(data.Creator)
				fmt.Println(data.Title)
				fmt.Println(data.Description)
				fmt.Println(data.URL)
				fmt.Println(data.ShortCutIconURL)
				fmt.Println(data.ImageURL)
			}
		})
	}
}
