package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

func RunProgram(path string, showContent bool) {
	err := WriteResultByPath(os.Stdout, path, showContent)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func WriteResultByPath(w io.Writer, path string, writeContent bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("WriteResultByPath: %w", err)
	}
	defer file.Close()

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		return fmt.Errorf("WriteResultByPath: invalid file: %v", path)
	}

	scanner := bufio.NewScanner(file)
	result, err := GetResult(scanner)
	if err != nil {
		return fmt.Errorf("WriteResultByPath: %w", err)
	}

	if writeContent {
		file.Seek(0, io.SeekStart)
		_, err = io.Copy(w, file)
		if err != nil {
			return fmt.Errorf("WriteResultByPath: %w", err)
		}
		fmt.Fprint(w, "\n\n# result\n")
	}
	result.WriteResult(w)

	return nil
}

func GetResult(scanner *bufio.Scanner) (*Result, error) {
	terrain := CreateAntFarm()
	var err error
	for scanner.Scan() {
		err = terrain.ReadDataFromLine(scanner.Text())
		if err != nil {
			return nil, errInvalidDataFormat(err)
		}
	}
	err = terrain.ValidateByFieldInfo()
	if err != nil {
		return nil, errInvalidDataFormat(err)
	}
	err = terrain.Match()
	if err != nil {
		return nil, errPaths(err)
	}
	return terrain.Result, nil
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

func (r *Result) WriteResult(w io.Writer) {
	sort.Slice(r.Paths, func(i, j int) bool { return r.Paths[i].Len < r.Paths[j].Len })
	steps, antsForEachPath := calcSteps(r.AntsCount, r.Paths)
	if steps == 1 {
		roomName := r.Paths[0].Front.Room.Name
		for ant := 1; ant <= antsForEachPath[0]; ant++ {
			fmt.Fprintf(w, "L%d-%s ", ant, roomName)
		}
		fmt.Fprintln(w)
	} else {
		paths := pathsOfListToSlice(r.Paths)
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
}
