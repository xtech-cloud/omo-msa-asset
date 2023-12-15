package cache

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
)

//补上缺失的代码
//* Clip 图片裁剪
//* 入参:图片输入、输出、缩略图宽、缩略图高、Rectangle{Pt(x0, y0), Pt(x1, y1)}，精度
//* 规则:如果精度为0则精度保持不变
//*
//* 返回:error
// */
func Clip(in io.Reader, out io.Writer, wi, hi, x0, y0, x1, y1, quality int) (data string, err error) {
	err = errors.New("unknow error")
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var origin image.Image
	var fm string
	origin, fm, err = image.Decode(in)
	if err != nil {
		log.Println(err)
		return data, err
	}

	if wi == 0 || hi == 0 {
		wi = origin.Bounds().Max.X
		hi = origin.Bounds().Max.Y
	}
	var canvas image.Image
	if wi != origin.Bounds().Max.X {
		//先缩略
		canvas = resize.Thumbnail(uint(wi), uint(hi), origin, resize.Lanczos3)
	} else {
		canvas = origin
	}

	switch fm {
	case "jpeg":
		img := canvas.(*image.YCbCr)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.YCbCr)
		buf := bytes.NewBuffer(nil)
		_ = png.Encode(buf, subImg)
		data = base64.StdEncoding.EncodeToString(buf.Bytes())
		return data, jpeg.Encode(out, subImg, &jpeg.Options{quality})
	case "png":
		switch canvas.(type) {
		case *image.NRGBA:
			img := canvas.(*image.NRGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.NRGBA)
			return "", png.Encode(out, subImg)
		case *image.RGBA:
			img := canvas.(*image.RGBA)
			subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
			return "", png.Encode(out, subImg)
		}
	case "gif":
		img := canvas.(*image.Paletted)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.Paletted)
		return "", gif.Encode(out, subImg, &gif.Options{})
	case "bmp":
		img := canvas.(*image.RGBA)
		subImg := img.SubImage(image.Rect(x0, y0, x1, y1)).(*image.RGBA)
		return "", bmp.Encode(out, subImg)
	default:
		return data, errors.New("ERROR FORMAT")
	}
	return data, err
}

/*
* Scale 缩略图生成
* 入参:图片输入、输出，缩略图宽、高，精度
* 规则: 如果width 或 hight其中有一个为0，则大小不变 如果精度为0则精度保持不变
* 返回:缩略图真实宽、高、error
 */
func Scale(in io.Reader, out io.Writer, width, height, quality int) (int, int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	var (
		w, h int
	)
	origin, fm, err := image.Decode(in)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	if width == 0 || height == 0 {
		width = origin.Bounds().Max.X
		height = origin.Bounds().Max.Y
	}
	if quality == 0 {
		quality = 100
	}
	canvas := resize.Thumbnail(uint(width), uint(height), origin, resize.Lanczos3)

	//return jpeg.Encode(out, canvas, &jpeg.Options{quality})
	w = canvas.Bounds().Dx()
	h = canvas.Bounds().Dy()
	switch fm {
	case "jpeg":
		return w, h, jpeg.Encode(out, canvas, &jpeg.Options{quality})
	case "png":
		return w, h, png.Encode(out, canvas)
	case "gif":
		return w, h, gif.Encode(out, canvas, &gif.Options{})
	//case "bmp":  //被我注释掉的是x/image/bmp
	//	return w, h, bmp.Encode(out, canvas)
	default:
		return w, h, errors.New("ERROR FORMAT")
	}
	return w, h, nil
}
