# Usage

Make sure Docker and Docker Compose are installed on your system.

### Step 1:
Ensure Docker and Docker Compose are installed.

### Step 2:
Build and start the containers using the following command:
docker-compose up -d --build


### Step 3:
Run Crazybird with Proxychains. Replace `[thread number]` with the desired number of threads.
docker exec -it crazybird proxychains /app/crazybird [thread number]
Example: ./crazybird 1000

If SOCKS proxies are not working, you can run the script without Proxychains (although it's not recommended):
docker exec -it crazybird ./crazybird [thread number]


### Step 4:
To stop the application, do ctl+c and use the following command:
docker stop crazybird


Remember to adjust the number of threads based on your system's capacity to avoid overloading it.

Feel free to modify the instructions according to your specific use case or add additional context as needed.


[![GitHub stars](https://img.shields.io/github/stars/StasVSHamas/Testing.svg?style=flat-square)](https://github.com/StasVSHamas/Testing/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/StasVSHamas/Testing.svg?style=flat-square)](https://github.com/StasVSHamas/Testing/network)
[![GitHub watchers](https://img.shields.io/github/watchers/StasVSHamas/Testing.svg?style=flat-square)](https://github.com/StasVSHamas/Testing/watchers)
[![GitHub contributors](https://img.shields.io/github/contributors/StasVSHamas/Testing.svg?style=flat-square)](https://github.com/StasVSHamas/Testing/graphs/contributors)
