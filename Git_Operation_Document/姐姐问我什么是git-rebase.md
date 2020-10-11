## 前言

> 哈喽，everybody，不知不觉8天的小长假也接近了尾声，玩耍了这么多天，今天也要收一收心，开始学习了呦～。明天就要上班啦，今天姐姐突然问我git-rebase指令是干什么的，怎么用？其实我是不想给他讲的，但是还是没有逃过姐姐的软磨硬泡，那么我们就一起来看一看什么是git-rebase吧！！！



## 缘起

话说，我和姐姐的缘分是在那一个月黑风高的晚上，啪，姐姐一巴掌打在了我的脸上并说了一句：能不能讲重点～～～。哈哈，不开玩笑了，直接说重点吧。我们先来看一个场景，我查看了一下我github上的个人仓库，`commit`提交次数很多，提交内容如下：

![]()



这么多的提交，有很多没有规范的命名，因为是自己使用，就随便开整了，这确实不好，还有一些没有必要的提交，其实是可以合并到一起的，这样会导致如下问题：

- 造成分支污染，项目中充满了许多`commit`记录，当出现紧急问题需要回滚代码时，就只能一条条的查看了。
- 代码`review`不方便，当你要做`code review`时，一个很小的功能却提交了很多次，看起来就不是很方便了。

这一篇文章我们先不讲`git`提交规范，我们先来解决一下如何合并多次提交记录。



## rebase作用一：合并提交记录

通过上面的场景，我们可以引申出`git-rebase`的第一个作用：合并提交记录。现在我们想合并最近5次的提交记录，执行：

```shell
$ git rebase -i HEAD~2
```

执行该指令后会自动弹出`vim`编辑模式：

```vim
pick e2c71c6 update readme
pick 3d2c660 wip: merge`

# Rebase 5f47a82..3d2c660 onto 5f47a82 (2 commands)
#
# Commands:
# p, pick <commit> = use commit
# r, reword <commit> = use commit, but edit the commit message
# e, edit <commit> = use commit, but stop for amending
# s, squash <commit> = use commit, but meld into previous commit
# f, fixup <commit> = like "squash", but discard this commit's log message
# x, exec <command> = run command (the rest of the line) using shell
# b, break = stop here (continue rebase later with 'git rebase --continue')
# d, drop <commit> = remove commit
# l, label <label> = label current HEAD with a name
# t, reset <label> = reset HEAD to a label
# m, merge [-C <commit> | -c <commit>] <label> [# <oneline>]
# .       create a merge commit using the original merge commit's
# .       message (or the oneline, if no original merge commit was
# .       specified). Use -c <commit> to reword the commit message.
#
# These lines can be re-ordered; they are executed from top to bottom.
#
# If you remove a line here THAT COMMIT WILL BE LOST.
#
# However, if you remove everything, the rebase will be aborted.
#
# Note that empty commits are commented out
```

从这里我们可以看出前面5行是我们要合并的记录，不过前面都带了一个相同的指令：`pick`，这是什么指令呢，不要慌，这不，下面已经给出了`commands`:

```shell
pick：保留该commit（缩写:p）

reword：保留该commit，但我需要修改该commit的注释（缩写:r）

edit：保留该commit, 但我要停下来修改该提交(不仅仅修改注释)（缩写:e）

squash：将该commit和前一个commit合并（缩写:s）

fixup：将该commit和前一个commit合并，但我不要保留该提交的注释信息（缩写:f）

exec：执行shell命令（缩写:x）

drop：我要丢弃该commit（缩写:d）

label：用名称标记当前HEAD(缩写:l)

reset：将HEAD重置为标签(缩写:t)

merge：创建一个合并分支并使用原版分支的commit的注释(缩写:m)
```

根据这些指令，我们可以进行修改，如下：

```shell
pick e2c71c6 update readme
s 3d2c660 wip: merge`
```

修改好后，我们点击保存退出，就会进入注释界面:

