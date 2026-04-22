# tech stack
- frontend: next.js
- backend: golang
- database: postgresql
- security: corazaWAF, cloudflare tunnel
- container: docker
- webserver: nginx

#run backend
go run ./server #running the server module from root

#docker rebuild or updating the content
sudo docker compose up -d --build

#check docker logs
sudo docker compose logs [service]

#check running container
sudo docker ps

#shutdown docker
sudo docker down

#check currently used resources by docker
sudo docker stats

#check ram
free -h

#check container size
docker ps -s

to do list:
1. implement backend |x|
2. implement nginx |x|
3. isolate nginx into docker container |x|
4. isolate webapp into docker container |x|
5. connect each container together with docker networking |x|
6. implement load balancer for webapp |x|
7. implement high availability for nginx | | (this need to use keepalived for simplicity)
8. implement content caching for html with nginx |x|
9. implement cloudflare tunneling |x| (put into docker container)
10. implement simple database using sql db |x|
11. implement CI/CD pipeline |x|
12. update frontend using next.js | |
13. vulnerability test (because we didn't use vps) | |
14. implement auth token on using API | |
15. fix my ci/cd incase previous steps are failed. | | (currently it do that, but gonna be fragile when doing new update with current logic.)
