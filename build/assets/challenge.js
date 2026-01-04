if (!navigator.cookieEnabled) return alert("Cookies are required to proceed.");
function t(t, n, r, e) {
  (this.c = t), (this.a = n), (this.b = r), (this.d = e), (this.f = 0);
}
this.solve = function (n, r) {
  var e,
    c = new t(
      0 |
        (e = (function (t) {
          for (var n = new Uint8Array(128), r = 0; r < 64; r = (r + 1) | 0)
            n[
              r < 26 ? r + 65 : r < 52 ? r + 71 : r < 62 ? r - 4 : 4 * r - 205
            ] = r;
          for (
            var e = (function (t) {
                for (var n = [], r = 0, e = t.length; e > r; r = (r + 1) | 0)
                  n.push(t.charCodeAt(r));
                return n;
              })(t),
              c = t.length,
              i = new Uint8Array(
                (3 * (c - (61 === e[(c - 1) | 0]) - (61 === e[(c - 2) | 0]))) /
                  4
              ),
              o = 0,
              a = 0;
            o < c;

          ) {
            var f = n[e[o++]],
              u = n[e[o++]],
              h = n[e[o++]],
              s = n[e[o++]];
            (i[a++] = (f << 2) | (u >> 4)),
              (i[a++] = ((15 & u) << 4) | (h >> 2)),
              (i[a++] = ((3 & h) << 6) | s);
          }
          return i;
        })(n))[0],
      e.slice(1, 33),
      e.slice(33, 45),
      e.slice(45)
    );
  (function (t) {
    for (
      var n = (4294967295 * Math.random()) | 0, r = 0, e = t.c;
      e > r;
      r = (r + 1) | 0
    )
      t.a[r >> 3] ^= ((n >>> r) & 1) << (7 & r);
  })(c),
    (function t(n, r) {
      crypto.subtle
        .importKey("raw", n.a, { name: "AES-GCM" }, !1, ["decrypt"])
        .then(function (t) {
          return crypto.subtle.decrypt({ name: "AES-GCM", iv: n.b }, t, n.d);
        })
        .then(function (t) {
          var n = new TextDecoder("utf-8"),
            e = new Uint8Array(t);
          r(n.decode(e));
        })
        .catch(function (e) {
          (function (t) {
            for (
              var n = 0, r = 0, e = t.c;
              e > r && ((n = 1 << (7 & r)), !((t.a[r >> 3] ^= n) & n));
              r = (r + 1) | 0
            );
          })(n),
            8191 & (n.f = (n.f + 1) | 0)
              ? t(n, r)
              : setTimeout(function () {
                  t(n, r);
                }, 1);
        });
    })(c, r);
};
this.solve("{{CHALLENGE}}", function (data) {
  document.cookie =
    "{{COOKIE_NAME}}=" +
    data +
    ";path=/;max-age={{PASSAGE_DURATION}};SameSite=Lax";
  document.body.setAttribute("v", "");
  // Dispatch a custom event to signal that the challenge has been solved
  // This can be useful for other scripts that may want to listen for this event
  // and perform actions accordingly
  // For example, analytics scripts or UI updates injected through the HEAD marker
  document.body.dispatchEvent(new Event("solved"));
  setTimeout(function () {
    location.reload();
  }, 500);
});
