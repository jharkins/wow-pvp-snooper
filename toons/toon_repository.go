package toons

import (
	"fmt"
	"sync"
	"time"
	"log"
)

var closeChan = make(chan bool)

/*
var classes []string = []string{
	"Demon Hunter",
	"Death Knight",
	"Paladin",
	"Priest",
	"Mage",
	"Rogue",
	"Warlock",
	"Warrior",
	"Hunter",
	"Shaman",
	"Monk",
	"Druid",
}
*/

type ToonRepository struct {
	toons        []*Toon
	numPages     int
	lastFetch    time.Time
	fetchTimeout time.Duration
	bNetURI      string
	lock         sync.Mutex
}

func (tr *ToonRepository) GetToons() []*Toon {
	return tr.toons
}

func (tr *ToonRepository) GetClassCount() map[string]int {
	out := make(map[string]int)
	for _, t := range tr.toons {
		if _, ok := out[t.Class]; !ok {
			out[t.Class] = 0
		}
		out[t.Class] += 1
	}
	return out
}

func (tr *ToonRepository) GetSpecCount() map[string]map[string]int {
	out := make(map[string]map[string]int)
	for _, t := range tr.toons {
		//cSpec := fmt.Sprintf("%s, %s", t.Spec, t.Class)
		if _, ok := out[t.Class]; !ok {
			out[t.Class] = make(map[string]int)
		}

		if _, ok := out[t.Class][t.Spec]; !ok {
			out[t.Class][t.Spec] = 0
		}
		out[t.Class][t.Spec] += 1
	}
	return out
}

func (tr *ToonRepository) loop(closeChan chan bool) {
	fetchTimeout := time.After(tr.fetchTimeout)
	for {
		select {
		case <-fetchTimeout:
			log.Println("Refreshing repository from bnet!")
			tr.updateRepository()
			fetchTimeout = time.After(tr.fetchTimeout)
		case <-closeChan:
			break
		}
	}
}

func (tr *ToonRepository) updateRepository() error {
	tr.lock.Lock()
	defer tr.lock.Unlock()

	// Reset it!
	tr.toons = []*Toon{}

	var wg sync.WaitGroup
	var toonChan chan *Toon = make(chan *Toon)

	for i := 1; i <= tr.numPages; i++ {
		wg.Add(1)
		go func(idx int) {
			uri := fmt.Sprintf("%s%d", tr.bNetURI, idx)
			td := GetToons(uri)
			for _, t := range td {
				toonChan <- t
			}
			wg.Done()
		}(i)
	}

	go func() {
		for t := range toonChan {
			tr.toons = append(tr.toons, t)
		}
	}()

	wg.Wait()
	close(toonChan)

	return nil
}

// Stop the loop
func (tr *ToonRepository) Close() error {
	closeChan <- true
	return nil
}

func NewToonRepository(bNetURI string, numPages int) (*ToonRepository, error) {
	tr := ToonRepository{
		bNetURI:      bNetURI,
		numPages:     numPages,
		fetchTimeout: 15 * time.Minute,
	}

	if err := tr.updateRepository(); err != nil {
		return nil, err
	}

	go tr.loop(closeChan)

	return &tr, nil
}
