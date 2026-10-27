(function (App) {
  var $ = App.$;
  var state = App.state;

  async function boot() {
    $("userId").textContent = state.user;
    state.items = [];
    state.next = null;
    await App.loadCatalog();
    await App.loadList(false);
    App.routeFromPath();
  }
  function switchUser(id) {
    state.user = id;
    App.saveUser(id);
    history.pushState({}, "", "/");
    boot();
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
  $("newUser").addEventListener("click", function () { switchUser(App.newUUID()); });
  $("useUser").addEventListener("click", function () {
    var id = $("pasteUser").value.trim().toLowerCase();
    if (!App.UUID.test(id)) { App.showError({ status: 400, data: "user id không hợp lệ", rid: "" }); return; }
    $("pasteUser").value = "";
    switchUser(id);
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
