if (!navigator.cookieEnabled) return alert("Cookies are required to proceed.");
// MODERN=1 SKEW_EXTRA_OPTIONS="--define:decrypt.wasm.DEFAULT_WASM_URL='{{ASSETS}}/w.wasm' --define:WEBLIB_EXPORT_NAME='solve'" make all
!(function (n) {
  function e(n) {
    for (
      var e = (4294967295 * Math.random()) | 0, t = 0, r = n[0];
      r > t;
      t = (t + 1) | 0
    )
      n.slice(1, 33)[t >> 3] ^= ((e >>> t) & 1) << (7 & t);
  }
  function t(n, e, r) {
    crypto.subtle
      .importKey("raw", n.slice(1, 33), { name: "AES-GCM" }, !1, ["decrypt"])
      .then(function (e) {
        return crypto.subtle.decrypt(
          { name: "AES-GCM", iv: n.slice(33, 45) },
          e,
          n.slice(45)
        );
      })
      .then(e)
      .catch(function () {
        (function (n) {
          for (
            var e = 0, t = 0, r = n[0];
            r > t && ((e = 1 << (7 & t)), !((n[(1 + (t >> 3)) | 0] ^= e) & e));
            t = (t + 1) | 0
          );
        })(n),
          t(n, e, (r = (r + 1) | 0));
      });
  }
  function r(n, v) {
    var d, m, y, p;
    "u" > typeof document && "u" > typeof Worker && l
      ? ((s = (s + 1) | 0),
        (d = f = (f + 1) | 0),
        null === c &&
          ((c = new Worker(
            "string" == typeof (p = u.src) && "" !== p
              ? p
              : "data:text/javascript;," + u.text
          )).onerror = function () {
            l = !1;
          }),
        (m = null),
        (y = null),
        (m = function (n) {
          var e = n.data;
          e.b ^ d ||
            (c.removeEventListener("message", m),
            c.removeEventListener("error", y),
            (s = (s - 1) | 0) || (null !== c && (c.terminate(), (c = null))),
            v(e.a));
        }),
        c.addEventListener("message", m),
        (y = function () {
          r(n, v);
        }),
        c.addEventListener("error", y),
        c.postMessage(new o(d, n, null)))
      : "u" > typeof WebAssembly && a
      ? (null === i &&
          (i = WebAssembly.instantiateStreaming(fetch("{{ASSETS}}/w.wasm"))),
        i).then(
          function (t) {
            !(function (n, t, r) {
              e(n);
              for (
                var o = n.length,
                  a = r.instance.exports,
                  i = new Uint8Array(a.memory.buffer),
                  u = 0;
                u < o;
                u = (u + 1) | 0
              )
                i[u] = n[u];
              a.decrypt(0, o, o + 1);
              for (
                var c = n.length - 61, s = new Uint8Array(c), f = 0;
                c > f;
                f = (f + 1) | 0
              )
                s[f] = i[(1 + o + f) | 0];
              t(s);
            })(n, v, t);
          },
          function () {
            (a = !1), e(n), t(n, v, 0);
          }
        )
      : (e(n), t(n, v, 0));
  }
  function o(n, e, t) {
    (this.b = n), (this.c = e), (this.a = t);
  }
  n.solve = function (n, e) {
    r(
      (function (n) {
        for (var e = new Uint8Array(128), t = 0; t < 64; t = (t + 1) | 0)
          e[t < 26 ? t + 65 : t < 52 ? t + 71 : t < 62 ? t - 4 : 4 * t - 205] =
            t;
        for (
          var r = (function (n) {
              for (var e = [], t = 0, r = n.length; r > t; t = (t + 1) | 0)
                e.push(n.charCodeAt(t));
              return e;
            })(n),
            o = n.length,
            a = new Uint8Array(
              (3 * (o - (61 === r[(o - 1) | 0]) - (61 === r[(o - 2) | 0]))) / 4
            ),
            i = 0,
            u = 0;
          o > i;

        ) {
          var c = e[r[i++]],
            s = e[r[i++]],
            f = e[r[i++]],
            l = e[r[i++]];
          (a[u++] = (c << 2) | (s >> 4)),
            (a[u++] = ((15 & s) << 4) | (f >> 2)),
            (a[u++] = ((3 & f) << 6) | l);
        }
        return a;
      })(n),
      function (n) {
        e(new TextDecoder().decode(new Uint8Array(n)));
      }
    );
  };
  var a = !0,
    i = null,
    u = "u" > typeof document ? document.currentScript : null,
    c = null,
    s = 0,
    f = 0,
    l = !0;
  "u" < typeof document &&
    addEventListener("message", function (n) {
      var e = n.data;
      r(e.c, function (n) {
        postMessage(new o(e.b, null, n));
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
