const mongoose = require('mongoose');
const express = require('express');
const mongoEndpoint = "mongodb://mongodb:27017/test"
const port = 80;
const {channelSchema, messageSchema} = require('./schemas');
const Channel = mongoose.model("Channel", channelSchema);
const Message = mongoose.model("Message", messageSchema);

const app = express();
app.use(express.json());

const connect = () => {
    mongoose.connect(mongoEndpoint).then(() => {
        console.log("success connect to mongo");
        }).catch((err) => {
            console.log("error connecting to mongo %v ", err);
        });
}

	//WEBSOCKET /////////////////////////////////////////////
	///////////////////////////////////////
	/////////////////////////////////////////////
	///////////////////////////////////////////////////
//connect to RabbitMQ
let channelInf;
let amqpAddr = process.env.AMQPADDR;

let amqp = require('amqplib/callback_api');


    amqp.connect("amqp://" + amqpAddr + ":5672/", (err, conn) => {
        if (!err && conn) {
            conn.createChannel((err, ch) => {
                
                var q = process.env.MQ;
                channelInf = ch;
                channelInf.assertQueue(q, {durable: false});
            });
        }
    });
	//WEBSOCKET /////////////////////////////////////////////
	///////////////////////////////////////
	/////////////////////////////////////////////
	///////////////////////////////////////////////////

app.get("/v1/channels", async(req, res) => {
    let user = req.get("X-User");
    if (user == null) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channles = await Channel.find({$or:[{members: user}, {private: false}]});
        res.set('Content-Type', 'application/json');
        res.json(channles);
        return
    } catch(e) {
        res.status(500).send("internal error when getting channels" + e);
        return;
    }
});

app.post('/v1/channels', async(req, res, next) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let {name, description, private, members} = req.body;
        if (!name) {
            res.status(400).send("name not provided");
            return;
        }
        res.set('Content-Type', 'application/json');
        if (!private){
            private = false;
        }

        let creator = user;
        if(!members) {
            members = [];
            members.push(user)
        }
        if (!members.includes(user)) {
            members.push(user);
        }

        let createdAt = new Date();
        let editedAt = new Date();
        let channel = {
            name,
            description,
            private,
            members,
            createdAt,
            creator,
            editedAt
        };

        let query = new Channel(channel);
        query.save((err, saveChannel) => {
            if (err) {
                res.status(500).send("fail to create channel");
                return;
            }
// send msg to queue
            if (!private || private === false) {
                members = [];
            }
            let memberID = [];
            for (let i = 0; i < members.length; i++) {
                memberID.push(members[i].id);
            }            
            let queueInfo = {msgType: "channel-new", msg: saveChannel, userIDs: memberID};
            channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
            res.status(201).json(saveChannel); 
        })
    } catch(e) {
        res.status(500).send("internal error when posting channels" + e);
        return;
    }
});

app.get("/v1/channels/:channelID", async(req, res) => {
    if (req.header("X-User") == ""){        
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        
        if (channel.private == true && !channel.members.includes(req.header("X-User"))){
            res.status(403).send("member not in this private channel");
            return;
        } else {
            if (req.query.before) {
                var messages = await Message.find({channelID: req.params.channelID, "_id": {"$lt": req.query.before}}).limit(100).sort({createdAt: -1});
            } else {
                var messages = await Message.find({channelID: req.params.channelID}).limit(100).sort({createdAt: -1});
            }
            res.set('Content-Type', 'application/json');
            res.status(200).json(messages);
        }
    } catch(e) {
        res.status(500).send("internal error when getting channel by id: " + e);
        return;
    }
});

app.post("/v1/channels/:channelID", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        if (channel.private == true && !channel.members.includes(req.header("X-User"))){
            res.status(403).send("member not in this private channel");
            return;
        } else {
            let {channelID, body, createdAt, creator, editedAt} = req.body;
            channelID = req.params.channelID;
            createdAt = new Date();
            body = req.body.body;
            creator = user;
            editedAt = new Date();
            let message = {
                channelID,
                body,
                createdAt,
                creator,
                editedAt
            }
            let query = new Message(message); 
            query.save((err, saveMessage) => {
                if (err) {
                    res.status(500).send(err);
                    return;
                }
// send msg to queue
if (!private || private === false) {
    members = [];
}
let memberID = [];
for (let i = 0; i < members.length; i++) {
    memberID.push(members[i].id);
}            
let queueInfo = {msgType: "message-new", msg: saveMessage, userIDs: memberID};
channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
                res.status(201).json(saveMessage); 
            });
        }
    } catch (err) {
        res.status(500).send("internal error when posting channel by id: " + err);
        return;
    }

});

