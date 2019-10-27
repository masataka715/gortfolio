package footprint

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fontmaker/core"
	"github.com/stretchr/objx"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	when := time.Now().Format("2006年01月02日 15時04分")
	Insert("あしあと", when)

	data := map[string]interface{}{}
	data["Footprint"] = GetAll()
	templates := template.Must(template.ParseFiles("templates/layout.html",
		"templates/footprint.html"))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("ipagp", "IPA_font/ipagp.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}
	fontSize := 15
	err = pdf.SetFont("ipagp", "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return
	}

	op := gopdf.CellOption{Align: gopdf.Center}
	rect := gopdf.Rect{W: 600, H: 100}
	text := "gortfolioの閲覧履歴(最新の20件)"
	pdf.CellWithOption(&rect, text, op)

	var parser core.TTFParser
	err = parser.Parse("IPA_font/ipagp.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	cap := float64(float64(parser.CapHeight()) * 1000.00 / float64(parser.UnitsPerEm()))
	realHeight := cap * (float64(fontSize) / 500.0)
	pdf.Br(realHeight)

	footprints := GetAll()
	var h float64
	for i, f := range footprints {
		pdf.Br(realHeight + 10.0)
		h = 90.0 + (realHeight+10.0)*(float64(i)-1.0)
		pdf.Line(150, h, 470, h)
		rect = gopdf.Rect{W: 600, H: h}
		// pdf.Cell(&rect, f.When+f.PageName)
		pdf.CellWithOption(&rect, f.When+" "+f.PageName, op)
	}

	pdf.Write(w)
}

func SetCookie(w http.ResponseWriter, uniqueID string, name string, avatarURL string) {
	FootprintCookie := objx.New(map[string]interface{}{
		"view_page": uniqueID,
		"when":      name,
	}).MustBase64()
	http.SetCookie(w, &http.Cookie{
		Name:  "footprint",
		Value: FootprintCookie,
		Path:  "/",
	})
}
