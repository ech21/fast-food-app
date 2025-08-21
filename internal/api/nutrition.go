package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type nutritionSvc struct {
	client              *http.Client
	NUTRITIONIX_APP_ID  string
	NUTRITIONIX_APP_KEY string
}

type nutritionixSearchInstantResult struct {
	Photo struct {
		Thumb string `json:"thumb"`
	} `json:"photo"`
	Nf_calories          uint   `json:"nf_calories"`
	Serving_unit         string `json:"serving_unit"`
	Serving_qty          uint   `json:"serving_qty"`
	Brand_type           uint   `json:"brand_type"`
	Region               uint   `json:"region"`
	Nix_brand_id         string `json:"nix_brand_id"`
	Brand_name           string `json:"brand_name"`
	Food_name            string `json:"food_name"`
	Brand_name_item_name string `json:"brand_name_item_name"`
	Nix_item_id          string `json:"nix_item_id"`
	Locale               string `json:"locale"`
}

func (svc *nutritionSvc) nutritionAutocomplete(query string) nutritionAutocompleteOutput {
	v := url.Values{}
	v.Add("query", query)
	url := "https://trackapi.nutritionix.com/v2/search/instant?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nutritionAutocompleteOutput{
			Results: nil,
			Err:     err,
		}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("x-app-id", svc.NUTRITIONIX_APP_ID)
	req.Header.Add("x-app-key", svc.NUTRITIONIX_APP_KEY)
	resp, err := svc.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nutritionAutocompleteOutput{
			Results: nil,
			Err:     err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nutritionAutocompleteOutput{
			Results: nil,
			Err:     err,
		}
	}
	if resp.StatusCode >= 300 {
		fmt.Println(resp.Status)
		fmt.Println(string(body))
	}
	var dat struct {
		Branded []nutritionixSearchInstantResult
	}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		fmt.Println(err)
		return nutritionAutocompleteOutput{
			Results: nil,
			Err:     err,
		}
	}
	out := make([]Item, len(dat.Branded))
	for i := 0; i < len(dat.Branded); i++ {
		v := dat.Branded[i]
		out[i] = Item{
			Name:         v.Food_name,
			LocationName: v.Brand_name,
			Price:        0,
		}
	}
	return nutritionAutocompleteOutput{
		Results: out,
		Err:     nil,
	}
}

func (svc *nutritionSvc) attach(mux *http.ServeMux) {
	fmt.Println("GET /api/nutritionAutocomplete")
	mux.HandleFunc("GET /api/nutritionAutocomplete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit: /api/nutritionAutocomplete")
		q := r.FormValue("query")
		if len(q) == 0 {
			fmt.Println("Missing in params 'query'")
			http.Error(w, "Missing in params 'query'", http.StatusBadRequest)
			return
		}
		out := svc.nutritionAutocomplete(q)
		jsonOut, err := json.Marshal(out)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Data could not be marshalled into json.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(jsonOut)
	})
}

func newNutritionSvc() *nutritionSvc {
	fmt.Println("New Nutrition Service")
	appId := os.Getenv("NUTRITIONIX_APP_ID")
	if len(appId) == 0 {
		panic("Missing env var: NUTRITIONIX_APP_ID")
	}
	fmt.Println("Found env var NUTRITIONIX_APP_ID: " + appId)

	appKey := os.Getenv("NUTRITIONIX_APP_KEY")
	if len(appKey) == 0 {
		panic("Missing env var: NUTRITIONIX_APP_KEY")
	}
	fmt.Println("Found env var NUTRITIONIX_APP_KEY: " + appKey)

	svc := &nutritionSvc{
		NUTRITIONIX_APP_ID:  appId,
		NUTRITIONIX_APP_KEY: appKey,
		client:              &http.Client{},
	}
	return svc
}
