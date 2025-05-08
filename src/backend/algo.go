package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Recipe struct {
	Result string   `json:"result"`
	Recipe []string `json:"recipe"`
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
		return nil
	}

	var recipes []Recipe
	if err := json.Unmarshal(file, &recipes); err != nil {
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

func findRecipes(ing1, ing2 string, graph map[string][][2]string) []string {
	var results []string

	for result, recipes := range graph {
		for _, recipe := range recipes {
			// Check both combinations: (ing1, ing2) and (ing2, ing1)
			if (recipe[0] == ing1 && recipe[1] == ing2) ||
				(recipe[0] == ing2 && recipe[1] == ing1) {
				results = append(results, result)
			}
		}
	}

	return results
}

func BFS(target string, graph map[string][][2]string) ([]RecipeStep, *TreeNode) {
	craftable := make(map[string]bool)
	recipeFor := make(map[string]RecipeStep)

	for base := range baseElements {
		craftable[base] = true
	}

	queue := []string{}

	for base := range baseElements {
		queue = append(queue, base)
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for ingredient := range craftable {
			possibleCombinations := findRecipes(current, ingredient, graph)

			for _, result := range possibleCombinations {
				if !craftable[result] {
					craftable[result] = true
					queue = append(queue, result)

					recipeFor[result] = RecipeStep{
						Ingredient1: current,
						Ingredient2: ingredient,
						Result:      result,
					}
				}
			}
		}
	}

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, nil
	}

	var craftingPath []RecipeStep
	visited := make(map[string]bool)

	pathQueue := []string{target}
	levelMap := make(map[string]int)
	levelMap[target] = 0

	for len(pathQueue) > 0 {
		current := pathQueue[0]
		pathQueue = pathQueue[1:]

		if baseElements[current] || visited[current] {
			continue
		}

		recipe, exists := recipeFor[current]
		if !exists {
			continue
		}

		visited[current] = true

		ing1, ing2 := recipe.Ingredient1, recipe.Ingredient2

		if !baseElements[ing1] && !visited[ing1] {
			pathQueue = append(pathQueue, ing1)
			levelMap[ing1] = levelMap[current] + 1
		}

		if !baseElements[ing2] && !visited[ing2] {
			pathQueue = append(pathQueue, ing2)
			levelMap[ing2] = levelMap[current] + 1
		}
	}

	type ElementLevel struct {
		Element string
		Level   int
	}

	var elementsByLevel []ElementLevel
	for element := range visited {
		elementsByLevel = append(elementsByLevel, ElementLevel{
			Element: element,
			Level:   levelMap[element],
		})
	}

	sort.Slice(elementsByLevel, func(i, j int) bool {
		return elementsByLevel[i].Level > elementsByLevel[j].Level
	})

	visitedForPath := make(map[string]bool)
	for _, el := range elementsByLevel {
		element := el.Element
		if !visitedForPath[element] {
			recipe := recipeFor[element]
			craftingPath = append(craftingPath, recipe)
			visitedForPath[element] = true
		}
	}

	treeRoot := buildCraftingTree(target, recipeFor, make(map[string][]string))

	return craftingPath, treeRoot
}

func DFS(target string, graph map[string][][2]string) ([]RecipeStep, *TreeNode) {
	craftable := make(map[string]bool)
	recipeFor := make(map[string]RecipeStep)

	for base := range baseElements {
		craftable[base] = true
	}

	visited := make(map[string]bool)
	stack := []string{}

	for base := range baseElements {
		stack = append(stack, base)
	}

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		current := stack[lastIdx]
		stack = stack[:lastIdx]

		if visited[current] {
			continue
		}

		visited[current] = true

		for ingredient := range craftable {
			possibleResults := findRecipes(current, ingredient, graph)

			for _, result := range possibleResults {
				if !craftable[result] {
					craftable[result] = true
					stack = append(stack, result)

					recipeFor[result] = RecipeStep{
						Ingredient1: current,
						Ingredient2: ingredient,
						Result:      result,
					}
				}
			}
		}
	}

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, nil
	}

	var craftingPath []RecipeStep
	visitedPath := make(map[string]bool)

	depthMap := make(map[string]int)

	explorePathDFS(target, recipeFor, visitedPath, depthMap, 0)

	type ElementDepth struct {
		Element string
		Depth   int
	}

	var elementsByDepth []ElementDepth
	for element := range visitedPath {
		if element != target && !baseElements[element] {
			elementsByDepth = append(elementsByDepth, ElementDepth{
				Element: element,
				Depth:   depthMap[element],
			})
		}
	}

	sort.Slice(elementsByDepth, func(i, j int) bool {
		return elementsByDepth[i].Depth > elementsByDepth[j].Depth
	})

	for _, el := range elementsByDepth {
		element := el.Element
		recipe := recipeFor[element]
		craftingPath = append(craftingPath, recipe)
	}

	if !baseElements[target] {
		craftingPath = append(craftingPath, recipeFor[target])
	}

	treeRoot := buildCraftingTreeDFS(target, recipeFor)

	return craftingPath, treeRoot
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

// func main() {
// 	graph := loadRecipes("recipes.json")

// 	target := "Mailbox"
// 	_, tree := DFS(target, graph)

// 	printTreeAsHeap(tree, "", true)

// }
