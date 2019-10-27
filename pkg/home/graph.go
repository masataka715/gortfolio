package home

import (
	"gortfolio/pkg/footprint"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	filename = "images/access_graph.png"
)

func MakeAccessGraph() {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// p.Title.Text = "Number of accesses per page"
	p.X.Label.Text = "page name"
	p.Y.Label.Text = "accesses"
	p.NominalX("home", "chat", "scraping", "shiritori", "task", "footprint")

	counts := footprint.GetCount()
	nums := plotter.Values{
		float64(counts[0].Count),
		float64(counts[1].Count),
		float64(counts[2].Count),
		float64(counts[3].Count),
		float64(counts[4].Count),
		float64(counts[5].Count),
	}
	breadth := vg.Points(25)
	bar, err := plotter.NewBarChart(nums, breadth)
	if err != nil {
		panic(err)
	}

	bar.LineStyle.Width = vg.Length(0)
	bar.Color = plotutil.Color(4)
	p.Add(bar)

	filename := "images/access_graph.png"
	if err := p.Save(4*vg.Inch, 4*vg.Inch, filename); err != nil {
		panic(err)
	}
}
