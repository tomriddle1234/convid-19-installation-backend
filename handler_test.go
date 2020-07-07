package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFuncAccumulateHandler(t *testing.T){
	ts := httptest.NewServer(GetMainEngine())
	defer ts.Close()
	client := &http.Client{}

	req, _ := http.NewRequest("POST", ts.URL + "/api/accumulatesignal", nil)

	// send out 11 requests within 10 seconds
	for i:=0;i<11;i++{
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			t.Fail()
		}


		//check response content
		body,_:= ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		resp.Body.Close()
		time.Sleep(800 * time.Millisecond)
	}

	// check status imediately
	req, _ = http.NewRequest("GET", ts.URL + "/api/statuscheck", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	//check response content
	body,_:= ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	defer resp.Body.Close()

	// check status after 10 seconds
	time.Sleep(10 * time.Second)
	req, _ = http.NewRequest("GET", ts.URL + "/api/statuscheck", nil)
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	//check response content
	body,_ = ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	defer resp.Body.Close()
}

func TestFuncAccumulateThenMultipleGetHandler(t *testing.T){
	ts := httptest.NewServer(GetMainEngine())
	defer ts.Close()
	client := &http.Client{}

	req, _ := http.NewRequest("POST", ts.URL + "/api/accumulatesignal", nil)

	// send out 2 requests within 2 seconds
	for i:=0;i<2;i++{
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			t.Fail()
		}


		//check response content
		body,_:= ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		resp.Body.Close()
		time.Sleep(800 * time.Millisecond)
	}

	// check status 100 times immediately
	req, _ = http.NewRequest("GET", ts.URL + "/api/statuscheck", nil)
	for i:=0;i<100; i++{
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
			t.Fail()
		}
		//check response content
		body,_:= ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		resp.Body.Close()
		time.Sleep(10 * time.Millisecond)

		// send another post signal among get requests
		if i == 15 {
			nreq, _ := http.NewRequest("POST", ts.URL + "/api/accumulatesignal", nil)
			nresp, err := client.Do(nreq)
			if err != nil {
				log.Fatal(err)
				t.Fail()
			}


			//check response content
			body,_:= ioutil.ReadAll(nresp.Body)
			log.Println(string(body))
			resp.Body.Close()
		}
	}

	// check status after 10 seconds
	time.Sleep(10 * time.Second)
	req, _ = http.NewRequest("GET", ts.URL + "/api/statuscheck", nil)
	fresp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}
	//check response content
	body,_:= ioutil.ReadAll(fresp.Body)
	log.Println(string(body))
	fresp.Body.Close()
}
