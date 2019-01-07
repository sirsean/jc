package commands

import (
    "log"
    "fmt"
	"github.com/sirsean/go-pool"
	p "github.com/sirsean/jc/path"
    "github.com/sirsean/jc/request"
	"sort"
)

func Ls() {
    requests, err := getRequests()
    if err != nil {
        log.Fatal(err)
    }

	for _, r := range requests {
		fmt.Printf("%s: %s %s\n", r.Id, r.Method, r.Url)
	}
}

type LoadRequestWorkUnit struct {
	Id    string
	OutCh chan request.Request
}

func (u LoadRequestWorkUnit) Perform() {
	r, err := request.LoadRequest(u.Id)
	if err == nil {
		u.OutCh <- r
	}
}

func getRequests() ([]request.Request, error) {
	files, err := p.ListFiles()
	if err != nil {
		return nil, err
	}

	p := pool.NewPool(100, 100) // we will read many files at a time
	p.Start()

	ch := make(chan request.Request, len(files)) // channel for each request
	rCh := make(chan []request.Request, 1)       // result channel for full requests list

	go func() {
		// we will build up the list of requests in this goroutine
		requests := make([]request.Request, 0)
		for r := range ch {
			requests = append(requests, r)
		}
		// all the requests are in the list, sort it
		sort.Slice(requests[:], func(i, j int) bool {
			return requests[i].Id < requests[j].Id
		})
		// send the list back to the main thread
		rCh <- requests
	}()

	for _, f := range files {
		p.Add(LoadRequestWorkUnit{
			Id:    f.Name(),
			OutCh: ch,
		})
	}

	p.Close()
	close(ch)

	requests := <-rCh

	return requests, nil
}
