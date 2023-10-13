package main

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"github.com/zocker-160/Berlin-Taxi-DCT/impl/dct"
)

func loadDCT(w fyne.Window, l *widget.Label, filename string) {
	dct, err := dct.DecodeFile(filename)
	if err != nil {
		l.SetText(err.Error())
	}

	image := canvas.NewImageFromImage(dct)
	image.FillMode = canvas.ImageFillContain

	w.SetContent(image)
}

func main() {
	a := app.New()
	w := a.NewWindow("DCT Viewer")

	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()

	sLabel := widget.NewLabel("Loading image...")
	w.SetContent(sLabel)

	if len(os.Args) != 2 {
		sLabel.SetText("No file specified!")
	} else {
		go loadDCT(w, sLabel, os.Args[1])
	}

	go w.SetOnDropped(func(p fyne.Position, u []fyne.URI) {
		fmt.Println("dropped:", u)

		file := u[0]

		if strings.HasSuffix(strings.ToLower(file.Name()), ".dct") {
			sLabel.SetText("Loading "+file.Name())
			loadDCT(w, sLabel, file.Path())
		}
	})

	w.ShowAndRun()
}
