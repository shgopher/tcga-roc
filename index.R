#install.packages("survival")
#install.packages("survminer")
#install.packages("timeROC")


library(survival)
library(survminer)
library(timeROC)

setwd("/Users/shgopher/Desktop/1/r/r1")

inputFile="data.csv"



# 使用read.table函数从指定的输入文件中读取数据，并将数据存储在名为rt的数据框中。

# header = TRUE表示第一行包含列名，sep = "\t"表示使用制表符作为列之间的分隔符，

# check.names = FALSE表示不检查列名的有效性，row.names = 1表示使用第一列作为行名。

# colnames(rt)[3]将第三列的列名存储在变量gene中。
rt=read.table(inputFile, header=T, sep=",", check.names=F, row.names=1)
gene=colnames(rt)[3]

head(rt)
head(gene)



ROC_rt=timeROC(T=rt$futime, delta=rt$fustat,
           marker=rt[,gene], cause=1,
           weighting='aalen',
           times=c(1,3,5), ROC=TRUE)

head(ROC_rt)

pdf(file="ROC.pdf", width=5, height=5)


plot(ROC_rt,time=1,col='green',title=FALSE,lwd=2)
plot(ROC_rt,time=3,col='blue',add=TRUE,title=FALSE,lwd=2)
plot(ROC_rt,time=5,col='red',add=TRUE,title=FALSE,lwd=2)

legend('bottomright',
   c(paste0('AUC at 1 years: ',sprintf("%.03f",ROC_rt$AUC[1])),
     paste0('AUC at 3 years: ',sprintf("%.03f",ROC_rt$AUC[2])),
     paste0('AUC at 5 years: ',sprintf("%.03f",ROC_rt$AUC[3]))),
     col=c("green",'blue','red'),lwd=2,bty = 'n')

dev.off()

