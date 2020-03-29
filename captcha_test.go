/*
 * @Package captcha
 * @Author Quan Chen
 * @Date 2020/3/29
 * @Description
 *
 */
package captcha

import (
	"bytes"
	"image"
	"image/color"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestNewCaptchaImage(t *testing.T) {
	type args struct {
		width   int
		height  int
		bgColor color.RGBA
	}
	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		{
			"1",
			args{
				width:   20,
				height:  20,
				bgColor: color.RGBA{R: 224, G: 224, B: 224},
			},
			&Image{
				image: image.NewNRGBA(image.Rect(0, 0, 20, 20)),
				width: 20, height: 20,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCaptchaImage(tt.args.width, tt.args.height, tt.args.bgColor)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCaptchaImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCaptchaImage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_SaveImage(t *testing.T) {
	type fields struct {
		image   *image.NRGBA
		width   int
		height  int
		Complex int
	}
	type args struct {
		imageFormat int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"1",
			fields{
				image:   image.NewNRGBA(image.Rect(0, 0, 20, 20)),
				width:   0,
				height:  0,
				Complex: 0,
			},
			args{imageFormat: 0},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captcha := &Image{
				image:   tt.fields.image,
				width:   tt.fields.width,
				height:  tt.fields.height,
				Complex: tt.fields.Complex,
			}
			w := &bytes.Buffer{}
			err := captcha.SaveImage(w, tt.args.imageFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			err = ioutil.WriteFile("image.png", w.Bytes(), 0777)
			if err != nil {
				t.Errorf("SaveImage() fail error: %v", err)
			}
		})
	}
}

func TestImage_DrawHollowLine(t *testing.T) {
	type fields struct {
		image   *image.NRGBA
		width   int
		height  int
		Complex int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1",
			fields{
				image:   image.NewNRGBA(image.Rect(0, 0, 200, 200)),
				width:   200,
				height:  200,
				Complex: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captcha := &Image{
				image:   tt.fields.image,
				width:   tt.fields.width,
				height:  tt.fields.height,
				Complex: tt.fields.Complex,
			}
			im := captcha.DrawHollowLine()
			w := &bytes.Buffer{}
			err := im.SaveImage(w, 0)
			if err != nil {
				t.Errorf("SaveImage() fail error: %v", err)
			}
			err = ioutil.WriteFile("image.png", w.Bytes(), 0777)
			if err != nil {
				t.Errorf("SaveImage() fail error: %v", err)
			}
		})
	}
}
func TestImage_DrawSineLine(t *testing.T) {
	type fields struct {
		image   *image.NRGBA
		width   int
		height  int
		Complex int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"1",
			fields{
				image:   image.NewNRGBA(image.Rect(0, 0, 200, 200)),
				width:   200,
				height:  200,
				Complex: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			captcha := &Image{
				image:   tt.fields.image,
				width:   tt.fields.width,
				height:  tt.fields.height,
				Complex: tt.fields.Complex,
			}
			im := captcha.DrawSineLine()
			w := &bytes.Buffer{}
			err := im.SaveImage(w, 0)
			if err != nil {
				t.Errorf("SaveImage() fail error: %v", err)
			}
			err = ioutil.WriteFile("sineImage.png", w.Bytes(), 0777)
			if err != nil {
				t.Errorf("SaveImage() fail error: %v", err)
			}
		})
	}
}
