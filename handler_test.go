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
