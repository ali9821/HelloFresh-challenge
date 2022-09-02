package cfg

import "github.com/spf13/viper"

type Config struct {
	DataFile                            string
	Output                              string
	MaxUniqueRecipeWorkersSize          int
	MaxMostPostCodeDeliveredWorkersSize int
	MaxSpecificPostCodeWorkersSize      int
	MaxRecipeListWorkersSize            int
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
		Output:                              v.GetString("output_format"),
		MaxUniqueRecipeWorkersSize:          v.GetInt("max_unique_recipe_workers_size"),
		MaxMostPostCodeDeliveredWorkersSize: v.GetInt("max_most_post_code_delivered_workers_size"),
		MaxSpecificPostCodeWorkersSize:      v.GetInt("max_specific_post_code_workers_size"),
		MaxRecipeListWorkersSize:            v.GetInt("max_recipe_list_workers_size"),
	}

}

func setDefaults(v *viper.Viper) {
	v.SetDefault("data_file", "data.json")
	v.SetDefault("output_format", "stdout")
	v.SetDefault("max_unique_recipe_workers_size", 100)
	v.SetDefault("max_most_post_code_delivered_workers_size", 100)
	v.SetDefault("max_specific_post_code_workers_size", 100)
	v.SetDefault("max_recipe_list_workers_size", 100)
}
