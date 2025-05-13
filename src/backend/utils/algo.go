package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	// "time"
)

type Recipe struct {
	Tier   int      `json:"tier"`
	Result string   `json:"result"`
	Recipe []string `json:"recipe"`
}

type RecipeStep struct {
	Ingredient1 string
	Ingredient2 string
	Result      string
}

type TreeNode struct {
	Element    string
	Children   []*TreeNode
	RecipeStep *RecipeStep
}

type TreeStats struct {
	NodeCount  int
	MaxDepth   int
	VisitCount int
}

type RecipePath struct {
	Steps    []RecipeStep
	TreeRoot *TreeNode
}

type JSONRecipeNode struct {
	Name   string               `json:"name"`
	Recipes [][2]*JSONRecipeNode `json:"recipes,omitempty"`
}

type JSONResponse struct {
	Data        *JSONRecipeNode `json:"data"`
	Errors      []string        `json:"errors"`
	Time        int64           `json:"time"`       
	NodeCount   int             `json:"nodeCount"`   
	RecipeFound int             `json:"recipeFound"` 
}

var baseElements = map[string]bool{
	"Air":   true,
	"Earth": true,
	"Fire":  true,
	"Water": true,
}

var (
	graph map[string][][2]string
	tiers map[string]int
)

func LoadRecipes(filename string) {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading recipes file: %v\n", err)
		return
	}

	var recipes []Recipe
	if err := json.Unmarshal(file, &recipes); err != nil {
		fmt.Printf("Error unmarshaling recipes: %v\n", err)
		return
	}

	graph = make(map[string][][2]string)
	tiers = make(map[string]int)

	for base := range baseElements {
		tiers[base] = 0
	}

	for _, r := range recipes {
		if len(r.Recipe) == 2 {
			graph[r.Result] = append(graph[r.Result], [2]string{r.Recipe[0], r.Recipe[1]})
		}
		if _, exists := tiers[r.Result]; !exists {
			tiers[r.Result] = r.Tier
		}
	}
}

func findRecipes(ing1, ing2 string, graph map[string][][2]string) []string {
	var results []string

	for result, recipes := range graph {
		for _, recipe := range recipes {
			if (recipe[0] == ing1 && recipe[1] == ing2) ||
				(recipe[0] == ing2 && recipe[1] == ing1) {
				results = append(results, result)
			}
		}
	}

	return results
}

func BFS(target string, graph map[string][][2]string, tiers map[string]int, maxRecipes int) ([]RecipePath, int) {
	craftable := make(map[string]bool)
	visited := make(map[string]bool)
	recipeVariants := make(map[string][]RecipeStep)
	visitCount := 0

	// Initialize base elements
	for base := range baseElements {
		craftable[base] = true
		visited[base] = true
	}

	queue := make([]string, 0, len(baseElements))
	for base := range baseElements {
		queue = append(queue, base)
	}

	for len(queue) > 0 && len(recipeVariants[target]) < maxRecipes {
		current := queue[0]
		queue = queue[1:]
		visitCount++

		for ingredient := range craftable {
			possibleResults := findRecipes(current, ingredient, graph)
			for _, result := range possibleResults {
				resultTier := tiers[result]
				currentTier := tiers[current]
				ingredientTier := tiers[ingredient]

				if resultTier > currentTier && resultTier > ingredientTier {
					newRecipe := RecipeStep{
						Ingredient1: current,
						Ingredient2: ingredient,
						Result:      result,
					}

					if len(recipeVariants[result]) < maxRecipes {
						isDuplicate := false
						for _, existing := range recipeVariants[result] {
							if (existing.Ingredient1 == current && existing.Ingredient2 == ingredient) ||
								(existing.Ingredient1 == ingredient && existing.Ingredient2 == current) {
								isDuplicate = true
								break
							}
						}
						if !isDuplicate {
							recipeVariants[result] = append(recipeVariants[result], newRecipe)
						}
					}

					if !craftable[result] {
						craftable[result] = true
						visited[result] = true
						queue = append(queue, result)
					}
				}
			}
		}
	}

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, visitCount
	}

	var allPaths []RecipePath
	processedCount := 0

	if maxRecipes > 1 {
		resultChan := make(chan RecipePath, maxRecipes)
		var wg sync.WaitGroup
		maxWorkers := 3 // gtw diatas >3 error

		sem := make(chan struct{}, maxWorkers)

		recipesToProcess := 0
		for range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++
			recipesToProcess++
		}

		processedCount = 0

		for _, recipeVariant := range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++

			wg.Add(1)
			sem <- struct{}{}
			go func(recipe RecipeStep) {
				defer wg.Done()
				defer func() { <-sem }()

				recipePaths := alternativeTree(recipe, recipeVariants, maxRecipes)
				for _, path := range recipePaths {
					resultChan <- path
				}
			}(recipeVariant)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		collectedCount := 0
		for path := range resultChan {
			allPaths = append(allPaths, path)
			collectedCount++
			if collectedCount >= recipesToProcess*maxRecipes {
				break
			}
		}
	} else {
		for _, recipeVariant := range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++

			recipePaths := alternativeTree(recipeVariant, recipeVariants, 1)
			allPaths = append(allPaths, recipePaths...)
		}
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > 0 {
		uniquePaths := []RecipePath{allPaths[0]}
		for i := 1; i < len(allPaths); i++ {
			isDuplicate := false
			for j := 0; j < len(uniquePaths); j++ {
				if isRecipeEqual(allPaths[i], uniquePaths[j]) {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				uniquePaths = append(uniquePaths, allPaths[i])
			}
		}
		allPaths = uniquePaths
	}

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
}

