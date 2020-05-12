function TwitchCollector(interval = 2e4) {
  const _hooks = {
    update: []
  };
  let _data = null;

  async function _getData() {
    const response = await fetch('/data');
    const data = await response.json();
    _update(data);
  }

  function _update(data) {
    _data = data;

    for (let hook of _hooks.update) {
      hook(data);
    }
  }

  function registerHooks(kind, callback) {
    if (kind in _hooks) {
      _hooks.kind.push(callback);
    }
  }

  function data() {
    return _data;
  }

  async function start() {
    await _getData();
    setInterval(_getData, interval);
  }

  return {
    data,
    registerHooks,
    start
  };
}
