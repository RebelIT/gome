import boto3
import random

access_key = "xxxx"
access_secret = "xxxx"
region ="us-east-2"
queue_url = "https://sqs.us-east-2.amazonaws.com/xxxxx/xxxx"

responsesOk = ["I sent the request for you",
               "done",
               "sure thing, sent it on",
               "no problem",
               "i can do that for you",
               "simple, done",
               "you got it",
               "on it",
               "sure thing, thats what i'm here for"
               ]

responsesIdk = ["I dont think a device like that exists",
                "what, nope, cant do it",
                "sorry, i'm not doing this on purpose",
                "bad request, five o three",
                "i wish i could help but i do not get it",
                "done hate me, but i cant help you"
                "check your code, this is not right"
                ]

def build_speechlet_response(title, output, reprompt_text, should_end_session):
    return {
        'outputSpeech': {
            'type': 'PlainText',
            'text': output
        },
        'card': {
            'type': 'Simple',
            'title': "SessionSpeechlet - " + title,
            'content': "SessionSpeechlet - " + output
        },
        'reprompt': {
            'outputSpeech': {
                'type': 'PlainText',
                'text': reprompt_text
            }
        },
        'shouldEndSession': should_end_session
    }

def build_response(session_attributes, speechlet_response):
    return {
        'version': '1.0',
        'sessionAttributes': session_attributes,
        'response': speechlet_response
    }

def post_message(client, message_body, url):
    response = client.send_message(QueueUrl = url, MessageBody= message_body)

def lambda_handler(event, context):
    client = boto3.client('sqs', aws_access_key_id = access_key, aws_secret_access_key = access_secret, region_name = region)
    intent_name = event['request']['intent']['name']
    slot_id = event['request']['intent']['slots']['NAME']['resolutions']['resolutionsPerAuthority'][0]['values'][0]['value']['id']
    slot_action = event['request']['intent']['slots']['ACTION']['value']
    m = intent_name+","+slot_id+","+slot_action

    # validation of intent
    if intent_name == "TuyaControl":
        post_message(client, m, queue_url)
        message = (random.choice(responsesOk))
    else:
        message = (random.choice(responsesIdk))

    speechlet = build_speechlet_response("message", message, "status", "true")
    return build_response({}, speechlet)