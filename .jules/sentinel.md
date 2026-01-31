## 2024-05-23 - Stored XSS via File Upload
**Vulnerability:** The application allowed uploading arbitrary files, including HTML and SVG, and served them inline. This allowed Stored XSS as the browser would execute scripts embedded in these files when viewed.
**Learning:** Even if the filename on disk is safe (ID-based), the original filename extension determines how the browser handles the file. Serving user content inline requires strict type validation or Content-Security-Policy.
**Prevention:** Force Content-Disposition: attachment for dangerous file types to prevent browser execution.
