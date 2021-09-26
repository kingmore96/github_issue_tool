# what is it?

*This tool is a command line tool used to create/update/get/lock/unlock an issue in a repository.*

*I use the github api to implement it.*

# how to use it ?
*You should use the commandline to start the tool.*

*After you start it,you will see some reminders you need to fill in.*

- *repo and owner: pro-client kingmore96 (this is my exaple)*
- *mode:(c u g l ul q) g (this is my example)*

*if you enter an read-mode g(get), there's no need to fill in the token for Authorization*

*otherwise you need to provide your token to the tool.*

*So you will see another reminder like below :*

- *token: XXX (your token)*

*When it show some results to you. You will see another remainder:*

- reuse *repo and owner?*

*if you print y, the tool will reuse the infos you last provided and ask the mode you need this round.*

*if you print n, the tool will begin a new round so you need to provide the new repo and owner info.*

## about mode

1. *the mode you can use in the tool.*

    *c: create*

    *u: update*

    *g: get*

    *l: lock the issue*

2. *write-mode and read-mode*

    *g is read-mode and the other 3 is write-mod.*

3. *if you choose c or u mode, there are somethingelse you need to provide.*

    *when you type c, you will see another remider*

    *please enter your issue title:*

    *if you entered "how to use it?", the tools will do the folowing things:*

    *1. create a md file named "how to use it?" in the filesystem(use the ongoing path).*

    *2. open the file with your default editor.*

    *3. After you finish editting the issue's body and save the file. you need to go to the*

    *command line and type yes to let the tool know you have finish your edit. (after a reminder finished body?)*

    *when you type u,you need to provide the same thing as you type c*

    *Besides,you need to provide the issue number to the tool,so it will know the exact issue you need to update*

    *when you type l,you don't need to provide the body,you only need to provide the issume number*
