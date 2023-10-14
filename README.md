#Usage:
#make sure docker & docker compose is installed 
#docker compose up -d --build
#docker exec -it crazybird proxychains /app/crazybird [thread number]
#Ex: ./crazybird 1000
#If socks are not working, you can run the script without proxychains(but not recommended):
#docker exec -it crazybird ./crazybird [thread number]
