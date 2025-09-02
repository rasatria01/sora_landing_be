package utils

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// SlugChecker defines a contract for checking slug uniqueness
type SlugChecker interface {
	SlugExists(ctx context.Context, slug string) (bool, error)
}

// Slugify creates a URL-safe slug from a string
func Slugify(s string) string {
	// lowercase
	s = strings.ToLower(s)

	// replace non-alphanumeric with "-"
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// trim leading/trailing "-"
	s = strings.Trim(s, "-")

	return s
}

// GenerateUniqueSlug ensures a slug is unique in the given repository
func GenerateUniqueSlug(ctx context.Context, repo SlugChecker, base string) (string, error) {
	slug := Slugify(base)
	uniqueSlug := slug
	counter := 1

	for {
		exists, err := repo.SlugExists(ctx, uniqueSlug)
		if err != nil {
			return "", err
		}
		if !exists {
			return uniqueSlug, nil
		}

		// append counter until unique
		uniqueSlug = fmt.Sprintf("%s-%d", slug, counter)
		counter++
	}
}
