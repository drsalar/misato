package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Conf struct {
	Pos       Position `json:"position"`
	Siz       Size     `json:"size"`
	Mode      int      `json:"mode"`
	Alpha     int      `json:"alpha"`
	Watermark string   `json:"watermark"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

var dir string
var conf Conf

func main() {
	confInit()

	var err error
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal("error in main:", err.Error())
		return
	}

	filesWalker()
	test()
}

func filesWalker() {
	err := filepath.Walk(dir+"/raw", run)
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func run(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}

	fmt.Println("open file:", path)

	if filepath.Ext(path) == ".jpg" {
		jpgHandler(path)
	} else if filepath.Ext(path) == ".png" {
		// pngHandler(path)
	}

	return nil
}

func jpgHandler(path string) {
	rmgb, _ := os.Open(path)
	rimg, _ := jpeg.Decode(rmgb)
	defer rmgb.Close()

	wmb, _ := os.Open(conf.Watermark)
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	//把水印写到右下角，并向0坐标各偏移10个像素
	offset := image.Pt(rimg.Bounds().Dx()-watermark.Bounds().Dx()-10, rimg.Bounds().Dy()-watermark.Bounds().Dy()-10)
	b := rimg.Bounds()
	m := image.NewNRGBA(b)

	yy := image.NewRGBA(image.Rect(0, 0, watermark.Bounds().Dx(), watermark.Bounds().Dy()))
	for x := 0; x < yy.Bounds().Dx(); x++ {
		for y := 0; y < yy.Bounds().Dy(); y++ {
			yy.SetRGBA(x, y, color.RGBA{255, 255, 255, uint8(conf.Alpha)})
		}
	}

	draw.Draw(m, b, rimg, image.ZP, draw.Src)
	// draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)
	draw.DrawMask(m, watermark.Bounds().Add(offset), watermark, image.ZP, yy, image.ZP, draw.Over)

	//生成新图片new.jpg，并设置图片质量..
	imgw, _ := os.Create(dir + "/new/" + filepath.Base(path))
	jpeg.Encode(imgw, m, &jpeg.Options{100})

	defer imgw.Close()

	fmt.Println("水印添加结束,请查看new.jpg图片...")
}

func confInit() {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("error in confInit:", err.Error())
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal("error in confInit:", err.Error())
	}
}

func test() {
	x, _ := os.Open("xxxx.png")
	xx, _ := png.Decode(x)
	defer x.Close()

	d, _ := os.Open("tttt.jpg")
	dd, _ := jpeg.Decode(d)
	defer d.Close()

	offset := image.Pt(dd.Bounds().Dx()-xx.Bounds().Dx()-10, dd.Bounds().Dy()-xx.Bounds().Dy()-10)

	mm := image.NewRGBA(dd.Bounds())
	yy := image.NewRGBA(image.Rect(0, 0, xx.Bounds().Dx(), xx.Bounds().Dy()))
	for x := 0; x < yy.Bounds().Dx(); x++ {
		for y := 0; y < yy.Bounds().Dy(); y++ {
			yy.SetRGBA(x, y, color.RGBA{255, 255, 255, uint8(conf.Alpha)})
		}
	}

	draw.Draw(mm, dd.Bounds(), dd, image.ZP, draw.Src)
	// draw.Draw(mm, xx.Bounds().Add(offset), xx, image.ZP, draw.Over)
	draw.DrawMask(mm, xx.Bounds().Add(offset), xx, image.ZP, yy, image.ZP, draw.Over)

	p, _ := os.Create("ssss.png")
	png.Encode(p, mm)
}
