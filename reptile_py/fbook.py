
import requests # requests 包 进入url访问
from bs4 import BeautifulSoup
from tqdm import tqdm #tqdm是一个快速、可扩展的python进度条，可以在python长循环中添加一个进度提示信息

def get_content(target):
    req = requests.get(url=target)
    req.encoding = 'utf-8'
    html = req.text
    bs = BeautifulSoup(html,'lxml')
    texts = bs.find('div',id='content')
    content = texts.text.strip().split('\xa0'*4)
    return content
if __name__ == '__main__':
    server = "https://www.xsbiquge.com"
    target = "https://www.xsbiquge.com/91_91600/"
    book_name = '三国之他们非要打种地的我.txt'
    req = requests.get(url=target)
    req.encoding = 'utf-8'
    html = req.text
    characters_name = BeautifulSoup(html,'lxml')
    characters = characters_name.find('div',id='list') # 寻找标签
    # 提出a标签
    characters = characters.find_all('a')
    #进行循环读取章节
    for character in tqdm(characters):
        # 得到href属性
        url = server + character.get('href')
        # 获取章节名字
        character_name = character.string
        # 获取文章
        content = get_content(url)
        with open(book_name,'a',encoding='utf-8') as f:
            f.write(character_name)
            f.write('\n')
            f.write('\n'.join(content))
            f.write('\n')