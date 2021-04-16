const cacheName = "app-" + "5e5b9e78a09a765c948b1c51c506bead285e7dac";

self.addEventListener("install", event => {
  console.log("installing app worker 5e5b9e78a09a765c948b1c51c506bead285e7dac");
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
  console.log("app worker 5e5b9e78a09a765c948b1c51c506bead285e7dac is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
