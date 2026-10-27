(function (App) {
  var $ = App.$;
  var el = App.el;
  var state = App.state;

  function modelName(id) {
    var found = id;
    if (state.catalog) state.catalog.providers.forEach(function (p) { p.models.forEach(function (m) { if (m.id === id) found = m.name; }); });
    return found;
  }
  function pickModel(id) {
    state.model = id;
    try { localStorage.setItem(App.MODEL_KEY, id); } catch (e) { void e; }
    var url = new URL(location.href);
    url.searchParams.set("model", id);
    history.replaceState({}, "", url);
    $("modelBtn").textContent = "Model: " + modelName(id) + " ▾";
    $("picker").hidden = true;
  }
  function renderPicker() {
    var box = $("picker");
    box.replaceChildren();
    var tabs = el("div", "tabs");
    state.catalog.providers.forEach(function (p) {
      var b = el("button", "btn small tab" + (p.id === state.tab ? " on" : ""), p.name);
      b.type = "button";
      b.addEventListener("click", function () { state.tab = p.id; renderPicker(); });
      tabs.append(b);
    });
    box.append(tabs);
    state.catalog.providers.filter(function (p) { return p.id === state.tab; }).forEach(function (p) {
      p.models.forEach(function (m) {
        var b = el("button", "opt");
        b.type = "button";
        b.disabled = !m.available;
        b.append(el("span", "", (m.id === state.model ? "✓ " : "") + m.name + (m.is_new ? " · new" : "")), el("small", "", m.available ? m.description : "Chưa cấu hình API key"));
        b.addEventListener("click", function () { pickModel(m.id); });
        box.append(b);
      });
    });
  }
  async function loadCatalog() {
    var cat = await App.api("GET", "/v1/models");
    if (!cat) return;
    state.catalog = cat;
    var wanted = new URL(location.href).searchParams.get("model");
    try { wanted = wanted || localStorage.getItem(App.MODEL_KEY); } catch (e) { void e; }
    var ok = {};
    cat.providers.forEach(function (p) { p.models.forEach(function (m) { if (m.available) ok[m.id] = p.id; }); });
    var id = wanted && ok[wanted] ? wanted : cat.default_model;
    state.tab = ok[id] || cat.providers[0].id;
    pickModel(id);
  }

  App.modelName = modelName;
  App.pickModel = pickModel;
  App.renderPicker = renderPicker;
  App.loadCatalog = loadCatalog;
})(App);
