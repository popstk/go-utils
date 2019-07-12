# lancher
提供xshell和SecureCRT启动参数agent

# dataSync
同步两个数据库的同一个表的数据，假定这种表都有一个自增且唯一的字段id，并且不会复用id和修改之前的行数据。
同步程序先查询两个表的最大id值，然后按照maxLines步进同步行数，批量插入目标数据库表。


## FAQ
###`/usr/bin/ld: cannot find -lclntsh`

zip不会创建符号链接，使用`ln libclntsh.so.11.1 libclntsh.so`创建

### `warning: libaio.so.1, needed by libclntsh.so, not found (try using -rpath or -rpath-link)`
```
sudo apt-get install libaio1 libaio-dev
sudo yum install libaio
```

