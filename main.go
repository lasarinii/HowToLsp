package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/lasarinii/howtolsp/analysis"
	"github.com/lasarinii/howtolsp/lsp"
	"github.com/lasarinii/howtolsp/rpc"
)

func main() {
    logger := getLogger("/home/euzinho/Coding/howtolsp/log.txt")
    logger.Println("Started")

    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(rpc.Split)

    state := analysis.NewState()
    writer := os.Stdout

    for scanner.Scan() {
        msg := scanner.Bytes()
        method, contents, err := rpc.DecodeMessage(msg)
        if err != nil {
            logger.Printf("Got error: %s", err)
            continue
        }
        handleMessage(logger, writer, state, method, contents)
    }
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
    logger.Printf("Got message with method: %s", method)

    switch method {
    case "initialize":
        var request lsp.InitializeRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("Couldn't unmarshal: %s", err)
        }

        logger.Printf("Connected to: %s, %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)

        msg := lsp.NewInitializeResponse(request.ID)
        writeResponse(writer, msg)

        logger.Print("Sent the reply")

    case "textDocument/didOpen":
        var request lsp.DidOpenTextDocumentNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didOpen: %s", err)
            return
        }

        logger.Printf("Opened: %s", request.Params.TextDocument.URI)
        diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
        writeResponse(writer, lsp.PublishDiagnosticsNotification{
            Notification: lsp.Notification{
                RPC: "2.0",
                Method: "textDocument/publishDiagnostics",
            },
            Params: lsp.PublishDiagnosticsParams{
                URI: request.Params.TextDocument.URI,
                Diagnostics: diagnostics,
            },
        })

    case "textDocument/didChange":
        var request lsp.TextDocumentDidChangeNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didChange: %s", err)
            return
        }

        logger.Printf("Changed: %s", request.Params.TextDocument.URI)
        for _, change := range request.Params.ContentChanges {
            diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
            writeResponse(writer, lsp.PublishDiagnosticsNotification{
            	Notification: lsp.Notification{
                    RPC: "2.0",
                    Method: "textDocument/publishDiagnostics",
                },
            	Params:       lsp.PublishDiagnosticsParams{
                    URI: request.Params.TextDocument.URI,
                    Diagnostics: diagnostics,
                },
            })
        }

    case "textDocument/hover":
        var request lsp.HoverRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/hover: %s", err)
            return
        }

        response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)

        writeResponse(writer, response)

    case "textDocument/definition":
        var request lsp.DefinitionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/definition: %s", err)
            return
        }

        response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)

        writeResponse(writer, response)

    case "textDocument/codeAction":
        var request lsp.CodeActionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/codeAction: %s", err)
            return
        }

        response := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)

        writeResponse(writer, response)

    case "textDocument/completion":
        var request lsp.CompletionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/completion: %s", err)
            return
        }

        response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)

        writeResponse(writer, response)
    }
}

func writeResponse(writer io.Writer, msg any) {
    reply := rpc.EncodeMessage(msg)
    writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
    logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
    if err != nil {
        panic("git me a good file")
    }

    return log.New(logfile, "[howtolsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
