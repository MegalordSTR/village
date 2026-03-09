package economy

// ProductionChainStep represents a single step in a production chain.
type ProductionChainStep struct {
	RecipeID string                `json:"recipe_id"`
	Recipe   *Recipe               `json:"recipe,omitempty"`
	Inputs   []ResourceRequirement `json:"inputs"`
	Outputs  []ResourceOutput      `json:"outputs"`
	Depth    int                   `json:"depth"` // depth from final product (0 = final step)
}

// ProductionChain returns the chain of recipes needed to produce the given resource type.
// It returns a slice of steps ordered from raw materials to final product.
func ProductionChain(output ResourceType) []ProductionChainStep { //nolint:gocognit
	// Build mapping from output type to recipe(s) that produce it.
	recipes := AllRecipes()
	outputToRecipe := make(map[ResourceType][]*Recipe)
	for i := range recipes {
		r := &recipes[i]
		for _, out := range r.Outputs {
			outputToRecipe[out.Type] = append(outputToRecipe[out.Type], r)
		}
	}

	// We'll perform a BFS backwards from target output.
	visited := make(map[string]bool)
	var steps []ProductionChainStep
	var queue []*Recipe

	// Find recipes that produce the target output.
	for _, r := range outputToRecipe[output] {
		if visited[r.ID] {
			continue
		}
		queue = append(queue, r)
		visited[r.ID] = true
	}

	// While queue not empty, process recipe and add its input recipes.
	for len(queue) > 0 {
		r := queue[0]
		queue = queue[1:]

		step := ProductionChainStep{
			RecipeID: r.ID,
			Recipe:   r,
			Inputs:   r.Inputs,
			Outputs:  r.Outputs,
			Depth:    0, // placeholder, will be computed later
		}
		steps = append(steps, step)

		// For each input resource type, find recipes that produce it.
		for _, inp := range r.Inputs {
			for _, prodRecipe := range outputToRecipe[inp.Type] {
				if visited[prodRecipe.ID] {
					continue
				}
				queue = append(queue, prodRecipe)
				visited[prodRecipe.ID] = true
			}
		}
	}

	// Compute depths using topological sort (simplified: assign depth based on longest chain).
	// For simplicity, we'll compute depth via recursion.
	// Build adjacency list: recipe -> dependent recipes (recipes that use its output).
	dependents := make(map[string][]string)
	for _, r := range recipes {
		for _, inp := range r.Inputs {
			for _, prodRecipe := range outputToRecipe[inp.Type] {
				dependents[prodRecipe.ID] = append(dependents[prodRecipe.ID], r.ID)
			}
		}
	}

	// Compute depth via memoized DFS.
	depthCache := make(map[string]int)
	var dfs func(string) int
	dfs = func(rid string) int {
		if d, ok := depthCache[rid]; ok {
			return d
		}
		maxDepth := 0
		for _, dep := range dependents[rid] {
			d := dfs(dep)
			if d > maxDepth {
				maxDepth = d
			}
		}
		depthCache[rid] = maxDepth + 1
		return maxDepth + 1
	}

	// Update step depths.
	for i := range steps {
		steps[i].Depth = dfs(steps[i].RecipeID)
	}

	// Sort steps by depth ascending (raw materials first).
	// Simple bubble sort for small n.
	for i := 0; i < len(steps); i++ {
		for j := i + 1; j < len(steps); j++ {
			if steps[i].Depth > steps[j].Depth {
				steps[i], steps[j] = steps[j], steps[i]
			}
		}
	}

	return steps
}

// BottleneckResources returns a list of resources that are required by multiple
// recipes in the chain and might become bottlenecks.
func BottleneckResources(chain []ProductionChainStep) []ResourceType {
	required := make(map[ResourceType]float64)
	for _, step := range chain {
		for _, inp := range step.Inputs {
			required[inp.Type] += inp.Quantity
		}
	}
	var bottlenecks []ResourceType
	for rt, total := range required {
		if total > 100 { // arbitrary threshold
			bottlenecks = append(bottlenecks, rt)
		}
	}
	return bottlenecks
}
