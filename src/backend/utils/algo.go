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
	Name     string              `json:"name"`
	Recipes  [][2]*JSONRecipeNode   `json:"recipes,omitempty"`
}

type JSONResponse struct {
	Data         *JSONRecipeNode `json:"data"`
	Errors       []string    `json:"errors"`
	Time         int64       `json:"time"`          // milliseconds
	NodeCount    int         `json:"nodeCount"`     // nodes visited
	RecipeFound  int         `json:"recipeFound"`   // recipes found
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
	visited := make(map[string]bool) // New visited map
	recipeVariants := make(map[string][]RecipeStep)
	visitCount := 0

	// Initialize base elements
	for base := range baseElements {
		craftable[base] = true
		visited[base] = true // Mark base as visited
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
				if !visited[result] {
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
						visited[result] = true // Mark as visited
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
		maxWorkers := 3
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

				recipeMap := buildIterativeRecipeMap(recipe, recipeVariants)
				craftingPath := make([]RecipeStep, 0, len(recipeMap))
				for _, step := range recipeMap {
					craftingPath = append(craftingPath, step)
				}

				treeRoot := buildCraftingTreeFromMap(target, recipeMap, make(map[string]bool))
				resultChan <- RecipePath{craftingPath, treeRoot}
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
			if collectedCount >= recipesToProcess {
				break
			}
		}
	} else {
		for _, recipeVariant := range recipeVariants[target] {
			if processedCount >= maxRecipes {
				break
			}
			processedCount++

			recipeMap := buildIterativeRecipeMap(recipeVariant, recipeVariants)
			craftingPath := make([]RecipeStep, 0, len(recipeMap))
			for _, step := range recipeMap {
				craftingPath = append(craftingPath, step)
			}

			treeRoot := buildCraftingTreeFromMap(target, recipeMap, make(map[string]bool))
			allPaths = append(allPaths, RecipePath{craftingPath, treeRoot})
			recipeMap = nil
		}
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
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
			validResults := make([]string, 0)

			for _, result := range possibleResults {
				resTier := tiers[result]
				currTier := tiers[current]
				ingTier := tiers[ingredient]

				if resTier > currTier || resTier > ingTier {
					validResults = append(validResults, result)
				}
			}

			for _, result := range validResults {
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

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, visitCount
	}

	var allPaths []RecipePath
	processedCount := 0

	if maxRecipes > 1 {
		resultChan := make(chan RecipePath, maxRecipes)
		var wg sync.WaitGroup
		maxWorkers := 3
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

				recipeMap := buildIterativeRecipeMap(recipe, recipeVariants)
				craftingPath := make([]RecipeStep, 0, len(recipeMap))
				for _, step := range recipeMap {
					craftingPath = append(craftingPath, step)
				}

				treeRoot := buildCraftingTreeFromMap(target, recipeMap, make(map[string]bool))
				resultChan <- RecipePath{craftingPath, treeRoot}
				recipeMap = nil
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

			recipeMap := buildIterativeRecipeMap(recipeVariant, recipeVariants)
			craftingPath := make([]RecipeStep, 0, len(recipeMap))
			for _, step := range recipeMap {
				craftingPath = append(craftingPath, step)
			}

			treeRoot := buildCraftingTreeFromMap(target, recipeMap, make(map[string]bool))
			allPaths = append(allPaths, RecipePath{craftingPath, treeRoot})
			recipeMap = nil
		}
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
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
				stack = append(stack, recipeVariants[current.Ingredient1][0])
			}
			if !baseElements[current.Ingredient2] {
				stack = append(stack, recipeVariants[current.Ingredient2][0])
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

func printTreeAsHeap(root *TreeNode, prefix string, isLast bool) {
	if root == nil {
		return
	}

	fmt.Print(prefix)
	if isLast {
		fmt.Print("└── ")
	} else {
		fmt.Print("├── ")
	}
	fmt.Println(root.Element)

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	for i, child := range root.Children {
		isLastChild := i == len(root.Children)-1
		printTreeAsHeap(child, childPrefix, isLastChild)
	}
}

func Search(target string, findShortest bool, useBFS bool, maxRecipes int) ([]RecipePath, int, int, error) {
	var recipePaths []RecipePath

	if useBFS {
		fmt.Println("Finding shortest recipe using BFS...")
		recipePaths, _ = BFS(target, graph, tiers, maxRecipes)
	} else {
		fmt.Println("Finding shortest recipe using DFS...")
		recipePaths, _ = DFS(target, graph, tiers, maxRecipes)
	}
	
	fmt.Printf("number of recipes: %d\n", len((recipePaths)))
	if len(recipePaths) > 0 {		
		if findShortest {
			sort.Slice(recipePaths, func(i, j int) bool {
				return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
			})

			path := recipePaths[0]
			stats := calculateTreeStats(path.TreeRoot)
			return []RecipePath{path}, stats.NodeCount, len(recipePaths), nil
		} else {
			var nodeCount int

			for _, path := range recipePaths {
				stats := calculateTreeStats(path.TreeRoot)
				nodeCount += stats.NodeCount
			}
			return recipePaths, nodeCount, len(recipePaths), nil
		}
	} else {
		return nil, 0, 0, fmt.Errorf("no recipes found")
	}
}

func _convertToJSONFormat(node *TreeNode) *JSONRecipeNode {
	if node == nil {
		return nil
	}

	jsonNode := &JSONRecipeNode{
		Name: node.Element,
	}

	// Base elements or nodes without children don't have recipes
	if len(node.Children) == 0 || baseElements[node.Element] {
		return jsonNode
	}

	// Create a recipe entry with 2 ingredients
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


// Convert the tree to JSON structure
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



// func main() {
// 	graph, tiers := loadRecipes("recipes.json")
// 	if graph == nil || tiers == nil {
// 		fmt.Println("Failed to load recipes")
// 		return
// 	}

// 	target := "Wood"
// 	findShortest := true
// 	useBFS := true
// 	maxRecipes := 8

// 	startTime := time.Now()

// 	if findShortest {
// 		if useBFS {
// 			fmt.Println("Finding shortest recipe using BFS...")
// 			recipePaths, visitCount := BFS(target, graph, tiers, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				sort.Slice(recipePaths, func(i, j int) bool {
// 					return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
// 				})

// 				path := recipePaths[0]

// 				fmt.Println("\nRecipe Tree:")
// 				printTreeAsHeap(path.TreeRoot, "", true)

// 				stats := calculateTreeStats(path.TreeRoot)
// 				fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 				fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 				fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 				fmt.Printf("Visited Nodes: %d\n", visitCount)
// 			}
// 		} else {
// 			fmt.Println("Finding shortest recipe using DFS...")
// 			recipePaths, visitCount := DFS(target, graph, tiers, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				sort.Slice(recipePaths, func(i, j int) bool {
// 					return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
// 				})

// 				path := recipePaths[0]

// 				fmt.Println("\nRecipe Tree:")
// 				printTreeAsHeap(path.TreeRoot, "", true)

// 				stats := calculateTreeStats(path.TreeRoot)
// 				fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 				fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 				fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 				fmt.Printf("Visited Nodes: %d\n", visitCount)
// 			}
// 		}
// 	} else {
// 		if useBFS {
// 			fmt.Printf("Finding up to %d recipes using BFS...\n", maxRecipes)
// 			recipePaths, visitCount := BFS(target, graph, tiers, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				fmt.Printf("\nFound %d recipes for %s using BFS:\n", len(recipePaths), target)
// 				for _, path := range recipePaths {
// 					printTreeAsHeap(path.TreeRoot, "", true)
// 				}
// 				stats := calculateTreeStats(recipePaths[0].TreeRoot)
// 				fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 				fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 				fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 				fmt.Printf("Visited Nodes: %d\n", visitCount)
// 			} else {
// 				fmt.Printf("No recipes found for %s using BFS\n", target)
// 			}
// 		} else {
// 			fmt.Printf("Finding up to %d recipes using DFS...\n", maxRecipes)
// 			recipePaths, visitCount := DFS(target, graph, tiers, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				fmt.Printf("\nFound %d recipes for %s using DFS:\n", len(recipePaths), target)
// 				for _, path := range recipePaths {
// 					printTreeAsHeap(path.TreeRoot, "", true)
// 				}
// 				stats := calculateTreeStats(recipePaths[0].TreeRoot)
// 				fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 				fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 				fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 				fmt.Printf("Visited Nodes: %d\n", visitCount)
// 			} else {
// 				fmt.Printf("No recipes found for %s using DFS\n", target)
// 			}
// 		}
// 	}
// 	elapsedTime := time.Since(startTime)
// 	fmt.Printf("Elapsed Time: %s\n", elapsedTime)
// }
