# Insta_api using golang and mongodb
Proper hashing algorithms have been implemented so that user password cannot be reverse engineered.

Features:</br>
User can create a new account</br>
User data is safe as hashing algorithms have been implemented</br>
User can see user details using userid</br>
User can post images with caption</br>
Automatic time stamp will be generated</br>
User can get all details of his posts</br>
User can search for a particular post using its id</br>


</br>
To get started:</br>
clone the repo</br>
Change the current working directory to the location where you cloned the directory</br>
Ensure you have golang installed(For this project, i have used 1.17)</br>
Ensure mongodb is up and running(I have used the default port number 27017)</br>
</br></br>
To run the application, type </br>
# go run app.go

Now the server starts on port 3000</br>
You can access it using Curl or postman

1. Create an User <- URL should be ‘/users'</br>
2. Get a user using id <- URL should be ‘/users/<id here>’
3. Create a Post <- URL should be ‘/posts'
4. Get a post using id <- URL should be ‘/posts/<id here>’
5. List all posts of a user <- URL should be ‘/posts/users/<Id here>'

                    
   
