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

## Technical Design & Decisions

### Project Structure

### Database Decision and ORM

For the database its been built against SQL primiarly to leverage the relational structure of the data that was being stored as the relationship can be mapped to one is many and vice vera such as a User can have many swipes but a swipe must belong to a single user.

When deciding what DB Tool to use with this, it was decided to be MYSql to levearge the built in functionality that would support with the distancing query rather than trying to implement my own solution. The only drawback to useing the built in functionality is that, according to documentation, its not the most performant at scale but for this solution and use case it wouldn't be a problem.

For the ORM in Go, I went with GORM because it allowed for the following features to be used:
- automigration for the table strucures and allowed a model first appraoch
- data parameterization to attempt to limit any potential SQL Injection
- Quicker startup due to it ease of setup

However there are some drawback, its automigration is a powerful tool but it only supports addiditve proccess so deletion and manual column changes are a challange.