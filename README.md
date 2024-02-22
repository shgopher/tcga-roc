<!--
 * @Author: shgopher shgopher@gmail.com
 * @Date: 2024-02-21 14:29:58
 * @LastEditors: shgopher shgopher@gmail.com
 * @LastEditTime: 2024-02-22 20:46:48
 * @FilePath: /r1/README.md
 * @Description: 
 * 
 * Copyright (c) 2024 by shgopher, All Rights Reserved. 
-->
# 算法分析
获取患者信息
- 还活着的病人计算其 days_to_last_follow_up 的日期，如果＞365 就是一年生存期 > 365*3 就是三年生存期,大于365*5 就是五年生存期
- 对于死去的病人，计算其 days_to_death

总而言之，现获取病人的生存状态，判断其生死，然后再计算其生存期

然后获取病人信息id，并且将其id寻找其携带基因为 tele1 的基因，因为我们只研究这个基因

得到了下面数据

|id|基因表达程度|生存期|
|:---:|:---:|:---:|
|1|0.3|3|

# 具体算法设计
1. 读取病例的信息， 返回三个字段：病人 case id，生存状态，days_to_last_follow_up,days_to_death

2. 读取基因表达信息，输入两个字段，病人case id，探索基因的类型，返回其基因的表达程度：fpkm_uq_unstranded + 病人id

3.生存期的计算，if 活着，算 last follow，else死了，算 to death，else 报错，再加入一个switch ，> 365 就是一年生存期，大于365*3 就是三年生存期，大于365*5 就是五年生存期，先从5年开始判断，如果不到5年再判断3年，然后再判断一年，所以不能用switch，用if else 

4.生成一组数据，病人id 基因表达程度 生存期


metadate.cart 有 caseid和gdc对应的数据



## 必有内容

- main.go
- go.mod
- go.sum
- index.R
- roc.pdf(后生成的)
- gen（cart）
- a.json （metadata.json）
- b.csv (cinitial.csv)





