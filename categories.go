package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Card struct {
	Number   string `json:"number"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	HasQImg  bool   `json:"hasqimg"`
	QImgName string `json:"QImgName"`
	HasAImg  bool   `json:"hasaimg"`
	AImgName string `json:"AImgName"`
}

type Category struct {
	Title  string `json:"title"`
	Number string `json:"number"`
	Cards  []Card `json:"cards"`
}

func readcategories() ([]Category, error) {
	var categories []Category
	for i := 1; i < 7; i++ {
		numberstring := strconv.Itoa(i)
		name := "category" + numberstring + ".json"
		namepath := "categories/" + name
		file, err := os.ReadFile(namepath)
		if err != nil {
			return categories, fmt.Errorf("error while reading the file")
		}

		var category Category
		if err = json.Unmarshal(file, &category); err != nil {
			return categories, fmt.Errorf("error while marshaling the category. %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}
