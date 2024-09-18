# Clone Project

To clone this ptoject you can do :
1. copy the url cloning project and than add open you directory you wanna use.
2. after opening the directory open your terminal and clone this project
    ```bash
    git clone url
    ```

# Rule Of Developing

1. Make a new branch with pattern `dev/your-name`
2. and than develop the fiture
3. after finishing the developt you can push you branch like
    ```bash
    git push origin dev/your-name
    ```
   *after you do `git clone` defaultly you'll have a remote named `origin`.
4. if you need other fiture branch, you just need to pull the branch exiting.
    ```bash
    git pull origin branch-name
    ```
5. after pulling the branch you can merge that new branch with yours.
    ```bash
    git merge branch-name
    ```

# Rule of commiting
1. Every commit should only include files from one package.
2. Identify every file you commit:
    
    a. You should only include `new` file with other `new` file in one commit with using flag `"[FEATURE] -> massage"`. 

    b. You should only include `modified` file with other `modified` file in one commit `"[IMPROVE] -> massage"`.
3. in massage commit, you need to identify what are you doin in that file or what are that file you commit doing
