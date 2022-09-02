package service

import (
	cfg "HelloFresh/config"
	"HelloFresh/pkg/model"
	"HelloFresh/tools"
	"sort"
	"sync"
)

type RecipeChecker struct {
	config          *cfg.Config
	recipeChecker   chan string
	uniqueRecipes   sync.Map
	rendererChannel *model.RendererChannel
}

func NewRecipeChecker(config *cfg.Config, recipeChecker chan string, rendererChannel *model.RendererChannel) model.Runner {
	return &RecipeChecker{
		config:          config,
		recipeChecker:   recipeChecker,
		uniqueRecipes:   sync.Map{},
		rendererChannel: rendererChannel,
	}
}

func (r *RecipeChecker) Run() error {
	var wg sync.WaitGroup
	for i := 0; i < r.config.MaxUniqueRecipeWorkersSize; i++ {
		wg.Add(1)
		go r.Worker(&wg)
	}
	wg.Wait()

	err := r.sendToRenderer()
	if err != nil {
		return err
	}
	return nil
}

func (r *RecipeChecker) Worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for recipe := range r.recipeChecker {
		numberOfOccurrence, ok := r.uniqueRecipes.Load(recipe)
		if !ok {
			r.uniqueRecipes.Store(recipe, 1)
		} else {
			r.uniqueRecipes.Store(recipe, numberOfOccurrence.(int)+1)
		}
	}

}

func (r *RecipeChecker) sendToRenderer() error {
	recipeNames := make([]string, 0)
	r.uniqueRecipes.Range(func(key, value interface{}) bool {
		recipeNames = append(recipeNames, key.(string))
		return true
	})
	sort.Strings(recipeNames)
	var countPerRecipes []model.CountPerRecipe
	var matchedRecipes []string
	for _, recipeName := range recipeNames {
		if tools.StringContains(recipeName, r.config.MatchedRecipes...) {
			matchedRecipes = append(matchedRecipes, recipeName)
		}
		recipeOccurrence, _ := r.uniqueRecipes.Load(recipeName)
		countPerRecipes = append(countPerRecipes, model.CountPerRecipe{
			Recipe: recipeName,
			Count:  recipeOccurrence.(int),
		})
	}
	for _, matchedRecipe := range matchedRecipes {
		r.rendererChannel.MatchByName <- matchedRecipe
	}
	close(r.rendererChannel.MatchByName)
	for _, countPerRecipe := range countPerRecipes {
		r.rendererChannel.CountPerRecipe <- countPerRecipe
	}
	close(r.rendererChannel.CountPerRecipe)
	r.rendererChannel.UniqueRecipeCount <- len(recipeNames)
	close(r.rendererChannel.UniqueRecipeCount)

	return nil
}
