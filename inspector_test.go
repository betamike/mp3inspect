package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func fixturePath(t *testing.T, fixture string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("problems recovering caller information")
	}

	return filepath.Join(filepath.Dir(filename), fixture)
}

func loadFixture(t *testing.T, fixture string) string {
	content, err := ioutil.ReadFile(fixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}

	return string(content)
}

func TestAllMp3s(t *testing.T) {

	files, err := ioutil.ReadDir("./test/files")
	if err != nil {
		fmt.Println("Be sure to read the README to generate the files needed for testing!")
		log.Fatal(err)
	}

	var mp3Files []string

	for _, file := range files {
		fileName := file.Name()
		if string(fileName[len(fileName)-3:]) == "mp3" {
			mp3Files = append(mp3Files, fileName)
		}
	}

	if len(mp3Files) == 0 {
		log.Fatal("Please read README, run mp3 and golden file generation scripts first")
	}

	for _, mp3File := range mp3Files {
		t.Run(mp3File, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			cmd := exec.Command(path.Join(dir, "/bin/mp3inspect"), path.Join(dir, "test/files", mp3File))
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error executing command")
				t.Fatal(err)
			}

			actual := string(output)

			re := regexp.MustCompile("\\s+")

			// I didn't want to have to deal with tabs and spaces
			// Just normalizing the actual and expected output to compare values easier
			actual = re.ReplaceAllLiteralString(actual, ",")

			goldenFileName := mp3File[0:len(mp3File)-3] + "golden"

			expected := loadFixture(t, path.Join("test/files", goldenFileName))
			expected = re.ReplaceAllLiteralString(expected, ",")

			if !strings.EqualFold(actual, expected) {
				t.Fatalf("\nactual   = %s\nexpected = %s", actual, expected)
			}
		})
	}
}

func TestMain(m *testing.M) {
	binaryName := "mp3inspect"
	make := exec.Command("make")
	err := make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v", binaryName, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
