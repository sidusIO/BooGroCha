package timeedit

import "fmt"

type room struct {
	Name   string `json:"fields.Lokalsignatur"`
	Id     string `json:"idAndType"`
	Seats  int    `json:"seats"`
	Campus string `json:"campus"`
}

type rooms []room

func (rs rooms) idFromName(name string) (string, error) {
	for _, room := range rs {
		if room.Name == name {
			return room.Id, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs rooms) nameFromId(id string) (string, error) {
	for _, room := range rs {
		if room.Id == id {
			return room.Name, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs rooms) removeAt(i int) rooms {
	return append(rs[:i], rs[i+1:]...)
}

func (rs rooms) remove(room room) rooms {
	for i, r := range rs {
		if r == room {
			return rs.removeAt(i)
		}
	}
	return rs
}

func (rs rooms) removeMany(rooms rooms) rooms {
	for _, n := range rooms {
		rs = rs.remove(n)
	}
	return rs
}