```shell
# This is a combination of 2 commits.
# This is the 1st commit message:

update readme

# This is the commit message #2:

wip: merge`

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# Date:      Thu Sep 17 22:03:52 2020 +0800
#
# interactive rebase in progress; onto 5f47a82
# Last commands done (2 commands done):
#    pick e2c71c6 update readme
#    squash 3d2c660 wip: merge`
# No commands remaining.
# You are currently rebasing branch 'master' on '5f47a82'.
#
# Changes to be committed:
#       new file:   hash/.idea/.gitignore
#       new file:   hash/.idea/hash.iml
#       new file:   hash/.idea/misc.xml
#       new file:   hash/.idea/modules.xml
#       new file:   hash/.idea/vcs.xml
#       new file:   hash/go.mod
#       new file:   hash/hash/main.go
#       modified:   snowFlake/Readme.md
```

上面把每一次的提交的`meassage`都列出了，因为我们要合并这两次的`commit`，所以提交注释可以修改成一条，最终编辑如下：

```shell
# This is a combination of 2 commits.
# This is the 1st commit message:

fix: merge update and wip: merge`

# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
#
# Date:      Thu Sep 17 22:03:52 2020 +0800
#
# interactive rebase in progress; onto 5f47a82
# Last commands done (2 commands done):
#    pick e2c71c6 update readme
#    squash 3d2c660 wip: merge`
# No commands remaining.
# You are currently rebasing branch 'master' on '5f47a82'.
#
# Changes to be committed:
#       new file:   hash/.idea/.gitignore
#       new file:   hash/.idea/hash.iml
#       new file:   hash/.idea/misc.xml
#       new file:   hash/.idea/modules.xml
#       new file:   hash/.idea/vcs.xml
#       new file:   hash/go.mod
#       new file:   hash/hash/main.go
#       modified:   snowFlake/Readme.md
```

编辑好后，保存退出就可以了。这样就完成了一次合并`commit`。我们来验证一下：

```shell
$ git log
15ace34 (HEAD -> master) fix: merge update and wip: merge`
5f47a82 update snowFlake code
```

从这里我们可以看到，两次提交变成了一次，减少了无用的提交信息。



## 作用二：分支合并

这个作用我们使用的很少，但是还是要知道，下面我们一起来看一下使用场景。

假设我们现在有一个新项目，现在我们要从`master`分支切出来一个`dev`分支，进行开发：

```shell
$ git checkout -b dev
```

这时候，你的同事完成了一次 `hotfix`，并合并入了 `master` 分支，此时 `master` 已经领先于你的 `dev` 分支了：

![]()

同事修复完事后，在群里通知了一声，正好是你需要的部分，所以我们现在要同步`master`分支的改动，使用`merge`进行合并：

```shell
$ git merge master
```

![]()

图中绿色的点就是我们合并之后的结果，执行`git log`就会在记录里发现一些 `merge` 的信息，但是我们觉得这样污染了 `commit` 记录，想要保持一份干净的 `commit`，怎么办呢？这时候，`git rebase` 就派上用场了。

所以现在我们来试一试使用`git rebase`，我们先回退到同事 `hotfix` 后合并 `master` 的步骤，我现在不使用`merge`进行合并了，直接使用`rebase`指令

```shell
$ git rebase master
```

这时，`git`会把`dev`分支里面的每个`commit`取消掉，然后把上面的操作临时保存成 `patch` 文件，存在 `.git/rebase` 目录下；然后，把 `dev` 分支更新到最新的 `master` 分支；最后，把上面保存的 `patch` 文件应用到 `dev` 分支上；

![]()

从 `commit` 记录我们可以看出来，`dev` 分支是基于 `hotfix` 合并后的 `master` ，自然而然的成为了最领先的分支，而且没有 `merge` 的 `commit` 记录，是不是感觉很舒服了。

我们在使用`rebase`合并分支时，也会出现`conflict`，在这种情况下，`git` 会停止 `rebase` 并会让你去解决冲突。在解决完冲突后，用 `git add` 命令去更新这些内容。然后再次执行`git rebase --continue`，这样`git` 会继续应用余下的 `patch` 补丁文件。

假如我们现在不想在执行这次`rebase`操作了，都可以通过`--abort`回到开始前状态：

```shell
git rebase --abort
```



## `rebase`是存在危险的操作 - 慎用

我们现在使用`rebase`操作看起来是完美的，但是他也是存在一定危险的，下面我们就一起来看一看。

现在假设我们在`dev`分支进行开发，执行了`rebase`操作后，在提交代码到远程之前，是这样的：

![]()

提交`dev`分支到远程代码仓库后，就变成了这样：

![]()

而此时你的同事也在 `dev` 上开发，他的分支依然还是以前的`dev`，并没有进行同步`master`：

