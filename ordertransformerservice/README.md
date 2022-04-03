# Order Transformer Service

1. Order Validator and Enrichment Application. Data is provided to us in form of a JSON file. We will have 4 functions :- 
  i.) One parses and stores multiple JSON files into a single struct to be processed further. It also takes into account a yaml file which defines the rules to validate the JSON. 
  ii.) Next, based on some yaml rules (taken from a file), we validate this struct.
  iii.) Based on the same yaml, we perform data enrichment using text templates.
  iv.) This is sent forward to another backend endpoint.
2. Understanding Mocks and Functional Testing 
  i.) Write Unit tests for the individual functions. 
  ii.) Show 2 ways we can mock the external http call. Normal `Patching` and `HTTPMOCKS`
  iii.) Write Functional Tests for the whole application (Only this wrapper code. We can extend this to master after this blog).
  iv.) Understand more about the domain of unit and functional tests.
3. Jenkins setup with Docker and connect to Github 
4. Code the CI/CD steps with Jenkins 
5. Continuous Deployment with Minikube, Jenkins and Github. 
6. Cloud Insights : Understanding Hosted Devops Pipelines and their power.

Keep all the Jenkins files and everything local to a single directory. Make sure Functional tests don't break the current application. 

# Rules

You can provide rules like 
1. Amount Filter
  i.) amountfilter: > 18000 
  ii.) amountfilter: < 18000
  iii.) amountfilter: = 18000
2. CreateAt Filter
  i.) createatfilter: > 1000000 
  ii.) createatfilter: < 1000000
  iii.) createatfilter: = 1000000
3. Remove the order with product ids from the last list.
  i.) blacklistproduct: - 1
                        - 2
4. Get orders from a single email id
  i.) emailfilter: abc@example.com
