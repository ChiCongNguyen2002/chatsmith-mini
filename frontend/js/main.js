(function (App) {
  var $ = App.$;
  var state = App.state;

  App.MODEL_KEY = "chatsmith_model";

  async function loadMe() {
    var me = await App.api("GET", "/v1/me");
    if (me) {
      $("userName").textContent = me.name || me.email;
      App.showApp();
    }
    return me;
  }
  async function boot() {
    state.items = [];
    state.next = null;
    if (!(await loadMe())) return;
    await App.loadCatalog();
    await App.loadList(false);
    App.routeFromPath();
  }

  function closeSidebar() {
    $("side").classList.remove("open");
    var overlay = document.querySelector(".sidebar-overlay");
    if (overlay) overlay.remove();
  }
  function openSidebar() {
    $("side").classList.add("open");
    var overlay = App.el("div", "sidebar-overlay");
    overlay.addEventListener("click", closeSidebar);
    document.body.append(overlay);
  }

  function closePicker() {
    $("picker").hidden = true;
    var overlay = document.querySelector(".picker-overlay");
    if (overlay) overlay.remove();
  }

  function autoResize(textarea) {
    textarea.style.height = "auto";
    textarea.style.height = Math.min(textarea.scrollHeight, 160) + "px";
  }

  $("newChat").addEventListener("click", function () {
    history.pushState({}, "", "/" + location.search);
    App.clearActive();
    App.renderHello();
    App.renderList();
    closeSidebar();
    $("input").focus();
  });
  $("more").addEventListener("click", function () { App.loadList(true); });
  $("logout").addEventListener("click", async function () {
    await fetch("/auth/logout", { method: "POST" });
    location.href = "/";
  });
  $("hamburger").addEventListener("click", function () {
    if ($("side").classList.contains("open")) closeSidebar();
    else openSidebar();
  });
  $("modelBtn").addEventListener("click", function () {
    var box = $("picker");
    if (!state.catalog) return;
    if (box.hidden) {
      App.renderPicker();
      box.hidden = false;
      var overlay = App.el("div", "picker-overlay");
      overlay.addEventListener("click", closePicker);
      box.parentElement.append(overlay);
    } else {
      closePicker();
    }
  });
  $("send").addEventListener("click", App.send);
  $("stop").addEventListener("click", App.stopStream);
  $("input").addEventListener("keydown", function (e) {
    if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); App.send(); }
  });
  $("input").addEventListener("input", function () { autoResize(this); });
  window.addEventListener("popstate", App.routeFromPath);

  App.closeSidebar = closeSidebar;
  App.closePicker = closePicker;
  boot();
})(App);
