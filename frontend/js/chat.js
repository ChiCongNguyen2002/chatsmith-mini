(function (App) {
  var $ = App.$;
  var state = App.state;

  function setBusy(on) {
    state.busy = on;
    $("send").hidden = on;
    $("stop").hidden = !on;
    $("input").disabled = on;
    if (!on) $("input").focus();
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
    state.active = c.id;
    history.pushState({}, "", "/conversation/" + c.id + location.search);
    App.renderConversation(c, []);
    return true;
  }
  function onDelta(ui, d) {
    ui.bubble.textContent += d.text;
    $("view").scrollTop = $("view").scrollHeight;
  }
  function onCompleted(ui, d) {
    ui.chip.textContent = d.model + (d.status !== "completed" ? " · " + d.status : "") + " · " + App.usd(d.cost_micro_usd) + " · TTFT " + d.first_token_ms + " ms · " + d.latency_ms + " ms";
    App.markStatus(ui, d.status);
  }
  function onFailed(ui, d) {
    App.markStatus(ui, "failed");
    ui.chip.textContent = ui.chip.textContent + " " + d.code + " · request_id " + d.request_id;
  }
  async function streamReply(text, ui, convId) {
    var started = performance.now();
    var first = 0;
    var res = await fetch("/v1/conversations/" + convId + "/messages", {
      method: "POST",
      headers: { "X-User-ID": state.user, "Content-Type": "application/json" },
      body: JSON.stringify({ content: text, model: state.model }),
      signal: state.abort.signal
    });
    var rid = res.headers.get("X-Request-Id") || "";
    if (!res.ok) {
      var j = null;
      try { j = await res.json(); } catch (e) { void e; }
      App.showError({ status: res.status, data: j && j.data, rid: rid });
      ui.row.remove();
      return;
    }
    App.showError(null);
    await readStream(res, function (name, d) {
      if (name === "message.delta") {
        if (!first) first = Math.round(performance.now() - started);
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
    if (!text || state.busy || !state.model) return;
    if (!(await ensureConversation())) return;
    var convId = state.active;
    var col = $("view").querySelector(".col");
    var empty = col.querySelector(".empty");
    if (empty) empty.remove();
    $("input").value = "";
    App.addMessage(col, { role: "user", content: text, status: "completed" });
    var ui = App.addMessage(col, { role: "assistant", content: "", status: "streaming", model: state.model });
    $("view").scrollTop = $("view").scrollHeight;
    state.abort = new AbortController();
    setBusy(true);
    try {
      await streamReply(text, ui, convId);
    } catch (e) {
      if (e.name !== "AbortError") App.showError({ status: 0, data: "mất kết nối khi stream (" + e.name + ")", rid: "" });
    } finally {
      state.abort = null;
      setBusy(false);
      await App.loadList(false);
      await refreshHeader(convId);
    }
  }

  App.setBusy = setBusy;
  App.send = send;
})(App);
