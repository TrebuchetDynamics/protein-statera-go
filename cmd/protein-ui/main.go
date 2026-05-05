//go:build !cgo

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/TrebuchetDynamics/protein-statera-go/internal/ui"
	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
	uiapp "github.com/gogpu/ui/app"
	"github.com/gogpu/ui/core/scrollview"
	"github.com/gogpu/ui/primitives"
	"github.com/gogpu/ui/render"
	uitheme "github.com/gogpu/ui/theme"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)

func main() {
	screenshotPath := flag.String("screenshot", "", "write an offscreen PNG screenshot and exit")
	flag.Parse()

	model := ui.DefaultModel()
	seed := widget.Hex(0x2F5D50)
	materialTheme := material3.New(seed)
	if *screenshotPath != "" {
		if err := saveScreenshot(*screenshotPath, model, materialTheme); err != nil {
			log.Fatal(err)
		}
		return
	}
	if !displayAvailable(os.Getenv) {
		path := headlessScreenshotPath(os.TempDir())
		if err := saveScreenshot(path, model, materialTheme); err != nil {
			log.Fatal(err)
		}
		log.Printf("no display detected; wrote offscreen screenshot: %s", path)
		return
	}

	gpuApp := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle(model.Spec.Title).
		WithSize(model.Spec.Width, model.Spec.Height).
		WithContinuousRender(false))

	appTheme := uitheme.DefaultLight()
	appTheme.Colors.Primary = seed
	appTheme.Colors.PrimaryDark = widget.Hex(0x1F463C)
	appTheme.Colors.PrimaryLight = widget.Hex(0x6F9C8D)
	app := uiapp.New(
		uiapp.WithWindowProvider(gpuApp),
		uiapp.WithPlatformProvider(gpuApp),
		uiapp.WithEventSource(gpuApp.EventSource()),
		uiapp.WithTheme(appTheme),
	)
	app.SetRoot(buildRoot(model, materialTheme))

	var canvas *ggcanvas.Canvas
	gpuApp.OnDraw(func(dc *gogpu.Context) {
		w, h := dc.Width(), dc.Height()
		if w <= 0 || h <= 0 {
			return
		}
		if canvas == nil {
			provider := gpuApp.GPUContextProvider()
			if provider == nil {
				return
			}
			var err error
			canvas, err = ggcanvas.New(provider, w, h)
			if err != nil {
				log.Printf("ggcanvas: %v", err)
				return
			}
		}

		app.Frame()
		cw, ch := canvas.Size()
		if cw != w || ch != h {
			if err := canvas.Resize(w, h); err != nil {
				log.Printf("resize: %v", err)
				return
			}
			cw, ch = w, h
		}

		if err := canvas.Draw(func(cc *gg.Context) {
			cc.SetRGBA(0.96, 0.97, 0.96, 1)
			cc.DrawRectangle(0, 0, float64(cw), float64(ch))
			cc.Fill()
			app.Window().DrawTo(render.NewCanvas(cc, cw, ch))
		}); err != nil {
			log.Printf("draw: %v", err)
			return
		}
		if err := canvas.Render(dc.RenderTarget()); err != nil {
			log.Printf("render: %v", err)
		}
	})
	gpuApp.OnClose(func() { gg.CloseAccelerator() })

	if err := gpuApp.Run(); err != nil {
		log.Fatal(err)
	}
}

func displayAvailable(getenv func(string) string) bool {
	return getenv("WAYLAND_DISPLAY") != "" || getenv("DISPLAY") != ""
}

func headlessScreenshotPath(tempDir string) string {
	return filepath.Join(tempDir, "protein-ui-headless.png")
}

