package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("world", "")
	want := "Hello, world"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestHelloArgs(t *testing.T) {
	got := Hello("sehyeong", "")
	want := "Hello, sehyeong"
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	t.Run("say hello to people", func(t *testing.T) {
		got := Hello("sehyeong", "")
		want := "Hello, sehyeong"

		if got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	})

	t.Run("say 'Hello world' when an empty string is supplied", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, world"

		if got != want {
			t.Errorf("got %q; want %q", got, want)
		}
	})
}

func TestHelloLang(t *testing.T) {
	t.Run("in specific language", func(t *testing.T) {
		got := Hello("sehyeong", "Spanish")
		want := "Hola, sehyeong"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q; want %q", got, want)
	}
}
