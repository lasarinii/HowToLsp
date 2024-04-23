package analysis

import (
	"fmt"
	"strings"

	"github.com/lasarinii/howtolsp/lsp"
)

type State struct {
    Documents map[string]string
}

func NewState() State {
    return State{Documents: map[string]string{}}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
    diagnostics := []lsp.Diagnostic{}
    for row, line := range strings.Split(text, "\n") {
        if strings.Contains(line, "VS Code") {
            i := strings.Index(line, "VS Code")
            diagnostics = append(diagnostics, lsp.Diagnostic{
            	Range:    LineRange(row, i, i+len("VS Code")),
            	Severity: 1,
            	Source:   "Common Sense",
            	Message:  "Make it NeoVim.",
            })
        }

        if strings.Contains(line, "NeoVim") {
            i := strings.Index(line, "NeoVim")
            diagnostics = append(diagnostics, lsp.Diagnostic{
            	Range:    LineRange(row, i, i+len("NeoVim")),
            	Severity: 4,
            	Source:   "Common Sense",
            	Message:  "Great Choice",
            })
        }
    }
    return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
    s.Documents[uri] = text

    return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
    s.Documents[uri] = text

    return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
    document := s.Documents[uri]

    return lsp.HoverResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: lsp.HoverResult{
            Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
        },
    }
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
    return lsp.DefinitionResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: lsp.Location{
            URI: uri,
            Range: lsp.Range{
                Start: lsp.Position{
                    Line: position.Line - 1,
                    Character: 0,
                },
                End: lsp.Position{
                    Line: position.Line - 1,
                    Character: 0,
                },
            },
        },
    }
}

func (s *State) TextDocumentCodeAction(id int, uri string) lsp.TextDocumentCodeActionResponse {
    text := s.Documents[uri]

    actions := []lsp.CodeAction{}
    for row, line := range strings.Split(text, "\n") {
        i := strings.Index(line, "VS Code")
        if i >= 0 {
            replaceChange := map[string][]lsp.TextEdit{}
            replaceChange[uri] = []lsp.TextEdit{
                {
                    Range: LineRange(row, i, i + len("VS Code")),
                    NewText: "NeoVim",
                },
            }

            actions = append(actions, lsp.CodeAction{
                Title: "Replace VS C*de with a superior editor.",
                Edit: &lsp.WorkspaceEdit{Changes: replaceChange},
            })

            censorChange := map[string][]lsp.TextEdit{}
            censorChange[uri] = []lsp.TextEdit{
                {
                    Range: LineRange(row, i, i + len("VS Code")),
                    NewText: "VS C*de",
                },
            }

            actions = append(actions, lsp.CodeAction{
                Title: "Censor to 'VS C*de'.",
                Edit: &lsp.WorkspaceEdit{Changes: censorChange},
            })
        }
    }
    response := lsp.TextDocumentCodeActionResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: actions,
    }

    return response
}

func (s *State) TextDocumentCompletion(id int, uri string) lsp.CompletionResponse {
    // analysis
    items := []lsp.CompletionItem{
        {
            Label: "NeoVim (BTW)",
            Detail: "The very cool editor",
            Documentation: "To use NeoVim, first you need to be very cool, second you need to hate VS C*de, thrid you need to use Neovim, and finally you need to repeat the rules.",
        },
    }

    response := lsp.CompletionResponse{
        Response: lsp.Response{
            RPC: "2.0",
            ID: &id,
        },
        Result: items,
    }

    return response
}

func LineRange(line, start, end int) lsp.Range {
    return lsp.Range{
        Start: lsp.Position{
            Line: line,
            Character: start,
        },
        End: lsp.Position{
            Line: line,
            Character: end,
        },
    }
}
