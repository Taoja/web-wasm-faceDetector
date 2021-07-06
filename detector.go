package main

import (
	"errors"
	"fmt"
	pigo "github.com/esimov/pigo/core"
	"math"
	"syscall/js"
)

var (
	faceClassifier   *pigo.Pigo // 脸部识别器
	imgParams        *pigo.ImageParams
	err              error
)

// unpackCascade 解压读取cascade
// @params uint8Array js.Uint8Array casacade文件buffer
func unpackCascade (uint8Array js.Value) error {
	p := pigo.NewPigo()
	
	buff := make([]byte, uint8Array.Get("length").Int())
	js.CopyBytesToGo(buff, uint8Array)
	faceClassifier, err = p.Unpack(buff)
	if err != nil {
		return errors.New("error unpacking the facefinder cascade file")
	}
	return nil
}

// faceDetector 人脸识别核心
// @params pixels []uint8 图片uint8格式
// @params width int 图片宽度
// @params height int 图片宽度
// @return 返回识别内容 []pigo.Detection
// type Detection struct {
//	Row   int
//	Col   int
//	Scale int
//	Q     float32
// }
func faceDetector (pixels []uint8, width, height int) []pigo.Detection {
	imgParams = &pigo.ImageParams{
		Pixels: pixels,
		Rows:   width,
		Cols:   height,
		Dim:    height,
	}
	cParams := pigo.CascadeParams{
		MinSize:     200,
		MaxSize:     640,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: *imgParams,
	}
	
	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := faceClassifier.RunCascade(cParams, 0.0)
	
	// Calculate the intersection over union (IoU) of two clusters.
	dets = faceClassifier.ClusterDetections(dets, 0.1)
	
	return dets
}

// detectFaces 人脸识别方法
// @params pixels []uint8 图片uint8格式
// @params width int 图片宽度
// @params height int 图片宽度
// @return x,y像素点阵 [][]int
func detectFaces(uint8Array js.Value, width, height int) []uint8 {
	var buff = make([]byte, uint8Array.Get("length").Int())
	
	js.CopyBytesToGo(buff, uint8Array)
	pixels := rgbaToGrayscale(buff, width, height)
	results := faceDetector(pixels, width, height)
	fmt.Printf("result = %v\n", results)
	//dets := make([][]int, len(results))
	//
	//for i := 0; i < len(results); i++ {
	//	dets[i] = append(dets[i], results[i].Row, results[i].Col, results[i].Scale, int(results[i].Q))
	//}
	return pixels
}

// rgbaToGrayscale rgba转灰度图片
func rgbaToGrayscale(data []uint8, width, height int) []uint8 {
	rows, cols := width, height
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// gray = 0.2*red + 0.7*green + 0.1*blue
			data[r*cols+c] = uint8(math.Round(
				0.2126*float64(data[r*4*cols+4*c+0]) +
					0.7152*float64(data[r*4*cols+4*c+1]) +
					0.0722*float64(data[r*4*cols+4*c+2])))
		}
	}
	return data
}