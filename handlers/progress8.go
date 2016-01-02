package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/andrewcharlton/school-dashboard/analysis"
	"github.com/andrewcharlton/school-dashboard/database"
	"github.com/andrewcharlton/school-dashboard/national"
)

type p8Slot struct {
	Att8        float64
	Att8PerSlot float64
	Entries     float64
	Prog8       float64
	Text        string
}

type p8Student struct {
	analysis.Student
	Slots [5]p8Slot
	Att   float64
}

type p8Graph struct {
	X    []float64
	Y    []float64
	Text []string
}

func Progress8(e database.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if redir := checkRedirect(e, queryOpts{true, true}, w, r); redir {
			return
		}

		Header(e, w, r)
		FilterPage(e, w, r, false)
		defer Footer(e, w, r)

		data := struct {
			Slots       [5]p8Slot
			Students    []p8Student
			Query       template.URL
			GraphNat    p8Graph
			GraphPupils [4]p8Graph
		}{
			[5]p8Slot{},
			[]p8Student{},
			template.URL(ShortenQuery(e, r.URL.Query())),
			p8Graph{},
			[4]p8Graph{},
		}

		f := GetFilter(e, r)
		g, err := e.DB.GroupByFilter(f)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		nat, exists := e.Nationals[f.NatYear]
		if !exists {
			fmt.Fprintf(w, "Error: %v", err)
		}

		keys := []string{}
		for key := range nat.Prog8 {
			keys = append(keys, key)
		}
		sort.Sort(sort.StringSlice(keys))

		for _, key := range keys {
			p8 := nat.Prog8[key]
			ks2, err := strconv.ParseFloat(key, 64)
			if err == national.ErrNoKS2 {
				continue
			}
			if err != nil {
				fmt.Printf("Error: %v", err)
				continue
			}

			data.GraphNat.X = append(data.GraphNat.X, ks2)
			data.GraphNat.Y = append(data.GraphNat.Y, p8.Att8)
		}

		totalN := 0
		for _, s := range g.Students {
			slots, err := p8StudentData(s, nat)
			if err != nil {
				continue
			}
			stdnt := p8Student{s, slots, 100.0 * s.Attendance.Latest()}
			data.Students = append(data.Students, stdnt)

			var key int
			switch {
			case s.PP && s.Gender == "Male":
				key = 0
			case s.Gender == "Male":
				key = 1
			case s.PP && s.Gender == "Female":
				key = 2
			case s.Gender == "Female":
				key = 3
			}
			data.GraphPupils[key].X = append(data.GraphPupils[key].X, s.KS2.APS/6.0)
			data.GraphPupils[key].Y = append(data.GraphPupils[key].Y, slots[4].Att8)
			data.GraphPupils[key].Text = append(data.GraphPupils[key].Text, fmt.Sprintf("<b>%v</b><br>%v", s.Name(), slots[4].Text))

			for i := 0; i < 5; i++ {
				data.Slots[i].Att8 += slots[i].Att8
				data.Slots[i].Att8PerSlot += slots[i].Att8PerSlot
				data.Slots[i].Entries += slots[i].Entries
				data.Slots[i].Prog8 += slots[i].Prog8
			}
			totalN++
		}

		for i := 0; i < 5; i++ {
			data.Slots[i].Att8 = data.Slots[i].Att8 / float64(totalN)
			data.Slots[i].Entries = data.Slots[i].Entries / float64(totalN)
			data.Slots[i].Prog8 = data.Slots[i].Prog8 / float64(totalN)
		}

		err = e.Templates.ExecuteTemplate(w, "progress8.tmpl", data)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
		}
	}
}

func p8StudentData(s analysis.Student, nat national.National) ([5]p8Slot, error) {

	exp, err := nat.Progress8(s.KS2.APS)
	if err != nil {
		return [5]p8Slot{}, err
	}

	slots := [5]p8Slot{}

	b := s.Basket()

	en := b.English(exp)
	slots[0] = p8Slot{en.Ach, en.Ach / 2.0, float64(en.EntN), en.Pts, p8Text(b.Slots[0:2])}

	ma := b.Mathematics(exp)
	slots[1] = p8Slot{ma.Ach, ma.Ach / 2.0, float64(ma.EntN), ma.Pts, p8Text(b.Slots[2:3])}

	eb := b.EBacc(exp)
	slots[2] = p8Slot{eb.Ach, eb.Ach / 3.0, float64(eb.EntN), eb.Pts, p8Text(b.Slots[4:7])}

	oth := b.Other(exp)
	slots[3] = p8Slot{oth.Ach, oth.Ach / 3.0, float64(oth.EntN), oth.Pts, p8Text(b.Slots[7:10])}

	p8 := b.Progress8(exp)
	slots[4] = p8Slot{p8.Ach, p8.Ach / 10.0, float64(p8.EntN), p8.Pts, p8Text(b.Slots[0:10])}

	return slots, nil
}

func p8Text(slots []analysis.Slot) string {

	text := ""
	for _, s := range slots {
		if s.Subj == "" {
			text += "<br>"
		} else {
			text += fmt.Sprintf("%v - %v<br>", s.Subj, s.Grade)
		}
	}
	return text
}
