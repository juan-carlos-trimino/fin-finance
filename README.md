# The vs-code-settings repo

VS Code provides different scopes for settings:<br>
**(a) User settings** - These settings are applied globally to any instance of VS Code you open.<br>
**(b) Workspace settings** - These settings are specific to a project and override User settings. VS Code stores workspace settings at the root of the project in a ***.vscode*** folder.<br>

VS Code stores setting values in a ***settings.json*** file:<br>
**(a) User settings.json location** - Depending on your platform, the user settings file is located here:<br>
&nbsp;&nbsp;&nbsp;&nbsp;**(a) Windows** - %AppData%\Code\User\settings.json<br>
&nbsp;&nbsp;&nbsp;&nbsp;**(b) macOS** - $HOME/Library/Application\ Support/Code/User/settings.json<br>
&nbsp;&nbsp;&nbsp;&nbsp;**(c) Linux** - $HOME/.config/Code/User/settings.json<br>
**(b) Workspace settings.json location** - The workspace settings file is located under the ***.vscode*** folder in your root folder. When you add a Workspace Settings ***settings.json*** file to your project or source control, the settings for the project will be shared by all users of that project.<br>

The Settings editor ([***File > Preferences > Settings***] or [***Ctrl + ,***]) is the user interface that enables you to review and modify setting values that are stored in a ***settings.json*** file.

## Using the settings
To use these settings on your project, you may add this repo as a dependency of your repo by using the ***git subtree tool***. This tool allows you to insert this repo as a sub-directory of yours. This is one of several ways Git projects can manage project dependencies.

To see the subtree commands, use:
```
$ git subtree -h
```

To add this repo as a subtree (i.e., as a subfolder) to your repository, run:
```
$ git subtree add --prefix=.vscode/ https://github.com/juan-carlos-trimino/vs-code-settings.git main --squash
```

To fetch content from a remote repo (i.e., this repo) and update the subtree, issue:
```
$ git subtree pull -P .vscode/ https://github.com/juan-carlos-trimino/vs-code-settings.git main --squash
```

If you make a change to the subtree, the commit will be stored in your repo and its logs.
```
$ git commit -am "add a comment"
$ git push origin main
```

If you want to update the remote repo (i.e., this repo) with the commit, run:
```
$ git subtree push -P .vscode https://github.com/juan-carlos-trimino/vs-code-settings.git main
```

To remove the subtree from your repo, execute:
```
$ git rm -r .vscode/
```
(The above command removes the subtree directory from your working directory and stages the removal for commit. The -r flag tells git to recursively delete the directory.)
```
$ git commit -m "remove subtree"
```
(The above command creates a new commit that records the removal of the subtree.)

The subtree is removed from your working directory and future commits, but the history of the subtree remains in your repository. You can still access it using tools like git log or git checkout.
```
$ git push origin main
```
