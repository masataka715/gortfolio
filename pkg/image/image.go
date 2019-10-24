package image

import (
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	qrCode, _ := qr.Encode("ポートフォリオをご覧頂き、ありがとうございます", qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 150, 150)
	png.Encode(w, qrCode)
}
