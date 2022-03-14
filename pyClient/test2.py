
# import urllib, urllib2, sys


# host = 'http://alirmcom2.market.alicloudapi.com'
# path = '/query/comp'
# method = 'GET'
# appcode = "76f77c394bd0401abb79d8151371c705"
# querys = 'p=1&ps=10&rout=CNST&sort=ZF&sorttype=0&where=where'
# bodys = {}
# url = host + path + '?' + querys

# request = urllib2.Request(url)
# request.add_header('Authorization', 'APPCODE ' + appcode)
# response = urllib2.urlopen(request)
# content = response.read()
# if (content):
#     print(content)

import requests

appcode = "76f77c394bd0401abb79d8151371c705"
headers = {    
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
    "Authorization": "APPCODE "+appcode
}

host = 'http://alirmcom2.market.alicloudapi.com'
path = '/query/comp'
type="CNST"
url = host+path
params = {
    "p": 1,
    "ps": 1000,
    "rout": type,
    "sort":"ZF",
    "sorttype":0,
    "where":"1"
}
r = requests.request('GET',url=url,headers=headers,params=params)
print(r)
print(r.text)