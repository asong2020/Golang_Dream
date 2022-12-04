import requests

if __name__ == '__main__':
    target = "http://www.baidu.com/"
    req = requests.get(url=target)
    req.encoding = 'utf-8'
    print(req.text)