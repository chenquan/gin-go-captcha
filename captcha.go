/*
 * @Package captcha
 * @Author Quan Chen
 * @Date 2020/3/28
 * @Description
 *
 */
package captcha

import (
	"errors"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	dpi        = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	r          = rand.New(rand.NewSource(time.Now().UnixNano()))
	fontFamily = make([]string, 0)
)

const (
	//图片格式
	ImageFormatPng  = iota // PNG
	ImageFormatJpeg        // JEP
	ImageFormatGif         //GIF
	//验证码噪点强度
	CaptchaComplexLower  = iota // 低
	CaptchaComplexMedium        // 中
	CaptchaComplexHigh          // 高
)

type Point struct {
	X int
	Y int
}

func NewPoint(x int, y int) *Point {
	return &Point{X: x, Y: y}
}

type Image struct {
	image   *image.NRGBA
	width   int
	height  int
	Complex int
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ReadFonts(dirPth string, suffix string) (err error) {
	files := make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	SetFontFamily(files...)
	return nil
}

//添加一个字体路径到字体库.
func SetFontFamily(fontPath ...string) {

	fontFamily = append(fontFamily, fontPath...)
}

//新建一个图片对象
func NewCaptchaImage(width int, height int, bgColor color.RGBA) (*Image, error) {

	m := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	return &Image{
		image:  m,
		height: height,
		width:  width,
	}, nil
}

//保存图片对象
func (captcha *Image) SaveImage(w io.Writer, imageFormat int) error {
	switch imageFormat {
	case ImageFormatPng:
		return png.Encode(w, captcha.image)
	case ImageFormatJpeg:
		return jpeg.Encode(w, captcha.image, &jpeg.Options{Quality: 100})
	case ImageFormatGif:
		return gif.Encode(w, captcha.image, &gif.Options{NumColors: 256})
	default:
		return errors.New("Not supported image format.")
	}
}

//添加一个较粗的空白线
func (captcha *Image) DrawHollowLine() *Image {

	first := captcha.width / 20
	end := first * 19

	// 线的颜色[随机]
	lineColor := RandLightColor()
	// 起始横坐标
	x1 := float64(r.Intn(first))
	x2 := float64(r.Intn(first) + end)

	multiple := float64(r.Intn(5)+3) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple = multiple * -1.0
	}

	w := captcha.height / 20

	for ; x1 < x2; x1++ {

		y := math.Sin(x1*math.Pi*multiple/float64(captcha.width)) * float64(captcha.height/3)

		if multiple < 0 {
			y = y + float64(captcha.height/2)
		}
		captcha.image.Set(int(x1), int(y), lineColor)

		for i := 0; i <= w; i++ {
			captcha.image.Set(int(x1), int(y)+i, lineColor)
		}
	}

	return captcha
}

//添加一个较粗的正弦线
func (captcha *Image) DrawSineLine() *Image {
	x := 0
	var y float64 = 0

	//振幅
	a := r.Intn(captcha.height / 2)

	//Y轴方向偏移量
	b := Random(int64(-captcha.height/4), int64(captcha.height/4))

	//X轴方向偏移量
	f := Random(int64(-captcha.height/4), int64(captcha.height/4))

	// 周期
	var t float64 = 0
	if captcha.height > captcha.width/2 {
		t = Random(int64(captcha.width/2), int64(captcha.height))
	} else {
		t = Random(int64(captcha.height), int64(captcha.width/2))
	}
	w := (2 * math.Pi) / t

	// 曲线横坐标起始位置
	px1 := 0
	px2 := int(Random(int64(float64(captcha.width)*0.8), int64(captcha.width)))
	// 随机颜色
	c := RandDarkColor()

	for x = px1; x < px2; x++ {
		if w != 0 {
			y = float64(a)*math.Sin(w*float64(x)+f) + b + (float64(captcha.width) / float64(5))
			i := captcha.height / 5
			for i > 0 {
				captcha.image.Set(x+i, int(y), c)
				i--
			}
		}
	}

	return captcha
}

//画一条直线
func (captcha *Image) DrawLine(num int) *Image {

	first := captcha.width / 10
	end := first * 9

	y := captcha.height / 3

	for i := 0; i < num; i++ {

		point1 := Point{X: r.Intn(first), Y: r.Intn(y)}
		point2 := Point{X: r.Intn(first) + end, Y: r.Intn(y)}

		if i%2 == 0 {
			point1.Y = r.Intn(y) + y*2
			point2.Y = r.Intn(y)
		} else {
			point1.Y = r.Intn(y) + y*(i%2)
			point2.Y = r.Intn(y) + y*2
		}

		captcha.drawBeeline(point1, point2, RandDarkColor())

	}
	return captcha
}

