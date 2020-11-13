Catalog of stolen bikes
====

![](https://images.unsplash.com/photo-1556316384-12c35d30afa4?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=3450&q=80)

Stolen bikes is a typical problem in Malmö, where the Docly HQ is. We need your help to build a service that will help the local police to solve bike thefts in the area.

# Requirements

- [ ] The police wants to able to add, edit and remove officers. (See data model suggestion below)
- [ ] Private citizens want to be able to report stolen bikes. (See data model suggestion below)
- [ ] The system should assign a police officer to handle stolen bike cases as they are being reported.
  - [ ] A police officer can only handle one case at a time.
- [ ] The police should be able to report bikes as found.
  - [ ] When the police finds a bike the case should be marked as resolved and the police officer in charge of the case should be marked as available to take a new case.
- [ ] The system should be able to assign unassigned cases to police officers as they become available.
- [ ] Police and private citizens both want to be able to list the reported bike thefts and their status.

# Tech requirements

- Go or Node.js (Typescript is a plus)
- Tests 
- You may use any database

# Your challenge

- Create an API that satisfies the requirments above.
- Create documentation on how to use Postman (or similar tool) to interact with the API.
- You can use any boilerplate and tools that you want to but we advice you to keep it simple. A clean, robust API is what we're looking for, and we're usually in favor of using all available tools and tricks to get things done.
- We prefer Docker to run our services, but as long as you have clear instructions on how to run yours you may use whatever you want.

# Instructions

- Fork this repo
- Build a clean and robust API
- Let us know that you've completed the challenge and how we can test it.

- police officer endpoint:
  - http://localhost:8080/officers                method: GET, description: fetch all the police officers 
  - http://localhost:8080/officers/{id}           method: GET, description: fetch one police officer by id, for example fetch officer whose id is 3: http://localhost:8080/officers/3
  - http://localhost:8080/officers/{id}           method: PUT, description: modify one police officer, go to Postman->Body->raw, enter: {"name": "Ellie Johansson"}
  - http://localhost:8080/officers                method: POST, description: create one police officer, go to Postman->Body->raw, enter: example: {"name": "Lars Nilsson"}
  - http://localhost:8080/officers/{id}           method: DELETE, description: remove one police officer

- bike thefts case management endpoint:
  - http://localhost:8080/bike-thefts             method: GET, description: to fetch all bike theft cases
  - http://localhost:8080/bike-thefts             method: POST, description: to create a bike theft case with one image file
                                                  example: go to Postman->Body->form-data, 
                                                  step 1: upload image file, under `KEY` enter `image` and choose `file`, under `VALUE` select an image file from the local computer.      
                                                  step 2: enter json data, in Postman->Body->form-data, under `KEY` enter `data`, under `VALUE` paste the json data like below:
                                                  {
                                                      "title": "people who are Joyriders",
                                                      "brand": "water soul",
                                                      "city": "malmö",
                                                      "description": "black, fint and well built, 28 model"
                                                  }
  - http://localhost:8080/bike-thefts-no-image    method: POST, description: to create a bike theft case without image file
                                                  example: go to Postman->Body->raw, enter the json data like below:
                                                  {
                                                      "title": "someone who is a Joyrider",
                                                      "brand": "apron kitten",
                                                      "city": "malmö",
                                                      "description": "green, old and well built, 24 model"
                                                  }
  - http://localhost:8080/bike-thefts/{id}        method: GET, description: to fetch a bike theft case
  - http://localhost:8080/bike-thefts/{id}        method: PUT, description: to modify a bike theft case, example: {"solved": true}
  - http://localhost:8080/bike-thefts/image/{id}  method: GET, description: to fetch a bike image by case id

- assign a case to an officer  
  - http://localhost:8080/case-to-officer         method: POST, example: {"case": 12,"officer": 3}, note: case id and officer id
# Suggested data model

### Police officers

| Column | Type    | Description     |
| ---    | ---     | ---             |
| id     | Integer |                 |
| name   | String  | Name of officer |

### Bike thefts

| Column      | Type     | Description                              |
| ---         | ---      | ---                                      |
| id          | Integer  |                                          |
| title       | String   | Title of report                          |
| brand       | String   | Brand of bike                            |
| city        | String   | City where theft occured                 |
| description | String   | Description of bike                      |
| reported    | DateTime | Date and time when theft was reported    |
| updated     | DateTime | Date and time when case was last updated |
| solved      | Boolean  | True when case has been solved           |
| officer     | Integer  | Officer in charge of the case            |
| image       | Blob     | Image of bike                            |


# License

This project is licensed under MIT. Feel free to use it anyway you see fit.
