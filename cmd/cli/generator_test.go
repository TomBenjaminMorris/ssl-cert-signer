package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_checkIfValidFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test*.csr")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	tests := []struct {
		name     string
		filename string
		want     bool
		wantErr  bool
	}{
		{"File does exist", tmpfile.Name(), true, false},
		{"File does not exist", "nowhere/test.csr", false, true},
		{"File is not csv", "test.txt", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkIfValidFile(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkIfValidFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkIfValidFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadCA(t *testing.T) {
	t.Run("load root ca", func(t *testing.T) {
		got, err := loadCA()
		notWant := CA{}

		if err != nil {
			t.Errorf("loadCA() error = %v", err)
			return
		}

		if got == notWant {
			t.Errorf("loadCA() returned a nil CA")
			return
		}
	})
}

func Test_getCSR(t *testing.T) {
	tests := []struct {
		name    string
		want    inputFile
		wantErr bool
		osArgs  []string
	}{
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Valid CSR no args", inputFile{"./exampleFiles/mydomain.com.csr", false, false}, false, []string{"cmd", "./exampleFiles/mydomain.com.csr"}},
		{"Valid CSR with args", inputFile{"./exampleFiles/mydomain.com.csr", true, true}, false, []string{"cmd", "-flag1", "-flag2", "./exampleFiles/mydomain.com.csr"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOsArgs := os.Args
			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) //flags are now reset
			}()

			os.Args = tt.osArgs
			got, err := getCSR()
			if (err != nil) != tt.wantErr {
				t.Errorf("getCSR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCSR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseCSR(t *testing.T) {
	tests := []struct {
		name  string
		input inputFile
		want  CSR
	}{
		{
			"name",
			inputFile{"../../exampleFiles/mydomain.com.csr", false, false},
			CSR{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCSR(tt.input); reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCSR() = %v, want %v", got, tt.want)
			}
		})
	}
}
