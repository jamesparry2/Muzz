# Muzz

The Muzz API is a minature api that is designed to handle the managament of users and provide a public interface for a consuming application to use to coordiante discovery of new users for a user to match with, and provide the user the opportunity to match with other users.

Below will be subsections denoting the requested functionality, the assumptions that were made that led to technical decisions and how the assumptions impacted the final endpoints. Also will be some setup steps on how to run the API locally, for manual verification of the endpoints. 

## Feature Requests
- [x] Allow the ability to create users
- [x] Allow users to login
- [x] Allow users to discover other users to swipe on
- [x] Allow users to swipe and potentially match with others
- [x] Allow users to filter discoveries based on age and gender
- [x] Order the results based on distance to the users
- [ ] Bonus: Apply a secondary order filter based on attractives

## Assumptions made
- For determining the distance, the assumption that the consumer would provide the lat and long of the users current location. This lead to a new location resource which would be a child of the user resource. 
- Regarding the variety of potential data integration options, the API would be built with the design of Ports/Adapters to allow easy switching of underlying technical stores. 

## Running the application locally
For running the application the assumption, outline in the document, is the user will have docker installed. So using this I've built a docker compose file that will use the local image along with a DB container and handle the networking between the two. There will be two ways to achieve this depending on the users setup:

- 1 Within the root there is a Makefile with a command which handles building the image and tagging the local image, it will then run the docker-compose in detached mode to free up the terminal. `make setup_and_run_local`
- 2 If the user does not have the 'Make' command installed, and do not wish to install it, they can follow the 3 commands to spin up the API: 
    `docker build --tag muzz-api .`
	`docker tag muzz-api:latest muzz-api:local`
	`docker-compose up -d`

## Technical Design Decisions
ToDo