//go sucks ga bisa compare struct
func isRecipeEqual(path1, path2 RecipePath) bool {
	if len(path1.Steps) != len(path2.Steps) {
		return false
	}

	stepsMap1 := make(map[string]struct{})
	stepsMap2 := make(map[string]struct{})

	for _, step := range path1.Steps {
		key := fmt.Sprintf("%s+%s=%s", step.Ingredient1, step.Ingredient2, step.Result)
		stepsMap1[key] = struct{}{}
	}

	for _, step := range path2.Steps {
		key := fmt.Sprintf("%s+%s=%s", step.Ingredient1, step.Ingredient2, step.Result)
		stepsMap2[key] = struct{}{}
	}

	if len(stepsMap1) != len(stepsMap2) {
		return false
	}

	for key := range stepsMap1 {
		if _, exists := stepsMap2[key]; !exists {
			return false
		}
	}

	return true
}

func DFS(target string, graph map[string][][2]string, tiers map[string]int, maxRecipes int) ([]RecipePath, int) {
	if maxRecipes <= 0 {
		maxRecipes = 1
	}

	craftable := make(map[string]bool)
	recipeVariants := make(map[string][]RecipeStep)
	visitCount := 0
	visited := make(map[string]bool)
	stack := []string{}

	for base := range baseElements {
		craftable[base] = true
		visited[base] = true
		stack = append(stack, base)
	}

	for len(stack) > 0 && len(recipeVariants[target]) < maxRecipes {
		lastIdx := len(stack) - 1
		current := stack[lastIdx]
		stack = stack[:lastIdx]
		visitCount++

		for ingredient := range craftable {
			possibleResults := findRecipes(current, ingredient, graph)

			for _, result := range possibleResults {
				resultTier := tiers[result]
				currentTier := tiers[current]
				ingredientTier := tiers[ingredient]

				if resultTier > currentTier && resultTier > ingredientTier {
					newRecipe := RecipeStep{
						Ingredient1: current,
						Ingredient2: ingredient,
						Result:      result,
					}

					if len(recipeVariants[result]) < maxRecipes {
						isDuplicate := false
						for _, existingRecipe := range recipeVariants[result] {
							if (existingRecipe.Ingredient1 == current && existingRecipe.Ingredient2 == ingredient) ||
								(existingRecipe.Ingredient1 == ingredient && existingRecipe.Ingredient2 == current) {
								isDuplicate = true
								break
							}
						}

						if !isDuplicate {
							recipeVariants[result] = append(recipeVariants[result], newRecipe)
						}
					}

					if !craftable[result] && !visited[result] {
						craftable[result] = true
						visited[result] = true
						stack = append(stack, result)
					}
				}
			}
		}
	}

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, visitCount
	}

	var allPaths []RecipePath
	processedCount := 0

	if maxRecipes > 1 {
		resultChan := make(chan RecipePath, maxRecipes)
		var wg sync.WaitGroup
		maxWorkers := 3 //same >3 error
		sem := make(chan struct{}, maxWorkers)

		for _, recipeVariant := range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++

			wg.Add(1)
			sem <- struct{}{}
			go func(recipe RecipeStep) {
				defer wg.Done()
				defer func() { <-sem }()

				recipePaths := alternativeTree(recipe, recipeVariants, maxRecipes)
				for _, path := range recipePaths {
					resultChan <- path
				}
			}(recipeVariant)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for path := range resultChan {
			allPaths = append(allPaths, path)
		}
	} else {
		for _, recipeVariant := range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++

			// Use the enhanced function that explores alternative recipes
			recipePaths := alternativeTree(recipeVariant, recipeVariants, 1)
			allPaths = append(allPaths, recipePaths...)
		}
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > 0 {
		uniquePaths := []RecipePath{allPaths[0]}
		for i := 1; i < len(allPaths); i++ {
			isDuplicate := false
			for j := 0; j < len(uniquePaths); j++ {
				if isRecipeEqual(allPaths[i], uniquePaths[j]) {
					isDuplicate = true
					break
				}
			}
			if !isDuplicate {
				uniquePaths = append(uniquePaths, allPaths[i])
			}
		}
		allPaths = uniquePaths
	}

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
}