func saveScreenshot(path string, model ui.AppModel, theme *material3.Theme) error {
	appTheme := uitheme.DefaultLight()
	app := uiapp.New(
		uiapp.WithWindowProvider(gpucontext.NullWindowProvider{W: model.Spec.Width, H: model.Spec.Height}),
		uiapp.WithTheme(appTheme),
	)
	app.SetRoot(buildRoot(model, theme))
	app.Frame()

	dc := gg.NewContext(model.Spec.Width, model.Spec.Height)
	dc.SetRGBA(0.96, 0.97, 0.96, 1)
	dc.DrawRectangle(0, 0, float64(model.Spec.Width), float64(model.Spec.Height))
	dc.Fill()
	app.Window().DrawTo(render.NewCanvas(dc, model.Spec.Width, model.Spec.Height))
	return dc.SavePNG(path)
}

func buildRoot(model ui.AppModel, theme *material3.Theme) widget.Widget {
	railItems := make([]widget.Widget, 0, len(model.Views)+1)
	railItems = append(railItems,
		primitives.Text("Views").FontSize(13).Bold().Color(widget.Hex(0x24483E)),
	)
	for _, view := range model.Views {
		railItems = append(railItems, railItem(view.Name))
	}
	rail := primitives.Box(railItems...).
		Width(200).
		Padding(14).
		Gap(8).
		Background(widget.Hex(0xE7EFEA)).
		Rounded(8)

	content := primitives.Box(
		header(model),
		summarySection(model, theme),
		confidenceSection(model, theme),
		structureSection(model, theme),
		clashSection(model, theme),
		evidenceSection(model, theme),
		methodologySection(model, theme),
	).Padding(20).Gap(12)

	return primitives.HBox(
		rail,
		scrollview.New(content, scrollview.PainterOpt(material3.ScrollbarPainter{Theme: theme})),
	).Padding(16).Gap(14).Background(theme.Colors.Surface)
}

func header(model ui.AppModel) widget.Widget {
	return primitives.Box(
		primitives.Text("Protein Statera Go").FontSize(24).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text("Evidence-first protein structure validation dashboard").FontSize(13).Color(widget.Hex(0x44504B)),
		primitives.Text(fmt.Sprintf("%d residues | %d atoms | confidence bands: VH=%d C=%d L=%d VL=%d | clashes: %d",
			model.Summary.ResidueCount, model.Summary.AtomCount,
			model.Summary.VeryHighConf, model.Summary.ConfidentConf,
			model.Summary.LowConf, model.Summary.VeryLowConf,
			model.Summary.StericClashes)).FontSize(11).Color(widget.Hex(0x52645C)),
	).Gap(6)
}

func summarySection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	s := model.Summary
	cards := []widget.Widget{
		summaryCard("Residues", fmt.Sprintf("%d", s.ResidueCount), "Total residues in structure", 200),
		summaryCard("Atoms", fmt.Sprintf("%d", s.AtomCount), "Total atoms across all residues", 200),
		summaryCard("Very High", fmt.Sprintf("%d (%.0f%%)", s.VeryHighConf, pct(s.VeryHighConf, s.ResidueCount)), "pLDDT >= 90", 200),
		summaryCard("Confident", fmt.Sprintf("%d (%.0f%%)", s.ConfidentConf, pct(s.ConfidentConf, s.ResidueCount)), "pLDDT 70-90", 200),
	}
	cards2 := []widget.Widget{
		summaryCard("Low", fmt.Sprintf("%d (%.0f%%)", s.LowConf, pct(s.LowConf, s.ResidueCount)), "pLDDT 50-70", 200),
		summaryCard("Very Low", fmt.Sprintf("%d (%.0f%%)", s.VeryLowConf, pct(s.VeryLowConf, s.ResidueCount)), "pLDDT < 50", 200),
		summaryCard("Steric Clashes", fmt.Sprintf("%d", s.StericClashes), "VDW overlap >= 0.4 Å", 200),
		summaryCard("Severe Clashes", fmt.Sprintf("%d", s.SevereClashes), "VDW overlap >= 0.9 Å", 200),
	}
	return primitives.Box(
		primitives.HBox(cards...).Gap(10),
		primitives.HBox(cards2...).Gap(10),
	).Gap(8)
}

