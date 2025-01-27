package main

type team struct {
	Name   string
	ID     int
	Points int
}

type teams struct {
	Teams []team
}

func createTeam(name string, id int) team {
	return team{
		Name:   name,
		ID:     id,
		Points: 0,
	}
}
