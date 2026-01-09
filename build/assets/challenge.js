// MODERN=1 SKEW_EXTRA_OPTIONS="--define:decrypt.wasm.DEFAULT_WASM_URL='{{ASSETS}}/w.wasm' --define:WEBLIB_EXPORT_NAME='solve'" make all
if (!navigator.cookieEnabled) return error("Cookies are required to proceed.");

if (!this.crypto || !this.crypto.subtle) {
  return error(
    "Your browser is not supported. Please update to a modern browser."
  );
}

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

function error(msg) {
  const error = document.querySelector(".e");
  error.hidden = false;
  error.textContent = msg;
}
