package file

import "os"

type Dirs struct {
	dir []os.DirEntry
}

func Dir(path string) (Dirs, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return Dirs{}, err
	}
	return Dirs{dir}, nil
}