func (captcha *Image) drawBeeline(point1 Point, point2 Point, lineColor color.RGBA) {
	dx := math.Abs(float64(point1.X - point2.X))

	dy := math.Abs(float64(point2.Y - point1.Y))
	sx, sy := 1, 1
	if point1.X >= point2.X {
		sx = -1
	}
	if point1.Y >= point2.Y {
		sy = -1
	}
	err := dx - dy
	for {
		captcha.image.Set(point1.X, point1.Y, lineColor)
		captcha.image.Set(point1.X+1, point1.Y, lineColor)
		captcha.image.Set(point1.X-1, point1.Y, lineColor)
		captcha.image.Set(point1.X+2, point1.Y, lineColor)
		captcha.image.Set(point1.X-2, point1.Y, lineColor)
		if point1.X == point2.X && point1.Y == point2.Y {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			point1.X += sx
		}
		if e2 < dx {
			err += dx
			point1.Y += sy
		}
	}
}

//画边框.
func (captcha *Image) DrawBorder(borderColor color.RGBA) *Image {
	for x := 0; x < captcha.width; x++ {
		captcha.image.Set(x, 0, borderColor)
		captcha.image.Set(x, captcha.height-1, borderColor)
	}
	for y := 0; y < captcha.height; y++ {
		captcha.image.Set(0, y, borderColor)
		captcha.image.Set(captcha.width-1, y, borderColor)
	}
	return captcha
}

//画噪点.
func (captcha *Image) DrawNoise(complex int) *Image {
	density := 18
	if complex == CaptchaComplexLower {
		density = 28
	} else if complex == CaptchaComplexMedium {
		density = 18
	} else if complex == CaptchaComplexHigh {
		density = 8
	}
	maxSize := (captcha.height * captcha.width) / density

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		captcha.image.Set(rw, rh, RandColor())
		size := r.Intn(maxSize)
		if size%3 == 0 {
			captcha.image.Set(rw+1, rh+1, RandColor())
		}
	}
	return captcha
}

//画文字噪点.
func (captcha *Image) DrawTextNoise(complex int) error {
	density := 1500
	if complex == CaptchaComplexLower {
		density = 2000
	} else if complex == CaptchaComplexMedium {
		density = 1500
	} else if complex == CaptchaComplexHigh {
		density = 1000
	}

	maxSize := (captcha.height * captcha.width) / density

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.image.Bounds())
	c.SetDst(captcha.image)
	c.SetHinting(font.HintingFull)
	rawFontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(10))

	for i := 0; i < maxSize; i++ {

		rw := r.Intn(captcha.width)
		rh := r.Intn(captcha.height)

		text := RandText(1)
		fontSize := rawFontSize/2 + float64(r.Intn(5))

		c.SetSrc(image.NewUniform(RandLightColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			log.Println(err)
			return err
		}
		c.SetFont(f)
		pt := freetype.Pt(rw, rh)

		_, err = c.DrawString(text, pt)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

//写字.
func (captcha *Image) DrawText(text string) error {
	c := freetype.NewContext()
	c.SetDPI(*dpi)

	c.SetClip(captcha.image.Bounds())
	c.SetDst(captcha.image)
	c.SetHinting(font.HintingFull)

	fontWidth := captcha.width / len(text)

	for i, s := range text {

		fontSize := float64(captcha.height) / (1 + float64(r.Intn(7))/float64(9))

		c.SetSrc(image.NewUniform(RandDarkColor()))
		c.SetFontSize(fontSize)
		f, err := RandFontFamily()

		if err != nil {
			log.Println(err)
			return err
		}
		c.SetFont(f)

		x := int(fontWidth)*i + int(fontWidth)/int(fontSize)

		y := 5 + r.Intn(captcha.height/2) + int(fontSize/2)

		pt := freetype.Pt(x, y)

		_, err = c.DrawString(string(s), pt)
		if err != nil {
			log.Println(err)
			return err
		}
		//pt.Y += c.PointToFixed(*size * *spacing)
		//pt.X += c.PointToFixed(*size);
	}
	return nil

}

//获取所有字体.
func RandFontFamily() (*truetype.Font, error) {
	fontFile := fontFamily[r.Intn(len(fontFamily))]

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return &truetype.Font{}, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return &truetype.Font{}, err
	}
	return f, nil
}

//随机生成深色系.
func RandDarkColor() color.RGBA {

	randColor := RandColor()

	increase := float64(30 + r.Intn(255))

	red := math.Abs(math.Min(float64(randColor.R)-increase, 255))

	green := math.Abs(math.Min(float64(randColor.G)-increase, 255))
	blue := math.Abs(math.Min(float64(randColor.B)-increase, 255))

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//随机生成浅色.
func RandLightColor() color.RGBA {

	red := r.Intn(55) + 200
	green := r.Intn(55) + 200
	blue := r.Intn(55) + 200

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

//生成随机颜色.
func RandColor() color.RGBA {

	red := r.Intn(255)
	green := r.Intn(255)
	blue := r.Intn(255)
	if (red + green) > 400 {
		blue = 0
	} else {
		blue = 400 - green - red
	}
	if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}
