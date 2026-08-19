package utils

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"net/http"

	"github.com/deepteams/webp"
)

const (
	EmbeddingDim   = 64
	pHashSize      = 32
	embeddingBlock = 8
)

func DecodeImage(data []byte) (image.Image, error) {
	contentType := http.DetectContentType(data)
	switch contentType {
	case "image/jpeg", "image/png":
		img, _, err := image.Decode(bytes.NewReader(data))
		return img, err
	case "image/webp":
		return webp.Decode(bytes.NewReader(data))
	default:
		return nil, errors.New("tipe gambar tidak didukung (hanya JPEG, PNG, WebP)")
	}
}

func ComputeImageEmbedding(img image.Image) ([]float64, error) {
	gray := resizeToGray(img)
	if gray == nil {
		return nil, errors.New("gambar tidak valid")
	}
	dct := dct2D(gray)

	embedding := make([]float64, EmbeddingDim)
	idx := 0
	for u := 0; u < embeddingBlock; u++ {
		for v := 0; v < embeddingBlock; v++ {
			embedding[idx] = dct[u][v]
			idx++
		}
	}
	normalizeL2(embedding)
	return embedding, nil
}

func CosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func resizeToGray(img image.Image) [][]float64 {
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return nil
	}

	out := make([][]float64, pHashSize)
	for y := 0; y < pHashSize; y++ {
		out[y] = make([]float64, pHashSize)
		sy := (float64(y)+0.5)*float64(srcH)/pHashSize - 0.5
		for x := 0; x < pHashSize; x++ {
			sx := (float64(x)+0.5)*float64(srcW)/pHashSize - 0.5
			out[y][x] = bilinearGray(img, bounds, sx, sy)
		}
	}
	return out
}

func bilinearGray(img image.Image, b image.Rectangle, x, y float64) float64 {
	x0 := int(math.Floor(x))
	y0 := int(math.Floor(y))
	tx := x - float64(x0)
	ty := y - float64(y0)

	g00 := grayAt(img, b, x0, y0)
	g10 := grayAt(img, b, x0+1, y0)
	g01 := grayAt(img, b, x0, y0+1)
	g11 := grayAt(img, b, x0+1, y0+1)

	return g00*(1-tx)*(1-ty) + g10*tx*(1-ty) + g01*(1-tx)*ty + g11*tx*ty
}

func grayAt(img image.Image, b image.Rectangle, x, y int) float64 {
	if x < b.Min.X {
		x = b.Min.X
	}
	if x > b.Max.X-1 {
		x = b.Max.X - 1
	}
	if y < b.Min.Y {
		y = b.Min.Y
	}
	if y > b.Max.Y-1 {
		y = b.Max.Y - 1
	}
	r, g, bl, _ := img.At(x, y).RGBA()
	return 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(bl>>8)
}

func dct2D(pixels [][]float64) [][]float64 {
	n := len(pixels)
	tmp := make([][]float64, n)
	for i := range tmp {
		tmp[i] = make([]float64, n)
	}

	for y := 0; y < n; y++ {
		for v := 0; v < n; v++ {
			sum := 0.0
			for x := 0; x < n; x++ {
				sum += pixels[y][x] * math.Cos((2*float64(x)+1)*float64(v)*math.Pi/(2*float64(n)))
			}
			tmp[y][v] = sum * dctScale(v, n)
		}
	}

	result := make([][]float64, n)
	for i := range result {
		result[i] = make([]float64, n)
	}

	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			sum := 0.0
			for y := 0; y < n; y++ {
				sum += tmp[y][v] * math.Cos((2*float64(y)+1)*float64(u)*math.Pi/(2*float64(n)))
			}
			result[u][v] = sum * dctScale(u, n)
		}
	}
	return result
}

func dctScale(k, n int) float64 {
	if k == 0 {
		return math.Sqrt(1.0 / float64(n))
	}
	return math.Sqrt(2.0 / float64(n))
}

func normalizeL2(v []float64) {
	var sum float64
	for _, x := range v {
		sum += x * x
	}
	if sum == 0 {
		return
	}
	norm := math.Sqrt(sum)
	for i := range v {
		v[i] /= norm
	}
}
