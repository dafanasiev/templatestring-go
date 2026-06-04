package templatestring

import (
	"fmt"
	"os"
	"testing"
)

func Test_NoTokens_NoPlugins_EmptyTemplate(t *testing.T) {
	ts := NewTemplateString("")
	expected := ""
	actual, err := ts.Render()
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_NoTokens_NoPlugins(t *testing.T) {
	ts := NewTemplateString("xyz")
	expected := "xyz"
	actual, err := ts.Render()
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_Tokens_NoPlugins(t *testing.T) {
	ts := NewTemplateString("${xyz}")
	_, err := ts.Render()
	if err == nil {
		t.Fatal("should produce error")
	}
}

func Test_DelegatePlugin(t *testing.T) {
	ts := NewTemplateString("${  xyz   }")

	expected := "XYZ"
	actual, err := ts.Render(NewDelegatePlugin(func(token string) (string, bool, error) {
		if token == "xyz" {
			return "XYZ", true, nil
		}

		return "", true, fmt.Errorf("unsupported token: %s", token)
	}))

	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_Token_With_Literals(t *testing.T) {
	ts := NewTemplateString("!${  xyz  } is ${  qwe  }!")

	expected := "!XYZ is QWE!"
	actual, err := ts.Render(NewDelegatePlugin(func(token string) (string, bool, error) {
		if token == "xyz" {
			return "XYZ", true, nil
		}

		if token == "qwe" {
			return "QWE", true, nil
		}

		return "", true, fmt.Errorf("unsupported token: %s", token)
	}))

	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_EnvPlugin(t *testing.T) {
	ts := NewTemplateString("${env:HOME}")

	expected := os.Getenv("HOME")
	actual, err := ts.Render(NewEnvPlugin())

	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}

func Test_EnvPlugin2(t *testing.T) {
	ts := NewTemplateString("${env:}")

	_, err := ts.Render(NewEnvPlugin())

	if err == nil {
		t.Fatal("should produce error")
	}
}

func Test_EnvPlugin3(t *testing.T) {
	ts := NewTemplateString("${p:1}")

	actual, err := ts.Render(NewEnvPlugin())
	t.Logf("%s", actual)

	if err == nil {
		t.Fatal("should produce error")
	}
}

func Test_StringMethod(t *testing.T) {
	expected := ".... ${env:HOME} ...."
	ts := NewTemplateString(expected)

	_, err := ts.Render(NewEnvPlugin())
	actual := ts.String()

	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("expected: %s, actual: %s", expected, actual)
	}
}
