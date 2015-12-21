package main

import (
"fmt"
 "encoding/json"
 "hash/crc64"
 "github.com/emirpasic/gods/maps/treemap"
 "strings"
 "log"
 "github.com/julienschmidt/httprouter"
 "net/http"
  )
type KeyValue struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

type GetAllKeysRes []KeyValue

type PutKeyValueResponse struct {
	Message string `json:"message"`
}
var circle *treemap.Map = treemap.NewWith(UInt64Comparator)

func main() {

		
	addToCircle("http://localhost:3000/keys")
	addToCircle("http://localhost:3001/keys")
	addToCircle("http://localhost:3002/keys")

	mux := httprouter.New()
	mux.PUT("/keys/:key_id/:value", storeKeyValue)
	mux.GET("/keys/:key_id", getKeyValue)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()

	
}

func getURLPut(key string, value string) string {

	index := getNodeForKey(key)
	inputUrl, _ := circle.Get(index)
	url := inputUrl.(string) + "/" + key + "/" + value
	fmt.Println("Input Url : ", url)
	return url
}

func getURLGet(key string) string {

	index := getNodeForKey(key)
	inputUrl, _ := circle.Get(index)
	url := inputUrl.(string) + "/" + key
	fmt.Println("Input Url : ", url)
	return url
}

func UInt64Comparator(a, b interface{}) int {
	aInt := a.(uint64)
	bInt := b.(uint64)
	switch {
	case aInt > bInt:
		return 1
	case aInt < bInt:
		return -1
	default:
		return 0
	}
}

func addToCircle(serverUrl string) {
	circle.Put(HashCode(serverUrl), serverUrl)
}


func HashCode(data string) uint64 {
	crcTable := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum([]byte(data), crcTable)
}

func getNodeForKey(key string) uint64 {

	keyHash := HashCode(key)
	keys := circle.Keys()

	return keys[getNodeIndex(keys, keyHash, 0)].(uint64)
}

func getNodeIndex(keys []interface{}, keyHash uint64, index int) int {

	if index < circle.Size() && keyHash > keys[index].(uint64) {
		return getNodeIndex(keys, keyHash, index+1)
	} else if index < circle.Size() && keyHash < keys[index].(uint64) {
		return index
	} else if index == circle.Size() {
		return 0
	} else {
		return 0
	}
}

func storeKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	key := p.ByName("key_id")
	value := p.ByName("value")
	url := getURLPut(key, value)

	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader(""))
	response, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
	}

	pkvr := PutKeyValueResponse{}

	json.NewDecoder(response.Body).Decode(&pkvr)

	resJson, _ := json.Marshal(pkvr)

	fmt.Println("Put Key Value Response ", pkvr)
	fmt.Fprintf(rw, "%s", resJson)

}

func getKeyValue(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	key := p.ByName("key_id")
	url := getURLGet(key)

	gkvr := KeyValue{}
	client := &http.Client{}
	response, err := client.Get(url)

	json.NewDecoder(response.Body).Decode(&gkvr)

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
	}

	resJson, _ := json.Marshal(gkvr)
	fmt.Println("Get Key Value Response ", gkvr)
	fmt.Fprintf(rw, "%s", resJson)

}

func getAllKeys(url string) {
	gakr := GetAllKeysRes{}
	client := &http.Client{}

	response, err := client.Get(url)
	json.NewDecoder(response.Body).Decode(gakr)

	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
	}

	fmt.Println("Get All Keys Response ", gakr)
}
