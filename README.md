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
## Q&A
**Q:** How should we run your solution?

**A:** See the `Run` section above.
##
**Q:** How long did you spend on the take home? What would you add to your solution if you
had more time and expected it to be used in a production setting?

**A:** ~10 hours because I was looking through the Prometheus Go client library for inspiration. See the `TODOs` section above.
##
**Q:** If you used any libraries not in the language’s standard library, why did you use them?

**A:** I used github.com/gorilla/mux and github.com/stretchr/testify. I used the HTTP router and URL matcher to get automatic `404 page not found` error handling - I could've used `net/http` instead but didn't realize until writing this answer. And I used `testify` as a force of habit.
##
**Q:** If you have any feedback, feel free to share your thoughts!

**A:** Awesome challenge!
##
