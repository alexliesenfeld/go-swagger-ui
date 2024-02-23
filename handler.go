package go_swagger_ui

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed swagger-ui/dist/*
var swaggerUIFS embed.FS

//go:embed swagger-ui/templates/*
var templatesFS embed.FS

var tplOverrides = map[string]*template.Template{
	"index.html":             template.Must(template.ParseFS(templatesFS, "swagger-ui/templates/index.html")),
	"swagger-initializer.js": template.Must(template.ParseFS(templatesFS, "swagger-ui/templates/swagger-initializer.js")),
}

var allFilePaths = Must(walkFS("swagger-ui/dist/", &swaggerUIFS, "."))

func NewHandler(opts ...Option) http.HandlerFunc {
	cfg := uiConfig{
		htmlTitle: "Swagger UI",
	}

	for idx := range opts {
		opts[idx](&cfg)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fileName := strings.TrimPrefix(strings.TrimSpace(path.Base(r.URL.Path)), "/")
		if fileName == "" {
			fileName = "index.html"
		}

		// Always serve "index.html" if a file is being asked for does not exist.
		// These cases are usually caused by http.Handler instances that are mounted on URL paths
		// that do not end with a slash (e.g., https://example.com/hello, in which case the
		// file name would be "hello", although "index.html" is what is expected to be returned).
		if _, exists := allFilePaths[fileName]; !exists {
			fileName = "index.html"
		}

		// We reload the spec file only for the CLI. In a normal production HTTP mode
		// "specFilePath" should be unset and not used at all. See WithSpecFilePath.
		if cfg.specFilePath != "" && fileName == "index.html" {
			newSpecContent, err := readSpecFile(cfg.specFilePath)
			if err != nil {
				slog.Error("error reading Swagger UI file", "err", err.Error())
				sendError(w, err)
				return
			}

			cfg.spec = newSpecContent
		}

		// We either load the requested file from the embed filesystem directly or rendering
		// a template instead.
		var responseBody []byte
		if tpl, ok := tplOverrides[fileName]; ok {
			var buf bytes.Buffer
			if err := replaceVars(&buf, tpl, &cfg); err != nil {
				slog.Error("failed to use Swagger UI template", "err", err.Error())
				sendError(w, err)
				return
			}
			responseBody = buf.Bytes()
		} else {
			var err error
			responseBody, err = fs.ReadFile(swaggerUIFS, "swagger-ui/dist/"+fileName)
			if err != nil {
				slog.Error("error reading file", "err", err.Error())
				sendError(w, err)
				return
			}
		}

		w.Header().Set("Content-Type", getContentType(fileName, responseBody))
		w.Write(responseBody)
	}
}

func sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if errors.Is(err, fs.ErrNotExist) {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(500)
	}
}

func replaceVars(w io.Writer, tpl *template.Template, cfg *uiConfig) error {
	urlsAsJSON, err := fromObject(cfg.urls)
	if err != nil {
		return fmt.Errorf("cannot marshal URLs: %w", err)
	}

	return tpl.Execute(w, struct {
		BasePath, Spec, URL, HTMLTitle, DocExpansion, DefaultModelExpandDepth, DefaultModelsExpandDepth,
		DefaultModelRendering, QueryConfigEnabled, SupportedSubmitMethods, DeepLinking,
		ShowMutatedRequest, ShowExtensions, ShowCommonExtensions, Filter, FilterString,
		DisplayOperationId, TryItOutEnabled, DisplayRequestDuration, PersistAuthorization, WithCredentials,
		OAuth2RedirectUrl, Layout, ValidatorURL, MaxDisplayedTags, PrimaryURL, ConfigURL, URLs string
	}{
		BasePath:                 cfg.basePath,
		ConfigURL:                fromStringConfigValue(cfg.configURL),
		Spec:                     strings.TrimSpace(escapeString(string(cfg.spec))),
		URL:                      fromStringConfigValue(cfg.url),
		HTMLTitle:                cfg.htmlTitle,
		DocExpansion:             fromDocExpansionConfigValue(cfg.docExpansion),
		DefaultModelExpandDepth:  fromIntConfigValue(cfg.defaultModelExpandDepth),
		DefaultModelsExpandDepth: fromIntConfigValue(cfg.defaultModelsExpandDepth),
		DefaultModelRendering:    fromModelRenderingConfigValue(cfg.defaultModelRendering),
		QueryConfigEnabled:       fromBoolConfigValue(cfg.queryConfigEnabled),
		SupportedSubmitMethods:   strings.TrimSpace(strings.Join(cfg.supportedSubmitMethods, ",")),
		DeepLinking:              fromBoolConfigValue(cfg.deepLinking),
		ShowMutatedRequest:       fromBoolConfigValue(cfg.showMutatedRequest),
		ShowExtensions:           fromBoolConfigValue(cfg.showExtensions),
		ShowCommonExtensions:     fromBoolConfigValue(cfg.showCommonExtensions),
		Filter:                   fromBoolConfigValue(cfg.filter),
		FilterString:             fromStringConfigValue(cfg.filterString),
		DisplayOperationId:       fromBoolConfigValue(cfg.displayOperationID),
		TryItOutEnabled:          fromBoolConfigValue(cfg.tryItOutEnabled),
		DisplayRequestDuration:   fromBoolConfigValue(cfg.displayRequestDuration),
		PersistAuthorization:     fromBoolConfigValue(cfg.persistAuthorization),
		WithCredentials:          fromBoolConfigValue(cfg.withCredentials),
		OAuth2RedirectUrl:        fromStringConfigValue(cfg.oauth2RedirectUrl),
		Layout:                   fromStringConfigValue(cfg.oauth2RedirectUrl),
		ValidatorURL:             fromStringConfigValue(cfg.validatorUrl),
		MaxDisplayedTags:         fromIntConfigValue(cfg.maxDisplayedTags),
		PrimaryURL:               fromStringConfigValue(cfg.urlsPrimary),
		URLs:                     urlsAsJSON,
	})
}

func fromStringConfigValue(v configValue[string]) string {
	if v.IsSet {
		return strings.ReplaceAll(v.Value, "\n", "\\n")
	}

	return ""
}

func fromDocExpansionConfigValue(v configValue[DocExpansion]) string {
	if v.IsSet {
		return string(v.Value)
	}

	return ""
}

func fromModelRenderingConfigValue(v configValue[ModelRendering]) string {
	if v.IsSet {
		return string(v.Value)
	}

	return ""
}

func fromIntConfigValue(v configValue[int]) string {
	if v.IsSet {
		return fmt.Sprintf("%d", v.Value)
	}

	return ""
}

func fromBoolConfigValue(v configValue[bool]) string {
	if v.IsSet {
		return strconv.FormatBool(v.Value)
	}

	return ""
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func readSpecFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading spec file: %w", err)
	}

	return data, nil
}

func getContentType(fileName string, content []byte) string {
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		// read a chunk to decide between utf-8 text and binary
		var buf [512]byte
		n, _ := io.ReadFull(bytes.NewReader(content), buf[:])
		contentType = http.DetectContentType(buf[:n])
	}

	return contentType
}

func fromObject(v any) (string, error) {
	if v == nil {
		return "", nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("cannot marshal object: %w", err)
	}

	return string(b), nil
}
