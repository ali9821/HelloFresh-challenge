package model

type Response struct {
	UniqueRecipeCount       int                     `json:"unique_recipe_count"`
	CountPerRecipe          []CountPerRecipe        `json:"count_per_recipe"`
	MatchByName             []string                `json:"match_by_name"`
	BusiestPostcode         BusiestPostcode         `json:"busiest_postcode"`
	CountPerPostcodeAndTime CountPerPostcodeAndTime `json:"count_per_postcode_and_time"`
}

type CountPerRecipe struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type BusiestPostcode struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type CountPerPostcodeAndTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

type RendererChannel struct {
	UniqueRecipeCount       chan int
	CountPerRecipe          chan CountPerRecipe
	BusiestPostcode         chan BusiestPostcode
	CountPerPostcodeAndTime chan CountPerPostcodeAndTime
	MatchByName             chan string
}
