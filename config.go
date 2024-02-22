package go_swagger_ui

type configValue[T any] struct {
	IsSet bool
	Value T
}

type uiConfig struct {
	htmlTitle                string
	spec                     []byte
	configURL                configValue[string]
	specFilePath             string
	url                      configValue[string]
	urls                     []SpecURL
	urlsPrimary              configValue[string]
	layout                   configValue[string]
	docExpansion             configValue[DocExpansion]
	defaultModelExpandDepth  configValue[int]
	defaultModelsExpandDepth configValue[int]
	defaultModelRendering    configValue[ModelRendering]
	queryConfigEnabled       configValue[bool]
	supportedSubmitMethods   []string
	showMutatedRequest       configValue[bool]
	deepLinking              configValue[bool]
	showExtensions           configValue[bool]
	showCommonExtensions     configValue[bool]
	filter                   configValue[bool]
	filterString             configValue[string]
	displayOperationID       configValue[bool]
	tryItOutEnabled          configValue[bool]
	displayRequestDuration   configValue[bool]
	persistAuthorization     configValue[bool]
	withCredentials          configValue[bool]
	oauth2RedirectUrl        configValue[string]
	maxDisplayedTags         configValue[int]
	validatorUrl             configValue[string]
}

type DocExpansion string

var (
	DocExpansionList DocExpansion = "list"
	DocExpansionFull DocExpansion = "full"
	DocExpansionNone DocExpansion = "none"
)

type ModelRendering string

var (
	ModelRenderingExample ModelRendering = "example"
	ModelRenderingModel   ModelRendering = "model"
)

type Layout string

var (
	LayoutBaseLayout Layout = "BaseLayout"
)

type Preset string

var (
	PresetAPIPreset Layout = "ApiPreset"
)

type SpecURL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Option is a function that takes a pointer to uiConfig and modifies it.
type Option func(*uiConfig)

// WithSpec sets the spec field of https://github.com/swagger-api/swagger-ui/blob/HEAD/docs/usage/configuration.md..

// WithSpec sets an OpenAPI specification document content. When used, the URL configuration setting will not be used.
// This is useful for testing manually-generated definitions without hosting them.
func WithSpec(value []byte) Option {
	return func(cfg *uiConfig) {
		cfg.spec = value
	}
}

// WithSpecURL sets the URL pointing to API definition (normally swagger.json or swagger.yaml).
// Will be ignored if WithSpecURLs or WithSpec is used.
func WithSpecURL(value string) Option {
	return func(cfg *uiConfig) {
		cfg.url = configValue[string]{
			IsSet: true,
			Value: value,
		}
	}
}

// WithSpecURLs sets the URLs array to multiple API definitions that are used by Topbar plugin.
// When used and Topbar plugin is enabled, the settings from WithSpecURL will not be used.
// Names and URLs must be unique among all items in this array, since they're used as identifiers.
// If the value of the 'primary' parameter matches the name of a spec provided in urls, that spec will
// be displayed when Swagger UI loads, instead of defaulting to the first spec in urls.
// Leave parameter 'primary' empty, if you do not want to set a preselected URL.
func WithSpecURLs(primary string, urls []SpecURL) Option {
	return func(cfg *uiConfig) {
		cfg.urls = urls
		if len(primary) > 0 {
			cfg.urlsPrimary = configValue[string]{
				IsSet: true,
				Value: primary,
			}
		}
	}
}

// WithSpecFilePath sets a file path to read from the OS file system.
// THIS OPTION IS NOT RECOMMENDED FOR PRODUCTION USE, because it reloads the file on every request.
// This option only exist to for testing purposes. Once file content is read, it will be used to set the spec field of
// https://github.com/swagger-api/swagger-ui/blob/HEAD/docs/usage/configuration.md and is equivalent to the
// WithSpec function.
func WithSpecFilePath(path string) Option {
	return func(cfg *uiConfig) {
		cfg.specFilePath = path
	}
}

// WithDocExpansion controls the default expansion setting for the operations and tags.
func WithDocExpansion(value DocExpansion) Option {
	return func(cfg *uiConfig) {
		cfg.docExpansion = configValue[DocExpansion]{
			IsSet: true,
			Value: value,
		}
	}
}

