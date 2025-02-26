package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type RequestData struct {
	ID int `json:"id"`
}

type Stat struct {
	BaseStat int               `json:"base_stat"`
	Effort   int               `json:"effort"`
	Stat     map[string]string `json:"stat"`
}

type PokemonResponse struct {
	Name    string                 `json:"name"`
	Stats   []Stat                 `json:"stats"`
	Sprites map[string]interface{} `json:"sprites"`
}

func fetchData(url string) (map[string]interface{}, error) {
	fmt.Println("กำลังดึงข้อมูลจาก:", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ไม่สามารถเรียก API ได้:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API ตอบกลับด้วย Status Code: %d", resp.StatusCode)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("อ่าน Response ไม่ได้:", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("JSON Decode ไม่สำเร็จ:", err)
		return nil, err
	}

	fmt.Println("ดึงข้อมูลสำเร็จ!")
	return data, nil
}

func getPokemonHandler(w http.ResponseWriter, r *http.Request) {
	var reqData RequestData

	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		fmt.Println("Error: อ่าน JSON ไม่ได้", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("📌 ได้รับค่า ID: %d\n", reqData.ID)

	pokemonURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d/", reqData.ID)
	pokemonFormURL := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-form/%d/", reqData.ID)

	pokemonData, err := fetchData(pokemonURL)
	if err != nil {
		http.Error(w, "Failed to fetch Pokemon data", http.StatusInternalServerError)
		return
	}

	pokemonFormData, err := fetchData(pokemonFormURL)
	if err != nil {
		http.Error(w, "Failed to fetch Pokemon form data", http.StatusInternalServerError)
		return
	}

	name, ok := pokemonFormData["name"].(string)
	if !ok {
		fmt.Println("Warning: ไม่มีค่า name")
		name = "unknown"
	}

	sprites := map[string]interface{}{
		"back_default":       nil,
		"back_female":        nil,
		"back_shiny":         nil,
		"back_shiny_female":  nil,
		"front_default":      nil,
		"front_female":       nil,
		"front_shiny":        nil,
		"front_shiny_female": nil,
	}
	if spriteData, ok := pokemonFormData["sprites"].(map[string]interface{}); ok {
		for key := range sprites {
			if val, exists := spriteData[key]; exists {
				sprites[key] = val
			}
		}
	} else {
		fmt.Println("Warning: ไม่มีค่า sprites")
	}

	statsData, ok := pokemonData["stats"].([]interface{})
	stats := []Stat{}
	if !ok {
		fmt.Println("Warning: ไม่มีค่า stats")
	} else {
		for _, stat := range statsData {
			statMap, ok := stat.(map[string]interface{})
			if !ok {
				continue
			}

			baseStat, ok := statMap["base_stat"].(float64)
			if !ok {
				baseStat = 0
			}

			effort, ok := statMap["effort"].(float64)
			if !ok {
				effort = 0
			}

			statInfo, ok := statMap["stat"].(map[string]interface{})
			if !ok {
				statInfo = make(map[string]interface{})
			}

			statName, ok := statInfo["name"].(string)
			if !ok {
				statName = "unknown"
			}

			statURL, ok := statInfo["url"].(string)
			if !ok {
				statURL = ""
			}

			stats = append(stats, Stat{
				BaseStat: int(baseStat),
				Effort:   int(effort),
				Stat: map[string]string{
					"name": statName,
					"url":  statURL,
				},
			})
		}
	}

	responseData := PokemonResponse{
		Name:    name,
		Sprites: sprites,
		Stats:   stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
	fmt.Println("Sprites ที่ดึงมา:", sprites)
	fmt.Println("Stats ที่ดึงมา:", stats)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get-pokemon", getPokemonHandler).Methods("POST")

	fmt.Println("Server is running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", r))

}
