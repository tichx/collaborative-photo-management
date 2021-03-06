const mongoose = require("mongoose")
const express = require("express")
const morgan = require("morgan")
let mysql = require("mysql2")
const {photoSchema, tagSchema} = require('./schemas')
const addr = process.env.ADDR || ":80"
const [host, port] = addr.split(":")
const mongoEndPoint = process.env.MONGOADDR
const mysqlEndPoint = process.env.MYSQLADDR.split(",")
const Photo = mongoose.model("Photo", photoSchema)
const Tag = mongoose.model("Tag", tagSchema)
const content = "Content-Type"
const appjson = "application/json"
const app = express();
app.use(express.json());
app.use(morgan("dev"));

const conn = mysql.createPool({
    host: mysqlEndPoint[0],
    user: mysqlEndPoint[1],
    password: mysqlEndPoint[2],
    database: 'Users',
    insecureAuth: true
});

const connect = () => {
    mongoose.connect(mongoEndPoint, {useNewUrlParser:true});
}
connect()
mongoose.connection.on('error', () => console.log("error connecting"))
const photoSample = {
    url: "http://example.com",
    description: "general image",
    members: [],
    createdAt: new Date()
}
const query = new Photo(photoSample);
query.save()
    .then(newChannel => {
        console.log(newChannel);
    }) 
    .catch(err => {
        console.log(err);
    });

app.listen(port, host, () => {
    console.log(`listening on ${port}`);
})

// helper to connect to mysql db
async function runQuery(query) {
    return new Promise(function (resolve, reject) {
        conn.query(query, async function (err, result, field) {
            if (err) {
                return reject(err)
            }
            resolve(result)
        })
    })
}

// Get all the photos one has access to (could see a photo for being as a member of a tag, 
// or being the creator of the photo).
// Params Usage: 
// 1. /v1/photos/ 
// 2. /v1/photos?year=2020
// 3. /v1/photos?month=12
// 4. /v1/photos?date=/2020/12/07
// Returns a json encode list of photo objects, defualt []
//
// 200: success
// 401: x-user is missing from header or user is not authenticated
// 500: internal errors
app.get("/v1/photos", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    tags =[]
    photos=[]
    try {
        // get all the tags the user is a member of
        tags = await Tag.find().or([{"members.id":userID}, {"creator.id":userID}]).select("id -_id")
        // select the tag id from those tags
        tags = tags.map(tag => tag.id);
        // get photos that associated with the tags, and the photos the user created.
        photos = await Photo.find({$or:[{"tags.id": {$in: tags}}, {"creator.id":userID}]})
    } catch(e) {
        res.status(500).send("Error: no photo found"+e)
        return;
    }
    try {
        res.setHeader(content, appjson)
        
        picSelection = []
        if (req.query.year) {
            // return images created within that yeara
            picSelection = photos.filter(function (el) {
                // console.log(el.createdAt.getYear())
                return el.createdAt.getFullYear() == req.query.year;
              });
        } else if (req.query.month) {
            // return images created on that month in this year
                picSelection = photos.filter(function (el) {
                return el.createdAt.getMonth() + 1 == req.query.month;
              });
        } else if (req.query.date) {
            // return images created on this date; client should format the date like 2020/03/26
            requestedDate = req.query.date.split("/")
            picSelection = photos.filter(function (el) {
                year = requestedDate[0]
                month = requestedDate[1]
                date = requestedDate[2]
                return el.createdAt.getFullYear() == year && el.createdAt.getMonth() + 1 == month && el.createdAt.getDate() == date;
              });

        } else {
            picSelection = photos
        }
        res.json(picSelection)
    } catch (e) {
        res.status(500).send("Error: no photo found"+e)
        return;
    }
});

