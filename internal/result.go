package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

func RunProgram(path string, showContent bool) {
	out := os.Stdout
	if showContent {
		err := WriteFileContent(out, path)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
	}

	err := WriteResult(out, path)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func WriteFileContent(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("WriteFileContent: %w", err)
	}
	defer file.Close()

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		return fmt.Errorf("WriteFileContent: invalid file: %v", path)
	}

	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("WriteFileContent: %w", err)
	}
	return nil
}

func WriteResult(w io.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("WriteResult: %w", err)
	}
	defer file.Close()

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		return fmt.Errorf("WriteResult: invalid file: %v", path)
	}

	scanner := bufio.NewScanner(file)
	terrain := CreateAntFarm()
	var errAnt error
	for scanner.Scan() {
		errAnt = terrain.ReadDataFromLine(scanner.Text())
		if errAnt != nil {
			return fmt.Errorf("WriteResult: %w", errInvalidDataFormat(errAnt))
		}
	}
	errAnt = terrain.ValidateByFieldInfo()
	if errAnt != nil {
		return fmt.Errorf("WriteResult: %w", errInvalidDataFormat(errAnt))
	}
	errAnt = terrain.Match()
	if errAnt != nil {
		return fmt.Errorf("WriteResult: %w", errPaths(errAnt))
	}

	sort.Slice(terrain.Result.Paths, func(i, j int) bool { return terrain.Result.Paths[i].Len < terrain.Result.Paths[j].Len })
	steps, antsForEachPath := calcSteps(terrain.Result.AntsCount, terrain.Result.Paths)

	if steps == 1 {
		roomName := terrain.Result.Paths[0].Front.Room.Name
		for ant := 1; ant <= antsForEachPath[0]; ant++ {
			fmt.Fprintf(w, "L%d-%s ", ant, roomName)
		}
		fmt.Fprintln(w)
	} else {
		paths := pathsOfListToSlice(terrain.Result.Paths)
		a, b := &AntQueue{}, &AntQueue{}
		antNum := 1
		for i := 1; i <= steps; i++ {
			cur := a.Dequeue()
			for cur != nil {
				w.Write([]byte(fmt.Sprintf("L%d-%s ", cur.Num, paths[cur.Path][cur.Pos].Name)))
				cur.Pos++
				if cur.Pos < len(paths[cur.Path]) {
					b.EnqueueAnt(cur)
				}
				cur = a.Dequeue()
			}
			for j, v := range antsForEachPath {
				if v > 0 {
					w.Write([]byte(fmt.Sprintf("L%d-%s ", antNum, paths[j][0].Name)))
					antsForEachPath[j]--
					b.Enqueue(antNum, j, 1)
					antNum++
				}
			}
			t := a
			a = b
			b = t
			fmt.Fprintln(w)
		}
	}
	return nil
}

func pathsOfListToSlice(paths []*List) [][]*Room {
	result := make([][]*Room, len(paths))
	for i, path := range paths {
		result[i] = make([]*Room, path.Len)
		j := 0
		for node := path.Front; node != nil; node = node.Next {
			result[i][j] = node.Room
			j++
		}
	}
	return result
}
