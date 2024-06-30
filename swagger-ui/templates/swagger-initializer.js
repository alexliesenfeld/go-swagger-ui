window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    dom_id: '#swagger-ui',
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    configUrl: blankToUndefined('{{ .ConfigURL }}'),
    spec: blankToUndefined(decodeHtmlEntities('{{ .Spec }}')),
    url: blankToUndefined('{{ .URL }}'),
    docExpansion: blankToUndefined('{{ .DocExpansion }}'),
    defaultModelExpandDepth: blankToUndefinedNumber('{{ .DefaultModelExpandDepth }}'),
    defaultModelsExpandDepth: blankToUndefinedNumber('{{ .DefaultModelsExpandDepth }}'),
    defaultModelRendering: blankToUndefined('{{ .DefaultModelRendering }}'),
    queryConfigEnabled: blankToUndefinedBool('{{ .QueryConfigEnabled }}'),
    supportedSubmitMethods: blankToUndefinedArray('{{ .SupportedSubmitMethods }}') || ["get", "put", "post", "delete", "options", "head", "patch", "trace"],
    deepLinking: blankToUndefinedBool('{{ .DeepLinking }}'),
    showMutatedRequest: blankToUndefinedBool('{{ .ShowMutatedRequest }}'),
    showExtensions: blankToUndefinedBool('{{ .ShowExtensions }}'),
    showCommonExtensions: blankToUndefinedBool('{{ .ShowCommonExtensions }}'),
    filter: blankToUndefined('{{ .Filter }}') || blankToUndefinedBool('{{ .FilterString }}'),
    displayOperationId: blankToUndefinedBool('{{ .DisplayOperationId }}'),
    tryItOutEnabled: blankToUndefinedBool('{{ .TryItOutEnabled }}'),
    displayRequestDuration: blankToUndefinedBool('{{ .DisplayRequestDuration }}'),
    persistAuthorization: blankToUndefinedBool('{{ .PersistAuthorization }}'),
    withCredentials: blankToUndefinedBool('{{ .WithCredentials }}'),
    oauth2RedirectUrl: blankToUndefined('{{ .OAuth2RedirectUrl }}'),
    layout: blankToUndefined('{{ .Layout }}') || "BaseLayout",
    validatorUrl: blankToUndefined('{{ .ValidatorURL }}'),
    maxDisplayedTags: blankToUndefinedNumber('{{ .MaxDisplayedTags }}'),
    urls: blankToUndefinedObject('{{ .URLs }}'),
    "urls.primaryName": blankToUndefined('{{ .PrimaryURL }}'),
  });

  //</editor-fold>
};

function blankToUndefined(input) {
  return (input || '').trim() === '' ? undefined : input
}

function blankToUndefinedNumber(input) {
  let cleanInput = blankToUndefined(input)
  if (!cleanInput) {
    return undefined
  }

  const parsed = parseInt(cleanInput, 10);
  if (isNaN(parsed)) {
    return undefined;
  }

  return parsed
}

function blankToUndefinedBool(input) {
  if (input === 'true') return true;
  if (input === 'false') return false;
  return undefined
}

function blankToUndefinedArray(input) {
  const arr = input.split(",").map(blankToUndefined).filter(e => !!e)
  if ((arr || []).length === 0) {
    return undefined
  }

  return arr
}

function blankToUndefinedObject(input) {
  if (!input) {
    return undefined
  }

  input = input.trim()
  if (input.length === 0) {
    return undefined
  }

  input = decodeHtmlEntities(input);

  return JSON.parse(input)
}

function decodeHtmlEntities(str) {
  if (!str) {
    return undefined
  }

  const textArea = document.createElement('textarea');
  textArea.innerHTML = str;
  return textArea.value;
}