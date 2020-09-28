package draw

import (
	"fmt"

	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

// FontCache is the custom font cache
type FontCache map[string]*truetype.Font

// Store stores a font in the cache
func (fc FontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

// Load loads a font from the cache
func (fc FontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]

	if !stored {
		return nil, fmt.Errorf("%s not in font cache", fd.Name)
	}

	return font, nil
}
