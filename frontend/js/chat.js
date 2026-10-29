(function (App) {
  var $ = App.$;
  var state = App.state;

  function refreshComposer() {
    var on = Boolean(state.active && state.streams[state.active]);
    $("send").hidden = on;
    $("stop").hidden = !on;
    $("input").disabled = on;
    if (!on) {
      $("input").focus();
      $("input").style.height = "auto";
    }
  }
  async function readStream(res, handle) {
    var reader = res.body.getReader();
    var dec = new TextDecoder();
    var buf = "";
    for (;;) {
      var r = await reader.read();
      if (r.done) break;
      buf += dec.decode(r.value, { stream: true });
      var cut = buf.indexOf("\n\n");
      while (cut >= 0) {
        var frame = buf.slice(0, cut);
        buf = buf.slice(cut + 2);
        readFrame(frame, handle);
        cut = buf.indexOf("\n\n");
      }
    }
  }
  function readFrame(frame, handle) {
    var name = "";
    var data = "";
    frame.split("\n").forEach(function (line) {
      if (line.indexOf("event: ") === 0) name = line.slice(7);
      else if (line.indexOf("data: ") === 0) data += line.slice(6);
    });
    if (name && data) handle(name, JSON.parse(data));
  }
  async function refreshHeader(id) {
    var c = await App.api("GET", "/v1/conversations/" + id);
    var view = $("view");
    var title = view.querySelector(".title");
    var meta = view.querySelector(".meta");
    if (!c || state.active !== id || !title || !meta) return;
    title.textContent = c.title || "Untitled";
    meta.textContent = "Cập nhật " + new Date(c.updated_at).toLocaleString() + " · " + view.querySelectorAll(".msg").length + " tin";
  }
  async function ensureConversation() {
    if (state.active) return true;
    var c = await App.api("POST", "/v1/conversations", {});
    if (!c) return false;
    state.wanted = c.id;
    state.active = c.id;
    history.pushState({}, "", "/conversation/" + c.id + location.search);
    App.renderConversation(c, []);
    return true;
  }
  function onDelta(ui, d) {
    ui.mdBuf += d.text;
    if (!ui.mdContent) {
      ui.mdContent = App.el("div", "md-content");
      ui.bubble.append(ui.mdContent);
    }
    App.renderMd(ui.mdContent, ui.mdBuf);
    $("view").scrollTop = $("view").scrollHeight;
  }
  function onCompleted(ui, d) {
    ui.chip.textContent = d.model + (d.status !== "completed" ? " · " + d.status : "") + " · " + App.usd(d.cost_micro_usd) + " · TTFT " + d.first_token_ms + " ms · " + d.latency_ms + " ms";
    App.markStatus(ui, d.status);
    if (!ui.row.isConnected) return;
    if (d.status === "truncated") App.showContinue(ui, state.active);
    else App.removeContinue();
  }
  function onFailed(ui, d) {
    App.markStatus(ui, "failed");
    ui.chip.textContent = ui.chip.textContent + " " + d.code + " · request_id " + d.request_id;
  }
  async function streamReply(url, body, ui, convId, signal) {
    var started = performance.now();
    var first = 0;
    var res = await fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
      signal: signal
    });
    var rid = res.headers.get("X-Request-Id") || "";
    if (!res.ok) {
      var j = null;
      try { j = await res.json(); } catch (e) { void e; }
      App.showError({ status: res.status, data: j && j.data, rid: rid });
      if (ui.continuing) App.showContinue(ui, convId);
      else ui.row.remove();
      return;
    }
    App.showError(null);
    var thinking = App.el("div", "thinking");
    thinking.innerHTML = '<span>Thinking</span><div class="thinking-dots"><span></span><span></span><span></span></div>';
    ui.bubble.append(thinking);
    var searchBlock = null;
    var sources = [];
    await readStream(res, function (name, d) {
      if (name === "search.started") {
        searchBlock = App.el("div", "search-block searching");
        searchBlock.innerHTML = '<span class="search-icon">&#128269;</span><span>Searching the web…</span>';
        ui.bubble.replaceChildren(searchBlock);
      } else if (name === "search.source") {
        sources.push(d);
      } else if (name === "search.completed") {
        if (searchBlock) {
          ui.bubble.replaceChildren(App.searchSummary(sources));
          sources = [];
        }
      } else if (name === "message.delta") {
        if (!first) {
          first = Math.round(performance.now() - started);
          var t = ui.bubble.querySelector(".thinking");
          if (t) t.remove();
        }
        onDelta(ui, d);
      } else if (name === "message.completed") {
        onCompleted(ui, d);
      } else if (name === "message.failed") {
        onFailed(ui, d);
      }
    });
    $("foot").textContent = "POST messages · TTFT " + first + " ms · tổng " + Math.round(performance.now() - started) + " ms · " + rid;
  }
  async function send() {
    var text = $("input").value.trim();
    if (!text || !state.model || state.wanted !== state.active || state.streams[state.active]) return;
    if (!(await ensureConversation())) return;
    var convId = state.active;
    var col = $("view").querySelector(".col");
    var empty = col.querySelector(".empty");
    if (empty) empty.remove();
    $("input").value = "";
    $("input").style.height = "auto";
    if (App.closeSidebar) App.closeSidebar();
    App.addMessage(col, { role: "user", content: text, status: "completed" });
    var ui = App.addMessage(col, { role: "assistant", content: "", status: "streaming", model: state.model });
    $("view").scrollTop = $("view").scrollHeight;
    await runStream("/v1/conversations/" + convId + "/messages", { content: text, model: state.model }, ui, convId);
  }
  async function continueTurn(convId, ui) {
    if (state.streams[convId]) return;
    App.removeContinue();
    ui.continuing = true;
    ui.mdContent = App.el("div", "md-content");
    ui.bubble.replaceChildren(ui.mdContent);
    App.renderMd(ui.mdContent, ui.mdBuf);
    $("view").scrollTop = $("view").scrollHeight;
    await runStream("/v1/conversations/" + convId + "/continue", {}, ui, convId);
  }
  async function runStream(url, body, ui, convId) {
    var ctrl = new AbortController();
    state.streams[convId] = ctrl;
    refreshComposer();
    try {
      await streamReply(url, body, ui, convId, ctrl.signal);
    } catch (e) {
      if (e.name !== "AbortError") App.showError({ status: 0, data: "mất kết nối khi stream (" + e.name + ")", rid: "" });
    } finally {
      delete state.streams[convId];
      state.views.delete(convId);
      refreshComposer();
      await App.loadList(false);
      await refreshHeader(convId);
    }
  }

  App.refreshComposer = refreshComposer;
  App.send = send;
  App.continueTurn = continueTurn;
  App.stopStream = function () { var c = state.streams[state.active]; if (c) c.abort(); };
})(App);
