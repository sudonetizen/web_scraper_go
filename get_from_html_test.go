package main

import "testing"

func TestGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := "Test Title"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetH1FromHTMLNoH1(t *testing.T) {
	inputBody := "<html><body><h2>Test Title</h2></body></html>"
	actual := getH1FromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	inputBody := "<html><body><p>outside p</p><main><p>main p</p></main></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "main p"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}


func TestGetFirstParagraphFromHTMLNoMainPriority(t *testing.T) {
	inputBody := "<html><body><p>outside p</p><p>outside p2</p></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := "outside p"

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLEmpty(t *testing.T) {
	inputBody := "<html><body>><main></main></body></html>"
	actual := getFirstParagraphFromHTML(inputBody)
	expected := ""

	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}