// ini ngakalin banget punten banget basically baseRecipeMap itu recipe tree dengan semua elemennya bounded ke 1 recipe akalin bikin generator alternative
func alternativeTree(startRecipe RecipeStep, recipeVariants map[string][]RecipeStep, maxAlternatives int) []RecipePath {
	baseRecipeMap := buildIterativeRecipeMap(startRecipe, recipeVariants)

	alternativeRecipes := make(map[string][]RecipeStep)
	for elem, recipe := range baseRecipeMap {
		if baseElements[elem] {
			continue
		}

		if variants, exists := recipeVariants[elem]; exists && len(variants) > 1 {
			for _, variant := range variants {
				isCurrentVariant := (variant.Ingredient1 == recipe.Ingredient1 && variant.Ingredient2 == recipe.Ingredient2) ||
					(variant.Ingredient1 == recipe.Ingredient2 && variant.Ingredient2 == recipe.Ingredient1)

				if !isCurrentVariant {
					alternativeRecipes[elem] = append(alternativeRecipes[elem], variant)
				}
			}
		}
	}

	if len(alternativeRecipes) == 0 {
		craftingPath := make([]RecipeStep, 0, len(baseRecipeMap))
		for _, step := range baseRecipeMap {
			craftingPath = append(craftingPath, step)
		}
		treeRoot := buildCraftingTreeFromMap(startRecipe.Result, baseRecipeMap, make(map[string]bool))
		return []RecipePath{{craftingPath, treeRoot}}
	}

	var elementsToTry []string
	for elem := range alternativeRecipes {
		elementsToTry = append(elementsToTry, elem)
	}

	sort.Slice(elementsToTry, func(i, j int) bool {
		return tiers[elementsToTry[i]] < tiers[elementsToTry[j]]
	})

	maxElementsToTry := 3
	if len(elementsToTry) > maxElementsToTry {
		elementsToTry = elementsToTry[:maxElementsToTry]
	}

	allPaths := []RecipePath{}

	craftingPath := make([]RecipeStep, 0, len(baseRecipeMap))
	for _, step := range baseRecipeMap {
		craftingPath = append(craftingPath, step)
	}
	treeRoot := buildCraftingTreeFromMap(startRecipe.Result, baseRecipeMap, make(map[string]bool))
	allPaths = append(allPaths, RecipePath{craftingPath, treeRoot})

	for _, elem := range elementsToTry {
		for _, altRecipe := range alternativeRecipes[elem] {
			if len(allPaths) >= maxAlternatives {
				break
			}

			altRecipeMap := make(map[string]RecipeStep)
			for k, v := range baseRecipeMap {
				altRecipeMap[k] = v
			}

			altRecipeMap[elem] = altRecipe

			depStack := []string{altRecipe.Ingredient1, altRecipe.Ingredient2}
			for len(depStack) > 0 {
				dep := depStack[len(depStack)-1]
				depStack = depStack[:len(depStack)-1]

				if baseElements[dep] {
					continue
				}

				if _, exists := altRecipeMap[dep]; exists {
					continue
				}

				if variants, exists := recipeVariants[dep]; exists && len(variants) > 0 {
					altRecipeMap[dep] = variants[0]
					depStack = append(depStack, variants[0].Ingredient1, variants[0].Ingredient2)
				}
			}

			altCraftingPath := make([]RecipeStep, 0, len(altRecipeMap))
			for _, step := range altRecipeMap {
				altCraftingPath = append(altCraftingPath, step)
			}

			altTreeRoot := buildCraftingTreeFromMap(startRecipe.Result, altRecipeMap, make(map[string]bool))

			allPaths = append(allPaths, RecipePath{altCraftingPath, altTreeRoot})
		}
	}

	if len(elementsToTry) >= 2 && len(allPaths) < maxAlternatives {
		for i := 0; i < len(elementsToTry); i++ {
			for j := i + 1; j < len(elementsToTry); j++ {
				elem1 := elementsToTry[i]
				elem2 := elementsToTry[j]

				for _, altRecipe1 := range alternativeRecipes[elem1] {
					for _, altRecipe2 := range alternativeRecipes[elem2] {
						if len(allPaths) >= maxAlternatives {
							break
						}

						altRecipeMap := make(map[string]RecipeStep)
						for k, v := range baseRecipeMap {
							altRecipeMap[k] = v
						}

						altRecipeMap[elem1] = altRecipe1
						altRecipeMap[elem2] = altRecipe2

						depStack := []string{
							altRecipe1.Ingredient1, altRecipe1.Ingredient2,
							altRecipe2.Ingredient1, altRecipe2.Ingredient2,
						}
						for len(depStack) > 0 {
							dep := depStack[len(depStack)-1]
							depStack = depStack[:len(depStack)-1]

							if baseElements[dep] {
								continue
							}

							if _, exists := altRecipeMap[dep]; exists {
								continue
							}

							if variants, exists := recipeVariants[dep]; exists && len(variants) > 0 {
								altRecipeMap[dep] = variants[0]
								depStack = append(depStack, variants[0].Ingredient1, variants[0].Ingredient2)
							}
						}

						altCraftingPath := make([]RecipeStep, 0, len(altRecipeMap))
						for _, step := range altRecipeMap {
							altCraftingPath = append(altCraftingPath, step)
						}

						altTreeRoot := buildCraftingTreeFromMap(startRecipe.Result, altRecipeMap, make(map[string]bool))

						allPaths = append(allPaths, RecipePath{altCraftingPath, altTreeRoot})
					}
				}
			}
		}
	}

	return allPaths
}

