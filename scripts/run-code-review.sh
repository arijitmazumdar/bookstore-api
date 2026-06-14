#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

prompt_file=".github/prompts/code-reviewer.prompt.md"
if [ ! -f "$prompt_file" ]; then
  echo "Prompt file not found: $prompt_file"
  exit 1
fi

copilot_cmd=""
if command -v copilot >/dev/null 2>&1; then
  copilot_cmd="copilot"
elif [ -x "$HOME/.vscode-remote/data/User/globalStorage/github.copilot-chat/copilotCli/copilot" ]; then
  copilot_cmd="$HOME/.vscode-remote/data/User/globalStorage/github.copilot-chat/copilotCli/copilot"
else
  echo "Copilot CLI not found. Please install it or add it to PATH."
  exit 1
fi

echo "Executing code-reviewer prompt via Copilot CLI"
echo "Prompt file: $prompt_file"

prompt_text=$(cat "$prompt_file")
$copilot_cmd -p "$prompt_text" --allow-all --silent

echo "Running golangci-lint static analysis"
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "Installing golangci-lint..."
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.56.0
fi

golangci-lint run ./...