![]()

那么当他 `pull` 远程 `master` 的时候，就会有丢失提交纪录。这就是为什么我们经常听到有人说 `git rebase` 是一个危险命令，因为它改变了历史，我们应该谨慎使用。

不过，如果你的分支上需要 `rebase` 的所有 `commits` 历史还没有被 `push` 过，就可以安全地使用 `git-rebase`来操作。



## 总结

在`asong`的细心讲解下，姐姐完全搞懂了怎么使用`git rebase`，我们来看一下姐姐的总结：

- 当我们在一个过时的分支上面开发的时候，执行 `rebase` 以此同步 `master` 分支最新变动；
- 假如我们要启动一个放置了很久的并行工作，现在有时间来继续这件事情，很显然这个分支已经落后了。这时候需要在最新的基准上面开始工作，所以 `rebase` 是最合适的选择。
- `git-rebase` 很完美，解决了我们的两个问题：
  - 合并 `commit` 记录，保持分支整洁；
  - 相比 `merge` 来说会减少分支合并的记录；
- 使用`rebase`操作要注意一个问题，如果你的分支上需要 `rebase` 的所有 `commits` 历史还没有被 `push` 过，就可以安全地使用 `git-rebase`来操作。

看来姐姐是真的学会了，那你们呢？

没有学会不要紧，亲自试验一下才能更好的理解呦～～～。

好啦这一篇文章到这里就结束了，我们下期见。

**结尾给大家发一个小福利吧，最近我在看[微服务架构设计模式]这一本书，讲的很好，自己也收集了一本PDF，有需要的小伙可以到自行下载。获取方式：关注公众号：[Golang梦工厂]，后台回复：[微服务]，即可获取。**

**我翻译了一份GIN中文文档，会定期进行维护，有需要的小伙伴后台回复[gin]即可下载。**

**我是asong，一名普普通通的程序猿，让我一起慢慢变强吧。我自己建了一个`golang`交流群，有需要的小伙伴加我`vx`,我拉你入群。欢迎各位的关注，我们下期见~~~**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/wx/qrcode_for_gh_efed4775ba73_258.jpg)

推荐往期文章：

- [手把手教姐姐写消息队列](https://mp.weixin.qq.com/s/0MykGst1e2pgnXXUjojvhQ)

- [详解Context包，看这一篇就够了！！！](https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg)

- [go-ElasticSearch入门看这一篇就够了(一)](https://mp.weixin.qq.com/s/mV2hnfctQuRLRKpPPT9XRw)

- [面试官：go中for-range使用过吗？这几个问题你能解释一下原因吗](https://mp.weixin.qq.com/s/G7z80u83LTgLyfHgzgrd9g)

- [学会wire依赖注入、cron定时任务其实就这么简单！](https://mp.weixin.qq.com/s/qmbCmwZGmqKIZDlNs_a3Vw)

- [听说你还不会jwt和swagger-饭我都不吃了带着实践项目我就来了](https://mp.weixin.qq.com/s/z-PGZE84STccvfkf8ehTgA)
- [掌握这些Go语言特性，你的水平将提高N个档次(二)](https://mp.weixin.qq.com/s/7yyo83SzgQbEB7QWGY7k-w)
- [go实现多人聊天室，在这里你想聊什么都可以的啦！！！](https://mp.weixin.qq.com/s/H7F85CncQNdnPsjvGiemtg)
- [grpc实践-学会grpc就是这么简单](https://mp.weixin.qq.com/s/mOkihZEO7uwEAnnRKGdkLA)
- [go标准库rpc实践](https://mp.weixin.qq.com/s/d0xKVe_Cq1WsUGZxIlU8mw)
- [2020最新Gin框架中文文档 asong又捡起来了英语，用心翻译](https://mp.weixin.qq.com/s/vx8A6EEO2mgEMteUZNzkDg)
- [基于gin的几种热加载方式](https://mp.weixin.qq.com/s/CZvjXp3dimU-2hZlvsLfsw)
- [boss: 这小子还不会使用validator库进行数据校验，开了～～～](https://mp.weixin.qq.com/s?__biz=MzIzMDU0MTA3Nw==&mid=2247483829&idx=1&sn=d7cf4f46ea038a68e74a4bf00bbf64a9&scene=19&token=1606435091&lang=zh_CN#wechat_redirect)