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
	Margin    Size   `json:"margin"`
	Siz       Size   `json:"size"`
	Mode      int    `json:"mode"`
	Alpha     int    `json:"alpha"`
	Watermark string `json:"watermark"`
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

	imgHandler(path)

	return nil
}

func imgHandler(path string) {
	rmgb, _ := os.Open(path)

	var rimg image.Image
	if filepath.Ext(path) == ".jpg" {
		rimg, _ = jpeg.Decode(rmgb)
	} else if filepath.Ext(path) == ".png" {
		rimg, _ = png.Decode(rmgb)
	}
	defer rmgb.Close()

	wmb, _ := os.Open(conf.Watermark)
	watermark, _ := png.Decode(wmb)
	defer wmb.Close()

	// offset := image.Pt(rimg.Bounds().Dx()-watermark.Bounds().Dx()-10, rimg.Bounds().Dy()-watermark.Bounds().Dy()-10)
	offset := calOffect(rimg.Bounds().Dx(), rimg.Bounds().Dy(), watermark.Bounds().Dx(), watermark.Bounds().Dy())
	b := rimg.Bounds()
	m := image.NewNRGBA(b)

	yy := image.NewRGBA(image.Rect(0, 0, watermark.Bounds().Dx(), watermark.Bounds().Dy()))
	for x := 0; x < yy.Bounds().Dx(); x++ {
		for y := 0; y < yy.Bounds().Dy(); y++ {
			yy.SetRGBA(x, y, color.RGBA{255, 255, 255, uint8(conf.Alpha)})
		}
	}

	draw.Draw(m, b, rimg, image.ZP, draw.Src)
	draw.DrawMask(m, watermark.Bounds().Add(offset), watermark, image.ZP, yy, image.ZP, draw.Over)

	imgw, _ := os.Create(dir + "/new/" + filepath.Base(path))
	jpeg.Encode(imgw, m, &jpeg.Options{100})

	defer imgw.Close()

	fmt.Println("水印添加结束,请查看new.jpg图片...")
}

func calOffect(x1, y1 int, x2, y2 int) (pt image.Point) {
	switch conf.Mode {
	case TOP_LEFT:
		pt = image.Pt(conf.Margin.W, conf.Margin.H)
		return
	case BOTTOM_RIGHT:
		pt = image.Pt(x1-x2-conf.Margin.W, y1-y2-conf.Margin.H)
		return
	case TOP_CENTER:
		pt = image.Pt((x1-x2)/2, conf.Margin.H)
		return
	case BOTTOM_CENTER:
		pt = image.Pt((x1-x2)/2, y1-y2-conf.Margin.H)
		return
	case TILE:
		//TODO
	case STRETCH:
		//TODO
	default:
		//TODO
	}
	pt = image.Pt(conf.Margin.W, conf.Margin.H)
	return
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
