package main

import (
	"bytes"
	"log"
	"text/template"
)

func addTeamTemplate(team *team) []byte {
	tmpl, err := template.ParseFiles("views/team.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, team)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to addteamtemplate, before sending the bytes to the host")
	return renderedMessage.Bytes()
}

type CardSelection struct { //its used to know if the card was browsed by a host or a client to set to their respective pages
	ClientStatus string
	ID           int
	Number       string
}

func showSelectedCard(card CardSelection) []byte {
	tmpl, err := template.ParseFiles("views/card.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, card)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to showselectedCard, before sending the bytes to the host")
	return renderedMessage.Bytes()
}

func showCardAnswer(card Card) []byte {
	tmpl, err := template.ParseFiles("views/cardanswer.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, card)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to showcardanswer, before sending the bytes to the host")
	return renderedMessage.Bytes()
}

func addpoints(team team) []byte {
	tmpl, err := template.ParseFiles("views/addpoints.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, team)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to addpoints, before sending the bytes to the host")
	log.Printf("this are the current team points: %v", team.Points)
	return renderedMessage.Bytes()
}

func resetQanimation() []byte {
	tmpl, err := template.ParseFiles("views/resetQanimation.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, nil)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to resetQanimation, before sending the bytes to the host")
	return renderedMessage.Bytes()
}

func removeQuestionCover() []byte {
	tmpl, err := template.ParseFiles("views/removequestioncover.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, nil)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to resetQanimation, before sending the bytes to the host")
	return renderedMessage.Bytes()

}
