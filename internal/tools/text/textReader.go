package text

import "image"

type TextReader interface {
	read(image image.Image) string
}
