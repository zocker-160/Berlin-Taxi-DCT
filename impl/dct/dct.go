package dct

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

type DCT struct {
	Header Header
	Data []byte
}

type Header struct {
	Magic [3]byte // always "DC2"
	Scale float32
	Xres uint32
	Yres uint32
	BPP uint8
	Unknown byte
	NumResolutions uint8
}

func ParseHeader(f io.Reader) (*Header, error) {
	var header = new(Header)

	err := binary.Read(f, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	if header.Magic != [3]byte{'D', 'C', '2'} {
		return nil, fmt.Errorf("unknown header: %v", header.Magic)
	}

	if header.Scale != 1 {
		return nil, fmt.Errorf("scale != 1.0")
	}

	if header.BPP != 24 && header.BPP != 32 {
		return nil, fmt.Errorf("unsupported BPP: %d", header.BPP)
	}

	return header, nil
}

func DecodeFile(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Decode(file)
}

func Decode(r io.Reader) (image.Image, error) {
	header, err := ParseHeader(r)
	if err != nil {
		return nil, err
	}

	size := calcBinarySize(int(header.Xres), int(header.Yres), int(header.BPP), int(header.NumResolutions))
	data := make([]byte, size)

	err = binary.Read(r, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}

	d := DCT{
		Header: *header,
		Data: data,
	}

	if d.Header.BPP == 24 {
		return &BGRImage{
			Width: int(d.Header.Xres),
			Height: int(d.Header.Yres),
			Data: d.Data,
		}, nil

	} else {
		return &BGRAImage{
			Width: int(d.Header.Xres),
			Height: int(d.Header.Yres),
			Data: d.Data,
		}, nil
	}
}

func Encode(w io.Writer, m image.Image) error {
	wb := bufio.NewWriter(w)
	width, height := m.Bounds().Max.X, m.Bounds().Max.Y

	dctHeader := Header{
		Magic: [3]byte{'D', 'C', '2'},
		Scale: 1,
		Xres: uint32(width), Yres: uint32(height),
		BPP: 32,
		Unknown: 0,
		NumResolutions: 1,
	}

	err := binary.Write(wb, binary.LittleEndian, dctHeader)
	if err != nil {
		return err
	}

	pixel := make([]byte, 4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := m.At(x, y).RGBA()

			pixel[0] = uint8(b >> 8)
			pixel[1] = uint8(g >> 8)
			pixel[2] = uint8(r >> 8)
			pixel[3] = uint8(a >> 8)

			wb.Write(pixel)
		}
	}

	wb.Flush()

	return nil
}

func calcBinarySize(xRes, yRes, bpp, numRes int) int {
	//size := 0
	//for i := 0; i < numRes; i++ {
	//	size += ( (xRes * yRes) / int(math.Pow(2, float64(2*i))) ) * (bpp / 8)
	//}

	// we only care about the first image
	size := (xRes * yRes * bpp) / 8

	return size
}

type BGRImage struct {
	Width, Height int
	Data []byte
}

func (img *BGRImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (img *BGRImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img *BGRImage) At(x, y int) color.Color {
	i := (y * img.Width + x) * 3
	b, g, r := img.Data[i], img.Data[i+1], img.Data[i+2]

	return color.NRGBA{r, g, b, 0xFF}
}


type BGRAImage struct {
	Width, Height int
	Data []byte
}

func (img *BGRAImage) ColorModel() color.Model {
	return color.NRGBAModel
}

func (img *BGRAImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.Width, img.Height)
}

func (img *BGRAImage) At(x, y int) color.Color {
	i := (y * img.Width + x) * 4
	b, g, r, a := img.Data[i], img.Data[i+1], img.Data[i+2], img.Data[i+3]

	return color.NRGBA{r, g, b, a}
}
