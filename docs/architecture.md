# retrocode architecture (summary)

Thin client (Rust + FLTK on antiX, 2006 PC) <-> one framed TCP socket (SSH-forwarded)
<-> headless Go server (Codespaces). Server owns: filesystem, edit engine with
rev-tagged ops, inotify, LSP proxy, PTYs, agent loop, permissions, LLM router.
Client renders and captures input only. Full spec: proto/messages.json.
