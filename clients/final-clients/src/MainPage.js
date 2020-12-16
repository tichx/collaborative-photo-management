import React, { useEffect, useState } from 'react';
import './Styles/MainPage.css';
import Alert from 'react-bootstrap/Alert';

const MainPage = (props) => {
    //variable to tell if there's a warning message to show
    const [show, setShow] = useState(false);
    //data structure for getting images
    const [imgDataList,setImgDataList] = useState([]);
    //data structure for getting tags
    const [tagDataList,setTagDataList] = useState([]);
    //map to store binding between tag id and object
    const [tagIDTable, setTagIDTable] = useState({});
    const [IDTagTable, setIDTagTable] = useState({});
    //variable to store selected tag
    const [tag,setTag] = useState(-1);
    const handleTagClick = (newTagId)=>{
        setTag(newTagId);
    }
    const handleShow = () => {
        setShow(true);
    }

    const [notifyTagChange, setNotifyTagChange] = useState(false);

    const NotifyTagUpdate = () =>{
        
        setNotifyTagChange(!notifyTagChange);
        
    }

    // TEST ONLY!!!backdoor code for getting authorization token.
    // console.log(localStorage.getItem("Authorization"));


    // Fetch all img data
    useEffect(() => {
        fetch("https://api.xutiancheng.me/v1/photos",  { 
            method: 'get', 
            headers: new Headers({
                'Authorization': localStorage.getItem("Authorization"), 
            })
          })
        .then(resp => resp.json())
        .then(data => {setImgDataList(data)})
    },[]);
    // Fetch all tag data
    useEffect(() => {
        fetch("https://api.xutiancheng.me/v1/tags",  { 
            method: 'get', 
            headers: new Headers({
                  'Authorization': localStorage.getItem("Authorization"), 
            })
          })
        .then(resp => resp.json())
        .then(data => 
        {   
            for (let index = 0; index < data.length; index++) {
                const element = data[index];
                let obj={};
                obj[element.id] = element.name;
                Object.assign(tagIDTable, obj);
    
                let obj2={};
                obj2[element.name] = element.id;
                Object.assign(IDTagTable, obj2);
            }
            setTagDataList(data);
        })


    },[notifyTagChange,show]);

    if (show) {
        return (
            <div>
            <TagTextField NotifyTagUpdate={NotifyTagUpdate} handleShow={handleShow}/>
            <TagButtonList tags={tagDataList} setTag={handleTagClick} NotifyTagUpdate={NotifyTagUpdate} handleShow={handleShow}/>
            <ImgCardList tag={tag} imgDataList={imgDataList} tagIDTable={tagIDTable} IDTagTable={IDTagTable} NotifyTagUpdate={NotifyTagUpdate} handleShow={handleShow}/>
            {/* <button onClick={BindTagImg} key="12345" IDTagTable={IDTagTable}>Add Tag</button> */}
            {/* <ImgCardList tags={tagDataList} tag={tag} imgDataList={imgDataList}/> */}
            <Alert variant="danger" onClose={() => setShow(false)} dismissible>
            <Alert.Heading>error for this request</Alert.Heading>
            <p>
              make sure you have a valid auth session, the image upload isn't 
              duplicate, also the tag is case sensitive. If it's not your problem
              then maybe there's an internal error.
            </p>
          </Alert>
        </div>

        );
    } else {
        return (
            <div>
                <TagTextField NotifyTagUpdate={NotifyTagUpdate}/>
                <TagButtonList tags={tagDataList} setTag={handleTagClick} NotifyTagUpdate={props.NotifyTagUpdate}/>
                <ImgCardList tag={tag} imgDataList={imgDataList} tagIDTable={tagIDTable} IDTagTable={IDTagTable} NotifyTagUpdate={NotifyTagUpdate}/>
                {/* <button onClick={BindTagImg} key="12345" IDTagTable={IDTagTable}>Add Tag</button> */}
                {/* <ImgCardList tags={tagDataList} tag={tag} imgDataList={imgDataList}/> */}
            </div>
            );
    }

    
}

//populates the tag buttons into a list
const TagButtonList = (props)=>{
    const TagLists = props.tags.map(data=><TagButton key={data.id} tag={data} setTag={props.setTag} NotifyTagUpdate={props.NotifyTagUpdate} handleShow={props.handleShow}/>);
    
    return (<div style={{display:"flex","flexDirection": "row"}}>{TagLists}</div>)
}

//gets the new tag name typed by user, after user press 'Enter', pass it to /v1/tags
//with payload: {"name": "tagname"}
const TagTextField = (props) => {
    return <input
    type="text"
    placeholder="create new tag"
    onKeyPress={event => {
                if (event.key === 'Enter') {
                    var newTag = event.target.value;
                        // POST request using fetch inside useEffect React hook
                        const requestOptions = {
                            method: 'post',
                            headers: new Headers({
                                'Content-Type': 'application/json',
                                'Authorization': localStorage.getItem("Authorization"), 
                              }),
                            body: JSON.stringify({"name": newTag})
                        };

                        fetch("https://api.xutiancheng.me/v1/tags", requestOptions)
                        .then(resp => resp.json())
                        .then(resp =>{
                            if (resp.status < 300){
                                setTimeout(() => {
                                    props.NotifyTagUpdate();
                                }, 1000)
                            }
                            //  else{
                            //         props.handleShow();
                            
                            // }
                        });
                        

                }
              }}
    />;

};

