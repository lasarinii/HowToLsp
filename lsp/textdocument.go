package lsp

type TextDocumentItem struct {
    URI string `json:"uri"`
    LanguageID string `json:"languageId"`
    Version int `json:"version"`
    Text string `json:"text"`
}

type TextDocumentIdentifier struct {
    URI string `json:"uri"`
}

type TextDocumentPositionParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position Position
}

type Position struct {
    Line int `json:"line"`
    Character int `json:"character"`
}

type Location struct {
    URI string `json:"uri"`
    Range Range `json:"range"`
}

type Range struct {
    Start Position `json:"start"`
    End Position `json:"end"`
}

type WorkspaceEdit struct {
    Changes map[string][]TextEdit `json:"changes"`
}

type TextEdit struct {
    Range Range `json:"range"`
    NewText string `json:"newText"`
}

type VersionTextDocumentIdentifier struct {
    TextDocumentIdentifier
    Version int `json:"version"`
}

type DidOpenTextDocumentNotification struct {
    Notification
    Params DidOpenTextDocumentParams `json:"params"`
}

type DidOpenTextDocumentParams struct {
    TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentDidChangeNotification struct {
    Notification
    Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
    TextDocument VersionTextDocumentIdentifier `json:"textDocument"`
    ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
    Text string `json:"text"`
}

type HoverRequest struct {
    Request
    Params HoverParams `json:"params"`
}

type HoverParams struct {
    TextDocumentPositionParams
}

type HoverResponse struct {
    Response
    Result HoverResult `json:"result"`
}

type HoverResult struct {
    Contents string `json:"contents"`
}

type DefinitionRequest struct {
    Request
    Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
    TextDocumentPositionParams
}

type DefinitionResponse struct {
    Response
    Result Location `json:"result"`
}

type CodeActionRequest struct {
    Request
    Params TextDocumentCodeActionParams `json:"params"`
}

type TextDocumentCodeActionResponse struct {
    Response
    Result []CodeAction `json:"result"`
}

type TextDocumentCodeActionParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Range Range `json:"range"`
    Context CodeActionContext `json:"context"`
}

type CodeActionContext struct {
    // ...
}

type CodeAction struct {
    Title string `json:"title"`
    Edit *WorkspaceEdit `json:"edit,omidempty"`
    Command *Command `json:"command,omitifempty"`
}

type Command struct {
    Title string `json:"title"`
    Command string  `json:"command"`
    Arguments []interface{} `json:"arguments,omitifempty"`
}

type CompletionRequest struct {
    Request
    Params CompletionParams `json:"params"`
}

type CompletionParams struct {
    TextDocumentPositionParams
}

type CompletionResponse struct {
    Response
    Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
    Label string `json:"label"`
    Detail string `json:"detail"`
    Documentation string `json:"documentation"`
}

type PublishDiagnosticsNotification struct {
    Notification
    Params PublishDiagnosticsParams `json:"params"`
}

type PublishDiagnosticsParams struct {
    URI string `json:"uri"`
    Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
    Range Range `json:"range"`
    Severity int `json:"severity"`
    Source string `json:"source"`
    Message string `json:"message"`
}
