# 网络爬虫 一本小说

import requests
from bs4 import BeautifulSoup
from tqdm import tqdm

def get_content(target):
    req = requests.get(url=target)
    req.encoding = 'utf-8'
    html = req.text
    bf = BeautifulSoup(html,'lxml')
    texts = bf.find('div',id='content')
    content = texts.text.strip().split('\xa0'*4)
    return content

if __name__ == '__main__':
    # target = "https://www.xsbiquge.com/15_15338/8549128.html"
    # req = requests.get(url=target)
    # req.encoding = 'utf-8'
    # html = req.text
    # bs = BeautifulSoup(html,'lxml')
    # texts = bs.find('div',id='content')
    # # texts.text 提取所有文字
    # # strip 方法去掉回车
    # # split 方法根据\ax0切分数据，每段开头都有4个空格
    # print(texts.text.strip().split('\xa0'*4))
    server = "https://www.xsbiquge.com"
    book_name = '诡秘之主.txt'
    target = 'https://www.xsbiquge.com/15_15338/'
    req = requests.get(url=target)
    req.encoding = 'utf-8'
    html = req.text
    chapter_bs = BeautifulSoup(html,'lxml')
    chapters = chapter_bs.find('div',id='list')
    chapters = chapters.find_all('a')
    for chapter in tqdm(chapters):
        chapter_name = chapter.string
        url = server + chapter.get('href')
        content = get_content(url)
        with open(book_name,'a',encoding='utf-8') as f:
           f.write(chapter_name)
           f.write('\n')
           f.write('\n'.join(content))
           f.write('\n')