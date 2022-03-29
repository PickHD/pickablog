# PickABlog
**RESTful API system blog using Golang Fiber**

## Setup Development
1. Install/Update/Delete packages, run :
    ```
    $ make deps
    ```
2. Creating docker container postgresql, run :
    ```
    $ make postup
    ```
3. Create new databases, run : 
    ```
    $ make dbup
    ```
4. Migrate table, run :
    ```
    $ make migrate up
    ```
5. Finally run _based on env_ :
    ```
    $ make run
    ```