// WithDefaultModelExpandDepth sets the default expansion depth for the model on the model-example section.
func WithDefaultModelExpandDepth(defaultModelExpandDepth int) Option {
	return func(cfg *uiConfig) {
		cfg.defaultModelExpandDepth = configValue[int]{Value: defaultModelExpandDepth, IsSet: true}
	}
}

// WithDefaultModelsExpandDepth sets the default expansion depth for models
// (set to -1 completely hide the models).
func WithDefaultModelsExpandDepth(defaultModelsExpandDepth int) Option {
	return func(cfg *uiConfig) {
		cfg.defaultModelsExpandDepth = configValue[int]{Value: defaultModelsExpandDepth, IsSet: true}
	}
}

// WithDefaultModelRendering controls how the model is shown when the API is first rendered.
// The user can always switch the rendering for a given model by clicking the 'Model' and 'Example Value' links.
func WithDefaultModelRendering(defaultModelRendering ModelRendering) Option {
	return func(cfg *uiConfig) {
		cfg.defaultModelRendering = configValue[ModelRendering]{Value: defaultModelRendering, IsSet: true}
	}
}

// WithQueryConfigEnabled enables overriding configuration parameters via URL search params.
func WithQueryConfigEnabled(queryConfigEnabled bool) Option {
	return func(cfg *uiConfig) {
		cfg.queryConfigEnabled = configValue[bool]{Value: queryConfigEnabled, IsSet: true}
	}
}

// WithSupportedSubmitMethods sets a list of HTTP methods that have the "Try it out" feature enabled.
// An empty array disables "Try it out" for all operations. This does not filter the operations from the display.
// Default is: ["get", "put", "post", "delete", "options", "head", "patch", "trace"].
func WithSupportedSubmitMethods(supportedSubmitMethods ...string) Option {
	return func(cfg *uiConfig) {
		cfg.supportedSubmitMethods = append(cfg.supportedSubmitMethods, supportedSubmitMethods...)
	}
}

// WithDeepLinking enables deep linking. See documentation at
// https://swagger.io/docs/open-source-tools/swagger-ui/usage/deep-linking/
// for more information.
func WithDeepLinking(deepLinking bool) Option {
	return func(cfg *uiConfig) {
		cfg.deepLinking = configValue[bool]{Value: deepLinking, IsSet: true}
	}
}

// WithShowExtensions controls the display of vendor extension (x-) fields and values for
// Operations, Parameters, Responses, and Schema.
func WithShowExtensions(showExtensions bool) Option {
	return func(cfg *uiConfig) {
		cfg.showExtensions = configValue[bool]{Value: showExtensions, IsSet: true}
	}
}

// WithShowCommonExtensions controls the display of extensions (pattern, maxLength, minLength, maximum, minimum)
// fields and values for Parameters.
func WithShowCommonExtensions(showCommonExtensions bool) Option {
	return func(cfg *uiConfig) {
		cfg.showCommonExtensions = configValue[bool]{Value: showCommonExtensions, IsSet: true}
	}
}

// WithFilter enables filtering. The top bar will show an edit box that you can use to filter the tagged
// operations that are shown. If enabled and a non-empty expression string is passed, then filtering
// will be enabled using that string as the filter expression. Filtering is case-sensitive matching
// the filter expression anywhere inside the tag. Leave the expression empty, if you only want to
// enable filtering but do not need a filter expression.
func WithFilter(enabled bool, expression string) Option {
	return func(cfg *uiConfig) {
		cfg.filter = configValue[bool]{Value: enabled, IsSet: true}
		if enabled && len(expression) > 0 {
			cfg.filterString = configValue[string]{Value: expression, IsSet: true}
		}
	}
}

// WithDisplayOperation controls the display of operationId in operations list. The default is false.
func WithDisplayOperation(displayOperationID bool) Option {
	return func(cfg *uiConfig) {
		cfg.displayOperationID = configValue[bool]{Value: displayOperationID, IsSet: true}

	}
}

