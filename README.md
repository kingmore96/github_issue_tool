# What is it?

 This tool is a command line tool used to create/update/get/lock/unlock an issue in a repository. 

 I use the github api to implement it. 

# Getting Started
 You should use the commandline to start the tool. 

 After you start it,you will see some reminders you need to fill in. 

-  repo and owner: pro-client kingmore96 (this is my exaple) 
-  mode:(c u g l ul q) g (this is my example) 

 If you enter an read-mode g(get), there's no need to fill in the token for Authorization.
 
 if you don't know the token,please see the Question section.
  
 Otherwise you need to provide your token to the tool. 

 When you type read-mode, you will see a reminder like below : 

-  token: XXX (your token) 

 When the tool shows some results to you. You will see another remainder: 

- reuse repo and owner(y n)? 

 if you print y, the tool will reuse the infos you last provided and ask the mode you need this round. 

 if you print n, the tool will begin a new round so you need to provide the new repo and owner info. 

## The Other modes

### the mode you can use in the tool

     c: create 

     u: update 

     g: get 

     l: lock the issue 
     
     ul : unlock the issue

### Write-Mode and Read-Mode 

g is read-mode and the others are write-mod. 

### Create Update Lock Or UnLock?

if you choose c or u mode, there are somethingelse you need to provide. 

when you type c, you will see another remider 

- please enter your issue title: 

if you enter "how to use it?", the tools will do the folowing things: 

1. create a md file named "how to use it?" in the filesystem(use the current path). 

2. open the file with your default editor. 

3. After you finish editting the issue's body and save the file. you need to go to the 

command line and type y to let the tool know you have finish your edit. (after a reminder finished body?) 

when you type u,you need to provide the same thing as you type c 

Besides,you need to provide the issue number to the tool,so it will know the exact issue you need to update 

when you type l,you don't need to provide the body,you only need to provide the issume number,so as to the ul mode.

## Question

### How can I get the Token?

Token shows who you are.

It's simple to get your token, you can visit the github settings page and find Developer Setting.

In it, you can choose generate personal token.
