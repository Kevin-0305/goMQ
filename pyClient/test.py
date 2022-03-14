import requests
data = {
    "channelName":"123",
    "content":"test",
    "messageType":1
}
r = requests.request('POST','http://0.0.0.0:9630/mq/messagePublish/',json=data)
print(r.text)