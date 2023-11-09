package internal

import (
	"errors"
	"strings"
)

func CreateAntFarm() *antFarm {
	return &antFarm{
		Rooms:     make(map[string]*Room),
		FieldInfo: &FieldInfo{UsingCoordinates: make(map[int]map[int]bool)},
		Result:    &Result{},
	}
}

func (a *antFarm) ValidateByFieldInfo() error {
	if a.FieldInfo.MODE != FIELD_PATHS {
		return errors.New("func Validate returns error")
	}
	if !a.FieldInfo.Start {
		return errors.New("please set ##start room")
	} else if !a.FieldInfo.End {
		return errors.New("please set ##end room")
	}
	return nil
}

func (a *antFarm) ReadDataFromLine(line string) error {
	if line == "" || strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
		return nil
	}

	switch a.FieldInfo.MODE {
	case FIELD_PATHS:
		return a.SetPathsFromLine(line)
	case FIELD_ROOMS:
		if strings.HasPrefix(line, "##") {
			if line == "##start" && !a.FieldInfo.Start && !a.FieldInfo.IsStart && !a.FieldInfo.IsEnd {
				a.FieldInfo.IsStart = true
				return nil
			} else if line == "##end" && !a.FieldInfo.End && !a.FieldInfo.IsEnd && !a.FieldInfo.IsStart {
				a.FieldInfo.IsEnd = true
				return nil
			}
			return errors.New("error with ## command")
		}

		if a.FieldInfo.IsStart || a.FieldInfo.IsEnd {
			err := a.SetMainRooms(line, a.FieldInfo.IsStart)
			if err != nil {
				return err
			}
			if a.FieldInfo.IsStart {
				a.FieldInfo.IsStart = false
				a.FieldInfo.Start = true
			} else {
				a.FieldInfo.IsEnd = false
				a.FieldInfo.End = true
			}
			return err
		}

		if len(strings.Split(line, " ")) != 3 {
			a.FieldInfo.MODE = FIELD_PATHS
			a.FieldInfo.UsingCoordinates = nil
			return a.ReadDataFromLine(line)
		}

		_, err := a.SetRoomFromLine(line)
		return err

	case FIELD_ANTS:
		err := a.SetAntsFromLine(line)
		if err != nil {
			return err
		}
		a.FieldInfo.MODE = FIELD_ROOMS
	}

	return nil
}

func (a *antFarm) Match() error {
	for {
		if !searchShortPath(a) {
			if a.StepsCount > 0 {
				return nil
			}
			return errors.New("path not found")
		}
		if !checkEffective(a) {
			return nil
		}
	}
}