app.patch("/v1/channels/:channelID", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        if (channel.private == true & channel.creator != req.header("X-User")){
            res.status(403).send("creator not in this private channel");
            return;
        } else {
            const {name ,description} = req.body;
            await Channel.update(channel, {name: name, description: description});
            // send msg to queue
if (!private || private === false) {
    members = [];
}
let memberID = [];
for (let i = 0; i < members.length; i++) {
    memberID.push(members[i].id);
}            
let queueInfo = {msgType: "channel-update", msg: await Channel.findById(req.params.channelID), userIDs: memberID};
channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
            res.set('Content-Type', 'application/json');
            res.json(await Channel.findById(req.params.channelID));
        }
    } catch(err) {
        res.status(500).send(err);
        return;
    }
});

app.delete("/v1/channels/:channelID", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        if (channel.private == true & channel.creator != req.header("X-User")){
            res.status(403).send("creator not in this private channel");
            return;
        } else {
            await Message.deleteMany({"channelID": req.params.channelID});
            await Channel.deleteMany({"_id": req.params.channelID});
                        // send msg to queue
if (!private || private === false) {
    members = [];
}
let memberID = [];
for (let i = 0; i < members.length; i++) {
    memberID.push(members[i].id);
}            
let queueInfo = {msgType: "channel-delete", msg: channelID, userIDs: memberID};
channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
            res.json("channel deleted");
        }
    } catch(err) {
        res.status(500).send(err);
        return;
    }
});

app.post("/v1/channels/:channelID/members", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        if (channel.private == true & channel.creator != req.header("X-User")){
            res.status(403).send("creator not in this private channel");
            return;
        } else {
            const body = req.body;
            await Channel.update({"_id": req.params.channelID}, {"$push": {"members": {id: body.id}}});
            res.status(201).json("user added as channel member");
        }
    } catch(err) {
        res.status(500).send(err);
        return;
    }
});

app.delete("/v1/channels/:channelID/members", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let channel = await Channel.findById(req.params.channelID);
        if (channel.private == true & channel.creator != req.header("X-User")){
            res.status(403).send("creator not in this private channel");
            return;
        } else {
            const body = req.body;
            await Channel.update({"_id": channelID}, {"$pull": { "members": {id: body.id}}});
            res.status(201).json("user deleted from channel member");
        }
    } catch(e) {
        res.status(500).send(e);
        return;
    }
});

app.patch("/v1/messages/:messageID", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let message = await Message.findById(req.params.messageID);
        if(message.creator != req.header("X-User")){
            res.status(403).send("not creator of the message");
            return;
        } else {
            const body = req.body;
            const patchMess = body.body;
            await Message.update({"_id": messageID}, {"body": patchMess});
                                    // send msg to queue
if (!private || private === false) {
    members = [];
}
let memberID = [];
for (let i = 0; i < members.length; i++) {
    memberID.push(members[i].id);
}            
let queueInfo = {msgType: "message-update", msg: await Message.findById(messageID), userIDs: memberID};
channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
            res.json(await Message.findById(messageID));
        }
    } catch(err) {
        res.status(500).send(err);
        return;
    }
});

app.delete("/v1/messages/:messageID", async(req, res) => {
    let user = req.get("X-User");
    if (!user) {
        res.status(401).send("unauthorized"); 
        return;
    }
    try {
        let message = await Message.findById(req.params.messageID);
        if(message.creator != req.header("X-User")){
            res.status(403).send("not creator of the message");
            return;
        } else {
            await Message.deleteOne({"_id": messageID});
// send msg to queue
if (!private || private === false) {
    members = [];
}
let memberID = [];
for (let i = 0; i < members.length; i++) {
    memberID.push(members[i].id);
}            
let queueInfo = {msgType: "message-delete", msg: messageID, userIDs: memberID};
channelInf.sendToQueue(q, Buffer.from(JSON.stringify(queueInfo)));
// send msg to queue
            res.json("message  deleted");
        }
    } catch(err) {
        res.status(500).send(err);
        return;
    }
});

connect();
mongoose.connection.on('error', console.error)
    .on('disconnected', connect)
    .once('open', main)

const name = "general";
const description = "";
const private = false;
const createdAt = new Date();
const channelInit = {
    name,
    description,
    private,
    createdAt
}
const newChannel = new Channel(channelInit);
newChannel.save((err, _) => {
    if (err) {
        console.log(err);
    } 
})

async function main() {
    app.listen(port, "", () => {
        console.log(`server listening $(port)`);
    })
} 