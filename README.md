# Task
Design and implement "Word of Wisdom" tcp server:

TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
The choice of the POW algorithm should be explained.
After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
Docker file should be provided both for the server and for the client that solves the POW challenge.

# Why hashcache? 
Strong Security: Based on cryptographic hashing, Hashcash provides robust protection against spam and denial-of-service attacks.
Simplicity: Its straightforward design makes integration smooth and minimizes development efforts.
Proven Track Record: Hashcash has been widely adopted and used successfully in various systems.
Adjustable Difficulty: We can fine-tune the mining difficulty for consistent block generation and network performance.
Energy Efficiency: Hashcash allows us to balance security and energy consumption by controlling the difficulty level.


