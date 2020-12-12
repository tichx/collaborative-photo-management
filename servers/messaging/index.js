const mongoose = require("mongoose")
const express = require("express")
const morgan = require("morgan")
let mysql = require("mysql2")
const {channelSchema, messageSchema} = require('./schemas')
const addr = process.env.ADDR || ":80"
const [host, port] = addr.split(":")
const mongoEndPoint = process.env.MONGOADDR
const mysqlEndPoint = process.env.MYSQLADDR.split(",")
const Channel = mongoose.model("Channel", channelSchema)
const Message = mongoose.model("Message", messageSchema)
const content = "Content-Type"
const appjson = "application/json"
const app = express();
app.use(express.json());
app.use(morgan("dev"));

const conn = mysql.createPool({
    host: mysqlEndPoint[0],
    user: mysqlEndPoint[1],
    password: mysqlEndPoint[2],
    database: mysqlEndPoint[3],
    insecureAuth: true
});

const connect = () => {
    mongoose.connect(mongoEndPoint, {useNewUrlParser:true});
}
connect()
mongoose.connection.on('error', () => console.log("error connecting"))
const generalChannel = {
    name: "general",
    description: "general channel",
    private: false,
    members: [],
    createdAt: new Date()
}
const query = new Channel(generalChannel);
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

// connects to db
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


app.get("/v1/channels", async (req, res) => {
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
        channels = await Channel.find().or([{"members.id":userID}, {"creator.id":userID}, {"private":false}, {"private":null}])
        res.json(channels)
    } catch (e) {
        res.status(500).send("Error: no channel found")
        return;
    }
});

