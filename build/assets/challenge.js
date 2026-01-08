if (!navigator.cookieEnabled) return alert("Cookies are required to proceed.");
// MODERN=1 SKEW_EXTRA_OPTIONS="--define:decrypt.wasm.DEFAULT_WASM_URL='{{ASSETS}}/w.wasm' --define:WEBLIB_EXPORT_NAME='solve'" make all
!(function (e) {
  function n(e) {
    for (
      var n = (4294967295 * Math.random()) | 0, t = 0, r = e[0];
      r > t;
      t = (t + 1) | 0
    )
      e.slice(1, 33)[t >> 3] ^= ((n >>> t) & 1) << (7 & t);
  }
  function t(e, f) {
    var l, d, v, y;
    "u" > typeof document && "u" > typeof Worker && s
      ? ((u = (u + 1) | 0),
        (l = c = (c + 1) | 0),
        null === a &&
          ((a = new Worker(
            "string" == typeof (y = i.src) && "" !== y
              ? y
              : "data:text/javascript;," + i.text
          )).onerror = function () {
            s = !1;
          }),
        (d = null),
        (v = null),
        (d = function (e) {
          var n = e.data;
          n.b ^ l ||
            (a.removeEventListener("message", d),
            a.removeEventListener("error", v),
            (u = (u - 1) | 0) || (null !== a && (a.terminate(), (a = null))),
            f(n.a));
        }),
        a.addEventListener("message", d),
        (v = function () {
          t(e, f);
        }),
        a.addEventListener("error", v),
        a.postMessage(new r(l, e, null)))
      : "u" > typeof WebAssembly
      ? (null === o &&
          (o = WebAssembly.instantiateStreaming(fetch("{{ASSETS}}/w.wasm"))),
        o).then(function (t) {
          !(function (e, t, r) {
            n(e);
            for (
              var o = e.length,
                i = r.instance.exports,
                a = new Uint8Array(i.memory.buffer),
                u = 0;
              u < o;
              u = (u + 1) | 0
            )
              a[u] = e[u];
            i.decrypt(0, o, o + 1);
            for (
              var c = e.length - 61, s = new Uint8Array(c), f = 0;
              c > f;
              f = (f + 1) | 0
            )
              s[f] = a[(1 + o + f) | 0];
            t(s);
          })(e, f, t);
        })
      : (n(e),
        (function e(n, t, r) {
          crypto.subtle
            .importKey("raw", n.slice(1, 33), { name: "AES-GCM" }, !1, [
              "decrypt",
            ])
            .then(function (e) {
              return crypto.subtle.decrypt(
                { name: "AES-GCM", iv: n.slice(33, 45) },
                e,
                n.slice(45)
              );
            })
            .then(t)
            .catch(function () {
              (function (e) {
                for (
                  var n = 0, t = 0, r = e[0];
                  r > t &&
                  ((n = 1 << (7 & t)), !((e[(1 + (t >> 3)) | 0] ^= n) & n));
                  t = (t + 1) | 0
                );
              })(n),
                e(n, t, (r = (r + 1) | 0));
            });
        })(e, f, 0));
  }
  function r(e, n, t) {
    (this.b = e), (this.c = n), (this.a = t);
  }
  e.solve = function (e, n) {
    t(
      (function (e) {
        for (var n = new Uint8Array(128), t = 0; t < 64; t = (t + 1) | 0)
          n[t < 26 ? t + 65 : t < 52 ? t + 71 : t < 62 ? t - 4 : 4 * t - 205] =
            t;
        for (
          var r = (function (e) {
              for (var n = [], t = 0, r = e.length; r > t; t = (t + 1) | 0)
                n.push(e.charCodeAt(t));
              return n;
            })(e),
            o = e.length,
            i = new Uint8Array(
              (3 * (o - (61 === r[(o - 1) | 0]) - (61 === r[(o - 2) | 0]))) / 4
            ),
            a = 0,
            u = 0;
          o > a;

        ) {
          var c = n[r[a++]],
            s = n[r[a++]],
            f = n[r[a++]],
            l = n[r[a++]];
          (i[u++] = (c << 2) | (s >> 4)),
            (i[u++] = ((15 & s) << 4) | (f >> 2)),
            (i[u++] = ((3 & f) << 6) | l);
        }
        return i;
      })(e),
      function (e) {
        n(new TextDecoder().decode(new Uint8Array(e)));
      }
    );
  };
  var o = null,
    i = "u" > typeof document ? document.currentScript : null,
    a = null,
    u = 0,
    c = 0,
    s = !0;
  "u" < typeof document &&
    addEventListener("message", function (e) {
      var n = e.data;
      t(n.c, function (e) {
        postMessage(new r(n.b, null, e));
      });
    });
})(this);

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
