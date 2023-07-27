cd ./staticv2;
npm run build;

cd ..;
docker compose down -v ;
docker compose build --no-cache; 
docker compose --env-file ./config.env up -d;