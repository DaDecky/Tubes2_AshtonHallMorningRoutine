package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func findOptimalCraftingPath(target string, graph map[string][][2]string) ([]RecipeStep, *TreeNode) {
	craftable := make(map[string]bool)
	for base := range baseElements {
		craftable[base] = true
	}

	usedInPath := make(map[string]bool)
	for base := range baseElements {
		usedInPath[base] = true
	}

	recipeFor := make(map[string]RecipeStep)

	dependencies := make(map[string][]string)

	queue := []string{}

	for result, recipes := range graph {
		for _, recipe := range recipes {
			ing1, ing2 := recipe[0], recipe[1]
			if baseElements[ing1] && baseElements[ing2] {
				craftable[result] = true
				step := RecipeStep{
					Ingredient1: ing1,
					Ingredient2: ing2,
					Result:      result,
				}
				recipeFor[result] = step
				dependencies[result] = []string{ing1, ing2}
				queue = append(queue, result)
			}
		}
	}

	for len(queue) > 0 {
		queue = queue[1:]

		for result, recipes := range graph {
			if craftable[result] {
				continue
			}

			for _, recipe := range recipes {
				ing1, ing2 := recipe[0], recipe[1]

				// If we can craft both ingredients
				if craftable[ing1] && craftable[ing2] {
					craftable[result] = true

					step := RecipeStep{
						Ingredient1: ing1,
						Ingredient2: ing2,
						Result:      result,
					}

					recipeFor[result] = step
					dependencies[result] = []string{ing1, ing2}
					queue = append(queue, result)

					if result == target {
						break
					}
				}
			}
		}
	}

	if !craftable[target] {
		fmt.Printf("Cannot craft %s from base elements\n", target)
		return nil, nil
	}

	var finalPath []RecipeStep

	var addRecipesToPath func(item string)
	addRecipesToPath = func(item string) {
		if baseElements[item] || usedInPath[item] {
			return
		}

		deps := dependencies[item]
		for _, dep := range deps {
			addRecipesToPath(dep)
		}

		finalPath = append(finalPath, recipeFor[item])
		usedInPath[item] = true
	}

	addRecipesToPath(target)

	treeRoot := buildCraftingTree(target, recipeFor, dependencies)

	return finalPath, treeRoot
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

func printCraftingPath(steps []RecipeStep) {
	if steps == nil || len(steps) == 0 {
		fmt.Println("No valid crafting path found.")
		return
	}

	fmt.Println("\nOptimized Crafting Path:")
	for i, step := range steps {
		fmt.Printf("Step %d: %s + %s = %s\n", i+1, step.Ingredient1, step.Ingredient2, step.Result)
	}
}

func printTreeAsHeap(root *TreeNode, prefix string, isLast bool) {
	if root == nil {
		return
	}

	fmt.Println(prefix + "└── " + root.Element)

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

// 	target := "Light"
// 	path, tree := findOptimalCraftingPath(target, graph)
// 	printCraftingPath(path)

// 	fmt.Println("\nCrafting Tree (Nested View):")
// 	printTreeAsHeap(tree, "", true)
// }
