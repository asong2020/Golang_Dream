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
$ git rebase -i HEAD~5
```

执行该指令后会自动弹出`vim`编辑模式：

```vim
pick 1195166 add oom demo
pick f4cd3ef update
pick 5803d8f update
pick 2f0b250 update
pick 6e4ad74 add git rebase

# Rebase 2bc0857..6e4ad74 onto 2f0b250 (5 commands)
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

从这里我们可以看出前面5行是我们要合并的记录，