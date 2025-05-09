package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

type Recipe struct {
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

var baseElements = map[string]bool{
	"Air":   true,
	"Earth": true,
	"Fire":  true,
	"Water": true,
}

func loadRecipes(filename string) map[string][][2]string {
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading recipes file: %v\n", err)
		return nil
	}

	var recipes []Recipe
	if err := json.Unmarshal(file, &recipes); err != nil {
		fmt.Printf("Error unmarshaling recipes: %v\n", err)
		return nil
	}

	graph := make(map[string][][2]string)
	for _, r := range recipes {
		if len(r.Recipe) == 2 {
			graph[r.Result] = append(graph[r.Result], [2]string{r.Recipe[0], r.Recipe[1]})
		}
	}

	return graph
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

func BFS(target string, graph map[string][][2]string, maxRecipes int) ([]RecipePath, int) {
	if maxRecipes <= 0 {
		maxRecipes = 1
	}

	craftable := make(map[string]bool)
	recipeVariants := make(map[string][]RecipeStep)
	visitCount := 0

	for base := range baseElements {
		craftable[base] = true
	}

	queue := []string{}
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

				if !craftable[result] {
					craftable[result] = true
					queue = append(queue, result)
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

	resultChan := make(chan RecipePath, maxRecipes)
	var wg sync.WaitGroup

	for _, recipeVariant := range recipeVariants[target] {
		if processedCount >= maxRecipes {
			break
		}
		processedCount++

		wg.Add(1)
		go func(recipe RecipeStep) {
			defer wg.Done()

			recipeMap := make(map[string]RecipeStep)
			visited := make(map[string]bool)

			buildRecursiveRecipeMap(recipe, recipeVariants, recipeMap, visited, 50)

			var craftingPath []RecipeStep
			for element, step := range recipeMap {
				if element != "" { // Skip any empty keys
					craftingPath = append(craftingPath, step)
				}
			}

			treeRoot := buildCraftingTreeFromMap(target, recipeMap)

			resultChan <- RecipePath{craftingPath, treeRoot}
		}(recipeVariant)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for path := range resultChan {
		allPaths = append(allPaths, path)
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
}

func buildRecursiveRecipeMap(recipe RecipeStep, recipeVariants map[string][]RecipeStep, recipeMap map[string]RecipeStep, visited map[string]bool, maxDepth int) {
	if maxDepth <= 0 {
		fmt.Println("Warning: Maximum recursion depth reached. Recipe chain might be incomplete.")
		return
	}

	if visited[recipe.Result] {
		return
	}

	visited[recipe.Result] = true
	recipeMap[recipe.Result] = recipe

	if !baseElements[recipe.Ingredient1] {
		if len(recipeVariants[recipe.Ingredient1]) > 0 {
			buildRecursiveRecipeMap(recipeVariants[recipe.Ingredient1][0], recipeVariants, recipeMap, visited, maxDepth-1)
		} else {
			fmt.Printf("Warning: No recipe found for intermediate ingredient %s\n", recipe.Ingredient1)
		}
	}

	if !baseElements[recipe.Ingredient2] {
		if len(recipeVariants[recipe.Ingredient2]) > 0 {
			buildRecursiveRecipeMap(recipeVariants[recipe.Ingredient2][0], recipeVariants, recipeMap, visited, maxDepth-1)
		} else {
			fmt.Printf("Warning: No recipe found for intermediate ingredient %s\n", recipe.Ingredient2)
		}
	}
}

func DFS(target string, graph map[string][][2]string, maxRecipes int) ([]RecipePath, int) {
	if maxRecipes <= 0 {
		maxRecipes = 1 // Default to at least one recipe
	}

	craftable := make(map[string]bool)
	recipeVariants := make(map[string][]RecipeStep)
	visitCount := 0

	for base := range baseElements {
		craftable[base] = true
	}

	visited := make(map[string]bool)
	stack := []string{}

	for base := range baseElements {
		stack = append(stack, base)
	}

	for len(stack) > 0 && len(recipeVariants[target]) < maxRecipes {
		lastIdx := len(stack) - 1
		current := stack[lastIdx]
		stack = stack[:lastIdx]

		visitCount++

		if visited[current] {
			continue
		}

		visited[current] = true

		for ingredient := range craftable {
			possibleResults := findRecipes(current, ingredient, graph)

			for _, result := range possibleResults {
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

				if !craftable[result] {
					craftable[result] = true
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

	resultChan := make(chan RecipePath, maxRecipes)
	var wg sync.WaitGroup

	for _, recipeVariant := range recipeVariants[target] {
		if processedCount >= maxRecipes {
			break
		}
		processedCount++

		wg.Add(1)
		go func(recipe RecipeStep) {
			defer wg.Done()

			recipeMap := make(map[string]RecipeStep)
			visited := make(map[string]bool)

			buildFixedDependencyGraph(recipe, recipeVariants, recipeMap, visited, 50)

			var craftingPath []RecipeStep
			for element, step := range recipeMap {
				if element != "" { // Skip any empty keys
					craftingPath = append(craftingPath, step)
				}
			}

			treeRoot := buildCraftingTreeFromMap(target, recipeMap)

			resultChan <- RecipePath{craftingPath, treeRoot}
		}(recipeVariant)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for path := range resultChan {
		allPaths = append(allPaths, path)
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i].Steps) < len(allPaths[j].Steps)
	})

	if len(allPaths) > maxRecipes {
		allPaths = allPaths[:maxRecipes]
	}

	return allPaths, visitCount
}

func buildFixedDependencyGraph(recipe RecipeStep, recipeVariants map[string][]RecipeStep, recipeMap map[string]RecipeStep, visited map[string]bool, maxDepth int) {
	if maxDepth <= 0 {
		fmt.Println("Warning: Maximum recursion depth reached. Recipe chain might be incomplete.")
		return
	}

	if _, exists := recipeMap[recipe.Result]; !exists {
		recipeMap[recipe.Result] = recipe
	}

	if !baseElements[recipe.Ingredient1] {
		if len(recipeVariants[recipe.Ingredient1]) > 0 {
			subRecipe := recipeVariants[recipe.Ingredient1][0]

			if !visited[recipe.Ingredient1] {
				visited[recipe.Ingredient1] = true
				buildFixedDependencyGraph(subRecipe, recipeVariants, recipeMap, visited, maxDepth-1)
				visited[recipe.Ingredient1] = false // Allow revisiting for different paths
			}
		} else {
			fmt.Printf("Warning: No recipe found for ingredient %s\n", recipe.Ingredient1)
		}
	}

	if !baseElements[recipe.Ingredient2] {
		if len(recipeVariants[recipe.Ingredient2]) > 0 {
			subRecipe := recipeVariants[recipe.Ingredient2][0]

			if !visited[recipe.Ingredient2] {
				visited[recipe.Ingredient2] = true
				buildFixedDependencyGraph(subRecipe, recipeVariants, recipeMap, visited, maxDepth-1)
				visited[recipe.Ingredient2] = false // Allow revisiting for different paths
			}
		} else {
			fmt.Printf("Warning: No recipe found for ingredient %s\n", recipe.Ingredient2)
		}
	}
}

func explorePathDFS(element string, recipeFor map[string]RecipeStep, visited map[string]bool, depthMap map[string]int, depth int) {
	if visited[element] || baseElements[element] {
		return
	}

	visited[element] = true
	depthMap[element] = depth

	recipe, exists := recipeFor[element]
	if !exists {
		return
	}

	ing1, ing2 := recipe.Ingredient1, recipe.Ingredient2

	explorePathDFS(ing1, recipeFor, visited, depthMap, depth+1)
	explorePathDFS(ing2, recipeFor, visited, depthMap, depth+1)
}

func buildCraftingTree(target string, recipeFor map[string]RecipeStep, dependencies map[string][]string) *TreeNode {
	if baseElements[target] {
		return &TreeNode{
			Element:  target,
			Children: nil,
		}
	}

	recipe, exists := recipeFor[target]
	if !exists {
		return &TreeNode{
			Element:  target,
			Children: nil,
		}
	}

	node := &TreeNode{
		Element:    target,
		RecipeStep: &recipe,
		Children:   make([]*TreeNode, 0),
	}

	ing1 := recipe.Ingredient1
	ing2 := recipe.Ingredient2

	node.Children = append(node.Children, buildCraftingTree(ing1, recipeFor, dependencies))
	node.Children = append(node.Children, buildCraftingTree(ing2, recipeFor, dependencies))

	return node
}

func buildCraftingTreeDFS(element string, bestRecipe map[string]RecipeStep) *TreeNode {
	if baseElements[element] {
		return &TreeNode{
			Element:  element,
			Children: nil,
		}
	}

	recipe, exists := bestRecipe[element]
	if !exists {
		return &TreeNode{
			Element:  element,
			Children: nil,
		}
	}

	node := &TreeNode{
		Element:    element,
		RecipeStep: &recipe,
		Children:   make([]*TreeNode, 0),
	}

	ing1 := recipe.Ingredient1
	ing2 := recipe.Ingredient2

	node.Children = append(node.Children, buildCraftingTreeDFS(ing1, bestRecipe))
	node.Children = append(node.Children, buildCraftingTreeDFS(ing2, bestRecipe))

	return node
}

func buildCraftingTreeFromMap(element string, recipeMap map[string]RecipeStep) *TreeNode {
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

	node := &TreeNode{
		Element:    element,
		RecipeStep: &recipe,
		Children:   make([]*TreeNode, 0),
	}

	ing1 := recipe.Ingredient1
	ing2 := recipe.Ingredient2

	node.Children = append(node.Children, buildCraftingTreeFromMap(ing1, recipeMap))
	node.Children = append(node.Children, buildCraftingTreeFromMap(ing2, recipeMap))

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

// debug
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

//claude from driver :3
// func main() {
// 	graph := loadRecipes("recipes.json")

// 	target := "Mailbox"
// 	findShortest := false
// 	useBFS := true
// 	maxRecipes := 6

// 	startTime := time.Now()

// 	if findShortest {
// 		if useBFS {
// 			fmt.Println("Finding shortest recipe using BFS...")
// 			recipePaths, visitCount := BFS(target, graph, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				sort.Slice(recipePaths, func(i, j int) bool {
// 					return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
// 				})

// 				path := recipePaths[0]
// 				fmt.Printf("\nFound shortest recipe for %s using BFS:\n", target)
// 				for i, step := range path.Steps {
// 					fmt.Printf("%d. Combine %s + %s = %s\n", i+1, step.Ingredient1, step.Ingredient2, step.Result)
// 				}

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
// 			recipePaths, visitCount := DFS(target, graph, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				sort.Slice(recipePaths, func(i, j int) bool {
// 					return len(recipePaths[i].Steps) < len(recipePaths[j].Steps)
// 				})

// 				path := recipePaths[0]
// 				fmt.Printf("\nFound shortest recipe for %s using DFS:\n", target)
// 				for i, step := range path.Steps {
// 					fmt.Printf("%d. Combine %s + %s = %s\n", i+1, step.Ingredient1, step.Ingredient2, step.Result)
// 				}

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
// 			recipePaths, visitCount := BFS(target, graph, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				fmt.Printf("\nFound %d recipes for %s using BFS:\n", len(recipePaths), target)
// 				for _, path := range recipePaths {
// 					fmt.Println("\nRecipe Path:")
// 					for i, step := range path.Steps {
// 						fmt.Printf("%d. Combine %s + %s = %s\n", i+1, step.Ingredient1, step.Ingredient2, step.Result)
// 					}

// 					fmt.Println("\nRecipe Tree:")
// 					printTreeAsHeap(path.TreeRoot, "", true)

// 					stats := calculateTreeStats(path.TreeRoot)
// 					fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 					fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 					fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 					fmt.Printf("Visited Nodes: %d\n", visitCount)
// 				}
// 			}
// 		} else {
// 			fmt.Printf("Finding up to %d recipes using DFS...\n", maxRecipes)
// 			recipePaths, visitCount := DFS(target, graph, maxRecipes)
// 			if len(recipePaths) > 0 {
// 				fmt.Printf("\nFound %d recipes for %s using DFS:\n", len(recipePaths), target)
// 				for _, path := range recipePaths {
// 					fmt.Println("\nRecipe Path:")
// 					for i, step := range path.Steps {
// 						fmt.Printf("%d. Combine %s + %s = %s\n", i+1, step.Ingredient1, step.Ingredient2, step.Result)
// 					}

// 					fmt.Println("\nRecipe Tree:")
// 					printTreeAsHeap(path.TreeRoot, "", true)

// 					stats := calculateTreeStats(path.TreeRoot)
// 					fmt.Printf("\nTree Statistics for crafting %s:\n", target)
// 					fmt.Printf("Total Nodes: %d\n", stats.NodeCount)
// 					fmt.Printf("Maximum Depth: %d\n", stats.MaxDepth)
// 					fmt.Printf("Visited Nodes: %d\n", visitCount)
// 				}
// 			}
// 		}
// 	}
// 	endTime := time.Now()
// 	elapsedTime := endTime.Sub(startTime)
// 	fmt.Printf("\nElapsed Time: %s\n", elapsedTime)
// }
