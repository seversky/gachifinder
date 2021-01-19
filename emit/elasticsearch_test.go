package emit

import (
	"testing"
	"time"

	"github.com/seversky/gachifinder"
)

const esURL = "http://192.168.56.105:9200"

func TestEmit_Connect(t *testing.T) {
	test := struct {
		name string
		e    gachifinder.Emitter
	}{
		name: "Elasticsearch connecting test",
		e: &Elasticsearch{
			URLs: []string{esURL},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		err := test.e.Connect()
		if err != nil {
			t.Error(err)
		}
		test.e.Close()
	})
}

func TestElasticsearch_Write(t *testing.T) {
	test := struct {
		name string
		e    gachifinder.Emitter
	}{
		name: "Elasticsearch writing test",
		e: &Elasticsearch{
			URLs: []string{esURL},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		err := test.e.Connect()
		if err != nil {
			t.Error(err)
		}
		defer test.e.Close()

		cd, done := func () (<-chan gachifinder.GachiData, <-chan bool) {
			done := make(chan bool)
			cd := make(chan gachifinder.GachiData)
			timestamp := time.Now()
	
			go func() {
				emitData := []gachifinder.GachiData{
					{
						Timestamp: timestamp,
						ShortCutIconURL: "https://ssl.pstatic.net/static.news/image/news/2014/favicon/favicon.ico",
						Title: "이재용 부회장 2년6개월 실형…남은 형량 1년반 다 채울까",
						URL: "https://news.naver.com/main/read.nhn?mode=LSD&mid=sec&oid=421&aid=0005114764&sid1=001",
						ImageURL: "https://imgnews.pstatic.net/image/421/2021/01/18/0005114764_001_20210118165914763.jpg",
						Creator: "뉴스1",
						Description: "(서울=뉴스1) 이세현 기자 = 이재용 삼성전자 부사장이 '국정농단' 사건 파기환송심에서 징역 2년6개월의 실형을 선고받았다. 통상 파기환송심은 대법원에서 그대로 확정되는 경우가 많은데다, 특별사면 절차도 녹록치 않",
					},
					{
						Timestamp: timestamp,
						ShortCutIconURL: "https://ssl.pstatic.net/static.news/image/news/2014/favicon/favicon.ico",
						Title: "文 ‘입양아 바꾸기’에 들끓은 여론…靑 “사전위탁제 말한 것”",
						URL: "https://news.naver.com/main/read.nhn?mode=LSD&mid=sec&oid=025&aid=0003070628&sid1=001",
						ImageURL: "https://imgnews.pstatic.net/image/025/2021/01/18/0003070628_001_20210118165511956.jpg",
						Creator: "중앙일보",
						Description: "18일 문재인 대통령이 신년 기자회견에서 입양에 대해 한 발언이 파문을 일으켰다. 발단은 양부모의 학대로 입양아가 사망한 ‘정인이 사건’에 대한 질문이었다. 문 대통령은 사건의 재발 방지 대책을 설명하던 중 “입양",
					},
				}
	
				for _, data := range emitData {
					cd <- data
				}
				done <- true
				close(cd)
				close(done)
			}()

			return cd, done
		}()

		test.e.Write(cd, done)

		
	})
}

func TestElasticsearch_ManualDelete(t *testing.T) {

}
