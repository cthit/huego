# huego

copy .env.example to .env

create a huego cilent in gamma at `http://localhost:3000` to get `GAMMA_SECRET` and `GAMMA_CLIENT_ID` 

to dev against the lights in Hubben you also need the `HUE_BASE_URL` and to be connected to digit network in Hubben

### run frontend with 

`docker-compose up`

runs at `http://localhost:3001`

### run backend

`cd backend`

`go run cmd/huego/main.go`