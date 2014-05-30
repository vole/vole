define(function(require) {

  function Memory() {
    this.size = 0;
    this._items = {};
  }

  Memory.prototype.get = function(key) {
    return this._items[key];
  };

  Memory.prototype.set = function(key, value) {
    if (!this._items.hasOwnProperty(key)) {
      this.size++;
    }

    this._items[key] = value;
  };

  Memory.prototype.del = function(key) {
    if (this._items.hasOwnProperty(key)) {
      this.size--;
    }

    delete this._items[key];
  };

  return {
    Memory: Memory
  };

});
