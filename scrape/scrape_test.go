package scrape

import (
	"fmt"
	"sync"
	"testing"

	"github.com/seversky/gachifinder"
)

func TestScrape_Do(t *testing.T) {
	tests := []struct {
		name	string
		p		*Scrape
		s		Scraper
	}{
		{
			name: 	"Scrape Test for Portal News",
			p: 	&Scrape {
					VisitDomains: []string {
						"https://" + NaverNews,
						"https://" + DaumNews,
					},
					// AllowedDomains: []string {
					// 	"news.naver.com",
					// 	"news.daum.net",
					// 	"news.v.daum.net",
					// },
				},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s = tt.p

			fs := []ParsingHandler {
				OnHTMLNaverHeadlineNews,
				OnHTMLDaumHeadlineNews,
			}
			dc := tt.s.Do(fs)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				emitData := make([]gachifinder.GachiData, 0, 20)
				for data := range dc {
					emitData = append(emitData, data)
				}

				length := len(emitData)
				if length == 0 {
					t.Error("There is not any collected data")
				}

				fmt.Println("The number of the collected data:", length)
				for _, data := range emitData {
					fmt.Println("")
					fmt.Println("Timestamp =", data.Timestamp)
					fmt.Println("VisitHost =", data.VisitHost)
					fmt.Println("Creator =", data.Creator)
					fmt.Println("Title =", data.Title)
					fmt.Println("Description =", data.Description)
					fmt.Println("URL =", data.URL)
					fmt.Println("ShortCutIconURL =", data.ShortCutIconURL)
					fmt.Println("ImageURL =", data.ImageURL)
				}

				wg.Done()
			}()
			wg.Wait()
		})
	}
}
