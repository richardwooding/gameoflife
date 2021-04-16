const cacheName = "app-" + "8262fc875b5a76de78f0df26f8221ea45025c3b1";

self.addEventListener("install", event => {
  console.log("installing app worker 8262fc875b5a76de78f0df26f8221ea45025c3b1");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/gameoflife",
        "/gameoflife/app.css",
        "/gameoflife/app.js",
        "/gameoflife/manifest.webmanifest",
        "/gameoflife/wasm_exec.js",
        "/gameoflife/web/app.css",
        "/gameoflife/web/app.wasm",
        "https://storage.googleapis.com/murlok-github/icon-192.png",
        "https://storage.googleapis.com/murlok-github/icon-512.png",
        
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker 8262fc875b5a76de78f0df26f8221ea45025c3b1 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
