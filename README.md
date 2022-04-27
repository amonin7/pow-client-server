# pow-client-server

## 0. Objective of this project

Design and implement _“Word of Wisdom”_ tcp server.                 
- TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.                 
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from _“word of wisdom”_ book or any other collection of the quotes.                 
- Docker file should be provided both for the server and for the client that solves the POW challenge.

## 1. Installation details

1. Clone this project
```shell
git clone git@github.com:amonin7/pow-client-server.git
```

2. Run the Docker on your machine
3. Go into the directory you downloaded the project
```shell
cd pow-client-server
```
4. Run the docker compose build command
```shell
docker-compose up --build
```
5. If everything was built and launched correctly you should see something like below in the logs:
```shell
server_1  | Starting tcp server
client_1  | Starting tcp client
server_1  | Successfully created the server - message listener on url: 0.0.0.0:3333
server_1  | Started listening to address{[::]:3333}
server_1  | Established new client: 172.19.0.3:46414
client_1  | Established connection to server{server:3333}

server_1  | received request for challenge from 172.19.0.3:46414
client_1  | received challenge from the server: ISRM{quadResidue=1142927, modulo=6552149}
client_1  | proof found:  2774
client_1  | challenge was successfully sent to server
server_1  | received request for resource from 172.19.0.3:46414, checking if proof is correct...
server_1  | proof, received from client is correct 172.19.0.3:46414
client_1  | Word of wisdom string received from server: Discipline is wisdom and vice versa. M. Scott Peck
```

## 2. Overview

#### 1. Below you could find the schema of the architecture of the developed service.
   ![alt text](https://upload.wikimedia.org/wikipedia/commons/5/55/Proof_of_Work_challenge_response.svg)

#### 2. I've chosen the [Quadratic Residue Modulo](https://en.wikipedia.org/wiki/Quadratic_residue#Complexity_of_finding_square_roots) implementation of Proof of Work algorithm.
Simply describing this implementation could be separated by stages:
1. Client: Sends the request for the challenge to server.
2. Server: Receives request for challenge.
3. Server: Generates random UInt16 `n`, random **prime** UInt32 `m <= 100_000_007`, which represents modulo.
   Then it calculates the quadratic residue modulo `q = n * n mod m`
4. Server: Sends the challenge to the client. Particularly challenge is an object containing 2 numbers: `q`, `m`.
5. Client: Receives challenge from the server.
6. Client: Calculates Proof `p`, such that `p * p = q mod m`. Particularly it starts with `p = 2` and increments `p` by  `1` if Proof is not correct.
7. Client: Sends the request for resource to the server. Particularly this is an object, containing 3 numbers: `p`, `q`, `m`.
8. Server: Receives request for resource.
9. Server: Calculates the `x = p * p mod m` and checks, whether `x = q mod m`.
   1. If the proof was found correctly it just sends the resource to the client. Particularly resource is just a single _Word-of-Wisdom_ quote.
   2. If the proof is incorrect it just writes the corresponding message about it to the console.
#### 3. I've chosen [Quadratic Residue Modulo](https://en.wikipedia.org/wiki/Quadratic_residue#Complexity_of_finding_square_roots) because:
1. It is simply to implement this algorithm on both client and server side
2. With the large prime numbers this algorithm works well - it requires client some real work (by checking all numbers until it finds correct one)
3. This algorithm is not used in some complex systems (such as e.g. [Hashcash](https://en.wikipedia.org/wiki/Hashcash) which is mainly used in blockchain), but this algorithm is time-proved!

#### 4. Prerequisites
1. [Golang](https://go.dev/)
2. [Docker](https://www.docker.com/)
3. [Docker-compose](https://docs.docker.com/compose/)