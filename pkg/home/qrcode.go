package home

import (
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func MakeQRcode() {
	filename := "images/qrcode.png"
	if _, err := os.Stat(filename); err == nil {
		return
	}
	qrCode, _ := qr.Encode("ポートフォリオをご覧頂き、ありがとうございます", qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 150, 150)
	file, _ := os.Create(filename)
	png.Encode(file, qrCode)
}