func buildIterativeRecipeMap(recipe RecipeStep, recipeVariants map[string][]RecipeStep) map[string]RecipeStep {
	stack := []RecipeStep{recipe}
	recipeMap := make(map[string]RecipeStep)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if _, exists := recipeMap[current.Result]; !exists {
			recipeMap[current.Result] = current
			if !baseElements[current.Ingredient1] {
				if variants, exists := recipeVariants[current.Ingredient1]; exists && len(variants) > 0 {
					stack = append(stack, variants[0])
				}
			}
			if !baseElements[current.Ingredient2] {
				if variants, exists := recipeVariants[current.Ingredient2]; exists && len(variants) > 0 {
					stack = append(stack, variants[0])
				}
			}
		}
	}
	return recipeMap
}

func buildCraftingTreeFromMap(element string, recipeMap map[string]RecipeStep, path map[string]bool) *TreeNode {
	if path[element] {
		return &TreeNode{
			Element:  element,
			Children: nil,
		}
	}

	if baseElements[element] {
		return &TreeNode{
			Element:  element,
			Children: nil,
		}
	}

	recipe, exists := recipeMap[element]
	if !exists {
		return &TreeNode{
			Element:  element,
			Children: nil,
		}
	}

	newPath := make(map[string]bool, len(path)+1)
	for k, v := range path {
		newPath[k] = v
	}
	newPath[element] = true

	node := &TreeNode{
		Element:    element,
		RecipeStep: &recipe,
		Children:   make([]*TreeNode, 0),
	}

	ing1 := recipe.Ingredient1
	ing2 := recipe.Ingredient2

	node.Children = append(node.Children, buildCraftingTreeFromMap(ing1, recipeMap, newPath))
	node.Children = append(node.Children, buildCraftingTreeFromMap(ing2, recipeMap, newPath))

	return node
}

