package app

import "github.com/revel/revel"
import "encoding/json"
import "io/ioutil"
import "fmt"
import "reflect"

type TPS struct {
    Index int `json:"index"`
    Location string `json:"location"`
    Kelurahan_id int `json:"kelurahan_id"`
    Jokowi int `json:"jokowi"`
    Prabowo int `json:"prabowo"`
}


var (
    Data   [][]interface{};
)

var Kabs = make(map[int][]TPS)
var Provinces = make(map[int][]TPS)
var National = make(map[int][]TPS)
var TotalTPS = 0

func Append(slice []TPS, elements []TPS) []TPS {
    n := len(slice)
    total := len(slice) + len(elements)
    if total > cap(slice) {
        // Reallocate. Grow to 1.5 times the new size, so we can still grow.
        newSize := total*3/2 + 1
        newSlice := make([]TPS, total, newSize)
        copy(newSlice, slice)
        slice = newSlice
    }
    slice = slice[:total]
    copy(slice[n:], elements)
    return slice
}

func InitTPS() {
    National[0] = make([]TPS, 0)
    /*
    NO_TPS = 0
    PRABOWO_TPS = 1
    JOKOWI_TPS = 2
    SAH_TPS = 3
    TIDAK_SAH_TPS = 4
    TERDATA_TPS = 5
    ERROR_TPS = 6

    TEMPAT_ID = 0
    ORTU_ID= 1
    NAMA = 2
    JUMLAH_TPS = 3
    ANAK = 4
    PRABOWO = 5
    JOKOWI = 6
    SAH = 7
    TIDAK_SAH = 8
    TPS_TERDATA = 9
    TPS_ERROR = 10
    */

    content, err := ioutil.ReadFile("tps.json")
    if err!=nil{
        fmt.Print("Error:",err)
    }
    err=json.Unmarshal(content, &Data)
    if err!=nil{
        fmt.Print("Error:",err)
    }

    maps := make(map[int][]interface{});

    for _, d := range Data {
        maps[int((d[0].(float64)))] = d
    }

    for _, d := range Data {
        arr := d[4].([]interface{})
        arr_len := int(d[3].(float64))
        is_tps := len(arr) == arr_len
        for _, tps := range arr {
            if reflect.TypeOf(tps).Kind() == reflect.Float64 {
                is_tps = false
            }
        }
        if is_tps {
            kec := maps[int(d[1].(float64))]
            kab_id := int(kec[1].(float64))
            kab := maps[kab_id]
            province_id := int(kab[1].(float64))
            province := maps[province_id]
            if _,ok := Kabs[kab_id]; !ok {
                Kabs[kab_id] = make([]TPS, 0)
            }
            if _,ok := Provinces[province_id]; !ok {
                //do something here
                Provinces[province_id] = make([]TPS, 0)
            }
            items := d[4].([]interface{})
            tpses := make([]TPS, len(items))
            for i, tps_raw := range items{
                raw := tps_raw.([]interface{})
                TotalTPS += 1
                tps := TPS{}
                tps.Index = int(i)
                tps.Location = province[2].(string) + " > " + kab[2].(string) + " > " + kec[2].(string) + " > " + d[2].(string) 
                tps.Kelurahan_id = int(d[0].(float64))
                tps.Prabowo = int(raw[1].(float64))
                tps.Jokowi = int(raw[2].(float64))
                tpses[i] = tps
            }
            Kabs[kab_id] = Append(Kabs[kab_id], tpses)
            Provinces[province_id] = Append(Provinces[province_id], tpses)
            National[0] = Append(National[0], tpses)
        }
    }
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}
    InitTPS();
	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB())
	// revel.OnAppStart(FillCache())
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
