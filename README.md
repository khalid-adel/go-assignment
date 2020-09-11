## Description
The program is a REST API that takes an alert object and returns an equivalent alarm object. 
The system in which this REST API works produces alerts with different formats from different resources. 
The task of this program is to convert any alert to a unified format called alarm. For that the program uses mapping files. 
Every module in the system has a mapping file that maps its alert fields to the unified alarm fields. 
For performance I added a feature that caches the content of mapping files i.e once the program reads the content of a mapping file 
it caches its content, so that when it receives another alert from the same module the program can use the cached mapping instead of reading the file again.

## How to use
In the folder “src” execute the command “go run main.go”, and the REST API will start. 
You can use Postman or any other software to make POST requests towards the API. 
The API is running on the URL: http://localhost:10000/convert

Under the folder “sample-alerts” you will find sample body requests that can be used in the POST requests. 
Under the folder “mapping” you will find sample mapping files that are used by the program.

## Notes:
1. The testing was done using requests from Postman.
1. You need Golang installed on the machine to run this code
1. You can create your own sample body requests, but you should notice the following
     1. There is a matching between “ServiceType” and the name of mapping file
     1. The alarm fields should keep the same names in any mapping file. You can find 
        their names in the Alarm struct in the code.
