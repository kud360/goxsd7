package loader

import (
	"context"
	"errors"
	"io"
	"testing"
	"testing/fstest"
)

func readAll(t *testing.T, doc *Document) string {
	t.Helper()
	defer func() { _ = doc.Body.Close() }()
	b, err := io.ReadAll(doc.Body)
	if err != nil {
		t.Fatalf("reading document body: %v", err)
	}
	return string(b)
}

func TestFS(t *testing.T) {
	fsys := fstest.MapFS{"schemas/a.xsd": &fstest.MapFile{Data: []byte("<schema-a/>")}}
	r := FS(fsys)

	doc, err := r.Resolve(context.Background(), Request{Location: "schemas/a.xsd"})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if got := readAll(t, doc); got != "<schema-a/>" {
		t.Errorf("body = %q", got)
	}

	_, err = r.Resolve(context.Background(), Request{Location: "missing.xsd"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("missing file: err = %v, want ErrNotFound", err)
	}
}

func TestFSRelativeResolution(t *testing.T) {
	fsys := fstest.MapFS{"schemas/b.xsd": &fstest.MapFile{Data: []byte("<schema-b/>")}}
	doc, err := FS(fsys).Resolve(context.Background(), Request{
		Location: "b.xsd",
		BaseURI:  "schemas/a.xsd",
	})
	if err != nil {
		t.Fatalf("Resolve relative: %v", err)
	}
	if doc.URI != "schemas/b.xsd" {
		t.Errorf("URI = %q, want schemas/b.xsd", doc.URI)
	}
}

func TestMapByNamespace(t *testing.T) {
	r := Map(map[string]string{"urn:example": "<schema-ns/>"})
	doc, err := r.Resolve(context.Background(), Request{Namespace: "urn:example"})
	if err != nil {
		t.Fatalf("Resolve by namespace: %v", err)
	}
	if got := readAll(t, doc); got != "<schema-ns/>" {
		t.Errorf("body = %q", got)
	}
}

func TestChain(t *testing.T) {
	miss := Map(nil)
	hit := Map(map[string]string{"a.xsd": "<from-hit/>"})
	doc, err := Chain(miss, hit).Resolve(context.Background(), Request{Location: "a.xsd"})
	if err != nil {
		t.Fatalf("Chain: %v", err)
	}
	if got := readAll(t, doc); got != "<from-hit/>" {
		t.Errorf("body = %q", got)
	}

	_, err = Chain(miss, miss).Resolve(context.Background(), Request{Location: "nope.xsd"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("exhausted chain: err = %v, want ErrNotFound", err)
	}

	abort := ResolverFunc(func(context.Context, Request) (*Document, error) {
		return nil, errors.New("io exploded")
	})
	_, err = Chain(abort, hit).Resolve(context.Background(), Request{Location: "a.xsd"})
	if err == nil || errors.Is(err, ErrNotFound) {
		t.Errorf("non-NotFound error must abort the chain, got %v", err)
	}
}
