package requests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	recipesBasePath = "recipes"
)

// BuildGetRecipeRequest builds an HTTP request for fetching a recipe.
func (b *Builder) BuildGetRecipeRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipesRequest builds an HTTP request for fetching a list of recipes.
func (b *Builder) BuildGetRecipesRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		recipesBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchForRecipesRequest builds an HTTP request for fetching a list of recipes.
func (b *Builder) BuildSearchForRecipesRequest(ctx context.Context, query string, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachToSpan(span, keys.SearchQueryKey, query)

	queryParams := filter.ToValues()
	queryParams.Set(types.SearchQueryKey, query)

	uri := b.BuildURL(
		ctx,
		queryParams,
		recipesBasePath,
		"search",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateRecipeRequest builds an HTTP request for creating a recipe.
func (b *Builder) BuildCreateRecipeRequest(ctx context.Context, input *types.RecipeCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildUpdateRecipeRequest builds an HTTP request for updating a recipe.
func (b *Builder) BuildUpdateRecipeRequest(ctx context.Context, recipe *types.Recipe) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipe == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipe.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipe.ID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	input := converters.ConvertRecipeToRecipeUpdateRequestInput(recipe)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildArchiveRecipeRequest builds an HTTP request for archiving a recipe.
func (b *Builder) BuildArchiveRecipeRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipeDAGRequest builds an HTTP request for fetching a recipe.
func (b *Builder) BuildGetRecipeDAGRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		"dag",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetRecipeMealPlanTasksRequest builds an HTTP request for fetching a recipe.
func (b *Builder) BuildGetRecipeMealPlanTasksRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		"prep_steps",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCloneRecipeRequest builds an HTTP request for fetching a recipe.
func (b *Builder) BuildCloneRecipeRequest(ctx context.Context, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	uri := b.BuildURL(
		ctx,
		nil,
		recipesBasePath,
		recipeID,
		"clone",
	)
	tracing.AttachToSpan(span, keys.RequestURIKey, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

const (
	imagePNG  = "image/png"
	imageJPEG = "image/jpeg"
	imageGIF  = "image/gif"
)

// BuildRecipeMediaUploadRequest builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildRecipeMediaUploadRequest(ctx context.Context, media []byte, extension, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if len(media) == 0 {
		return nil, ErrNilInputProvided
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("upload", fmt.Sprintf("media.%s", extension))
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating form file")
	}

	if _, err = io.Copy(part, bytes.NewReader(media)); err != nil {
		return nil, observability.PrepareError(err, span, "copying file contents to request")
	}

	if err = writer.Close(); err != nil {
		return nil, observability.PrepareError(err, span, "closing media writer")
	}

	uri := b.BuildURL(ctx, nil, recipesBasePath, recipeID, "images")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building media upload request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func (b *Builder) prepareUploads(ctx context.Context, files map[string][]byte) (io.Reader, string, error) {
	_, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if files == nil {
		return nil, "", ErrNilInputProvided
	}

	var ct string
	for filename := range files {
		switch strings.ToLower(strings.TrimSpace(filepath.Ext(filename))) {
		case ".jpeg":
			if ct == "" {
				ct = imageJPEG
			} else if ct != "" && ct != imageJPEG {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		case ".png":
			if ct == "" {
				ct = imagePNG
			} else if ct != "" && ct != imagePNG {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		case ".gif":
			if ct == "" {
				ct = imageGIF
			} else if ct != "" && ct != imageGIF {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		}
	}

	if ct == "" {
		return nil, "", ErrInvalidPhotoEncodingForUpload
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for filename, content := range files {
		part, err := writer.CreateFormFile("upload", filename)
		if err != nil {
			return nil, "", observability.PrepareError(err, span, "creating form file")
		}

		if _, err = io.Copy(part, bytes.NewReader(content)); err != nil {
			return nil, "", observability.PrepareError(err, span, "copying file contents to request")
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", observability.PrepareError(err, span, "closing media writer")
	}

	return body, writer.FormDataContentType(), nil
}

// BuildMultipleRecipeMediaUploadRequest builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildMultipleRecipeMediaUploadRequest(ctx context.Context, files map[string][]byte, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	body, formDataContentType, err := b.prepareUploads(ctx, files)
	if err != nil {
		return nil, observability.PrepareError(err, span, "preparing upload request")
	}

	uri := b.BuildURL(ctx, nil, recipesBasePath, recipeID, "images")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building media upload request")
	}

	req.Header.Set("Content-Type", formDataContentType)

	return req, nil
}

// BuildMultipleRecipeMediaUploadRequestForRecipeStep builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildMultipleRecipeMediaUploadRequestForRecipeStep(ctx context.Context, files map[string][]byte, recipeID, recipeStepID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeStepIDKey, recipeStepID)

	body, formDataContentType, err := b.prepareUploads(ctx, files)
	if err != nil {
		return nil, observability.PrepareError(err, span, "preparing upload request")
	}

	uri := b.BuildURL(ctx, nil, recipesBasePath, recipeID, recipeStepsBasePath, recipeStepID, "images")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building media upload request")
	}

	req.Header.Set("Content-Type", formDataContentType)

	return req, nil
}
