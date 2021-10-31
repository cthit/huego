# huego

copy .env.exapmle to .env

get `GAMMA_SECRET` and `GAMMA_CLIENT_ID` from gamma at `http://localhost:3000`
to dev against the lights in hubben you also need the `HUE_BASE_URL` and to be connected to digit network in Hubben

### run frontend with 

`docker-compose up`

runs at `http://localhost:3001`

### run backend

`cd backend`

`go run cmd/huego/main.go`