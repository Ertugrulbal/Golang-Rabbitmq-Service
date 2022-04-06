cd # The Number Prediction Game with Golang and RabbitMQ- (Run in Kuberetes)

## Dependencies used in the project

- Go 1.17.3
- RabbitMQ 3.9.14
- Docker 20.10.4
- DevSpace 5.18.4
- Kubernetes-Kompose 1.26.1

## Project Setup
- Clone the project with `git clone https://github.com/Ertugrulbal/Number_Prediction_Game_With_Golang_RabbitMQ.git`
- cd `RabbitMQ-In-Kubernetes`
- 
## Run 
- The test for application with `docker-compose up`
- kubectl apply -f rabbitmq, my-service-a, my-service-b
- Also the enter a code for po tests `kubectl get po`
- ![pods2](https://user-images.githubusercontent.com/92356291/161960621-46e7d727-f958-4b78-a60c-53800f7d4a58.jpeg)


## What is this project aims to?
![image](https://user-images.githubusercontent.com/92356291/161713615-1afb6712-65ad-46cd-9e0d-a8dbc03a2ec2.png)

* Program A generates a number between 0 and 9 and adds it to the NuerStoreA queue.
* Program B also expects data from the NumberStoreB queue
* If data comes from the queue, it makes 5 number predictions and if one of the prediction is correct, it gets 1 point to PointStore queue.
* This process will continue continuously for both programs.



## The Structure of Code For Kubernetes

![image](https://user-images.githubusercontent.com/92356291/161961195-bc24b7ac-ec99-4100-89a2-319ca5d70fb6.png)
* ProgramA and ProgramB have a Publisher Function and Consumer Function.
* Each of the Program have a Docker Image
* There is a docker-compose.yml for the whole docker image work together. 
* In order to run project in kubernetes, we had to convert docker-compose.yml file to app/pod.yaml and app/service.yaml files via kubectl and kompose tools.
* All of the services have a `pod.yaml` and `service.yaml` files.


