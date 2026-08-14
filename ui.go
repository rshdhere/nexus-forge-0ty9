package main

import (
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Nexus Chat</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@400;600;700&family=IBM+Plex+Sans:wght@400;500;600&display=swap" rel="stylesheet" />
  <style>
    :root {
      --ink: #14212b;
      --muted: #5a6b78;
      --paper: #eef3f6;
      --panel: rgba(255, 255, 255, 0.82);
      --line: rgba(20, 33, 43, 0.1);
      --accent: #0f8f7b;
      --accent-deep: #0a6b5c;
      --mine: #d7f3ee;
      --theirs: #ffffff;
      --glow: rgba(15, 143, 123, 0.22);
      color-scheme: light;
    }
    * { box-sizing: border-box; }
    html, body { height: 100%; }
    body {
      margin: 0;
      font-family: "IBM Plex Sans", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(1200px 600px at 10% -10%, #c8ebe4 0%, transparent 55%),
        radial-gradient(900px 500px at 100% 0%, #d5e4f5 0%, transparent 50%),
        linear-gradient(165deg, #f7fafb 0%, var(--paper) 45%, #e4ecef 100%);
      min-height: 100vh;
    }
    .shell {
      width: min(920px, 100%);
      margin: 0 auto;
      min-height: 100vh;
      display: grid;
      grid-template-rows: auto 1fr;
      padding: 1.25rem 1rem 1rem;
      gap: 1rem;
    }
    header {
      display: flex;
      align-items: end;
      justify-content: space-between;
      gap: 1rem;
      animation: rise 0.55s ease both;
    }
    .brand {
      font-family: Outfit, sans-serif;
      font-weight: 700;
      font-size: clamp(1.8rem, 4vw, 2.4rem);
      letter-spacing: -0.03em;
      line-height: 1;
      margin: 0;
    }
    .brand span { color: var(--accent-deep); }
    .tag {
      margin: 0.35rem 0 0;
      color: var(--muted);
      font-size: 0.95rem;
    }
    .status {
      font-size: 0.85rem;
      color: var(--muted);
      display: flex;
      align-items: center;
      gap: 0.45rem;
      white-space: nowrap;
    }
    .dot {
      width: 0.55rem;
      height: 0.55rem;
      border-radius: 50%;
      background: #9aa8b3;
      box-shadow: 0 0 0 4px rgba(154, 168, 179, 0.2);
    }
    .dot.live {
      background: var(--accent);
      box-shadow: 0 0 0 4px var(--glow);
      animation: pulse 1.8s ease infinite;
    }
    .app {
      display: grid;
      grid-template-columns: 220px 1fr;
      gap: 0.9rem;
      min-height: 0;
      height: calc(100vh - 7.5rem);
      animation: rise 0.7s 0.08s ease both;
    }
    aside, .chat {
      background: var(--panel);
      border: 1px solid var(--line);
      backdrop-filter: blur(10px);
      border-radius: 1.1rem;
      min-height: 0;
    }
    aside {
      padding: 1rem;
      display: flex;
      flex-direction: column;
      gap: 0.85rem;
    }
    aside h2 {
      margin: 0;
      font-family: Outfit, sans-serif;
      font-size: 0.95rem;
      letter-spacing: 0.02em;
      text-transform: uppercase;
      color: var(--muted);
      font-weight: 600;
    }
    #online {
      list-style: none;
      margin: 0;
      padding: 0;
      display: flex;
      flex-direction: column;
      gap: 0.45rem;
      overflow: auto;
    }
    #online li {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.45rem 0.55rem;
      border-radius: 0.65rem;
      background: rgba(15, 143, 123, 0.08);
      font-size: 0.92rem;
      font-weight: 500;
    }
    #online li::before {
      content: "";
      width: 0.45rem;
      height: 0.45rem;
      border-radius: 50%;
      background: var(--accent);
    }
    .chat {
      display: grid;
      grid-template-rows: 1fr auto;
      overflow: hidden;
    }
    #messages {
      overflow: auto;
      padding: 1.1rem;
      display: flex;
      flex-direction: column;
      gap: 0.7rem;
      scroll-behavior: smooth;
    }
    .empty {
      margin: auto;
      text-align: center;
      color: var(--muted);
      max-width: 18rem;
      line-height: 1.45;
    }
    .msg {
      max-width: min(78%, 34rem);
      padding: 0.7rem 0.85rem;
      border-radius: 1rem;
      border: 1px solid var(--line);
      animation: pop 0.28s ease both;
    }
    .msg.mine {
      align-self: flex-end;
      background: var(--mine);
      border-bottom-right-radius: 0.35rem;
    }
    .msg.theirs {
      align-self: flex-start;
      background: var(--theirs);
      border-bottom-left-radius: 0.35rem;
    }
    .msg .meta {
      display: flex;
      justify-content: space-between;
      gap: 0.75rem;
      font-size: 0.75rem;
      color: var(--muted);
      margin-bottom: 0.25rem;
      font-weight: 600;
    }
    .msg .text {
      white-space: pre-wrap;
      word-break: break-word;
      line-height: 1.4;
    }
    form.composer {
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 0.55rem;
      padding: 0.85rem;
      border-top: 1px solid var(--line);
      background: rgba(255,255,255,0.55);
    }
    input, button {
      font: inherit;
    }
    input[type="text"] {
      width: 100%;
      border: 1px solid var(--line);
      border-radius: 0.8rem;
      padding: 0.8rem 0.95rem;
      background: #fff;
      color: var(--ink);
      outline: none;
    }
    input[type="text"]:focus {
      border-color: rgba(15, 143, 123, 0.55);
      box-shadow: 0 0 0 3px var(--glow);
    }
    button {
      border: 0;
      border-radius: 0.8rem;
      background: var(--accent);
      color: #fff;
      font-weight: 600;
      padding: 0.8rem 1.15rem;
      cursor: pointer;
      transition: transform 0.15s ease, background 0.15s ease;
    }
    button:hover { background: var(--accent-deep); }
    button:active { transform: translateY(1px) scale(0.98); }
    button:disabled { opacity: 0.55; cursor: not-allowed; }
    .gate {
      position: fixed;
      inset: 0;
      display: grid;
      place-items: center;
      padding: 1rem;
      background:
        radial-gradient(800px 400px at 20% 20%, rgba(200, 235, 228, 0.9), transparent 60%),
        rgba(238, 243, 246, 0.92);
      backdrop-filter: blur(8px);
      z-index: 10;
    }
    .gate[hidden] { display: none; }
    .card {
      width: min(420px, 100%);
      background: #fff;
      border: 1px solid var(--line);
      border-radius: 1.25rem;
      padding: 1.5rem;
      box-shadow: 0 20px 50px rgba(20, 33, 43, 0.08);
      animation: rise 0.45s ease both;
    }
    .card h1 {
      font-family: Outfit, sans-serif;
      font-size: 1.75rem;
      margin: 0 0 0.35rem;
      letter-spacing: -0.03em;
    }
    .card p {
      margin: 0 0 1.1rem;
      color: var(--muted);
      line-height: 1.45;
    }
    .card form {
      display: grid;
      gap: 0.65rem;
    }
    @keyframes rise {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: none; }
    }
    @keyframes pop {
      from { opacity: 0; transform: translateY(6px) scale(0.98); }
      to { opacity: 1; transform: none; }
    }
    @keyframes pulse {
      0%, 100% { box-shadow: 0 0 0 4px var(--glow); }
      50% { box-shadow: 0 0 0 7px rgba(15, 143, 123, 0.08); }
    }
    @media (max-width: 760px) {
      .app { grid-template-columns: 1fr; height: calc(100vh - 8.2rem); }
      aside { display: none; }
      header { align-items: start; flex-direction: column; }
    }
  </style>
