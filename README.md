# Backend Engineer Take-home Challenge

This is a solution for Oden Technologies's Cable Diameter Challenge.

The solution calculates a one minute Simple Moving Average(SMA) of the [cable diameter](http://takehome-backend.oden.network/?metric=cable-diameter) that is updated once per second. 

## TODOs
1. Handle errors - my (lack of) error handling had more to do with pragmatism than oversight
2. Add authentication and authorization
3. Improve/Add unit tests - I added just enough tests case to give insight into how I test code
4. Manage application secrets through environment variables - I like https://cloud.google.com/secret-manager
5. Document packages better. Ideally by adding a README.md into the subfolders
6. Add more metric types
7. Better abstraction for metric registration
8. Eventually CI/CD, IaC, collect application metrics, document endpoints(OpenAPI Specification)
9. Load IoT metrics into a time series database and visualize the data through a interactive application. 

## Run
```
docker-compose up
```

## Test
```
make test
```
