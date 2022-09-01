package cfg

import "github.com/spf13/viper"

type Config struct {
	JsonFile                            string
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
		JsonFile:                            v.GetString("json_file"),
		Output:                              v.GetString("output_format"),
		MaxUniqueRecipeWorkersSize:          v.GetInt("max_UniqueRecipeWorkers_size"),
		MaxMostPostCodeDeliveredWorkersSize: v.GetInt("max_mostPostCodeDeliveredWorkers_size"),
		MaxSpecificPostCodeWorkersSize:      v.GetInt("max_specificPostCodeWorkers_size"),
		MaxRecipeListWorkersSize:            v.GetInt("max_recipeListWorkers_size"),
	}

}

func setDefaults(v *viper.Viper) {
	v.SetDefault("json_file", "data.json")
	v.SetDefault("output_format", "stdout")
	v.SetDefault("max_UniqueRecipeWorkers_size", 100)
	v.SetDefault("max_mostPostCodeDeliveredWorkers_size", 100)
	v.SetDefault("max_specificPostCodeWorkers_size", 100)
	v.SetDefault("max_recipeListWorkers_size", 100)
}
