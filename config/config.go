package cfg

import "github.com/spf13/viper"

type Config struct {
	DataFile                            string
	Output                              string
	MaxUniqueRecipeWorkersSize          int
	MaxMostPostCodeDeliveredWorkersSize int
	MaxSpecificPostCodeWorkersSize      int
	MaxRecipeListWorkersSize            int
	MatchedRecipes                      []string
}

func NewConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	v.AutomaticEnv()

	setDefaults(v)

	return &Config{
		DataFile:                            v.GetString("data_file"),
		Output:                              v.GetString("output"),
		MaxUniqueRecipeWorkersSize:          v.GetInt("max_unique_recipe_workers_size"),
		MaxMostPostCodeDeliveredWorkersSize: v.GetInt("max_most_post_code_delivered_workers_size"),
		MaxSpecificPostCodeWorkersSize:      v.GetInt("max_specific_post_code_workers_size"),
		MaxRecipeListWorkersSize:            v.GetInt("max_recipe_list_workers_size"),
		MatchedRecipes:                      v.GetStringSlice("matched_recipes"),
	}

}

func setDefaults(v *viper.Viper) {
	v.SetDefault("data_file", "data2.json")
	v.SetDefault("output", "stdout")
	v.SetDefault("max_unique_recipe_workers_size", 100)
	v.SetDefault("max_most_post_code_delivered_workers_size", 100)
	v.SetDefault("max_specific_post_code_workers_size", 100)
	v.SetDefault("max_recipe_list_workers_size", 100)
	v.SetDefault("matched_recipes", []string{"Potato", "Veggie", "Mushroom"})
}
