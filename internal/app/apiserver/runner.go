package apiserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
)

type listURL struct {
	URL []string `json:"urls"`
}

func (l *listURL) validate() error {
	if len(l.URL) > 20 {
		return errors.New("invalid number of urls in request")
	}

	return nil
}

func (l *listURL) String() string {
	result := ""
	for _, str := range l.URL {
		result += str
	}
	return result
}

//curl -X POST -H "Content-Type: application/json" -d '{"urls":["http://localhost:8000/test1", "http://localhost:8000/test2","http://localhost:8000/test3","http://localhost:8000/test4","http://localhost:8000/test5","http://localhost:8000/test6","http://localhost:8000/test7","http://localhost:8000/test8"]}' http://localhost:8080/run

func (l *listURL) Start() ([]string, error) {
	errChan := make(chan error)
	outChan := make(chan string)
	in := make(chan string)
	result := []string{}
	wg := sync.WaitGroup{}

	go func() {
		pushURL(l.URL, outChan)
	}()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range outChan {
				err := getRequest(url, in)
				if err != nil {
					errChan <- err
					break
				}
			}
		}()
	}

	for i := 0; i < len(l.URL); i++ {
		select {
		case err := <-errChan:
			wg.Wait()
			return nil, err
		case str := <-in:
			result = append(result, str)
		}
	}

	return result, nil
}

func pushURL(urls []string, out chan<- string) {
	defer close(out)

	for _, url := range urls {
		out <- url
	}
}

func getRequest(url string, in chan<- string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	in <- string(body)
	return nil
}

func (a *APIServer) runner(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	var urls listURL

	err := json.NewDecoder(r.Body).Decode(&urls)
	if err != nil {
		return err
	}

	err = urls.validate()
	if err != nil {
		return err
	}

	text, err := urls.Start()
	if err != nil {
		return err
	}

	return jsonResponse(w, text)
}

func jsonResponse(w http.ResponseWriter, text []string) error {
	resp := struct {
		Response struct {
			Msg []string `json:"url_response"`
		} `json:"responses"`
	}{}
	resp.Response.Msg = text

	bytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, string(bytes))
	return err
}
