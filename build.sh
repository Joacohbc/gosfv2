cd ./staticv2
npm run build
cp ./sw.js ../static

cd ..

cd ./src/
env GOOS=linux GOARCH=amd64 go build -o ../bin/app-gosf.bin
env GOOOS=windows GOARCH=amd64 go build -o ../bin/app-gosf.exe

cd ..
docker compose down -v;
docker compose build --no-cache;
docker compose --env-file ./config.env up -d;
docker logs -f app-gosf;