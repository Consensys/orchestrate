package provider

import (
	"context"
	"strings"
)

type providerCtxKey int

const (
	providerKey providerCtxKey = iota
)

// WithProviderName attachs provider name to context
func WithName(ctx context.Context, qualifiedName string) context.Context {
	providerName := GetName(qualifiedName)
	return context.WithValue(ctx, providerKey, providerName)
}

func NameFromContext(ctx context.Context) string {
	name, _ := ctx.Value(providerKey).(string)
	return name
}

func GetName(qualifiedName string) string {
	parts := strings.Split(qualifiedName, "@")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

// QualifyName Creates a qualified name for an element
func QualifyName(providerName, elementName string) string {
	parts := strings.Split(elementName, "@")
	if len(parts) == 1 {
		return elementName + "@" + providerName
	}
	return elementName
}

// QualifyNameFromContext Gets the fully qualified name.
func QualifyNameFromContext(ctx context.Context, elementName string) string {
	providerName := NameFromContext(ctx)
	return QualifyName(providerName, elementName)
}