// Post a new photo url to the user's account. 
// Expected  {"url":"string", "description":"some description"} url is a required field.
// Returns a json encoded object of inserted item.
// 201: success
// 401: x-user is missing from header or user is not authenticated
// 500: internal errors
app.post("/v1/photos/", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    const {url, description} = req.body;
    if (!url) {
        res.status(400).send("Must have the url field")
        return;
    }
    try {
        names=""
        let results = await runQuery("select email from users where id=" + mysql.escape(userID))
        if(results[0] && results[0].email) {
            names = results[0].email
        } else {
            res.status(500).send("Cannot find email")
            return;
        }
    } catch(e) {
        res.status(400).send("Error finding email")
    }
    const query = new Photo({
        "url": url,
        "description": description,
        "createdAt": new Date(),
        "creator": {"id":userID, "email":names}
    });
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send("Error when creating new photo: "+ err.toString());
            return;
        }
        res.status(201).json(newChannel);
        return;
    })
});

// Post: Add a tag to a photo
// Constraints:
// - client must be the creator of the tag or a member of the tag (you can't access tags someone else created, unless you are a member of the tag)
// - the target photo's creator must also be the the member of the tag (you can't tag stranger's photos)
//
// 
// Returns the updated photo object with tag attached to the tag attribute.
// 200: success
// 400: user's email is not found / error occurred when retrieving the tag
// 401: x-user is missing from header or user is not authenticated
// 403: user is not the creator of the photo
// 404: the target photo is not found
// 500: internal errors
app.post("/v1/photos/:photoID/tag/:tagID", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    names=""
    try {
        let results = await runQuery("select email from users where id=" + mysql.escape(userID))
        if(results[0] && results[0].email) {
            names = results[0].email
        } else {
            res.status(400).send("Cannot find email")
            return;
        }
    } catch(e) {
        res.status(400).send("Error finding email")
    }
    photo = await Photo.findOne({"id":req.params.photoID})
    if (!photo) {
        res.status(404).send("photo not found with given photoID")
        return;
    }
    tag = await Tag.findOne({"id":req.params.tagID})
    if (!tag) {
        res.status(404).send("no tag found with given id");
        return;
    }

    try {
        let updates = {"tags":[{"id":tag.id}]}
        Photo.findOneAndUpdate({id:req.params.photoID}, {$push: updates}, {new:true}, function(err, photo){
            if (err) {
                res.status = 400;
                res.send("Error finding channel. "+ err);
                return;
            }
            if(userID != photo.creator.id){
                res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
            }
            // if(names+"" != tag.creator.email+"" || names+"" != tag.members.email+""){
            //     res.status(403);
            //     res.send("Unauthorized: this user is not the tag's creator or member"+tag.creator.email + names);
            //     return;
            // }
            res.setHeader(content, appjson)
            res.send(photo)
            return;
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any photos")
        return;
    }
});


// Get all the tags this user has access to
// returns the encoded json obejcts in a list
// 200: success
// 401: x-user is missing from header or user is not authenticated
// 500: internal errors
app.get("/v1/tags", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    try {
        res.setHeader(content, appjson)
        tags = await Tag.find().or([{"members.id":userID}, {"creator.id":userID}])
        res.json(tags)
    } catch (e) {
        res.status(500).send("Error: no photo found"+e)
        return;
    }
});

// Creates a new tag
// Requires {"name": "some tag name"}
// Returns a json encoded tag object
// 201: success
// 401: x-user is missing from header or user is not authenticated
// 500: internal errors
app.post("/v1/tags", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    const {name} = req.body;
    if (!name) {
        res.status(400).send("Must have the url field")
        return;
    }
    try {
        try {
            names=""
            let results = await runQuery("select email from users where id=" + mysql.escape(userID))
            if(results[0] && results[0].email) {
                names = results[0].email
            } else {
                res.status(400).send("Cannot find email")
                return;
            }
        } catch(e) {
            res.status(400).send("Error finding email")
        }
        const query = new Tag({
            "name": name,
            "members" : {"id":userID, "email":names},
            "createdAt": new Date(),
            "creator": {"id":userID, "email":names}
        });
        query.save((err, tag) => {
            if (err) {
                res.status(400).send("Error when creating new tag: "+ err.toString());
                return;
            }
            res.status(201).json(tag);
            return;
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to create tags")
        return;
    }
});


// delete a tag by tagID
// reutrns a message
// 200: success
// 401: x-user is missing from header
// 403 user is not authenticated
// 500: internal errors
app.delete("/v1/tags/:tagID", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    try {
        const tag = await Tag.findOne({id:req.params.tagID});
        if (tag.creator.id != userID) {
            res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
        } else {
            Tag.findOneAndDelete({id:req.params.tagID}, function(err, tags){
                if (err) {
                    res.status = 400;
                    res.send("Error finding tag"+ err);
                    return;
                }
                // let photos = Photo.find({"tags.id":req.params.tagID})
                // console.log(photos)
                // Photo.deleteMany({tags:req.params.channelID}, function(err, message){
                //     if (err) {
                //         res.status = 400;
                //         res.send("Error:"+ err);
                //         return;
                //     }
                // })
                res.send("Tags successfully deleted")
                return;
            })
        }
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any tags")
        return;
    }
});