// WithTryItOutEnabled controls whether the "Try it out" section should be enabled by default.
func WithTryItOutEnabled(tryItOutEnabled bool) Option {
	return func(cfg *uiConfig) {
		cfg.tryItOutEnabled = configValue[bool]{Value: tryItOutEnabled, IsSet: true}

	}
}

// WithDisplayRequestDuration controls the display of the request duration (in milliseconds) for "Try it out" requests.
func WithDisplayRequestDuration(displayRequestDuration bool) Option {
	return func(cfg *uiConfig) {
		cfg.displayRequestDuration = configValue[bool]{Value: displayRequestDuration, IsSet: true}
	}
}

// WithPersistAuthorization configures Swagger UI to persist authorization data, so that it is not lost
// on browser close/refresh.
func WithPersistAuthorization(persistAuthorization bool) Option {
	return func(cfg *uiConfig) {
		cfg.persistAuthorization = configValue[bool]{Value: persistAuthorization, IsSet: true}
	}
}

// WithCredentials enables passing credentials, as defined in the Fetch standard, in CORS requests that are
// sent by the browser. Note that Swagger UI cannot currently set cookies cross-domain (see swagger-js#1163) -
// as a result, you will have to rely on browser-supplied cookies (which this setting enables sending)
// that Swagger UI cannot control.
func WithCredentials(withCredentials bool) Option {
	return func(cfg *uiConfig) {
		cfg.withCredentials = configValue[bool]{Value: withCredentials, IsSet: true}
	}
}

// WithOauth2RedirectUrl sets the OAuth redirect URL.
func WithOauth2RedirectUrl(oauth2RedirectUrl string) Option {
	return func(cfg *uiConfig) {
		cfg.oauth2RedirectUrl = configValue[string]{Value: oauth2RedirectUrl, IsSet: true}
	}
}

// WithHTMLTitle sets the index HTML page htmlTitle.
func WithHTMLTitle(title string) Option {
	return func(cfg *uiConfig) {
		cfg.htmlTitle = title
	}
}

// WithLayout sets the name of a component available via the plugin system to use as the top-level
// layout for Swagger UI.
func WithLayout(layout Layout) Option {
	return func(cfg *uiConfig) {
		cfg.layout = configValue[string]{Value: string(layout), IsSet: true}
	}
}

// WithPresets sets the list of presets to use in Swagger UI. Usually, you'll want to
// include PresetAPIPreset if you use this option.
func WithPresets(preset Preset) Option {
	return func(cfg *uiConfig) {
		cfg.layout = configValue[string]{Value: string(preset), IsSet: true}
	}
}

// WithMaxDisplayedTags limits the number of tagged operations displayed to at most this many.
// The default is to show all operations.
func WithMaxDisplayedTags(maxTags int) Option {
	return func(cfg *uiConfig) {
		cfg.maxDisplayedTags = configValue[int]{Value: maxTags, IsSet: true}
	}
}

// WithValidatorURL sets the validator URL to use to validate specification files. By default, Swagger UI
// attempts to validate specs against swagger.io's online validator. You can use this parameter to set a
// different validator URL, for example for locally deployed validators (e.g., Validator Badge,
// see https://github.com/swagger-api/validator-badge). Disabling it or setting the URL to 127.0.0.1 or localhost
// will disable validation.
func WithValidatorURL(enabled bool, validatorUrl string) Option {
	return func(cfg *uiConfig) {
		var valURL string
		if enabled {
			valURL = validatorUrl
		}

		cfg.validatorUrl = configValue[string]{Value: valURL, IsSet: enabled}
	}
}

// WithShowMutatedRequest configures Swagger UI to use the mutated request returned from a
// requestInterceptor to produce the curl command in the UI, otherwise the request before
// the requestInterceptor was applied is used.
// Refer to https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/ for more information.
func WithShowMutatedRequest(showMutatedRequest bool) Option {
	return func(cfg *uiConfig) {
		cfg.showMutatedRequest = configValue[bool]{Value: showMutatedRequest, IsSet: true}
	}
}

// WithConfigURL sets the URL to fetch external configuration document from.
func WithConfigURL(configURL string) Option {
	return func(cfg *uiConfig) {
		cfg.configURL = configValue[string]{Value: configURL, IsSet: true}
	}
}
