# lancher
提供xshell和SecureCRT启动参数agent

# dataSync

## FAQ
###`/usr/bin/ld: cannot find -lclntsh`

zip不会创建符号链接，使用`ln libclntsh.so.11.1 libclntsh.so`创建

### `warning: libaio.so.1, needed by libclntsh.so, not found (try using -rpath or -rpath-link)`
```
sudo apt-get install libaio1 libaio-dev
sudo yum install libaio
```

