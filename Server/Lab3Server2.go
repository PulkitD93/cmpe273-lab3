package main
import (
"fmt"
"encoding/json"
  "github.com/julienschmidt/httprouter"
  "net/http"
  //"os"
  "strconv")

type KeyVal struct {
  ID int
  Value string
}
var m map[int]string = make(map[int]string)
func create(rw http.ResponseWriter, req *http.Request, p httprouter.Params ){
    keyID:= p.ByName("key_id")
    value:= p.ByName("value")
    keyInt,_:=strconv.Atoi(keyID)
    m[keyInt]=value


    rw.WriteHeader(201)
    
}
func getkeyval(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
   
   key_id:=p.ByName("key_id")
   keyIdInt,_:=strconv.Atoi(key_id)
   value := m[keyIdInt]
   getRes:= KeyVal {}
   getRes.ID=keyIdInt
   getRes.Value = value

   resJson, _ := json.Marshal(getRes)
    rw.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(rw, "%s", resJson) 
}

func getAll(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    getRes:= []KeyVal {}
    mapStruct:= KeyVal{}
    for i,j:= range m{
    mapStruct.ID= i
    mapStruct.Value=j

    getRes= append(getRes, mapStruct)
}
    resJson, _ := json.Marshal(getRes)
    rw.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(rw, "%s", resJson) 
}

func main() {
  mux := httprouter.New()
  mux.PUT("/keys/:key_id/:value", create)
  mux.GET("/keys/:key_id", getkeyval )
  mux.GET("/keys", getAll )
  server := http.Server{
    Addr:    "0.0.0.0:3001",
    Handler: mux,
  }
  server.ListenAndServe()
}

