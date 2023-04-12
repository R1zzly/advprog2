package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_readUserInput(t *testing.T) {
	doneChan := make(chan bool)

	var stdin bytes.Buffer

	stdin.Write([]byte("5\nq\n"))

	go readUserInput(&stdin, doneChan)
	<-doneChan
	close(doneChan)
}

func Test_checkNumbers(t *testing.T) {
	tests := []struct {
		input    string
		msg      string
		wantQuit bool
	}{
		{input: "q\n", msg: "", wantQuit: true},
		{input: "5\n", msg: "5 is a prime number!", wantQuit: false},
		{input: "abc\n", msg: "Please enter a whole number!", wantQuit: false},
	}

	for _, tc := range tests {
		scanner := bufio.NewScanner(strings.NewReader(tc.input))
		msg, quit := checkNumbers(scanner)

		if msg != tc.msg {
			t.Errorf("checkNumbers(%q) msg = %q; want %q", tc.input, msg, tc.msg)
		}

		if quit != tc.wantQuit {
			t.Errorf("checkNumbers(%q) quit = %v; want %v", tc.input, quit, tc.wantQuit)
		}
	}
}

func Test_intro(t *testing.T) {
	fstout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	intro()
	_ = w.Close()

	os.Stdout = fstout
	out, _ := io.ReadAll(r)

	if string(out) != `Is it Prime?
------------
Enter a whole number, and we'll tell you if it is a prime number or not. Enter q to quit.
-> ` {
		t.Errorf("incorect prompt() want %q", string(out))
	}
}

func Test_prompt(t *testing.T) {
	fstout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	prompt()
	_ = w.Close()

	os.Stdout = fstout
	out, error := io.ReadAll(r)
	if error != nil {
		t.Error(error)
	}

	if string(out) != "-> " {
		t.Errorf("incorect prompt() want %q", string(out))
	}
}

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}
