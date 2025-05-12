package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Name string `json:"name"`
	Tier int    `json:"tier"`
}

type ElementRecipe struct {
	Tier   int      `json:"tier"`
	Result string   `json:"result"`
	Recipe []string `json:"recipe"`
}

func InitializeData() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	url2 := "https://little-alchemy.fandom.com/wiki/Elements_(Myths_and_Monsters)"

	// Get filters (myth & monsters recipes)
	filters := getFilters(url2)
	filters = append(filters, "Time")

	// Get recipes minus filters
	recipes, elements := getRecipesAndElements(url, filters)

	// Write recipes to file
	writeToFile("recipes.json", recipes)
	
	// Write elements to file
	writeToFile("elements.json", elements)

	fmt.Printf("Amount of recipes found: %d\n", len(recipes))
	fmt.Printf("Amount of elements found: %d\n", len(elements))
}

func getFilters(url string) []string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var filters []string

	doc.Find("table.list-table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, tr *goquery.Selection) {
			if j == 0 {
				return // skip header row
			}

			tds := tr.Find("td")
			if tds.Length() == 0 {
				return
			}

			result := tds.Eq(0).Find("a[href^='/wiki/']").First().Text()
			filters = append(filters, result)
		})
	})

	return filters
}

func getRecipesAndElements(url string, filters []string) ([]ElementRecipe, []Element) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var recipes []ElementRecipe
	var elements []Element
	tier := 0

	doc.Find("table.list-table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, tr *goquery.Selection) {
			if j == 0 {
				return // skip header row
			}

			tds := tr.Find("td")
			if tds.Length() == 0 {
				return
			}

			result := tds.Eq(0).Find("a[href^='/wiki/']").First().Text()
			
			// Get element name and tier only
			element := Element{
				Name: result,
				Tier: adjustTier(tier),
			}

			if !contains(filters, result) {
				elements = append(elements, element)
			}

			// Get recipes
			if tds.Eq(1).Find("ul").Length() > 0 {
				var currentRecipe ElementRecipe
				var recipeItems []string

				tds.Eq(1).Find("a[href^='/wiki/']").Each(func(k int, a *goquery.Selection) {
					item := a.Text()
					recipeItems = append(recipeItems, item)

					if k%2 == 1 {
						currentRecipe = ElementRecipe{
							Tier:   adjustTier(tier),
							Result: result,
							Recipe: recipeItems,
						}

						valid := true
						for _, filter := range filters {
							if filter == currentRecipe.Result || contains(currentRecipe.Recipe, filter) {
								valid = false
								break
							}
						}

						if valid {
							recipes = append(recipes, currentRecipe)
						}

						recipeItems = []string{}
					}
				})
			}
		})
		tier++
	})

	return recipes, elements
}

func adjustTier(tier int) int {
	if tier > 1 {
		return tier - 1
	}
	return tier
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func writeToFile(filename string, data interface{}) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}