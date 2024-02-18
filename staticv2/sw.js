self.addEventListener("install", (event) => {
    event.waitUntil(
        caches
            .open("v1")
            .then((cache) =>
                cache.addAll([
                    "/files",
                    "/me",
                    "/notes",
                    "/"
                ]),
            ),
    );
});

self.addEventListener("fetch", (event) => {
    event.respondWith(
        // caches.match() always resolves
        // but in case of success response will have value
        caches.match(event.request).then((response) => {
            if(navigator.onLine) {
                return fetch(event.request)
                .then((response) => {
                    let responseClone = response.clone();
                    caches.open("v1").then((cache) => {
                        cache.put(event.request, responseClone);
                    });
                    return response;
                })
            }

            if (response !== undefined) {
                return response;
            }
        }),
    );
});