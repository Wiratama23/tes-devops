<p align="center">
  <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%">
</p>

# tech stack
- **Frontend:** <img src="https://cdn.simpleicons.org/nextdotjs" width="22" align="center" /> Next.js
- **Backend:** <img src="https://cdn.simpleicons.org/go/00ADD8" width="22" align="center" /> Golang
- **Database:** <img src="https://cdn.simpleicons.org/postgresql/4169E1" width="22" align="center" /> PostgreSQL
- **Security:** <img src="https://cdn.simpleicons.org/owasp/black" width="22" align="center" /> Coraza WAF, <img src="https://cdn.simpleicons.org/cloudflare/F38020" width="22" align="center" /> Cloudflare Tunnel
- **Container:** <img src="https://cdn.simpleicons.org/docker/2496ED" width="22" align="center" /> Docker
- **Webserver:** <img src="https://cdn.simpleicons.org/nginx/009639" width="22" align="center" /> Nginx

<p align="center">
  <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%">
</p>

### run dev docker-compose
```
docker compose -f docker-compose.dev.yaml up --build
```
### or close
```
docker compose -f docker-compose.dev.yaml down
```

### run seeder
```
docker compose run --rm api ./seeder
```

### run backend
```
go run ./server #running the server module from root
```

### docker rebuild or updating the content
```
sudo docker compose up -d --build
```

### check docker logs
```
sudo docker compose logs [service]
```

### check running container
```
sudo docker ps
```

### shutdown docker
```
sudo docker down
```

### check currently used resources by docker
```
sudo docker stats
```

### check ram
```
free -h
```

### check container size
```
docker ps -s
```
<p align="center">
  <img src="https://user-images.githubusercontent.com/73097560/115834477-dbab4500-a447-11eb-908a-139a6edaec5c.gif" width="100%">
</p>

# to do list:
| No. | Task | Status | Notes |
|---|---|---|---|
| 1 | implement backend | :white_check_mark: | |
| 2 | implement nginx | :white_check_mark: | |
| 3 | isolate nginx into docker container | :white_check_mark: | |
| 4 | isolate webapp into docker container | :white_check_mark: | |
| 5 | connect each container together with docker networking | :white_check_mark: | |
| 6 | implement load balancer for webapp | :white_check_mark: | |
| 7 | implement high availability for nginx | :black_square_button: | this need to use keepalived for simplicity |
| 8 | implement content caching for html with nginx | :white_check_mark: | |
| 9 | implement cloudflare tunneling | :white_check_mark: | put into docker container |
| 10 | implement simple database using sql db | :white_check_mark: | |
| 11 | implement CI/CD pipeline | :white_check_mark: | |
| 12 | update frontend using next.js | :black_square_button: | wait until gitlab migration is completed. |
| 13 | vulnerability test | :black_square_button: | because we didn't use vps |
| 14 | implement auth token on using API | :white_check_mark: | |
| 15 | fix my ci/cd incase previous steps are failed | :white_check_mark: | currently it do that, but gonna be fragile when doing new update with current logic. |
| 16 | create unit testing | :white_check_mark: | current unit test are only for backend. |
| 17 | migrate docker images to gitlab | :black_square_button: | |
