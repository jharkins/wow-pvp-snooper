package toons

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

type Toon struct {
	Name    string
	Level   int
	Class   string
	Spec    string
	Rank    int
	Link    string
	Talents []Talent
}

type Talent struct {
	Name       string
	PictureURI string
}

type ToonData struct {
	toonCount map[string]uint
	toonList  []Toon
}

func (t *ToonData) GetCount() map[string]uint {
	return t.toonCount
}

func getTalents() []Talent {
	out := []Talent{}

	return out
}

func getToon(s string) (Toon, error) {
	parts := strings.Split(s, " ")
	l, err := strconv.Atoi(parts[0])
	if err != nil {
		return Toon{}, err
	}

	t := Toon{
		Level: l,
	}

	t.Spec = parts[1]


	switch parts[2] {
	case "Mastery":
		t.Class = "Hunter"
	case "Demon":
		t.Class = "Demon Hunter"
	case "Death":
		t.Class = "Death Knight"
	default:
		t.Class = parts[2]

	}
	return t, nil
}

func getPage(pageURL string) *goquery.Document {
	doc, err := goquery.NewDocument(pageURL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getToons(doc *goquery.Document) []*Toon {
	toons := []*Toon{}
	doc.Find(".SortTable-row").Each(func(i int, s *goquery.Selection) {
		tData := s.Find(".Character-level").Text()
		if tData == "" {
			return
		}
		t, err := getToon(tData)
		if err != nil {
			panic(err)
		}
		if t.Level != 110 {
			return
		}
		foo := s.Find(".Link")
		link, exists := foo.Attr("href")
		if exists {
			t.Link = link
		}
		t.Name = s.Find(".Character-name").Text()
		t.Rank, err = strconv.Atoi(s.Children().First().Text())
		if err != nil {
			panic(err)
		}
		toons = append(toons, &t)
	})
	return toons
}

func parseToons(toons []*Toon) ToonData {
	var tD = ToonData{}
	tD.toonCount = make(map[string]uint)
	tD.toonList = make([]Toon, len(toons))
	for i, t := range toons {
		tD.toonList[i] = *t
		classSpec := fmt.Sprintf("%s, %s", t.Class, t.Spec)
		if _, ok := tD.toonCount[classSpec]; ok {
			tD.toonCount[classSpec]++
		} else {
			tD.toonCount[classSpec] = 1
		}
	}

	return tD
}

func ParseURL(url string) ToonData {
	doc := getPage(url)
	rawToons := getToons(doc)
	toonData := parseToons(rawToons)
	return toonData
}

func GetToons(url string) []*Toon {
	doc := getPage(url)
	return getToons(doc)
}
