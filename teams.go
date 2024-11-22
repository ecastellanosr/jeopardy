package main

type team struct {
	Name   string
	id     int
	Points int
}

type teams struct {
	Teams []team
}

func createTeam(name string, id int) team {
	return team{
		Name:   name,
		id:     id,
		Points: 0,
	}
}

func NewTeams() *teams {
	return &teams{
		Teams: []team{},
	}
}
