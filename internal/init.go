package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (a *antFarm) SetAntsFromLine(line string) error {
	countAnts, err := strconv.Atoi(line)
	if err != nil || countAnts < 1 {
		return errors.New("invalid number of Ants")
	}
	a.AntsCount = countAnts
	a.Result.AntsCount = countAnts
	return nil
}

func (a *antFarm) SetRoomFromLine(line string) (*Room, error) {
	splited := strings.Split(line, " ")
	if len(splited) != 3 || len(splited[0]) < 1 {
		return nil, errors.New("invalid format of room")
	}

	name := splited[0]
	if strings.HasPrefix(name, "L") || strings.Contains(name, "-") {
		return nil, errors.New("room name is invalid")
	}

	x, errX := strconv.Atoi(splited[1])
	y, errY := strconv.Atoi(splited[2])
	if errX != nil || errY != nil {
		return nil, errors.New("room coords can only be numbers")
	}

	if a.Rooms[name] != nil {
		return nil, fmt.Errorf("room name duplicated: '%v'", name)
	}

	if _, ok := a.FieldInfo.UsingCoordinates[x]; !ok {
		a.FieldInfo.UsingCoordinates[x] = make(map[int]bool)
	}
	if a.FieldInfo.UsingCoordinates[x][y] {
		return nil, fmt.Errorf("room coords must be unique; room name: '%v'", name)
	}

	a.FieldInfo.UsingCoordinates[x][y] = true
	room := &Room{
		Name:   name,
		X:      x,
		Y:      y,
		Paths:  make(map[*Room]int),
		Weight: [2]int{0, 0},
	}
	a.Rooms[name] = room
	return room, nil
}

func (a *antFarm) SetMainRooms(line string, startOrEnd bool) error {
	room, err := a.SetRoomFromLine(line)
	if err != nil {
		return err
	}
	if startOrEnd {
		a.Start = room.Name
	} else {
		a.End = room.Name
	}
	return nil
}

func (a *antFarm) SetPathsFromLine(line string) error {
	splited := strings.Split(line, "-")
	if len(splited) != 2 || len(splited[0]) < 1 || len(splited[1]) < 1 {
		return errors.New("invalid format of path")
	}
	name1, name2 := splited[0], splited[1]
	if name1 == name2 {
		return fmt.Errorf("rooms can't link themselves. Line: '%v'", line)
	}

	room1 := a.Rooms[name1]
	room2 := a.Rooms[name2]
	if room1 == nil || room2 == nil {
		return fmt.Errorf("path contains unknown room. Line: '%v'", line)
	}

	room1.Paths[room2] = STABLE
	room2.Paths[room1] = STABLE
	return nil
}
