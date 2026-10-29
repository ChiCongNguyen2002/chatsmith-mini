var App = (window.App = {});

(function (App) {
  var PAGE = 20;
  var MSG_PAGE = 200;
  var VIEW_LIMIT = 10;
  var TRUNCATED_NOTICE = "Câu trả lời bị cắt vì giới hạn độ dài";
  var $ = function (id) { return document.getElementById(id); };

  marked.setOptions({ breaks: true, gfm: true });
  function renderMd(bubble, text) {
    bubble.innerHTML = marked.parse(text);
  }

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
      var actions = el("span", "item-actions");
      var ren = el("span", "item-rename", "✎");
      ren.addEventListener("click", function (e) {
        e.stopPropagation();
        startRename(c, b, title);
      });
      var del = el("span", "item-delete", "×");
      del.addEventListener("click", function (e) {
        e.stopPropagation();
        confirmDelete(c.id, c.title || "Untitled");
      });
      actions.append(ren, del);
      b.append(title, actions);
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
    App.state.views.delete(id);
    if (App.state.wanted === id) {
      clearActive();
      history.pushState({}, "", "/" + location.search);
      renderHello();
    }
    renderList();
  }
  function startRename(conv, btn, titleSpan) {
    var input = el("input", "rename-input");
    input.type = "text";
    input.value = conv.title || "";
    input.maxLength = 200;
    titleSpan.replaceWith(input);
    input.focus();
    input.select();
    var done = false;
    function finish() {
      if (done) return;
      done = true;
      var val = input.value.trim();
      if (val && val !== conv.title) {
        api("PATCH", "/v1/conversations/" + conv.id, { title: val }).then(function (res) {
          if (res) {
            conv.title = val;
            btn.classList.remove("untitled");
          }
          renderList();
        });
      } else {
        renderList();
      }
    }
    input.addEventListener("blur", finish);
    input.addEventListener("keydown", function (e) {
      if (e.key === "Enter") { e.preventDefault(); finish(); }
      if (e.key === "Escape") { done = true; renderList(); }
    });
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
    if (status !== "truncated" && ui.notice) {
      ui.notice.remove();
      ui.notice = null;
    }
  }
  function removeContinue() {
    $("view").querySelectorAll(".continue-btn").forEach(function (b) { b.remove(); });
  }
  function showContinue(ui, convId) {
    ui.row.parentNode.querySelectorAll(".continue-btn").forEach(function (b) { b.remove(); });
    var btn = el("button", "btn small continue-btn", "Tiếp tục");
    btn.type = "button";
    btn.addEventListener("click", function () { App.continueTurn(convId, ui); });
    ui.row.insertBefore(btn, ui.chip);
  }
  function searchSummary(sources) {
    var block = el("div", "search-block done");
    var list = el("div", "search-sources");
    list.hidden = true;
    var toggle = el("span", "search-toggle", " ▸");
    sources.forEach(function (s) {
      var a = el("a", "search-link", s.title || s.url);
      if (/^https?:\/\//.test(s.url)) { a.href = s.url; a.target = "_blank"; a.rel = "noopener"; }
      a.addEventListener("click", function (e) { e.stopPropagation(); });
      list.append(a);
    });
    block.append(el("span", "search-icon", "✓"), el("span", "", "Searched " + sources.length + " sources"), toggle, list);
    block.addEventListener("click", function () {
      list.hidden = !list.hidden;
      toggle.textContent = list.hidden ? " ▸" : " ▾";
    });
    return block;
  }
  function addMessage(col, m) {
    var row = el("div", "msg " + m.role);
    var bubble = el("div", "bubble");
    if (m.role === "assistant" && m.content) {
      renderMd(bubble, m.content);
    } else {
      bubble.textContent = m.content || "";
    }
    if (m.sources && m.sources.length) bubble.prepend(searchSummary(m.sources));
    var chip = el("div", "chip", chipText(m));
    var ui = { row: row, bubble: bubble, chip: chip, notice: null, mdBuf: m.content || "", mdContent: null, model: m.model };
    row.append(bubble, chip);
    col.append(row);
    markStatus(ui, m.status);
    return ui;
  }
  function buildEntry(c, msgs) {
    var col = el("div", "col");
    col.append(el("h2", "title", c.title || "Untitled"), el("div", "meta", "Cập nhật " + new Date(c.updated_at).toLocaleString() + " · " + msgs.length + " tin"));
    if (msgs.length === 0) col.append(el("p", "empty", "Cuộc trò chuyện chưa có tin nhắn."));
    var last = null;
    msgs.forEach(function (m) { last = addMessage(col, m); });
    if (last && msgs[msgs.length - 1].status === "truncated") showContinue(last, c.id);
    return { c: c, msgs: msgs, col: col };
  }
  function showEntry(entry) {
    var view = $("view");
    view.replaceChildren(entry.col);
    view.scrollTop = view.scrollHeight;
  }
  function remember(entry) {
    var views = App.state.views;
    views.delete(entry.c.id);
    views.set(entry.c.id, entry);
    if (views.size <= VIEW_LIMIT) return;
    var evict = Array.from(views.keys()).find(function (key) { return !App.state.streams[key] && key !== App.state.active; });
    if (evict) views.delete(evict);
  }
  function renderConversation(c, msgs) {
    var entry = buildEntry(c, msgs);
    remember(entry);
    showEntry(entry);
  }
  function closeConversation() {
    clearActive();
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
  function showLoading() {
    var col = el("div", "col");
    col.append(el("p", "empty", "Đang tải cuộc trò chuyện…"));
    $("view").replaceChildren(col);
  }
  async function fetchConversation(id) {
    var c = await api("GET", "/v1/conversations/" + id);
    if (!c) return null;
    var msgs = await loadMessages(id);
    return msgs ? { c: c, msgs: msgs } : null;
  }
  function sameData(a, b) {
    return a.c.title === b.c.title && JSON.stringify(a.msgs) === JSON.stringify(b.msgs);
  }
  async function openConversation(id, push) {
    App.state.wanted = id;
    if (push) history.pushState({}, "", "/conversation/" + id + location.search);
    var cached = App.state.views.get(id);
    App.state.active = cached ? id : null;
    if (cached) {
      showEntry(cached);
      remember(cached);
    } else {
      showLoading();
    }
    renderList();
    App.refreshComposer();
    var fresh = await fetchConversation(id);
    if (App.state.wanted !== id) return;
    if (!fresh) { if (!cached) closeConversation(); return; }
    App.state.active = id;
    var keep = cached && (App.state.streams[id] || sameData(cached, fresh));
    if (!keep) renderConversation(fresh.c, fresh.msgs);
    renderList();
    App.refreshComposer();
  }
  function clearActive() {
    App.state.wanted = null;
    App.state.active = null;
    App.refreshComposer();
  }
  function routeFromPath() {
    var m = location.pathname.match(/^\/conversation\/([0-9a-f-]{36})$/);
    if (m) openConversation(m[1], false); else { clearActive(); renderHello(); renderList(); }
  }

  App.state = { items: [], next: null, active: null, wanted: null, views: new Map(), streams: {}, catalog: null, model: "", tab: "" };
  App.clearActive = clearActive;
  App.searchSummary = searchSummary;
  App.$ = $;
  App.el = el;
  App.showLogin = showLogin;
  App.showApp = showApp;
  App.showError = showError;
  App.api = api;
  App.usd = usd;
  App.renderMd = renderMd;
  App.showContinue = showContinue;
  App.removeContinue = removeContinue;
  App.chipText = chipText;
  App.markStatus = markStatus;
  App.addMessage = addMessage;
  App.renderConversation = renderConversation;
  App.renderHello = renderHello;
  App.renderList = renderList;
  App.loadList = loadList;
  App.routeFromPath = routeFromPath;
})(App);
