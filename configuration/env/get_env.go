package env

import "os"

func GetNewsTokenApi() string {
	return os.Getenv("NEWS_API_KEY")
}
