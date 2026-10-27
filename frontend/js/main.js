(function (App) {
  var $ = App.$;
  var state = App.state;

  async function loadMe() {
    var me = await App.api("GET", "/v1/me");
    if (me) $("userName").textContent = me.name || me.email;
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

  $("newChat").addEventListener("click", function () {
    if (state.busy) return;
    history.pushState({}, "", "/" + location.search);
    state.active = null;
    App.renderHello();
    App.renderList();
    $("input").focus();
  });
  $("more").addEventListener("click", function () { App.loadList(true); });
  $("logout").addEventListener("click", async function () {
    await fetch("/auth/logout", { method: "POST" });
    location.href = "/";
  });
  $("modelBtn").addEventListener("click", function () {
    var box = $("picker");
    if (!state.catalog) return;
    if (box.hidden) App.renderPicker();
    box.hidden = !box.hidden;
  });
  $("send").addEventListener("click", App.send);
  $("stop").addEventListener("click", function () { if (state.abort) state.abort.abort(); });
  $("input").addEventListener("keydown", function (e) {
    if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); App.send(); }
  });
  window.addEventListener("popstate", App.routeFromPath);
  boot();
})(App);
