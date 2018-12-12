import boto3
import os
import time
from roku import Roku

access_key = "<access_key>"
access_secret = "<Super_Secret>"
region ="us-east-2"
queue_url = "https://sqs.us-east-2.amazonaws.com/<id>/account"
tv = Roku('10.0.0.xxx<tv_ip>')

def pop_message(client, url):
    response = client.receive_message(QueueUrl = url, MaxNumberOfMessages = 10)

    #last message posted becomes messages
    message = response['Messages'][0]['Body']
    receipt = response['Messages'][0]['ReceiptHandle']
    client.delete_message(QueueUrl = url, ReceiptHandle = receipt)
    return message

client = boto3.client('sqs', aws_access_key_id = access_key, aws_secret_access_key = access_secret, region_name = region)

waittime = 20
client.set_queue_attributes(QueueUrl = queue_url, Attributes = {'ReceiveMessageWaitTimeSeconds': str(waittime)})

time_start = time.time()
while (time.time() - time_start < 60):
    print("Checking...")
    try:
        message = pop_message(client, queue_url)
        print(message)
        if message == "powerButton":
            tv._post('/keypress/Power')

        elif message == "hdmiOn":
            os.system("/home/pi/automation/local_hdmi_control/rpi-hdmi.sh on")

        elif message == "hdmiOff":
            os.system("/home/pi/automation/local_hdmi_control/rpi-hdmi.sh off")

        elif message == "down":
            os.system("sudo /sbin/shutdown -h now")

        elif message == "reboot":
            os.system("sudo /sbin/shutdown -r now")

    except:
        pass