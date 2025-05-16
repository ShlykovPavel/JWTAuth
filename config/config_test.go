package config

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnv(t *testing.T) {
	//Устанавливаем переменные окружаения
	t.Setenv("ENV", "production")
	t.Setenv("AUTH_USERNAME", "testUser")
	t.Setenv("AUTH_PASSWORD", "testPass")
	t.Setenv("AUTH_RETRY_COUNT", "3")
	//Загружаем переменные окружения
	config := LoadConfig("")

	if config.Username != "testUser" || config.Password != "testPass" || config.RetryCount != 3 {
		t.Fatal("Error reading env config")
	}

}

func TestLoadConfigFromFile(t *testing.T) {
	//Создаём временную папку
	tempDir := t.TempDir()
	// Создаём временный тестовый файл
	content := `
ENV=production
AUTH_USERNAME=testUser
AUTH_PASSWORD=testPass
AUTH_RETRY_COUNT=3

`
	tmpFile, err := os.CreateTemp(tempDir, "test_config_*.env")
	if err != nil {
		t.Fatal("Error creating temp file", err.Error())
	}
	defer os.Remove(tmpFile.Name())
	//Записываем в файл наше заготовленное содержимое
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatal("Error writing to temp file", err.Error())
	}
	//Закрываем файл после создания
	tmpFile.Close()
	if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) {
		t.Error("File was not created")
	}
	// Тестируем загрузку
	cfg := LoadConfig(tmpFile.Name())

	if cfg.Username != "testUser" || cfg.Password != "testPass" || cfg.RetryCount != 3 {
		t.Fatal("Error reading env config")
	}

}
