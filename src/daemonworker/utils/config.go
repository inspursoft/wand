package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadConfig(path string) (map[string]string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0640)
	if err != nil {
		log.Printf("Failed to open file: %s, error: %v\n", path, err)
		return nil, err
	}
	items := make(map[string]string)
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) > 0 && strings.Index(line, "=") != -1 {
			parts := strings.Split(line, "=")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			items[key] = value
		}
	}
	return items, nil
}

func ListConfig(config map[string]string) {
	log.Printf(">>>>>> Current Config <<<<<<\n")
	for key, val := range config {
		log.Printf("%s=%s\n", key, val)
	}
}
