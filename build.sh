cd ./staticv2
npm run build
cp ./sw.js ../static

cd ..
docker compose down -v;
docker compose build --no-cache;
docker compose --env-file ./config.env up -d;
docker logs -f app-gosf;