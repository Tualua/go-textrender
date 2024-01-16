package render_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/shaping"

	"golang.org/x/image/math/fixed"
)

func Test_Render(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 425, 250))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	f1, _ := font.ParseTTF(data)

	r := &render.Renderer{
		FontSize: 48,
		Color:    color.Black,
	}
	str := "Hello! ± ज्या"
	r.DrawString(str, img, f1)
	r.DrawStringAt(str, img, 0, 100, f1)

	r.PixScale = 2
	r.Color = color.Gray{Y: 0xcc}
	r.DrawStringAt("baseline", img, 0, 180, f1)

	data, _ = os.Open("testdata/NotoSans-Bold.ttf")
	f2, _ := font.ParseTTF(data)
	r.FontSize = 36
	r.Color = color.NRGBA{R: 0xcc, G: 0, B: 0x33, A: 0x99}
	x := r.DrawStringAt("Red", img, 60, 140, f2)
	r.DrawStringAt("Bold", img, x, 140, f2)

	// from https://github.com/adobe-fonts/emojione-color, MIT license
	data, _ = os.Open("testdata/EmojiOneColor.otf")
	f3, _ := font.ParseTTF(data)
	r.FontSize = 36
	r.DrawStringAt("🚀🖥️", img, 270, 80, f3)

	data, _ = os.Open("testdata/Greybeard-22px.ttf")
	f4, _ := font.ParseTTF(data)
	r.FontSize = 22
	r.Color = color.NRGBA{R: 0xcc, G: 0x66, B: 0x33, A: 0xcc}
	r.DrawStringAt("\uE0A2░", img, 366, 164, f4)

	str = "Hello ज्या 😀! 🎁 fin."
	rs := []rune(str)
	sh := &shaping.HarfbuzzShaper{}
	in := shaping.Input{
		Text:     rs,
		RunStart: 0,
		RunEnd:   len(rs),
		Size:     fixed.I(int(r.FontSize)),
	}
	faces := []font.Face{f1, f2, f3}
	fontRuns := shaping.SplitByFontGlyphs(in, faces)
	seg := shaping.Segmenter{}

	x = 0
	for _, funtRun := range fontRuns {
		runs := seg.Split(funtRun, singleFontMap{funtRun.Face})

		line := make(shaping.Line, len(runs))
		r.Color = color.NRGBA{R: 0x33, G: 0x99, B: 0x33, A: 0xcc}

		for i, run := range runs {
			line[i] = sh.Shape(run)
		}

		for _, run := range line {
			x = r.DrawShapedRunAt(run, img, x, 232)
		}
	}

	w, _ := os.Create("testdata/out.png")
	png.Encode(w, img)
	w.Close()
}

func TestRender_PixScaleAdvance(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 350, 180))

	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	f, _ := font.ParseTTF(data)

	r := &render.Renderer{
		FontSize: 48,
		Color:    color.Black,
	}
	str := "Testing"
	adv0 := r.DrawString(str, img, f)

	r.PixScale = 1 // instead of the zero value
	adv1 := r.DrawString(str, img, f)
	if adv0 != adv1 {
		t.Error("unscaled font did not advance as default")
	}

	r.PixScale = 2
	adv2 := r.DrawString(str, img, f)
	if adv2 <= int(float32(adv1)*1.9) || adv2 >= int(float32(adv1)*2.1) {
		t.Error("scaled font did not advance proportionately")
	}
}

func TestRenderHindi(t *testing.T) {
	text := "नमस्ते"
	r := &render.Renderer{
		FontSize: 30,
		Color:    color.Black,
	}

	img := image.NewNRGBA(image.Rect(0, 0, 120, 50))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	data, _ := os.Open("testdata/NotoSans-Regular.ttf")
	face, _ := font.ParseTTF(data)

	r.DrawString(text, img, face)

	w, _ := os.Create("testdata/out_hindi.png")
	png.Encode(w, img)
	w.Close()
}

type singleFontMap struct {
	face font.Face
}

func (sf singleFontMap) ResolveFace(rune) font.Face { return sf.face }
