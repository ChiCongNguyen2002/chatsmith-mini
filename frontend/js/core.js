var App = (window.App = {});

(function (App) {
  var PAGE = 20;
  var MSG_PAGE = 200;
  var TRUNCATED_NOTICE = "Câu trả lời bị cắt vì giới hạn độ dài";
  var $ = function (id) { return document.getElementById(id); };

  function el(tag, cls, text) {
    var node = document.createElement(tag);
    if (cls) node.className = cls;
    if (text !== undefined) node.textContent = text;
    return node;
  }
  function showError(err) {
    var box = $("banner");
    if (!err) { box.hidden = true; return; }
    box.hidden = false;
    box.textContent = "Lỗi " + err.status + " " + (err.data || "") + (err.rid ? " · request_id " + err.rid : "");
  }
  async function api(method, path, body) {
    var started = performance.now();
    var res;
    try {
      res = await fetch(path, { method: method, headers: { "Content-Type": "application/json" }, body: body ? JSON.stringify(body) : undefined });
    } catch (e) {
      showError({ status: 0, data: "không kết nối được máy chủ (" + e.name + ")", rid: "" });
      return null;
    }
    var ms = Math.round(performance.now() - started);
    var rid = res.headers.get("X-Request-Id") || "";
    var json = null;
    try { json = await res.json(); } catch (e) { json = null; }
    $("foot").textContent = method + " " + path.split("?")[0] + " · " + ms + " ms · " + res.status + (rid ? " · " + rid : "");
    if (res.status === 401) {
      showLogin();
      return null;
    }
    if (!res.ok) {
      showError({ status: res.status, data: json && json.data, rid: rid });
      return null;
    }
    showError(null);
    return json ? json.data : true;
  }
  function showLogin() {
    $("side").hidden = true;
    document.querySelector(".topbar").hidden = true;
    document.querySelector(".composer").hidden = true;
    document.querySelector(".foot").hidden = true;
    var view = $("view");
    view.replaceChildren();
    var box = el("div", "hello");
    var link = el("a", "login-btn", "Sign in with Google");
    link.href = "/auth/google/login";
    box.append(el("p", "", "Welcome to"), el("h1", "", "Chat Smith Mini"), el("p", "subtitle", "Your AI assistant, powered by the latest models"), link);
    view.append(box);
  }
  function showApp() {
    $("side").hidden = false;
    $("userSection").hidden = false;
    document.querySelector(".topbar").hidden = false;
    document.querySelector(".composer").hidden = false;
    document.querySelector(".foot").hidden = false;
  }
  function groupOf(iso) {
    var d = new Date(iso);
    var now = new Date();
    var today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    var day = new Date(d.getFullYear(), d.getMonth(), d.getDate());
    var days = Math.round((today - day) / 86400000);
    if (days <= 0) return "Today";
    if (days === 1) return "Yesterday";
    if (days <= 7) return "Previous 7 days";
    if (days <= 30) return "Previous 30 days";
    return "Older";
  }
  function renderList() {
    var list = $("list");
    list.replaceChildren();
    var last = "";
    App.state.items.forEach(function (c) {
      var g = groupOf(c.updated_at);
      if (g !== last) { list.append(el("div", "group", g)); last = g; }
      var b = el("button", "item" + (c.id === App.state.active ? " active" : "") + (c.title ? "" : " untitled"));
      b.type = "button";
      var title = el("span", "item-title", c.title || "Untitled");
      var del = el("span", "item-delete", "×");
      del.addEventListener("click", function (e) {
        e.stopPropagation();
        confirmDelete(c.id, c.title || "Untitled");
      });
      b.append(title, del);
      b.addEventListener("click", function () { openConversation(c.id, true); });
      list.append(b);
    });
    $("more").hidden = !App.state.next;
  }
  function confirmDelete(id, title) {
    var overlay = el("div", "confirm-overlay");
    var box = el("div", "confirm-box");
    box.append(el("p", "", "Delete \"" + title + "\"?"));
    var actions = el("div", "confirm-actions");
    var cancel = el("button", "btn small", "Cancel");
    cancel.type = "button";
    cancel.addEventListener("click", function () { overlay.remove(); });
    var confirm = el("button", "btn small btn-danger", "Delete");
    confirm.type = "button";
    confirm.addEventListener("click", async function () {
      overlay.remove();
      await deleteConversation(id);
    });
    actions.append(cancel, confirm);
    box.append(actions);
    overlay.append(box);
    document.body.append(overlay);
  }
  async function deleteConversation(id) {
    var res = await api("DELETE", "/v1/conversations/" + id);
    if (!res) return;
    App.state.items = App.state.items.filter(function (c) { return c.id !== id; });
    if (App.state.active === id) {
      App.state.active = null;
      history.pushState({}, "", "/" + location.search);
      renderHello();
    }
    renderList();
  }
  async function loadList(more) {
    var path = "/v1/conversations?limit=" + PAGE + (more && App.state.next ? "&cursor=" + encodeURIComponent(App.state.next) : "");
    var page = await api("GET", path);
    if (!page) return;
    App.state.items = more ? App.state.items.concat(page.items) : page.items;
    App.state.next = page.next_cursor;
    renderList();
  }
  function renderHello() {
    var view = $("view");
    view.replaceChildren();
    var box = el("div", "hello");
    box.append(el("p", "", "Welcome"), el("h1", "", "How can I help you today?"));
    view.append(box);
  }
  function usd(micro) { return "$" + (micro / 1e6).toFixed(micro < 1000 ? 6 : 4); }
  function chipText(m) {
    if (m.role !== "assistant" || !m.model) return "";
    var t = m.model;
    if (m.status !== "completed") t += " · " + m.status;
    if (m.cost_micro_usd !== null && m.cost_micro_usd !== undefined) t += " · " + usd(m.cost_micro_usd);
    if (m.first_token_ms) t += " · TTFT " + m.first_token_ms + " ms";
    if (m.latency_ms) t += " · " + m.latency_ms + " ms";
    return t;
  }
  function markStatus(ui, status) {
    ui.row.classList.remove("failed", "canceled", "truncated");
    if (status === "failed" || status === "canceled" || status === "truncated") ui.row.classList.add(status);
    if (status === "truncated" && !ui.notice) {
      ui.notice = el("div", "notice", TRUNCATED_NOTICE);
      ui.row.insertBefore(ui.notice, ui.chip);
    }
  }
  function addMessage(col, m) {
    var row = el("div", "msg " + m.role);
    var bubble = el("div", "bubble", m.content);
    var chip = el("div", "chip", chipText(m));
    var ui = { row: row, bubble: bubble, chip: chip, notice: null };
    row.append(bubble, chip);
    col.append(row);
    markStatus(ui, m.status);
    return ui;
  }
  function renderConversation(c, msgs) {
    var view = $("view");
    view.replaceChildren();
    var col = el("div", "col");
    col.append(el("h2", "title", c.title || "Untitled"), el("div", "meta", "Cập nhật " + new Date(c.updated_at).toLocaleString() + " · " + msgs.length + " tin"));
    if (msgs.length === 0) col.append(el("p", "empty", "Cuộc trò chuyện chưa có tin nhắn."));
    msgs.forEach(function (m) { addMessage(col, m); });
    view.append(col);
    view.scrollTop = view.scrollHeight;
  }
  function closeConversation() {
    App.state.active = null;
    renderHello();
    renderList();
  }
  async function loadMessages(id) {
    var msgs = [];
    var after = 0;
    for (;;) {
      var batch = await api("GET", "/v1/conversations/" + id + "/messages?after_seq=" + after + "&limit=" + MSG_PAGE);
      if (!batch) return null;
      msgs = msgs.concat(batch);
      if (batch.length < MSG_PAGE) return msgs;
      after = batch[batch.length - 1].seq;
    }
  }
  async function openConversation(id, push) {
    if (App.state.busy) return;
    var c = await api("GET", "/v1/conversations/" + id);
    var msgs = c && (await loadMessages(id));
    if (!msgs) { closeConversation(); return; }
    App.state.active = id;
    if (push) history.pushState({}, "", "/conversation/" + id + location.search);
    renderConversation(c, msgs);
    renderList();
  }
  function routeFromPath() {
    var m = location.pathname.match(/^\/conversation\/([0-9a-f-]{36})$/);
    if (m) openConversation(m[1], false); else { App.state.active = null; renderHello(); renderList(); }
  }

  App.state = { items: [], next: null, active: null, catalog: null, model: "", tab: "", busy: false, abort: null };
  App.$ = $;
  App.el = el;
  App.showLogin = showLogin;
  App.showApp = showApp;
  App.showError = showError;
  App.api = api;
  App.usd = usd;
  App.chipText = chipText;
  App.markStatus = markStatus;
  App.addMessage = addMessage;
  App.renderConversation = renderConversation;
  App.renderHello = renderHello;
  App.renderList = renderList;
  App.loadList = loadList;
  App.routeFromPath = routeFromPath;
})(App);