</head>
<body>
  <div class="gate" id="gate">
    <div class="card">
      <h1>Nexus Chat</h1>
      <p>Pick a display name and jump into the shared room. Messages stream live over SSE.</p>
      <form id="joinForm">
        <input id="nameInput" type="text" maxlength="24" placeholder="Your name" autocomplete="nickname" required />
        <button type="submit">Enter chat</button>
      </form>
    </div>
  </div>

  <div class="shell">
    <header>
      <div>
        <p class="brand">Nexus <span>Chat</span></p>
        <p class="tag">Realtime room powered by Go</p>
      </div>
      <div class="status"><span class="dot" id="liveDot"></span><span id="liveLabel">Connecting…</span></div>
    </header>

    <div class="app">
      <aside>
        <h2>Online</h2>
        <ul id="online"></ul>
      </aside>
      <section class="chat" aria-label="Chat">
        <div id="messages"><p class="empty">No messages yet. Say hello.</p></div>
        <form class="composer" id="sendForm">
          <div style="display:grid;gap:0.35rem">
            <input id="textInput" type="text" maxlength="1000" placeholder="Write a message…" autocomplete="off" disabled />
            <span id="charCount" style="font-size:0.75rem;color:var(--muted);padding-left:0.2rem">0 / 1000</span>
          </div>
          <button type="submit" id="sendBtn" disabled>Send</button>
        </form>
      </section>
    </div>
  </div>

  <script>
    const state = { user: "", messages: [], online: [], es: null };

    const gate = document.getElementById("gate");
    const joinForm = document.getElementById("joinForm");
    const nameInput = document.getElementById("nameInput");
    const sendForm = document.getElementById("sendForm");
    const textInput = document.getElementById("textInput");
    const sendBtn = document.getElementById("sendBtn");
    const messagesEl = document.getElementById("messages");
    const onlineEl = document.getElementById("online");
    const liveDot = document.getElementById("liveDot");
    const liveLabel = document.getElementById("liveLabel");

    const saved = localStorage.getItem("nexus-chat-user");
    if (saved) nameInput.value = saved;

    joinForm.addEventListener("submit", (e) => {
      e.preventDefault();
      const name = nameInput.value.trim();
      if (!name) return;
      state.user = name;
      localStorage.setItem("nexus-chat-user", name);
      gate.hidden = true;
      textInput.disabled = false;
      sendBtn.disabled = false;
      textInput.focus();
      connect();
      heartbeat();
      setInterval(heartbeat, 15000);
    });

    textInput.addEventListener("input", () => {
      document.getElementById("charCount").textContent = textInput.value.length + " / 1000";
    });

    sendForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const text = textInput.value.trim();
      if (!text || !state.user) return;
      textInput.value = "";
      sendBtn.disabled = true;
      try {
        const res = await fetch("/api/messages", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ user: state.user, text }),
        });
        if (!res.ok) throw new Error("send failed");
      } catch (err) {
        textInput.value = text;
        alert("Could not send message.");
      } finally {
        sendBtn.disabled = false;
        textInput.focus();
      }
    });

    async function heartbeat() {
      if (!state.user) return;
      try {
        await fetch("/api/presence", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ user: state.user }),
        });
      } catch (_) {}
    }

    function connect() {
      if (state.es) state.es.close();
      const es = new EventSource("/api/stream");
      state.es = es;
      es.onopen = () => {
        liveDot.classList.add("live");
        liveLabel.textContent = "Live";
      };
      es.onerror = () => {
        liveDot.classList.remove("live");
        liveLabel.textContent = "Reconnecting…";
      };
      es.onmessage = (ev) => {
        try {
          const payload = JSON.parse(ev.data);
          if (payload.type === "snapshot") {
            state.messages = payload.data.messages || [];
            state.online = payload.data.online || [];
            renderMessages();
            renderOnline();
          } else if (payload.type === "message") {
            state.messages.push(payload.data);
            if (state.messages.length > 500) state.messages = state.messages.slice(-500);
            renderMessages(true);
          } else if (payload.type === "presence") {
            state.online = payload.data || [];
            renderOnline();
          }
        } catch (_) {}
      };
    }

    function renderMessages(stick) {
      const nearBottom = messagesEl.scrollHeight - messagesEl.scrollTop - messagesEl.clientHeight < 80;
      if (!state.messages.length) {
        messagesEl.innerHTML = '<p class="empty">No messages yet. Say hello.</p>';
        return;
      }
      messagesEl.innerHTML = state.messages.map((m) => {
        const mine = m.user === state.user;
        const when = new Date(m.createdAt).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
        const who = mine ? escapeHtml(m.user) + " (you)" : escapeHtml(m.user);
        return '<article class="msg ' + (mine ? "mine" : "theirs") + '">' +
          '<div class="meta"><span>' + who + '</span><span>' + when + '</span></div>' +
          '<div class="text">' + escapeHtml(m.text) + '</div></article>';
      }).join("");
      if (stick || nearBottom) messagesEl.scrollTop = messagesEl.scrollHeight;
    }

    function renderOnline() {
      const names = [...new Set(state.online)].sort((a, b) => a.localeCompare(b));
      if (!names.length) {
        onlineEl.innerHTML = '<li style="background:transparent;color:var(--muted)">Waiting for people…</li>';
        return;
      }
      onlineEl.innerHTML = names.map((n) => "<li>" + escapeHtml(n) + "</li>").join("");
    }

    function escapeHtml(s) {
      return String(s)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;");
    }
  </script>
</body>
</html>`
