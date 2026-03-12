# tts-mcp

A Go MCP server that gives Claude the ability to speak. Uses xAI or OpenAI for text-to-speech. Streams PCM audio directly to your speakers in real-time.

## Setup

**1. Get an API key**
- xAI: [console.x.ai](https://console.x.ai/team/default/api-keys)
- OpenAI: [platform.openai.com](https://platform.openai.com/api-keys)

**2. Build**

```bash
build.bat
```

**3. Add to Claude Code**

```bash
claude mcp add -s user tts-mcp -e API_KEY=YOUR_KEY -e PROVIDER=xai -- build/tts-mcp.exe -voice ara
```

## Configuration

| Source | Name | Description |
|--------|------|-------------|
| Env / Flag | `PROVIDER` / `-provider` | `xai` or `openai` (default: `xai`) |
| Env / Flag | `API_KEY` / `-key` | API key for the selected provider |
| Flag | `-voice` | Default voice (see voices below) |
| Flag | `-lang` | Language code, xAI only (default: `en`) |

Flags override env vars. `API_KEY` is required.

## MCP Tools

### `talk`

| Parameter | Required | Description |
|-----------|----------|-------------|
| `text` | yes | Text to speak aloud. Long text is automatically chunked. |
| `voice` | no | Per-call voice override. |

## Voices

### xAI

| Voice | Tone |
|-------|------|
| `eve` | Energetic, upbeat |
| `ara` | Warm, friendly |
| `rex` | Confident, clear |
| `sal` | Smooth, balanced |
| `leo` | Authoritative, strong |

### OpenAI

| Voice | Tone |
|-------|------|
| `alloy` | Neutral, balanced |
| `ash` | Soft, warm |
| `coral` | Clear, expressive |
| `echo` | Smooth, resonant |
| `fable` | Warm, narrative |
| `nova` | Bright, friendly |
| `onyx` | Deep, authoritative |
| `sage` | Calm, measured |
| `shimmer` | Light, upbeat |

## License

MIT
