# Candlelight-Backend

This should be the main source of the backend from now on. Branch structure is as follows:

### main
The branch that is currently deployed to production

### dev
The "staging" branch that work branches come from. The latest (and possibly unstable) code is found here

### Work Item Branches
Individual branches are created for each Issue on the board to isolate features for testing purposes. If you want to look at a specific piece of functionality, go to the branch for that Issue

To run the code, you can simply run `docker compose build` then `docker compose up` from the root directory, or if you want to run it manually, ensure that redis is running on your machine, then run `go run ./candlelight-api` from the root directory
