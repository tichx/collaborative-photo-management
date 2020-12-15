import React, { useEffect, useState } from 'react';
import './Styles/MainPage.css';

import ReactDOM from 'react-dom';
import { checkPropTypes } from 'prop-types';

const MainPage = (props) => {
    //data structure for getting images
    const [imgDataList,setImgDataList] = useState([]);
    //data structure for getting tags
    const [tagDataList,setTagDataList] = useState([]);
    //variable to store user imput tag name
    const [tagText, setTagText] = useState("");

    //map to store binding between tag id and object
    const [tagIDTable, setTagIDTable] = useState({});
    const [IDTagTable, setIDTagTable] = useState({});
    //variable to store selected tag
    const [tag,setTag] = useState(-1);
    const handleTagClick = (newTagId)=>{
        setTag(newTagId);
    }

    // let tagIDTable = {};

    // Fetch all img data
    console.log(localStorage.getItem("Authorization"));
    useEffect(() => {
        // console.log(localStorage.getItem("Authorization"));
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
        // console.log(localStorage.getItem("Authorization"));
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


    },[]);

    
    return (
    <div>
        <TagTextField/>
        <TagButtonList tags={tagDataList} setTag={handleTagClick}/>
        <ImgCardList tag={tag} imgDataList={imgDataList} tagIDTable={tagIDTable} IDTagTable={IDTagTable}/>
        {/* <button onClick={BindTagImg} key="12345" IDTagTable={IDTagTable}>Add Tag</button> */}
        {/* <ImgCardList tags={tagDataList} tag={tag} imgDataList={imgDataList}/> */}
    </div>
    );
    
}

//populates the tag buttons into a list
const TagButtonList = (props)=>{
    const TagLists = props.tags.map(data=><TagButton key={data.id} tag={data} setTag={props.setTag}/>);
    
    return (<div style={{display:"flex","flexDirection": "row"}}>{TagLists}</div>)
}

//TODO:
//gets the new tag name typed by user, after user press 'Enter', pass it to /v1/tags
//with payload: {"name": "tagname"}
const TagTextField = () => {
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
                                'Authorization': localStorage.getItem("Authorization"), 
                              }),
                            body: JSON.stringify({ newTag: "Selected" })
                        };
                        fetch("https://api.xutiancheng.me/v1/tags", requestOptions)
                        // .then(resp => resp.json())
                        // .then(console.log(resp));
                        // var xhr = new XMLHttpRequest();
                        // xhr.open('POST', "https://api.xutiancheng.me/v1/tags", true);
                        // xhr.setRequestHeader('Authorization', localStorage.getItem("Authorization"));
                        // xhr.send(JSON.stringify({ name: "tagText" }));
                        // xhr.onreadystatechange = processRequest;
                        // console.log(xhr);
                        // function processRequest(e) {
                        // if (xhr.readyState == 4 && xhr.status == 200) {
                        // var response1 = JSON.parse(xhr.responseText);
                        // }}
                    // empty dependency array means this effect will only run once (like componentDidMount in classes)
                }
              }}
    />;
};

//individual button showing each tag
const TagButton = (props)=>{

    const handleTagClick = (e)=>{
        props.setTag(props.tag.id)
    }
        return <button onClick={handleTagClick}>{props.tag.name}</button>;

    // return <button onClick={props.onClick(props.tag.id)}>props.tag.name</button>;
}

//populates the image cards into a list
const ImgCardList = (props) => {
    var photoListItems;
    if (props.tag === -1) {
        photoListItems = props.imgDataList.map((data) => {
            // return <ImgCard key={data.id} img={data} tags={props.tags}/>

            return <ImgCard key={data.id} img={data} style={{display:"flex","flex-wrap":"wrap","flexDirection": "row"}}  tagIDTable={props.tagIDTable} IDTagTable={props.IDTagTable}/>
        });
    } else {
        photoListItems = props.imgDataList.map((data) => {
            if (data.tags.includes(props.tag)){
                // return <ImgCard key={data.id} img={data} tags={props.tags}/>

                return <ImgCard key={data.id} img={data} style={{display:"flex","flex-wrap":"wrap","flexDirection": "row"}} tagIDTable={props.tagIDTable} IDTagTable={props.IDTagTable}/>
            }
        });
    }
    
    return (photoListItems)
}

//individual image card, still missing function that when clicked,
//lets user add a tag for this photo
const ImgCard = (props) =>{
    //TODO:
    // let displayTags = props.img.tags.map((item, i) => (
    //     <p>{tagIDTable[item.id].name}</p>
    // ));
    

    let tagNameList = [];
    for (let index = 0; index < props.img.tags.length; index++) {
        let tagIDFromImage = props.img.tags[index].id;
        let tagID = props.tagIDTable[tagIDFromImage];
        tagNameList.push(tagID);
    }
        // console.log(cardComponentList);
    let displayResult = tagNameList.map((item, i) => (
          <p key={i} className="font-size-0-8">
            {item}
          </p>)
    );
    
    return(
    <div>
        <img src={props.img.url} style={{width:"150px",height:"150px"}}/>
        {displayResult}
        <BindTagImg imgID={props.img.id} IDTagTable={props.IDTagTable}/>
         {/* <button onClick={BindTagImg} key={props.img.id} IDTagTable={props.IDTagTable}>Add Tag</button> */}
    </div>
     );
    
}

//TODO: bind image to a new tag name input by user and post 
//it to /v1/photos/:photoID/tag/:tagID 
const BindTagImg = (props)=>{
    var imageID = props.imgID;
    // console.log(imageID);
    return <input
    type="text"
    placeholder="bind photo with tag"
    onKeyPress={event => {
                if (event.key === 'Enter') {
                    var tagName = event.target.value;
                    console.log("tagName is: " + tagName);
                    var tagID = props.IDTagTable[tagName];
                    console.log("tagID is: " + tagID);
                    
                        // POST request using fetch inside useEffect React hook
                        const requestOptions = {
                            method: 'post',
                            headers: new Headers({
                                'Authorization': localStorage.getItem("Authorization"), 
                              }),
                        };
                        fetch("https://api.xutiancheng.me/v1/photos/"+imageID+"/tag/"+tagID, requestOptions)
                    
                }
              }}
    />;
}


export default MainPage;