// debug
func calculateTreeStats(root *TreeNode) TreeStats {
	if root == nil {
		return TreeStats{0, 0, 0}
	}

	var stats TreeStats
	stats.NodeCount = 1

	maxChildDepth := 0
	for _, child := range root.Children {
		childStats := calculateTreeStats(child)
		stats.NodeCount += childStats.NodeCount
		stats.VisitCount += childStats.VisitCount

		if childStats.MaxDepth > maxChildDepth {
			maxChildDepth = childStats.MaxDepth
		}
	}

	stats.MaxDepth = maxChildDepth + 1
	return stats
}

func Search(target string, findShortest bool, useBFS bool, maxRecipes int) ([]RecipePath, int, int, error) {

	if useBFS {
		fmt.Println("Finding recipes using BFS...")
		recipePaths, _ := BFS(target, graph, tiers, maxRecipes)

		fmt.Printf("Number of recipes found: %d\n", len(recipePaths))
		if len(recipePaths) == 0 {
			return nil, 0, 0, fmt.Errorf("no recipes found for %s", target)
		}

		if findShortest {
			sort.Slice(recipePaths, func(i, j int) bool {
				return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
			})

			path := recipePaths[0]
			stats := calculateTreeStats(path.TreeRoot)
			return []RecipePath{path}, stats.NodeCount, 1, nil
		}

		var totalNodeCount int
		for _, path := range recipePaths {
			stats := calculateTreeStats(path.TreeRoot)
			totalNodeCount += stats.NodeCount
		}
		return recipePaths, totalNodeCount, len(recipePaths), nil
	} else {
		fmt.Println("Finding recipes using DFS...")
		recipePaths, _ := DFS(target, graph, tiers, maxRecipes)

		fmt.Printf("Number of recipes found: %d\n", len(recipePaths))
		if len(recipePaths) == 0 {
			return nil, 0, 0, fmt.Errorf("no recipes found for %s", target)
		}

		if findShortest {
			sort.Slice(recipePaths, func(i, j int) bool {
				return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
			})

			path := recipePaths[0]
			stats := calculateTreeStats(path.TreeRoot)
			return []RecipePath{path}, stats.NodeCount, 1, nil
		}

		var totalNodeCount int
		for _, path := range recipePaths {
			stats := calculateTreeStats(path.TreeRoot)
			totalNodeCount += stats.NodeCount
		}
		return recipePaths, totalNodeCount, len(recipePaths), nil
	}
}

func _convertToJSONFormat(node *TreeNode) *JSONRecipeNode {
	if node == nil {
		return nil
	}

	jsonNode := &JSONRecipeNode{
		Name: node.Element,
	}

	if len(node.Children) == 0 || baseElements[node.Element] {
		return jsonNode
	}

	var recipe [2]*JSONRecipeNode
	for i := 0; i < 2 && i < len(node.Children); i++ {
		recipe[i] = _convertToJSONFormat(node.Children[i])
	}

	jsonNode.Recipes = append(jsonNode.Recipes, recipe)
	return jsonNode
}

func ConvertToJSONFormat(recipes []RecipePath) *JSONRecipeNode {
	if recipes == nil {
		return nil
	}

	jsonNode := &JSONRecipeNode{
		Name: recipes[0].TreeRoot.Element,
	}

	for _, recipe := range recipes {
		jsonNode.Recipes = append(jsonNode.Recipes, _convertToJSONFormat(recipe.TreeRoot).Recipes...)
	}

	return jsonNode
}

func WriteTreeToJSONFile(recipes []RecipePath, filename string) error {
	jsonRoot := ConvertToJSONFormat(recipes)

	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a JSON encoder with indentation
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 2-space indentation

	// Write the JSON to file
	if err := encoder.Encode(jsonRoot); err != nil {
		return err
	}

	return nil
}
