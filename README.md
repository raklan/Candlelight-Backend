# Candlelight-Backend

The backend code for my Senior Capstone project, written over the course of 2024. A majority of the code found in this repository has been written solely by hand, by me. In Spring of 2024, [Zachary Gust](https://github.com/ZSGust) wrote most of the netcode, but over the course of the summer and fall, I ended up refactoring and rewriting most of it as needs arose.

During development, I also created a wiki outlining a bit more about how the code works. These wiki pages have been copied into the `wiki` folder with only alterations to fix links between pages

To run the code, you can simply run `docker compose build` then `docker compose up` from the root directory, or if you want to run it manually, ensure that redis is running on your machine, then run `go run ./candlelight-api` from the root directory
