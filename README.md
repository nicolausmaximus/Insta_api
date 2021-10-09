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
Proper pagination algorithms have been implemented to keep a check on the user data(Internet consumption while using the system). 
Used mutex to make the server thread safe
Unit testing has been deployed to test the app.


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

                    
User creation</br>
![image](https://user-images.githubusercontent.com/63350417/136671832-e5faadfb-3594-4308-a48f-c2d332a21a2f.png)
</br>
Updated Database</br>
![image](https://user-images.githubusercontent.com/63350417/136671852-894e7236-2b5a-461a-bde2-25d8dbf4a3a5.png)
As we can see, the password is properly hashed</br>
                    
Get user details using userid</br>
![image](https://user-images.githubusercontent.com/63350417/136671882-fb0c7e27-db75-4a0c-8183-4eaa9522fe67.png)


![image](https://user-images.githubusercontent.com/63350417/136671285-ce08f28e-c147-45fe-837b-e6bda56d5aae.png)</br>
![image](https://user-images.githubusercontent.com/63350417/136671290-41e033f9-8846-4e6b-8b6a-7112a596ca33.png)</br>

get post according to id
![image](https://user-images.githubusercontent.com/63350417/136671952-6598da4a-e4cf-4cfa-af14-fb903dae7360.png)

Used mutex to make the server thread safe

                  
