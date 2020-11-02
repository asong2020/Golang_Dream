# 工作中常用Git使用指南

`author:asong` `公众号：Golang梦工厂`

> 因为自己工作经常与git打交道，所以就想着整理一份自己工作中常用的git，只整理指令，不讲解原理，因为原理这个东西，官方文档讲的太明白了，我也没必要在这里再整理一遍。这里只介绍我工作中常用的命令，因为命令这个东西比较多，但是真正使用的却没有那么多，所以这个主打常使用指令，针对工作没几年的、或者初入职场的朋友。有一些不好的指令，我会特别表明，请慎用！！！
> 我自己做这个肯定考虑是不全的，有兴趣的小伙伴可以一快加入进来，我们共同维护。因为我也是一个新手，有错误欢迎指出。
>
> 这个我会最后整理成PDF，可以到我的github上获取，也可以到我的公众号获取，目前还没有放上去。



### 基础操作

#### 查看提交历史

```shell
# 不带任何参数，git log以相反的时间顺序列出在该存储库中所做的提交，
$ git log 

# -p 或者 --patch，它显示了每次提交中引入的差异(补丁输出). -2 限制显示的日志条目的数量。例如-2用于仅显示最后两个条目。
$ git log -p -2

# 上面查看的是比较详细的信息 如果想查看每次提交的一些简短统计信息，可以使用如下命令
$ git log --stat

# 还有一个真正有用的选项是 --pretty.此选项将日志输出更改为默认格式以外的其他格式。一些预建的选项值可供您使用。oneline此选项的值将每条提交打印在一行上，如果您要查看大量的提交，这将很有用。 short，full，和fuller值显示在大致相同的格式，但分别与更少或更多的信息，输出：
$ git log --pretty=oneline

# 最有趣的选项值是format，它允许您指定自己的日志输出格式。当生成用于机器解析的输出时，这特别有用-因为您明确指定了格式，所以您知道它不会随着Git更新而改变：
$ git log --pretty=format:"%h - %an, %ar : %s"
```

### 撤销操作

#### 撤销`git add`
```shell
$ git reset HEAD . # 这里撤销的是已经执行git add的文件
$ git reset HEAD -filename # 撤销某个文件或文件夹
```


### 关联github仓库

这里省略github上创建仓库的步骤，新创建一个仓库后，在本地关联步骤：

```shell
$ echo "# go-algorithm" >> README.md
$ git init
$ git add README.md
$ git commit -m "first commit"
$ git branch -M master
$ git remote add origin git@github.com:<github名字>/<仓库名字>.git
$ git push -u origin master
```

若是一个已存在的仓库，在本地关联步骤:

```shell
$ git branch -M master
$ git remote add origin git@github.com:<github名字>/<仓库名字>.git
$ git push -u origin master
```




### 分支操作指令

#### 查看分支

- 查看所有的分支

```shell
$ git branch
```

#### 创建分支
```shell
$ git branch <branch name> 
$ git checkout -b <branch name> # 这个我平常用的比较多,创建分支并切换到该分支
```

#### 切换分支

```shell
$ git checkout <branch name> # 切换分支
$ git checkout -b <branch name> # 这个我平常用的比较多,创建分支并切换到该分支
```



#### 本地分支关联远程分支

```shell
$ git push origin <branch name>:<branch name> # 将本地xn分支推送至远程xn分支

$ git push --set-upstream origin <branch name> # 切换到你要推送到远程分支的本地分支 进行关联
```


#### 远程分支

##### 查看远程分支

- 查看远程分支

```shell
$ git branch -a # 查看本地分支和远程分支
$ git branch -v
$ git branch -r # 查看远程仓库的分支
$ git branch -vv # 查看本地分支关联(跟踪)的远程分支之间的对应关系，本地分支对应那个远程分支
```

- 查看远程分支列表

```shell
$ git ls-remote
```

##### 删除远程分支

```shell
$ git push origin :<branch name> //将一个空分支推送到远程即为删除
$ git push origin --delete <branchName>
```


#### 合并分支

合并分支都是现在本地分支进行合并，然后推到远程分支的，具体步骤看接下来的指令。

##### 1. 合并分支

这里[branchName] 是指要合并的分支。
```shell
$ git merge <branchName>
```

简便写法：
```shell
$ git merge - 
```
短横线代表上一个切换过来的分支。

##### 2. 合并冲突

合并分支并不是一番风顺的，所以当合并分支发生冲突时，我需要手动解决冲突。

- 冲突标志

```text
这些行与普通行没有什么不同
祖先或彻底解决，因为只有一侧发生了变化。
<<<<<<< yours：sample.txt
解决冲突很难。
我们去买东西吧。
=======
Git使冲突解决变得容易。
>>>>>>>他们的：sample.txt
这是另一行经过完全解析或未修改的行。
```
其中一对相互矛盾的变化发生的区域标有标记 <<<<<<<，=======和>>>>>>>。之前的部分======= 通常是您的一面，之后的部分通常是他们的一面。

- 查看分支状态，可以找到发生冲突的地方

```shell
$ git status
```

- 手动解决冲突后，重新提交文件

```shell
$ git add <file>
```

- 冲突停止后，结束合并指令

```shell
$ git merge --continue
```

- 本地合并分支结束后，推送到远程即可，记得在推送时，拉取一下远程分支的代码

```shell
$ git pull
$ git push
```

- 如果分支合错了，想退回到原来版本，执行如下指令

```shell
$ git reset --hard
```

#### 根据commit_id查看对应分支

说一下这个使用的场景，我们通常在团队开发中，团队中队员有事临时不再，却又不知道对方代码所在位置的时候，可以通过一些`commit_id`来确定其代码的分支。

```shell
$ git branch -r --contains COMMIT_ID
```