// delete a photo by photoID
// 200: success
// 401: x-user is missing from header
// 403: user is not authenticated
// 500: internal errors
app.delete("/v1/photos/:photoID", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    try {
        const photo = await Photo.findOne({id:req.params.photoID});
        if (photo.creator.id != userID) {
            res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
        } else {
            Photo.findOneAndDelete({id:req.params.photoID}, function(err, photos){
                if (err) {
                    res.status = 400;
                    res.send("Error finding tag"+ err);
                    return;
                }
                // Message.deleteMany({channelID:req.params.channelID}, function(err, message){
                //     if (err) {
                //         res.status = 400;
                //         res.send("Error:"+ err);
                //         return;
                //     }
                // })
                res.send("Photo successfully deleted")
                return;
            })
        }
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any photos")
        return;
    }
});

// add a member to a tag
// expects {"id" : 4}, a user id in integer
// returns a message "user added"
//
// 201: success
// 400: user id is not supplied
// 401: x-user is missing from header or user is not authenticated
// 500: internal errors
app.post("/v1/tags/:tagID/members", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    let user = req.body.id
        if (!user) {
            res.status(400).send("user id must be supplied.")
            return;
        }
    let names=""
    try {
        
        let results = await runQuery("select email from users where id=" + mysql.escape(req.body.id))
        if(results[0] && results[0].email) {
            names = results[0].email
        } else {
            res.status(500).send("Cannot find email")
            return;
        }
    } catch(e) {
        res.status(400).send("Error finding email")
    }
        
    try {  
        Tag.findOne({id:req.params.tagID}, function(err, tag){
            if(userID != tag.creator.id){
                res.status(403).send("Unauthorized: this user is not the creator");
                return;
            } else {
                Tag.findOneAndUpdate({id:req.params.tagID}, { $push: {members:{id:req.body.id, email:names}} }, {new:true}, function(err, tag){
                    if (err) {
                        res.status(400).send("Error adding new user "+err);
                        return;
                    } 
                    res.status(201).send("User added");
                    return;
                });
            }
        })        
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any tags")
        return;
    }
});

// delete a user from members of a tag
// expect a {"id":8}, a user id in integer format
// 200: success
// 401: x-user is missing from header
// 403: user is not authenticated
// 500: internal errors
app.delete("/v1/tags/:tagID/members", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    let user = req.body.id
        if (!user) {
            res.status(400).send("user id must be supplied.")
            return;
        }
    try {  
        Tag.findOne({id:req.params.tagID}, function(err, tag){
            if(userID != tag.creator.id){
                res.status(403).send("Unauthorized: this user is not the creator")
                return;
            } else {
                let members = tag.members
                let myArray = members.filter(function( member ) {
                    return member.id !== req.body.id;
                  });
                Tag.findOneAndUpdate({id:req.params.tagID}, ({members:myArray}), {new:true}, function(err, tag){
                    if (err) {
                        res.status(400).send("Error deletings new user "+err)
                        return;
                    } 
                    res.status(200).send("User removed");
                    return;
                });
            }
        })        
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any tags")
        return;
    }
});
