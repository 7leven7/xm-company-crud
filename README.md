# xm-company-crud

I have made an effort to complete everything on the checklist.
To run the project, it will be sufficient to execute "docker compose up." 
You can find the routes in the "router.go" folder. 
If you want to run Litner, you can first enter the "xm-company-crud-app-1" container,
execute the following commands:
docker exec -id xm-company-crud-app-1 sh
golangci-lint run

In the docker-compose.yaml file, you can find the ports configuration. 
The app port is set to 8888, the PostgreSQL port is set to 1234, and Kafka is using the default port 9092. 
