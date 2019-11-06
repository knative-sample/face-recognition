package manager

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"path/filepath"

	"image/draw"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
)

func Mark(fontPath, orginPath, targetPath string, fa *FaceAttribute) error {
	//需要加水印的图片
	imgfile, err := os.Open(orginPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer imgfile.Close()

	jpgimg, err := jpeg.Decode(imgfile)
	if err != nil {
		log.Println(err)
		return err
	}
	img := image.NewNRGBA(jpgimg.Bounds())

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, jpgimg.At(x, y))
		}
	}
	//拷贝一个字体文件到运行目录
	fontBytes, err := ioutil.ReadFile(filepath.Join(fontPath, "simsun.ttc"))
	if err != nil {
		log.Println(err)
		return err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return err
	}
	maleCount := 0
	femaleCount := 0
	for i := 0; i < fa.FaceNum; i++ {
		f := freetype.NewContext()
		f.SetDPI(72)
		f.SetFont(font)
		f.SetFontSize(120)
		f.SetClip(jpgimg.Bounds())
		f.SetDst(img)
		pos := i * 4
		x := fa.FaceRect[pos]
		y := fa.FaceRect[pos+1]

		pt := freetype.Pt(x+10, y+20)
		sex := ""
		if fa.Gender[i] == 0 {
			sex = "F"
			f.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 255, B: 0, A: 255}))
			femaleCount++
		} else if fa.Gender[i] == 1 {
			sex = "M"
			f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 255, B: 0, A: 255}))
			maleCount++
		}
		_, err = f.DrawString(fmt.Sprintf("%s", sex), pt)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	imgRGBA, err := GetHeadImageRGBA(filepath.Join(fontPath, "bg.jpg"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	x0, y0 := jpgimg.Bounds().Dx()-1600, 10
	draw.DrawMask(img, image.Rect(x0, y0, x0+1300, y0+500), imgRGBA, image.ZP, nil, image.ZP, draw.Over)
	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(150)
	f.SetClip(jpgimg.Bounds())
	f.SetDst(img)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
	dx := jpgimg.Bounds().Dx() - 1500
	if dx < 0 {
		dx = 0
	}
	pt := freetype.Pt(dx, 150)
	_, err = f.DrawString(fmt.Sprintf("Male:Female=%v:%v", maleCount, femaleCount), pt)
	if err != nil {
		log.Println(err)
		return err
	}
	//保存到新文件中
	newFile, err := os.Create(targetPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer newFile.Close()

	err = jpeg.Encode(newFile, img, &jpeg.Options{100})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func GetHeadImageRGBA(iamgePath string) (*image.RGBA, error) {
	img, err := imaging.Open(iamgePath)
	if err != nil {
		return nil, err
	}
	imgRGBA := image.NewRGBA(image.Rect(0, 0, 2000, 190))
	err = graphics.Scale(imgRGBA, img)
	if err != nil {
		return nil, err
	}
	return imgRGBA, nil
}
