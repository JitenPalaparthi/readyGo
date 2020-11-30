# readyGo 

# Overview

A Simple configuration gives you a working project.

- readyGo is a command line interface( probably the name of readyGo CLI would be rgo) application, it is designed to scaffold creation of different types of go based projects.readyGo is designed for developers in mind. Ideally readyGo should provide ready to use application code. The code is generated based on configurations provided by the end user i.e "The great developer :)".
- By version 1 release, it will support http, grpc, CloudEvents template engines with various databases (sql/no sql), pub-sub and CloudEvents plugins and probably even more.

*readGo cooks for you.As your business logic varies , you have to add the required spices according your taste but one thing .. readyGo is not a template engine, it gives you a working project.*

## Reason behind this project

- A brief history of my work in general: I get ideas like anyone else. Start developing that idea, create a workspace, working on the project , by the time I am done with structure ,CRUD operations and even more repetitive logic/code, I loose most my time and then park the project aside. This made me understand that we should not loose time for repetitive tasks.So readyGo.

## The current state of the project

- Give a config file(json | yaml) with fewer details as mentioned in the below yaml file, readyGo gives you a working REST API with MongoDb as the back end.
- The config file is self explanatory. Here is a sample one .Note: This will change in future

``` yaml
    ---
version: '0.1'
project: example
type: http
port: '50054'
db: mongo
models:
- name: person
  fields:
  - name: name
    type: string
    isKey: true
  - name: email
    type: string
    validateExp: "[a-zA-Z0-9]"
    isKey: true
  - name: mobile
    type: string
  - name: status
    type: string
  - name: last_Modified
    type: string

```
 ### From the above config, it says to readyGo engine that,

- Create a http restapi that runs on port 50054.
- Use the mongodb as database.
- Create a model mamed person with (name,email, mobile,status and last_updated) as fields.
- Validate email with the given validateExp (regular expression).
    
### As a result , the output is,

- A root directory with the name example is created.
- A database directory with all database related methods are created.
- An interface directory with CreatePerson,UpdatePersonByID,DeletePersonByID,GetPersonByID,GetAllPersons,GetAllPersonsBy  definitions are created.
- Models(person),handlers(CreatePerson,UpdatePersonByID,DeletePersonByID,GetPersonByID,GetAllPersons(with skip and limit for pagination),GetAllPersonsBy (various search parameters)) with required validations are created.
- A Dockerfile and docker-compose files are created.
- An application.json configuration file is created.
- The restful service is ready to be started without making a single change to the project.

