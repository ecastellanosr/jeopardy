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
	ID       int    `json:"id"`
}

type Category struct {
	Title string `json:"title"`
	Cards []Card `json:"cards"`
}

func readcategories() ([]Category, error) {
	var categories []Category
	category_id := 0
	for i := 1; i < 7; i++ {
		card_id := 1
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
		for e := range category.Cards {
			id := category_id + card_id
			category.Cards[e].ID = id
			card_id++
		}
		category_id = category_id + 10
		categories = append(categories, category)
	}
	return categories, nil
}

func FindCard(categories []Category, question_id int) Card {
	Looked_card := Card{}
	category_id := 0

	switch {
	case question_id < 10:
		break
	case question_id < 20:
		category_id = 1
	case question_id < 30:
		category_id = 2
	case question_id < 40:
		category_id = 3
	case question_id < 50:
		category_id = 4
	default:
		category_id = 5
	}

	card_id := question_id
	for _, card := range categories[category_id].Cards {
		if card_id == card.ID {
			Looked_card = card
		}
	}
	return Looked_card
}

func addpointstoteam(categories []Category, teams *teams, question_id, team_id int) team {
	winnerteam := team{}
	points := 0

	card := FindCard(categories, question_id)
	card_points, _ := strconv.Atoi(card.Number)
	points = card_points

	for _, team := range teams.Teams {
		if team_id == team.ID {
			team.Points += points
			winnerteam = team
		}
	}
	return winnerteam
}