//individual button showing each tag
const TagButton = (props)=>{

    const handleTagClick = (e)=>{
        props.setTag(props.tag.id)
    }
        return <div style={{display:"flex",flexWrap:"wrap",flexDirection: "row", margin:"15px"}}>
            <button onClick={handleTagClick}>{props.tag.name}</button>
            <BindTagMember tagID = {props.tag.id} NotifyTagUpdate={props.NotifyTagUpdate} handleShow={props.handleShow}/>
        </div>
}

//populates the image cards into a list
const ImgCardList = (props) => {
    var photoListItems;
    if (props.tag === -1) {
        photoListItems = props.imgDataList.map((data) => {

            return <ImgCard key={data.id} img={data} style={{display:"flex","flex-wrap":"wrap","flexDirection": "row"}}  tagIDTable={props.tagIDTable} IDTagTable={props.IDTagTable} NotifyTagUpdate={props.NotifyTagUpdate} handleShow={props.handleShow}/>
        });
    } else {

        photoListItems = props.imgDataList.map((data) => {
            var tagIds = []
            data.tags.forEach(e => {
                tagIds.push(e.id)
            });
            if (tagIds.includes(props.tag)){
                // return <ImgCard key={data.id} img={data} tags={props.tags}/>

                return <ImgCard key={data.id} img={data} style={{display:"flex","flex-wrap":"wrap","flexDirection": "row"}} tagIDTable={props.tagIDTable} IDTagTable={props.IDTagTable} NotifyTagUpdate={props.NotifyTagUpdate} handleShow={props.handleShow}/>
            }
        });
    }
    
    return (photoListItems)
}

//individual image card, still missing function that when clicked,
//lets user add a tag for this photo
const ImgCard = (props) =>{
    
    let tagNameList = [];
    for (let index = 0; index < props.img.tags.length; index++) {
        let tagIDFromImage = props.img.tags[index].id;
        let tagID = props.tagIDTable[tagIDFromImage];
        tagNameList.push(tagID);
    }
    let displayResult = tagNameList.map((item, i) => (
          <p key={i} className="font-size-0-8">
            {item}
          </p>)
    );
    
    return(
    <div>
        <img src={props.img.url} alt="wrong image url" style={{width:"150px",height:"150px"}}/>
        {displayResult}
        <BindTagImg imgID={props.img.id} IDTagTable={props.IDTagTable} NotifyTagUpdate={props.NotifyTagUpdate} handleShow={props.handleShow}/>
         {/* <button onClick={BindTagImg} key={props.img.id} IDTagTable={props.IDTagTable}>Add Tag</button> */}
    </div>
     );
    
}

//bind image to a new tag name input by user and post 
//it to /v1/photos/:photoID/tag/:tagID 
const BindTagImg = (props)=>{
    var imageID = props.imgID;
    return <input
    type="text"
    placeholder="bind photo with old tag"
    onKeyPress={event => {
                if (event.key === 'Enter') {
                    var tagName = event.target.value;
                    var tagID = props.IDTagTable[tagName];
                    
                        // POST request using fetch inside useEffect React hook
                        const requestOptions = {
                            method: 'post',
                            headers: new Headers({
                                'Authorization': localStorage.getItem("Authorization"), 
                              }),
                        };
                        fetch("https://api.xutiancheng.me/v1/photos/"+imageID+"/tag/"+tagID, requestOptions)
                        .then(resp =>{
                            if (resp.status < 300){
                                setTimeout(() => {
                                    props.NotifyTagUpdate();
                                }, 1000)
                            } 
                            // else{
                            //     props.handleShow();
                            // }
                        });
                }
              }}
    />;
}

//bind tag to a user id input by user and post 
//it
const BindTagMember = (props)=>{
    var tagID = props.tagID;
    return <input
    type="text"
    placeholder="userID"
    style={{width:"50px"}}
    onKeyPress={event => {
                if (event.key === 'Enter') {
                    var userID = event.target.value;
                    
                        // POST request using fetch inside useEffect React hook
                        const requestOptions = {
                            method: 'post',
                            headers: new Headers({
                                'Content-Type': 'application/json',
                                'Authorization': localStorage.getItem("Authorization"), 
                              }),
                            body: JSON.stringify({"id": userID})
                        };
                        fetch("https://api.xutiancheng.me/v1/tags/"+tagID+"/members", requestOptions)
                        .then(resp =>{
                            if (resp.status < 300){
                                setTimeout(() => {
                                    props.NotifyTagUpdate();
                                }, 1000)
                            } 
                            // else{
                            //     props.handleShow();
                            // }
                        });

                }
              }}
    />;
}


export default MainPage;