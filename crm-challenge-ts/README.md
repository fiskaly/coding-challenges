## Instructions

Welcome challenger!  
When you see this, you have successfully passed the first interview at fiskaly. Now 
comes the time to prove yourselves. We have the following challenge prepared for you, 
containing elements from frontend development, backend development as well as database
handling. Now let's jump right into it! 

### Prerequisites:
- docker

### Setup: 
- run `setup.sh`. This will start 3 docker containers, one each for frontend, backend, and a postgres server, respectively. 
- make sure the ports 3000, 3001, and 5432 are not occupied on your machine. 
- you can find the frontend at `http://localhost:3000`

### The Challenge:   
We want you to prepare / implement a simple CRM tool for us. The basic structure already 
exists. Your task is to add some features to it.    

- add a form in the frontend, which lets you add new customers to the database. On submission, the data from the form will be sent to the backend, which in turn handles the database entry. The form should allow you to enter information for the first-name, last-name and mail of the customer. 
- add an api-endpoint in the backend, which lets you create a new customer entry in the database. Return a message which indicates success.
- add a table in the frontend, which displays all the customers currently in the database. 
- add a simple filter to the table, which lets you filter customers by their last name

Please do not focus on styling in this challenge, but rather functionality! You are free to add any additional 
api-endpoints you might need. 

Optional features: 
 - for each new customer, create a random uuid and add this in the `customer_id` field in the database
 - you might have noticed that some customers have multiple entries in the database, whereas the difference is the tss_id. Find a solution to display the list of tss along with customer information
 - restrict the input of the form to valid input: i.e. only email-addresses are allowed in the mail-field; no numbers are allowed in the name fields; etc. 
 - add an option to create a new tss (as indicated by the tss_id) for an existing customer. Handle cases properly, where the customer does not exist already.

We wish you all the best and are looking forward to your results!
Cheers