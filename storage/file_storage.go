package storage

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"
)

type fileStorage struct {
	quotes []string
	size   int
}

func (s *fileStorage) RandomQuote() (string, error) {
	index := rand.Intn(s.size)

	return s.quotes[index], nil
}

func NewFileStorage(fileName string) (*fileStorage, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, errors.New("storage is empty")
	}

	rand.Seed(time.Now().UnixNano())

	fs := &fileStorage{
		quotes: lines,
		size:   len(lines),
	}

	return fs, nil
}
