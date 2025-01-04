_(When this was being worked on, it was hosted in Gitlab, where I had set up a CI/CD pipeline to manage testing/deployment of the code. This was written for anyone else on the team looking to access my deployment automation. Original text of the wiki page has been edited to remove sensitive info such as IP addresses. Alterations have been surrounded by double-underscores)_

Everything needed for deploying the backend is completely automated via bash scripts and the Gitlab Runner CI functionality. However, since we don't have access to the normal Gitlab runners, I (Ryan) have one locally hosted on my desktop. This means that pipelines will only get picked up/complete if I have it running on my machine. However, since I'm the only person that works in the backend, this hasn't been an issue thus far.

Typically, (again, assuming the Runner is running) any checkin to this Repo will kick off the CI Pipeline. Every branch will automatically run the test suite, and if the commit was to `main`, it will also build a docker image from the most recent code, push said image to Docker Hub, SSH into the remote EC2 instance, and trigger the redeploy script I've saved on that machine - that script will simply stop the currently running container on the machine, pull down the most recent Docker image for the backend (which should be the one you just pushed), and restart it using the new image.

However, if for some reason you need to manually deploy the code, do the following (after running the Go test suite to ensure it all works, of course)
1. Run `docker login -u candlelightdevteam -p [__CLI PASSWORD HERE__]`
2. Run `docker build --no-cache -t candlelightdevteam/candlelight-backend:latest .`
3. Run `docker push candlelightdevteam/candlelight-backend:latest`
4. Run `ssh [__IP of EC2 instance__] -i ~/.ssh/candlelight-dev-login.pem -l ec2-user 'sudo bash ./candlelight/candlelight-backend/rebuildbackend.sh'` - This requires you to have the candlelight-dev-login.pem keyfile locally downloaded at the location `~/.ssh/candlelight-dev-login.pem`. If you need this file, reach out to Ryan. 
    1. This step can also be done manually if preferred. Just SSH into the EC2 instance and do the following:
    2. Run `sudo docker ps` and find the container ID of "candlelight-backend-web-1"
    3. Run `sudo docker stop [WEB CONTAINER ID HERE] && sudo docker rm [WEB CONTAINER ID HERE]`
    4. Run `sudo docker images` and find the Image ID of "candlelightdevteam/candlelight-backend"
    5. Run `sudo docker rmi [IMAGE ID HERE]`
    6. Run `cd candlelight/candlelight-backend`
    7. Run `docker compose up -d` - This will run the backend + the Redis database as 2 separate containers in a Docker Compose Stack

Once those steps are done, the backend should be deployed to the EC2 instance. The Database should also persist as long as you don't remove the container its running in
