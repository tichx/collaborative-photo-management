

# Collaborative Photo Management
Group Members: Nicholas Xu, Danfeng Yang, Sailesh Sivakumar, Martin Zhang
## Project Description
We are proposing to build a _Collaborative Photo Management_ website for businesses/teams. Today millions of organizations and teams are working collaboratively online. Images assets are mission critical for many different business functions, like marketing, product development, disease diagnosis, photo editing, etc. However, there is no software that enables teams to collaboratively import, manage, consume photo media effectively.  So we believe any team that needs image-focused collaboration will be our target audience.

For example, a startup could be our client. Say they have some photographer vendors to shoot a product launch event. They want to collectively select photos for editing.  After editing they want to select the best ones to send to the media and to post on their social media for marketing purposes. Currently they might use Google Drive/Dropbox to share photos, but google Drive doesn’t allow you to tag or group photos, or browse them easily with gallery view, or view them by month/years, or to publish them publicly so that anyone can view the gallery. 

As developers, our team has someone who’s a photographer working for UW student organizations. He often frets about archiving and finding event photos that are stored in Google Drive. It’s very difficult to find a specific photo of a certain event because of the lack of tagging/grouping or image search features. It is also immensely ineffective to browse many photos together because Google Drive is essentially built for documents collaboration. So as a team, we find motivation to make an Image Collaboration Tool for teams so that in case of events with photo media, any teams would find this useful to organize and manage their image assets.

## Technical Descriptioin

### Architecture Diagram

![graph.png](/static/graph.png)  

#### User Stories

| Priority | User | Description | Technical Implementation Strategy |
| --- | --- | --- | --- |
|P0| As an usesr| I want to be able to create new account, log in, log out, change name of  my account | Store log in credentials in database, authenticate user and store sessions inside Redis database. |
|P0|As a user|I want to be able to upload  photos from a local device.|Use an endpoint to insert new photos and photos will be stored on AWS S3 buckets.|
|P0|As a contributor|I want to tag photos and filter/search them by tag name|Will have to store parameters of the pictures and use a framework to show a processed picture in the client. |
|P0|As a contributor|I want to be able to glance photo gallery by tag,  by day, month, year to quickly find useful photos|Store the tags, etc. to the database. Use and endpoint to get the gallery then implement a filter to filter out tags, etc.|
|P1|As a contributor|Stretch goal: I want to add title and description to photos|Will have the image struct be defined with fields to hold title and description information. And associate the images with users.|
|P1|As a contributor|Stretch goal: I want to discuss with other contributors live as browsing through the photos|Will implement a chat feature through websockets to enable client side messaging. Authenticated users can find other users by email to chat and collaborate|
|P2|As a contributor|Stretch goal: I want to share the photos with other contributors|The user can use an endpoint to get to the resources (maybe by album id/photo id)|

### Endpoints
#### Sessions
- /v1/session: begin a new session for the user
    - POST
        - 201 application/json: created a new user sessions
        - 403: invalid username
        - 415: unsupported media
        - 500: internal server error
- /v1/sessions/mine: getting the current session
    - GET
        - 400: bad request
        - 403: forbidden request if not authenticated
    - DELETE
        - 400: bad request
        - 403: forbidden request if not user's session
#### Users
- /v1/user : registering new users
    - POST; application/json: Create new user account with email, username, password, first name, last name.
        - 201; application/json: Successfully gets the information. Return the user profile in json. 
        - 404: Cannot find the user
        - 500: Internal server error.
- /v1/user/{UserID | me} : getting existing user profiles
    - GET; application/json: get the user profiles, name, avatar, number of albums, number of pictures, etc.
        - 201; application/json: Successfully gets the information. Return the user profile in json.
        - 404: Cannot find the user
        - 500: Internal server error.
#### Photo Management
- /v1/photo: import photos
    - POST; application/json: Import a new photo under the user’s account.
        - 201; application/json: Successfully adds the photo to the gallery; returns a json object of the imported photo information.
        - 400: Some parameters (unsupported format etc.) are wrong.
        - 500: Internal server error.
- /v1/photo?{year|month|datetime}=
    - GET; application/json: get all the photos, then filter by year inside the sql query
        - 201; application/json: Successfully gets the information of all photos, returns an array of json objects of the photo information. 
        - 400: wrong format for year, etc.
        - 500: Internal server error.
#### Tags
- /v1/photo?tag= gets the PhotoID[] from the Tag table with this tag name(or tag id)
    - GET; application/json: get all the photos, then filter by tags
        - 201; application/json: Successfully gets the information of all photos, returns an array of json objects of the photo information. 
        - 404: Cannot find any photo for this tag
        - 500: Internal server error.
- /v1/tags: create a new tag with a selection of photos
    - POST: create a new tag with a name
        - 200  tagged
        - 422 : invalid user input
        - 415: unsupported media
        - 500: internal server error
- /v1/tags/: gets all the tags created by the user
    - GET; application/json: return all tags with a specific userid

### Appendix
#### Database Schema
```
# Photos
Photo: {
    PhotoID int
    Title string
    Note string
    Tag String
    DateCreated datetime
    LastModifiedTime datetime
    Link string
    PhotoURL string
}

# Users
User: {
    UserID int
    Username string
    PassHash string
    FirstName string
    LastName string
    Album AlbumID[]
    Tag TagID[]
    Email string
}
 
# Albums
Album: {
    AlbumID int
    AlbumName string
    DateCreated datetime
    Photos PhotoID[]
}

#Tags
Tags: {
    TagID int
    TagName string
    DateCreated datetime
    Photos PhotoID[]
    UserID int
}
```
