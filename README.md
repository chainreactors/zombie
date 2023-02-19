# Zombie 

一个轻量级的服务口令爆破工具, 继承了hydra的命令行设计, hashcat的字典生成, 以及红队向的功能设计. 

## QuickStart

参考了hydra的命令行设计, 小写为命令行输出, 大写为文件输入, 留空为使用默认值.

使用默认字典爆破ssh口令

zombie -I targets.txt -u root -s ssh

打开debug,判断是否存在漏报

zombie -I targets.txt -u root -p password123 -s ssh --debug