app.post("/v1/channels", async (req, res) => {
    if (!("x-user" in req.headers)) {
        res.status(401).send("user is not authenticated");
        return;
    }
    const {userID} = JSON.parse(req.headers['x-user'])
    if (!userID) {
        res.status(401).send("user id missing from x-user");
        return;
    }
    const {name, description, private} = req.body;
    if (!name) {
        res.status(400).send("Must have the name field")
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
    const query = new Channel({
        "name": name,
        "description": description,
        "private": private,
        "members": [{"id":userID, "email":names}],
        "createdAt": new Date(),
        "creator": {"id":userID, "email":names}
    });
    query.save((err, newChannel) => {
        if (err) {
            res.status(400).send("Error when creating new channel: "+ err.toString());
            return;
        }
        res.status(201).json(newChannel);
        return;
    })
});

app.get("/v1/channels/:channelID", async (req, res) => {
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
        const channel = await Channel.findOne({"id":req.params.channelID});
        if(!channel) {
            res.status(404).send("Error: Channel does not exist.")
            return;
        }
        if ( channel.private && !channel['members'].some(user => user.id == userID)) {
            res.status(403).send("User is not authorized to see");
                return;
        }
        let before = req.query.before
        const msgs = await Message.find({channelID: req.params.channelID});
        msgs.sort((a, b) => b.createdAt - a.createdAt);
        if (before) {
            let index = msgs.findIndex(message => message.id == before);
            if (index === -1) {
                index = 0;
            }
            let myArray = []
            for (let i = index; i < msgs.length && i < index + 100; i++) {
                myArray.push(msgs[i]);
            }
            res.status(201).json(myArray);
        } else {
            let myArray = []
            for (let i = 0; i < msgs.length && i < 100; i++) {
                myArray.push(msgs[i]);
            }
            res.status(201).json(myArray);
        }
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.post("/v1/channels/:channelID", async (req, res) => {
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

        const channel = await Channel.findOne({"id":req.params.channelID});
        if(!channel) {
            res.status(404).send("Error: Channel does not exist.")
            return;
        }
        if (!channel['members'].some(user => user.id == userID) && channel.private) {
            res.status(403).send("User is not authorized to see");
                return;
        }
        let channelID = req.params.channelID
        let row = await runQuery("select email from users where id=" + mysql.escape(userID))
        names = ""
        if(row[0] && row[0].email) {
            names = row[0].email
        } else {
            res.status(404).send("Error: Cannot find email")
            return;
        }
        let msg = Message({
            channelID: channelID,
            body: req.body.body,
            createdAt: Date.now(),
            creator: {id:userID, email:names}
        })
        msg.save(function(err) {
            if (err) {
                res.status = 500;
                res.send("Error saving new message " + err);
                return;
            }
            res.status = 201
            res.setHeader(content, appjson)
            res.send(msg)
            return;
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.patch("/v1/channels/:channelID", async (req, res) => {
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
        let updates = {}
        if (req.body.name && req.body.description) {
            updates = {name:req.body.name,description:req.body.description }
        } else if (req.body.name) {
            updates = {name:req.body.name }
        } else if (req.body.description) {
            updates = {description:req.body.description }
        } else {
            res.status(401).send("update object missing or misformed.");
            return;
        }
        Channel.findOneAndUpdate({id:req.params.channelID}, {$set: updates}, {new:true}, function(err, channel){
            if (err) {
                res.status = 400;
                res.send("Error finding channel. "+ err);
                return;
            }
            if(userID != channel.creator.id){
                res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
            }
            res.setHeader(content, appjson)
            res.send(channel)
            return;
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.delete("/v1/channels/:channelID", async (req, res) => {
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
        const ch = await Channel.findOne({id:req.params.channelID});
        const general = await Channel.findOne({"name":"general"});
        if(general.id == req.params.channelID) {
            res.status(403).send("Error: general channel cannot be deleted")
            return;
        } else if (ch.creator.id != userID) {
            res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
        } else {
            Channel.findOneAndDelete({id:req.params.channelID}, function(err, channel){
                if (err) {
                    res.status = 400;
                    res.send("Error finding channel"+ err);
                    return;
                }
                // if(userID != channel.creator.id){
                //     res.status(403);
                //     res.send("Unauthorized: this user is not the creator");
                //     return;
                // }
                Message.deleteMany({channelID:req.params.channelID}, function(err, message){
                    if (err) {
                        res.status = 400;
                        res.send("Error:"+ err);
                        return;
                    }
                })
                res.send("Channel successfully deleted")
                //next();
                return;
            })
        }
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.put("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.copy("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.head("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.options("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.link("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.purge("/v1/channels/:channelID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.get("/v1/channels/:channelID/members", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.put("/v1/channels/:channelID/members", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.patch("/v1/channels/:channelID/members", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.copy("/v1/channels/:channelID/members", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.copy("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.get("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.post("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.put("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.link("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.head("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});
app.options("/v1/messages/:messageID", async (req, res) => {
    res.status(405).send("Method not allowed")
    return;
});

app.post("/v1/channels/:channelID/members", async (req, res) => {
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
        // let newUID = user.id
        // if (!newUID) {
        //     res.status(400).send("user id must be supplied.")
        //     return;
        // }
        
    try {  
        Channel.findOne({id:req.params.channelID}, function(err, channel){
            if(userID != channel.creator.id){
                res.status(403).send("Unauthorized: this user is not the creator");
                return;
            } else {
                Channel.findOneAndUpdate({id:req.params.channelID}, { $push: {members:req.body} }, {new:true}, function(err, channel){
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
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.delete("/v1/channels/:channelID/members", async (req, res) => {
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
        // let newUID = user.id
        // if (!newUID) {
        //     res.status(400).send("user id must be supplied.")
        //     return;
        // }
        
    try {  
        Channel.findOne({id:req.params.channelID}, function(err, channel){
            if(userID != channel.creator.id){
                res.status(403).send("Unauthorized: this user is not the creator")
                return;
            } else {
                let members = channel.members
                let myArray = members.filter(function( member ) {
                    return member.id !== req.body.id;
                  });
                Channel.findOneAndUpdate({id:req.params.channelID}, ({members:myArray}), {new:true}, function(err, channel){
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
        res.status(500).send("Unable to find any channels")
        return;
    }
});

app.patch("/v1/messages/:messageID", async (req, res) => {
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
        
        Message.findOneAndUpdate({id:req.params.messageID}, {$set: {body:req.body.body}}, {new:true}, function(err, message){
            if (err) {
                res.status = 400;
                res.send("Error finding message. "+ err);
                return;
            }
            if (message) {
                if(userID != message.creator.id){
                    res.status(403);
                    res.send("Unauthorized: this user is not the creator");
                    return;
                }
            } else {
                res.status = 400;
                res.send("Error finding message.");
                return;
            }
            res.setHeader(content, appjson)
            res.send(message)
            return;
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any messages")
        return;
    }
});

app.delete("/v1/messages/:messageID", async (req, res) => {
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
        Message.findOneAndDelete({id:req.params.messageID}, function(err, message){
            if (err) {
                res.status = 400;
                res.send("Error finding message"+ err);
                return;
            }
            if(userID != message.creator.id){
                res.status(403);
                res.send("Unauthorized: this user is not the creator");
                return;
            }
            res.send("Message successfully deleted")
        })
    } catch (e) {
        console.log(e);
        res.status(500).send("Unable to find any messages")
        return;
    }
});


