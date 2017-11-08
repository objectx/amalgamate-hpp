package main

// FileList represents non-duplicated list of files.
// The order is preserved.
type FileList struct {
	files []string
	dict  map[string]int
}

// Items Retrieves registered files
func (l *FileList) Items() []string {
	return l.files
}

// FindIndex finds the index of `p` in the list.
// Returns -1 if missing.
func (l *FileList) FindIndex(p string) int {
	if idx, ok := l.dict[p]; ok {
		return idx
	}
	return -1
}

// Register records `filePath` as a file
func (l *FileList) Register(filePath string) int {
	if l.dict == nil {
		l.dict = make(map[string]int)
	}
	idx, ok := l.dict[filePath]
	if ok {
		return idx // Already appeared
	}
	idx = len(l.files)
	l.dict[filePath] = idx
	l.files = append(l.files, filePath)
	return idx
}

func (l *FileList) Clear() {
	l.dict = nil
	l.files = nil
}
