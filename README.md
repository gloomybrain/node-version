#node-version  

Yet another nodejs version manager, written completely in [Go](https://github.com/golang)  
This projects purpose is to produce something working to better understand the language and it's standart library  
However if you find this project useful any PRs are welcome =)

To start using node-version do the following:
- install the binary somewhere in your path
- add the following lines to your .bash_profile

```bash
export NODE_VERSION_PATH=~/.node-version
export PATH=$PATH:$NODE_VERSION_PATH/default/bin
```
