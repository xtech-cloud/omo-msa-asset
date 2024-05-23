package cache

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"omo.msa.asset/proxy"
)

func clipFaces(asset, key, owner, url, operator string, info *FaceResponse) error {
	size, buf, err := downloadAsset(url)
	if err != nil {
		return err
	}
	if size < 100 {
		return errors.New("the data is empty of url = " + url)
	}
	if info.Result == nil {
		return errors.New("not found the face of url = " + url)
	}
	for _, face := range info.Result.List {
		bs64, er := clipImageFace(buf, face.Location)
		if er != nil {
			return er
		}
		_, er1 := CreateThumb(asset, key, owner, bs64, operator, face)
		if er1 != nil {
			return er1
		}
	}
	return nil
}

func downloadAsset(url string) (int64, *bytes.Buffer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	l, er := buf.ReadFrom(resp.Body)
	if er != nil {
		return 0, nil, er
	}
	return l, buf, nil
}

func clipImageFace(buf *bytes.Buffer, loc proxy.LocationInfo) (string, error) {
	if buf == nil {
		return "", errors.New("the buf is nil")
	}
	origin, err := decodeImage(buf.Bytes(), int(loc.Width), int(loc.Height))
	if err != nil {
		return "", err
	}
	wid := int(loc.Width)
	hei := int(loc.Height)
	if wid == 0 || hei == 0 {
		wid = origin.Bounds().Max.X
		hei = origin.Bounds().Max.Y
	}
	//img := origin.(*image.YCbCr)
	//subImg := img.SubImage(image.Rect(int(loc.Left), int(loc.Top), wid, hei))
	//subImg := img.SubImage(image.Rect(100, 300, wid, hei))
	left := int(loc.Left) - 20
	if left < 1 {
		left = 1
	}
	top := int(loc.Top) - 60
	if top < 1 {
		top = 1
	}
	subImg := imaging.Crop(origin, image.Rect(left, top, left+wid+50, top+hei+50))
	subBuf := bytes.NewBuffer(nil)
	err = jpeg.Encode(subBuf, subImg, &jpeg.Options{100})
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(subBuf.Bytes())
	return data, nil
}

func decodeImage(bts []byte, wid, hei int) (image.Image, error) {
	reader := bytes.NewReader(bts)
	cfg, format, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, err
	}
	if cfg.Width < wid || cfg.Height < hei {
		return nil, errors.New("the image width or height is limited")
	}
	reader.Seek(0, 0)
	var img image.Image
	if format == "png" {
		img, err = png.Decode(reader)
	} else if format == "jpeg" || format == "jpg" {
		img, err = jpeg.Decode(reader)
	} else if format == "bmp" {
		img, err = bmp.Decode(reader)
	} else if format == "gif" {
		img, err = gif.Decode(reader)
	} else if format == "webp" {
		img, err = webp.Decode(reader)
	} else {
		err = errors.New("the image format not support of " + format)
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}

//* Clip 图片裁剪
//* 入参:图片输入、输出、缩略图宽、缩略图高、Rectangle{Pt(x0, y0), Pt(x1, y1)}，精度
//* 规则:如果精度为0则精度保持不变
//*
//* 返回:error
// */
func ClipImage(in io.Reader, out io.Writer, wi, hi, x0, y0, x1, y1, quality int) (data string, err error) {
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
func ScaleImage(in io.Reader, out io.Writer, width, height, quality int) (int, int, error) {
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
