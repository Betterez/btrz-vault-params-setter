Service files
===============
The service files are json files contains:
1. service name. repeat twice, in the section name and in the service name.
2. environments: in case of several environments multiple values can bve used.
3. arns: AWS arn for the policies used by the service
4. mongodb: will create a user and a password for this service for each environment

service files should be in the services folder. A sample service file is in this folder under the name of `services.json`
