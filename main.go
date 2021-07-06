package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	jsLoad := js.FuncOf(load)
	jsRender := js.FuncOf(render)
	js.Global().Set("load", jsLoad)
	js.Global().Set("render", jsRender)
	select {}
}

func load(_ js.Value, inputs []js.Value) interface{} {
	u8a := inputs[0]
	err := unpackCascade(u8a)
	if err != nil {
		fmt.Printf("初始化失败：%v\n", err)
	}
	return nil
}

func render(_ js.Value, inputs []js.Value) interface{} {
	u8a := inputs[0]
	width := inputs[1].Int()
	height := inputs[2].Int()
	imgDataMap := detectFaces(u8a, width, height)
	result := js.Global().Get("Uint8Array").New(len(imgDataMap))
	js.CopyBytesToJS(result, imgDataMap)
	return result
}