### Part 1

I used the python's built-in webserver, as it was the simplest thing to do: </br>
Dockerfile

``` Dockerfile
FROM python:3-slim

WORKDIR /app

# Create a simple HTML file
RUN echo "Hello World" > index.html

# Expose port 8080
EXPOSE 8080

# Run a simple HTTP server on port 8080
CMD ["python3", "-m", "http.server", "8080"]
```

docker build -t hello-world . </br>
docker run -d -p 8080:8080 hello-world

Verify </br>
curl 127.0.0.1:8080 </br>
From elsewhere on the network, it will be available on the ipv4/6 address of the host machine. Assuming the container is on the host/bridge network. </br>
From my home server I can curl my laptop: </br>
curl 192.168.1.82:8080 </br>
Hello World


I have used Flask before to serve a webhook, in this case, there is a "Hello world" example already done on [Flask website](https://flask.palletsprojects.com/en/stable/quickstart/#a-minimal-application)
or [on github](https://github.com/docker/awesome-compose/blob/master/flask/app/Dockerfile).
</br>My approach seemed simpler, however.

If you need some more nginx sites, I run a few locally, which you can access:
[My vault](https://vault.dgekwbxjjww.online) </br>
[My Homeassistant](https://ha.dgekwbxjjww.online/) </br>
Or my Audiobook instance, which uses OpenID connect: [Audiobookshelf](https://audiobooks.dgekwbxjjww.online/audiobookshelf/login) </br>
All of these are based on docker (compose).


### Parts 2 and 3 
I will not attempt, I do not have that much experience with k8s/cloud environments. Rather than turn this into generated AI slop.


### Part 4
Task list this complex does not make sense to be a playbook, but rather should be a role.  

My role structure:

``` bash
roles/
    common/               # this hierarchy represents a "role"
        tasks/            #
            main.yml      #  <-- tasks file can include smaller files if warranted
        handlers/         #
            main.yml      #  <-- handlers file
        templates/        #  <-- files for use with the template resource
            ntp.conf.j2   #  <------- templates end in .j2
        files/            #
            bar.txt       #  <-- files for use with the copy resource
            foo.sh        #  <-- script files for use with the script resource
        vars/             #
            main.yml      #  <-- variables associated with this role
        defaults/         #
            main.yml      #  <-- default lower priority variables for this role

```
This is the "default" recommended structure based on [ansible documentation](https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_reuse_roles.html) I trimmed it of things that are not neccessary for a simple role like this.


tasks/main.yml - In this case, we have two distinct server groups; RHEL based ones and ubuntu. Therefore I will devide the group-specific tasks, while keeping main.yml with tasks for both (all) groups:

``` yaml


```