func summaryCard(name, metric, desc string, width float32) widget.Widget {
	return primitives.Box(
		primitives.Text(name).FontSize(12).Bold().Color(widget.Hex(0x183D34)),
		primitives.Text(metric).FontSize(16).Bold().Color(widget.Hex(0x246B45)),
		primitives.Text(desc).FontSize(10).Color(widget.Hex(0x52645C)),
	).Width(width).Padding(10).Gap(3).Background(widget.Hex(0xFFFFFF)).Rounded(6).BorderStyle(1, widget.Hex(0xD8E1DC))
}

func pct(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func confidenceSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Confidence Analysis (AlphaFold pLDDT)").FontSize(16).Bold().Color(widget.Hex(0x183D34)),
		boundary("EMBL-EBI standard pLDDT bands. Very High >= 90 | Confident 70-90 | Low 50-70 | Very Low < 50"),
	}
	for _, rec := range model.Confidence {
		color := widget.Hex(0x246B45)
		switch rec.Band {
		case "confident":
			color = widget.Hex(0x3A7CA5)
		case "low":
			color = widget.Hex(0xC9A037)
		case "very-low":
			color = widget.Hex(0xC96037)
		}
		children = append(children, card(
			primitives.Text(fmt.Sprintf("%s %d (chain %s)", rec.ResidueName, rec.Index, rec.ChainID)).FontSize(12).Bold().Color(color),
			primitives.Text(fmt.Sprintf("pLDDT: %.1f | band: %s", rec.PLDDT, rec.Band)).FontSize(10).Color(widget.Hex(0x44504B)),
		))
	}
	if len(model.VeryLowSegs) > 0 {
		for _, seg := range model.VeryLowSegs {
			children = append(children, card(
				primitives.Text(fmt.Sprintf("Very-low segment: residues %d-%d (%d residues)", seg.Start, seg.End, seg.Count)).FontSize(11).Color(widget.Hex(0xC96037)),
			))
		}
	}
	return section("Confidence", children, theme, widget.Hex(0xE8F0EC))
}

func structureSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	s := model.Structures[0]
	children := []widget.Widget{
		primitives.Text("Structure Overview").FontSize(16).Bold().Color(widget.Hex(0x183D34)),
		card(
			primitives.Text(fmt.Sprintf("Structure: %s", s.ID)).FontSize(13).Bold(),
			primitives.Text(fmt.Sprintf("Source: %s", s.Source)).FontSize(10).Color(widget.Hex(0x44504B)),
			primitives.Text(fmt.Sprintf("Residues: %d | Atoms: %d | pLDDT range: %s", s.ResidueCount, s.AtomCount, s.PLDDTRange)).FontSize(10).Color(widget.Hex(0x52645C)),
		),
	}
	for _, note := range s.Notes {
		children = append(children, card(
			primitives.Text(note).FontSize(11).Color(widget.Hex(0x246B45)),
		))
	}
	return section("Structure", children, theme, widget.Hex(0xEBF5F0))
}

func clashSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Steric Clash Detection (VDW Radii)").FontSize(16).Bold().Color(widget.Hex(0x183D34)),
		boundary("MolProbity-style VDW overlap detection. Severe: >= 0.9 Å overlap. Non-severe: >= 0.4 Å overlap."),
	}
	if len(model.Clashes) == 0 {
		children = append(children, card(
			primitives.Text("No steric clashes detected").FontSize(12).Color(widget.Hex(0x246B45)),
		))
	} else {
		for _, clash := range model.Clashes {
			severity := "non-severe"
			color := widget.Hex(0xC9A037)
			if clash.Severe {
				severity = "severe"
				color = widget.Hex(0xC96037)
			}
			children = append(children, card(
				primitives.Text(fmt.Sprintf("Clash: %s/%s <-> %s/%s", clash.ResidueA, clash.AtomA, clash.ResidueB, clash.AtomB)).FontSize(11).Bold().Color(color),
				primitives.Text(fmt.Sprintf("distance: %.2f Å | %s (chain %s <-> %s)", clash.Distance, severity, clash.ChainA, clash.ChainB)).FontSize(10).Color(widget.Hex(0x44504B)),
			))
		}
	}
	return section("Geometry", children, theme, widget.Hex(0xF5E8DC))
}

