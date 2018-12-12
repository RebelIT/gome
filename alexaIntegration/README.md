# Controlling rPi with Alexa
Credit where credit is due, I used this blog post as well as his code framework to set this up, re-documenting for my own purposes and automating in my playbooks to incorporate it into my piProjects. http://www.cyber-omelette.com/2017/01/alexa-run-script.html

* Create an AWS account
  ```
  https://us-east-2.console.aws.amazon.com/console/home?region=us-east-2
  ```

* Register your amazon developer account (Must be the same account registered with our echo)

  ```
  https://developer.amazon.com/edw/home.html
  ```
---

* Create your messaging queue
1. In your AWS console search for the "Simple Queue Service (SQS)" service
2. Create a new Queue
3. Give it a descriptive name
4. Standard Queue, Quick-Create
5. Once created note the following:

  ```
  URL:	https://sqs.us-east-2.amazonaws.com/<some number>/<queue name>
  ARN:	arn:aws:sqs:us-east-2:<some number>:<queue name>
  ```

---
* Create an API user for your queue
1. In your AWS console search for the "IAM" service
2. Go to users
3. Create a new user and assign permissions

---
* Create your lambda function to do an action
1. In your AWS console search for the "lambda" service
2. IMPORTANT: Switch to US East N. Virginia region (this region supports alexa skills)
3. Create a Function
4. Name it
5. Use Python 2.7 as the runtime language
6. Use a role from template
7. Name the role (read more about roles and templates on your own)
8. Leave Policy template blank
9. Create Function
10. In the next screen click on your functions
11. Scroll down to edit the functions
12. Copy contents from lambda_function.py into the code block editing if statement on line 42 for your intents and message to write.
13. Click on the blank dotted box in the tree view, add Alexa Skills Kit

---
* Create your skill
1. Follow the link to the Alexa skills developer portal
2. Click get started to the Alexa Skills Kit
3. Use custom interaction model
4. Fill out the name and invocation name (what you want Alexa to trigger on) and Save
5. Fill out the interaction model (what the intents are and what you want Alexa to respond to) - Use the examples (alexa_intent_model.json and alexa_utterances.ini)
6. Configuration, use AWS Lambda ARN endpoint and paste the lambda ARN url from the function you created. The rest can be default.
7. Test, here you can type in a phrase and test it against your utternaces and intents, this will actually write to your queue, so log in and check it when you test.
8. Publishing, pick a category that fits your skill
9. Privacy and compliance, be honest to the questions...
10. Save it, if you log into your Alexa account it should now appear in your list of installed skills with a [dev] tag.
11. Follow the Alexa ansible role (ansible/roles/alexa_consumer) to configure the consumer.
