cd # The Number Prediction Game with Golang and RabbitMQ


https://user-images.githubusercontent.com/92356291/161639046-955bbb14-967b-4398-99be-79c9c5ca525c.mp4




## Dependencies used in the project

- Go 1.17.3
- RabbitMQ 3.9.14
- Docker 20.10.4

## Project Setup
- Clone the project with `git clone https://github.com/Ertugrulbal/Number_Prediction_Game_With_Golang_RabbitMQ.git`
## Run 
- Run the project with `docker-compose up`

## What is this project aims to?
Picture of Project 
* Program A generates a number between 0 and 9 and adds it to the NuerStoreA queue.
* Program B also expects data from the NumberScoreB queue
* If data comes from the queue, it makes 5 number predictions and if one of the prediction is correct, it gets 1 point to PointStore queue.
* This process will continue continuously for both programs.





## The Structure of Code 

    -image of structure
* Each of the Program have a Docker Image
* There is a docker-compose.yml for the whole docker image work together. 
* ProgramA and ProgramB have a Publisher Function and Consumer Function. 
