import {useParams, Link} from "react-router-dom"
import {useState, useEffect} from "react"
import type {Post, Comment} from "../Api"
import {getPostById, getPostComments} from "../Api"
import PostCard from "../components/PostCard"
import calcDuration from "../functions/calcDuration"
import CommentIcon from '../assets/icons/CommentIcon'
import UpvoteIcon from '../assets/icons/UpvoteIcon'
import DownvoteIcon from '../assets/icons/DownvoteIcon'

import "../css/Forum.css"
export default function PostPage(){
    const {id} = useParams();
    const [post, setPost] = useState<Post|null> (null)
    const [comments, setComments] = useState<Comment[]|null> (null)

    useEffect(() => {
        async function loadPost(){
            try {
                const data = await getPostById(Number(id))
                setPost(data)
            } catch (err){
                console.log("failed retrieving posts", err)
            }
        }
        loadPost();
    }, []);

    useEffect(() =>{
        async function loadComments(){
            if (!post) return
            try{
                const data = await getPostComments(post.post_id)
                setComments(data)
            } catch(err){
                console.log("failed to retrieve comments", err)
            }
        }
        
        loadComments()
        console.log("loadComments attempted")
        console.log(comments)
    },[post]);
    

    if (post) return (
        <div style={{display:"flex", flexDirection:"column", margin:"1%", alignItems:"center",}}>
        
        <PostCard post={post}/>

        <div style={{marginBottom:"10px"}}/>

        {comments && comments.map(comment => (
            <div style={{ marginBottom: '15px', width:"35vw", border: "1px solid grey"}}>
                <div style={{display:"flex", flexDirection:"row", marginBottom:"5px", gap:"5px"}}>
                    <p style={{margin:0,marginLeft:"3px", color:"black", fontSize: "13px"}}>{comment.username}</p>
                    <p style={{margin:0, fontSize: "12px", color: "gray", marginLeft:"8px"}}> {calcDuration(comment.created_at)}</p>
                </div>

                <p className="text-14-scalable">{comment.text_body}</p> 

                <div style={{display:"flex", flexDirection:"row", alignSelf: "flex-start", alignItems:"center", fontSize: "12px",  color: "gray", height:"20px", gap:"8px",  marginTop:"2px", marginLeft:"5px" }}>
                    <div style={{display:"flex", flexDirection:"row",alignItems:"center", gap:"2px"}}>
                            <p style={{fontSize: "13px",padding:"0", margin:"0"}}>0</p>
                            <UpvoteIcon/>
                        </div>
                        
                        <div style={{display:"flex", flexDirection:"row",alignItems:"center", gap:"2px"}}>
                            <p style={{fontSize: "13px",padding:"0", margin:"0"}}>0</p>
                            <DownvoteIcon/>
                        </div>

                        <div style={{display:"flex", flexDirection:"row", alignItems:"center", gap:"2px"}}>
                            <p style={{fontSize: "13px",padding:"0", margin:"0"}}>0</p>
                            <CommentIcon/>
                    </div>
                </div>

            </div>
        ))}

        {!comments && 
            <div>
                <h3> Post has no comments! </h3>

            </div>}   

        </div>
    )




    //rendered if failed to retrieve post
    return (
        <div style={{display:"flex", flexDirection:"column", height:"100vh", justifyContent:"center", alignContent:"center"}}>
            <h1> Page not found 🥲 </h1>
   
            <Link to="/" style={{fontSize:"20px", marginTop:"20px"}}>
                    Back to home
                </Link>
        </div>
    )
}