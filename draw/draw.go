package draw

import (
	"bytes"
	"fmt"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"

	"image"

	"get.cutie.cafe/levelup/fetch"
	"get.cutie.cafe/levelup/types"
)

func openFont(font string) (*truetype.Font, error) {
	file, err := os.Open(font)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fontBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return truetype.Parse(fontBytes)
}

func calcRuneWidth(graphics *draw2dimg.GraphicContext, font *truetype.Font, r rune) float64 {
	return float64((font.HMetric(fixed.Int26_6(graphics.GetFontSize()*float64(graphics.DPI)*(64.0/72.0)), font.Index(r)).AdvanceWidth / 64))
}

func drawCapsule(graphics *draw2dimg.GraphicContext, capsule image.Image) {
	graphics.MoveTo(0, 0)
	graphics.DrawImage(capsule)
}

func calcTextWidth(graphics *draw2dimg.GraphicContext, font *truetype.Font, text string) float64 {
	fw := float64(0)

	for i := 0; i < len(text); i++ {
		fw += calcRuneWidth(graphics, font, rune(text[i]))
	}

	return fw
}

func drawText(graphics *draw2dimg.GraphicContext, x float64, y float64, text string, preferredFontSize float64, imageWidth float64) {
	fw := float64(9999)
	font, _ := graphics.FontCache.Load(draw2d.FontData{Name: "font"})

	preferredFontSize++

	for fw >= imageWidth && preferredFontSize > 0 {
		preferredFontSize--
		graphics.SetFontSize(preferredFontSize)

		fw = calcTextWidth(graphics, font, text)
	}

	graphics.SetFontSize(preferredFontSize - 1)

	graphics.FillStringAt(text, x, y)
}

// Stolen from draw2dkit
func fillRectangle(path draw2d.GraphicContext, x1, y1, x2, y2 float64) {
	path.BeginPath()
	path.MoveTo(x1, y1)
	path.LineTo(x2, y1)
	path.LineTo(x2, y2)
	path.LineTo(x1, y2)
	path.Close()

	pth := path.GetPath()

	path.Fill(&pth)
}

// Draw an app widget.
func Draw(app *types.App) (*image.RGBA, error) {
	// establish font cache, see https://github.com/llgcode/draw2d/issues/127#issuecomment-267845074
	font, err := openFont("font.ttf")

	fontCache := FontCache{}
	fontCache.Store(draw2d.FontData{Name: "font"}, font)
	draw2d.SetFontCache(fontCache)

	// create image
	image := image.NewRGBA(image.Rect(0, 0, 460, 335))

	capsuleBytes, err := fetch.HTTPGet(fmt.Sprintf("https://cdn.cloudflare.steamstatic.com/steam/apps/%d/header.jpg?t=%d", app.AppID, time.Now().Unix()))
	if err != nil {
		return nil, err
	}

	capsule, err := jpeg.Decode(bytes.NewReader(capsuleBytes))
	if err != nil {
		return nil, err
	}

	graphics := draw2dimg.NewGraphicContext(image)
	defer graphics.Close()

	// load font
	if err != nil {
		return nil, err
	}

	graphics.SetFontData(draw2d.FontData{
		Name: "font",
	})

	// initial fill/stroke ops
	graphics.SetFillColor(color.RGBA{255, 255, 255, 255})
	graphics.SetStrokeColor(color.RGBA{0, 0, 0, 255})

	graphics.SetFillRule(draw2d.FillRuleWinding)

	// draw capsule
	drawCapsule(graphics, capsule)

	// draw text
	imgWidth := float64(image.Bounds().Size().X)

	drawText(graphics, 0, 250, app.Name, 27, imgWidth)

	nextOffset := 0

	if len(app.Developers) > 0 || len(app.Publishers) > 0 {
		nextOffset += 25

		str := ""

		if len(app.Developers) > 0 && app.Developers[0] != app.Publishers[0] {
			str += app.Developers[0] + "; "
		}

		if len(app.Publishers) > 0 {
			str += app.Publishers[0]
		}

		drawText(graphics, float64(0), float64(250+nextOffset), str, 15, imgWidth)
	}

	prevWidth := float64(0)

	// genres text
	for _, g := range app.Genres {
		graphics.SetFillColor(color.RGBA{0x17, 0xa2, 0xb8, 255})

		graphics.SetFontSize(10)

		width := calcTextWidth(graphics, font, g)

		fillRectangle(graphics, prevWidth, float64(250+nextOffset+8), prevWidth+width+6, float64(250+nextOffset+20+8))

		graphics.SetFillColor(color.RGBA{255, 255, 255, 255})

		drawText(graphics, prevWidth+4, float64(250+nextOffset+20)+2, g, 10, 9999)

		prevWidth += width + 10
	}

	// price text
	graphics.SetFillColor(color.RGBA{0x4c, 0x6b, 0x22, 255})

	graphics.SetFontSize(10)

	priceText := ""
	if app.Prices.Steam["us"].OriginalPrice == 0 {
		priceText = "Free to Play"
	} else {
		priceText = "$" + fmt.Sprintf("%.2f", (float64(app.Prices.Steam["us"].OriginalPrice)/100.0)) + " USD"
	}

	width := calcTextWidth(graphics, font, priceText)

	fillRectangle(graphics, prevWidth, float64(250+nextOffset+8), prevWidth+width+8, float64(250+nextOffset+20+8))

	graphics.SetFillColor(color.RGBA{0xa4, 0xd0, 0x07, 255})
	drawText(graphics, prevWidth+4, float64(250+nextOffset+20)+2, priceText, 10, 9999)

	// draw bottom footer
	graphics.SetFillColor(color.RGBA{255, 255, 255, 255})
	drawText(graphics, 0, float64(250+nextOffset+50), "View more recommendations at recommendations.steamsal.es - by Cutie Café", 10, imgWidth)

	return image, nil
}
