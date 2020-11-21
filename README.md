# readyGo 

## What is readyGo project?

- readyGo is a project that gives ready to use project structure  based on project type .
  The following project types are suppported      
- 
    1.http
    2.grpc
    3.cloudEvents
    4.cli

## Various components in this project

- There are mainly three components.

  1. Templates --> All templates reside in templates/ directory. This contains templates for generating files.Each template file name is taken as a key for the template engine to keep templates in memory. 
    To create a new template, simply add a file to template directory with go templatting syntax. The data feed for the template as of now is only Module data that comes from configs/config.json or any user given configs file.

  1. Configurators --> Each project type has a configurator. Its again a json file resides in configs/template_config.json. This file contains three things. 1) Directories to create in DS section. 2) Static files to copy with src and destination in FS section and 3) Templates to consider for a project type that is in TS section.
  So what directories to be created , which static files to be copied from source to destination and the third one is which all templates/ to be used for this particular project.
  
  example:

  `{
	"FC": {
		"http_mongo": [{
			"src": "static/databases/mongo/database.static",
			"dst": "database/database.go"
		}, {
			"src": "static/containers/Dockerfile",
			"dst": "Dockerfile"
		}, {
			"src": "static/helper/helper.static",
			"dst": "helper/helper.go"
		}]
	},
	"TC": {
		"http_mongo": ["models", "interfaces","database", "handlers"]
	},
	"DC": {
		"http_mongo": ["models", "interfaces","database", "handlers","helper"]
	}
}`

from above json template example, the project type is http_mongo (that means http service with mongodb as backend).


- 3 files to be copied from source to destination. In general all source files be in the readyGo project static/ directory
- 4 templates to be used they are frmo templates/ directory.
- 5 directories to be created " models,interfaces,database,handlers,helper

To configure another project , add a new entry in configs/template_config.json file. Create needed templates. configure required directories. If required copy static files

1. Generate -- > Takes the input of templates , static files and configurators, Generate generates all required files with the given directory and entity details in the configs/config.json file