func evidenceSection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Evidence Class").FontSize(16).Bold().Color(widget.Hex(0x183D34)),
		boundary("All metrics derived from parsed structure data. pLDDT values are AlphaFold model confidence estimates, not experimental measurements."),
		card(
			primitives.Text("Evidence class: predicted-structure").FontSize(12).Bold().Color(widget.Hex(0x3A7CA5)),
			primitives.Text("AlphaFold prediction confidence mapped to EMBL-EBI standard bands").FontSize(10).Color(widget.Hex(0x44504B)),
			primitives.Text("Steric clashes flagged by VDW radii overlap detection").FontSize(10).Color(widget.Hex(0x44504B)),
			primitives.Text("This tool validates structures; it does not predict them").FontSize(10).Color(widget.Hex(0x52645C)),
		),
	}
	return section("Evidence", children, theme, widget.Hex(0xE8EEF0))
}

func methodologySection(model ui.AppModel, theme *material3.Theme) widget.Widget {
	children := []widget.Widget{
		primitives.Text("Methodology").FontSize(16).Bold().Color(widget.Hex(0x183D34)),
		card(
			primitives.Text("PDB Parsing").FontSize(12).Bold(),
			primitives.Text("ATOM record parsing with residue grouping, chain ID tracking, and multi-model support").FontSize(10).Color(widget.Hex(0x44504B)),
		),
		card(
			primitives.Text("pLDDT Extraction").FontSize(12).Bold(),
			primitives.Text("AlphaFold writes per-residue confidence to B-factor columns. Averaged across atoms per residue.").FontSize(10).Color(widget.Hex(0x44504B)),
		),
		card(
			primitives.Text("Steric Clash Detection").FontSize(12).Bold(),
			primitives.Text("Van der Waals radii overlap check (MolProbity-style). Element-specific radii from standard tables.").FontSize(10).Color(widget.Hex(0x44504B)),
		),
		card(
			primitives.Text("RMSD Comparison").FontSize(12).Bold(),
			primitives.Text("Root-mean-square deviation for already-aligned structures. Uses Kabsch algorithm for superposition.").FontSize(10).Color(widget.Hex(0x44504B)),
		),
	}
	return section("Methodology", children, theme, widget.Hex(0xEEF0F2))
}

func section(title string, children []widget.Widget, theme *material3.Theme, bg widget.Color) widget.Widget {
	return primitives.Box(
		append([]widget.Widget{children[0]}, children[1:]...)...,
	).Padding(14).Gap(8).Background(bg).Rounded(8).BorderStyle(1, widget.Hex(0xD4DED8))
}

func railItem(name string) widget.Widget {
	return primitives.Box(
		primitives.Text(name).FontSize(13).Bold().Color(widget.Hex(0x183D34)),
	).Padding(8).Gap(4).Background(widget.Hex(0xF6FAF7)).Rounded(6)
}

func card(children ...widget.Widget) widget.Widget {
	return primitives.Box(children...).
		Padding(10).
		Gap(4).
		Background(widget.Hex(0xFFFFFF)).
		Rounded(6).
		BorderStyle(1, widget.Hex(0xD8E1DC))
}

func boundary(text string) widget.Widget {
	return primitives.Box(
		primitives.Text(text).FontSize(11).Color(widget.Hex(0x5C3B00)),
	).Padding(10).Background(widget.Hex(0xFFF4D8)).Rounded(6).BorderStyle(1, widget.Hex(0xE7C66A))
}
