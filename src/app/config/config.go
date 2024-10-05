package config

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 配置信息
 * @Date: 2024-10-04 19:37
 */
type Config struct {
	DB struct {
		MySqlDataSource string
	}
	Redis struct {
		Host     string
		Password string
		DB       int
	}
}
