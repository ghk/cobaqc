package controllers

import "github.com/revel/revel"
import "cobaqc/app"
import "math/rand"
import "math"

type App struct {
	*revel.Controller
}

type QcResult struct {
    Used []app.TPS `json:"used"`
    Prabowo int `json:"prabowo"`
    Jokowi int `json:"jokowi"`

}
func Pick(m int, n int) map[int]bool {
    set := make(map[int]bool)
    found := 0
    for ;found < n; {
        i := rand.Intn(m)
        if _,ok := set[i]; !ok {
            set[i] = true
            found += 1
        }
    }
    return set;
}
func Sample(tpses []app.TPS, n int) []app.TPS {
    result := make([]app.TPS, n)
    picked := Pick(len(tpses), n)
    i := 0
    for key, value := range picked {
        if value {
            result[i] = tpses[key]
            i = i+1
        }
    }
    return result
}

func DoQc(sample_type int, sample int) QcResult {
    result := QcResult{}
    result.Used = make([]app.TPS, 0)
    result.Jokowi = 0
    result.Prabowo = 0
    used_count := 0
    source := app.Kabs
    if sample_type == 1 {
        source = app.Provinces
    }
    if sample_type == 2 {
        source = app.National
    }
    for _, tpses := range source {
        sampled_count := int(math.Floor(float64(sample) * float64(len(tpses)) / float64(app.TotalTPS) + 0.5))
        if sampled_count > 0 {
            s_tpses := Sample(tpses, sampled_count)
            result.Used = app.Append(result.Used, s_tpses)
            for _,tps := range s_tpses {
                result.Prabowo += tps.Prabowo
                result.Jokowi += tps.Jokowi
                used_count += 1
            }
        }
    }
    return result
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Result(typ int, count int) revel.Result {
    qc := DoQc(typ, count)
	return c.RenderJson(qc)
}
