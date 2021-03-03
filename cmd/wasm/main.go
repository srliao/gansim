package main

import (
	"encoding/json"
	"log"
	"syscall/js"
)

func main() {
	done := make(chan struct{}, 0)
	global := js.Global()
	global.Set("wasmPrint", js.FuncOf(runsim))
	<-done
}

type probconfig struct {
}

func runsim(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		log.Println("error, invalid number of args")
		return "ERROR: invalid number of args"
	}
	var ps = args[0].String()
	var p profile
	err := json.Unmarshal([]byte(ps), &p)
	if err != nil {
		log.Println(err)
		return "ERROR parsing"
	}
	calc(p, artifactSet{}, true)
	return ""
}
