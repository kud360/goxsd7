// Package loader is the IO seam of goxsd7: how schema documents are
// obtained given a namespace and/or location hint. Everything that needs
// a schema document — xs:import/include/redefine/override during
// compilation, and xsi:schemaLocation hints during instance validation —
// goes through a Resolver, so multi-schema loading has exactly one home
// (docs/LESSONS.md: goxsd5 deferred this across an architectural boundary
// and paid for it).
//
// Users implement Resolver for custom transports/catalogs, or compose the
// provided helpers: Dir, FS, HTTP, Map, Chain.
package loader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strings"
)

// Request identifies a schema document to load.
type Request struct {
	// Namespace is the expected target namespace; "" means absent
	// (no-namespace schema). Resolvers may use it as a lookup key
	// (catalogs) or to disambiguate.
	Namespace string
	// Location is the schemaLocation hint: a URL or path, possibly
	// relative to BaseURI. May be "" (e.g. xs:import without location).
	Location string
	// BaseURI is the URI of the referencing document, for resolving a
	// relative Location. May be "".
	BaseURI string
}

// resolved returns Location resolved against BaseURI when both permit it.
func (r Request) resolved() string {
	if r.BaseURI == "" || r.Location == "" {
		return r.Location
	}
	base, err := url.Parse(r.BaseURI)
	if err != nil {
		return r.Location
	}
	ref, err := url.Parse(r.Location)
	if err != nil {
		return r.Location
	}
	return base.ResolveReference(ref).String()
}

// Document is a loaded schema document. URI is the resolved absolute
// identity used for deduplication (loading the same URI twice must yield
// the same components). Callers own closing Body.
type Document struct {
	URI  string
	Body io.ReadCloser
}

// ErrNotFound reports that a resolver does not have the requested
// document. Chain treats it as "try the next resolver"; any other error
// aborts resolution.
var ErrNotFound = errors.New("schema document not found")

// Resolver obtains schema documents. Implementations must be safe for
// concurrent use.
type Resolver interface {
	Resolve(ctx context.Context, req Request) (*Document, error)
}

// ResolverFunc adapts a function to Resolver.
type ResolverFunc func(ctx context.Context, req Request) (*Document, error)

func (f ResolverFunc) Resolve(ctx context.Context, req Request) (*Document, error) {
	return f(ctx, req)
}

// FS resolves location hints as paths within fsys. BaseURI-relative
// resolution applies before lookup; a leading "./" and "/" are trimmed.
func FS(fsys fs.FS) Resolver {
	return ResolverFunc(func(_ context.Context, req Request) (*Document, error) {
		loc := strings.TrimPrefix(strings.TrimPrefix(req.resolved(), "./"), "/")
		if loc == "" {
			return nil, fmt.Errorf("resolving namespace %q: no location hint: %w", req.Namespace, ErrNotFound)
		}
		f, err := fsys.Open(loc)
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("resolving %q: %w", loc, ErrNotFound)
		}
		if err != nil {
			return nil, fmt.Errorf("opening schema %q: %w", loc, err)
		}
		return &Document{URI: loc, Body: f}, nil
	})
}

// HTTP fetches http(s) location hints. A nil client uses
// http.DefaultClient. Non-http(s) locations report ErrNotFound so Chain
// can try elsewhere; HTTP status 404/410 also map to ErrNotFound.
func HTTP(client *http.Client) Resolver {
	if client == nil {
		client = http.DefaultClient
	}
	return ResolverFunc(func(ctx context.Context, req Request) (*Document, error) {
		loc := req.resolved()
		if !strings.HasPrefix(loc, "http://") && !strings.HasPrefix(loc, "https://") {
			return nil, fmt.Errorf("resolving %q: not an http(s) location: %w", loc, ErrNotFound)
		}
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, loc, nil)
		if err != nil {
			return nil, fmt.Errorf("building request for %q: %w", loc, err)
		}
		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("fetching schema %q: %w", loc, err)
		}
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
			resp.Body.Close()
			return nil, fmt.Errorf("fetching schema %q: status %s: %w", loc, resp.Status, ErrNotFound)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("fetching schema %q: unexpected status %s", loc, resp.Status)
		}
		return &Document{URI: loc, Body: resp.Body}, nil
	})
}

// Map serves fixed documents from memory. Keys are matched against the
// resolved location first, then the namespace — so it doubles as a tiny
// catalog. Values are the document bytes.
func Map(docs map[string]string) Resolver {
	return ResolverFunc(func(_ context.Context, req Request) (*Document, error) {
		if body, ok := docs[req.resolved()]; ok {
			return &Document{URI: req.resolved(), Body: io.NopCloser(strings.NewReader(body))}, nil
		}
		if body, ok := docs[req.Namespace]; ok {
			return &Document{URI: req.Namespace, Body: io.NopCloser(strings.NewReader(body))}, nil
		}
		return nil, fmt.Errorf("resolving (ns %q, loc %q): %w", req.Namespace, req.Location, ErrNotFound)
	})
}

// Chain tries each resolver in order; ErrNotFound falls through to the
// next, any other error aborts. All resolvers exhausted → ErrNotFound
// decorated with every attempt.
func Chain(resolvers ...Resolver) Resolver {
	return ResolverFunc(func(ctx context.Context, req Request) (*Document, error) {
		var attempts []error
		for _, r := range resolvers {
			doc, err := r.Resolve(ctx, req)
			if err == nil {
				return doc, nil
			}
			if !errors.Is(err, ErrNotFound) {
				return nil, err
			}
			attempts = append(attempts, err)
		}
		return nil, fmt.Errorf("all %d resolvers exhausted: %w", len(resolvers), errors.Join(attempts...))
	})
}
